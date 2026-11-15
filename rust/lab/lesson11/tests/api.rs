use axum::Router;
use axum::body::Body;
use axum::http::{Method, Request, StatusCode, header};
use http_body_util::BodyExt;
use serde_json::{Value, json};
use tower::ServiceExt;

async fn send(app: &Router, method: Method, uri: &str, body: Option<&str>) -> (StatusCode, Value) {
    let mut builder = Request::builder().method(method).uri(uri);
    let body = match body {
        Some(text) => {
            builder = builder.header(header::CONTENT_TYPE, "application/json");
            Body::from(text.to_string())
        }
        None => Body::empty(),
    };
    let response = app
        .clone()
        .oneshot(builder.body(body).unwrap())
        .await
        .unwrap();
    let status = response.status();
    let bytes = response.into_body().collect().await.unwrap().to_bytes();
    let value = serde_json::from_slice(&bytes)
        .unwrap_or_else(|_| Value::String(String::from_utf8_lossy(&bytes).into_owned()));
    (status, value)
}

#[tokio::test]
async fn health_returns_ok() {
    let app = lesson11::app();
    let (status, body) = send(&app, Method::GET, "/health", None).await;
    assert_eq!(status, StatusCode::OK);
    assert_eq!(body, json!({ "status": "ok" }));
}

#[tokio::test]
async fn create_then_get_task() {
    let app = lesson11::app();

    let (status, created) = send(
        &app,
        Method::POST,
        "/tasks",
        Some(r#"{"title": "  axum 공부하기  "}"#),
    )
    .await;
    assert_eq!(status, StatusCode::CREATED);
    assert_eq!(created["id"], 1);
    assert_eq!(created["title"], "axum 공부하기");
    assert_eq!(created["done"], false);
    let created_at = created["created_at"].as_str().unwrap();
    assert!(created_at.contains('T') && created_at.ends_with('Z'));

    let (status, fetched) = send(&app, Method::GET, "/tasks/1", None).await;
    assert_eq!(status, StatusCode::OK);
    assert_eq!(fetched, created);
}

#[tokio::test]
async fn create_sets_location_header() {
    let app = lesson11::app();
    let request = Request::post("/tasks")
        .header(header::CONTENT_TYPE, "application/json")
        .body(Body::from(r#"{"title": "위치 확인"}"#))
        .unwrap();
    let response = app.oneshot(request).await.unwrap();
    assert_eq!(response.status(), StatusCode::CREATED);
    assert_eq!(response.headers()[header::LOCATION], "/tasks/1");
}

#[tokio::test]
async fn create_rejects_empty_title() {
    let app = lesson11::app();
    let expected = json!({ "error": "title은 비어 있을 수 없습니다" });

    for body in [r#"{"title": "   "}"#, "{}", r#"{"title": null}"#] {
        let (status, error) = send(&app, Method::POST, "/tasks", Some(body)).await;
        assert_eq!(status, StatusCode::BAD_REQUEST, "본문: {body}");
        assert_eq!(error, expected, "본문: {body}");
    }
}

#[tokio::test]
async fn create_limits_title_to_100_characters() {
    let app = lesson11::app();

    let exactly_100 = "가".repeat(100);
    let body = json!({ "title": exactly_100 }).to_string();
    let (status, _) = send(&app, Method::POST, "/tasks", Some(&body)).await;
    assert_eq!(status, StatusCode::CREATED);

    let body = json!({ "title": "가".repeat(101) }).to_string();
    let (status, error) = send(&app, Method::POST, "/tasks", Some(&body)).await;
    assert_eq!(status, StatusCode::BAD_REQUEST);
    assert_eq!(error, json!({ "error": "title은 100자 이하여야 합니다" }));
}

#[tokio::test]
async fn create_rejects_broken_json_and_unknown_fields() {
    let app = lesson11::app();

    let (status, error) = send(&app, Method::POST, "/tasks", Some(r#"{"title": "#)).await;
    assert_eq!(status, StatusCode::BAD_REQUEST);
    let message = error["error"].as_str().unwrap();
    assert!(message.starts_with("요청 본문을 읽을 수 없습니다: "));

    let (status, error) = send(&app, Method::POST, "/tasks", Some(r#"{"titel": "x"}"#)).await;
    assert_eq!(status, StatusCode::BAD_REQUEST);
    let message = error["error"].as_str().unwrap();
    assert!(message.starts_with("요청 본문을 읽을 수 없습니다: "));
    assert!(message.contains("unknown field `titel`"));
}

#[tokio::test]
async fn missing_content_type_returns_415() {
    let app = lesson11::app();
    let request = Request::post("/tasks")
        .body(Body::from(r#"{"title": "x"}"#))
        .unwrap();
    let response = app.oneshot(request).await.unwrap();
    assert_eq!(response.status(), StatusCode::UNSUPPORTED_MEDIA_TYPE);
}

#[tokio::test]
async fn list_is_empty_array_then_ordered_by_id() {
    let app = lesson11::app();

    let (status, empty) = send(&app, Method::GET, "/tasks", None).await;
    assert_eq!(status, StatusCode::OK);
    assert_eq!(empty, json!([]));

    send(&app, Method::POST, "/tasks", Some(r#"{"title": "첫째"}"#)).await;
    send(&app, Method::POST, "/tasks", Some(r#"{"title": "둘째"}"#)).await;
    let (_, tasks) = send(&app, Method::GET, "/tasks", None).await;
    let ids: Vec<_> = tasks.as_array().unwrap().iter().map(|t| &t["id"]).collect();
    assert_eq!(ids, [1, 2]);
    assert_eq!(tasks[1]["title"], "둘째");
}

#[tokio::test]
async fn unknown_task_returns_404_json() {
    let app = lesson11::app();
    let (status, body) = send(&app, Method::GET, "/tasks/99", None).await;
    assert_eq!(status, StatusCode::NOT_FOUND);
    assert_eq!(body, json!({ "error": "할 일을 찾을 수 없습니다" }));
}

#[tokio::test]
async fn invalid_id_returns_400_json() {
    let app = lesson11::app();
    let expected = json!({ "error": "id는 1 이상의 정수여야 합니다" });

    for uri in [
        "/tasks/abc",
        "/tasks/0",
        "/tasks/-1",
        "/tasks/9223372036854775808",
    ] {
        let (status, body) = send(&app, Method::GET, uri, None).await;
        assert_eq!(status, StatusCode::BAD_REQUEST, "경로: {uri}");
        assert_eq!(body, expected, "경로: {uri}");
    }
}

#[tokio::test]
async fn patch_checks_body_before_looking_up_the_task() {
    let app = lesson11::app();
    send(&app, Method::POST, "/tasks", Some(r#"{"title": "운동"}"#)).await;

    let empty = json!({ "error": "수정할 필드(title, done)가 없습니다" });
    let (status, body) = send(&app, Method::PATCH, "/tasks/1", Some("{}")).await;
    assert_eq!(status, StatusCode::BAD_REQUEST);
    assert_eq!(body, empty);
    let (status, body) = send(&app, Method::PATCH, "/tasks/9999", Some("{}")).await;
    assert_eq!(status, StatusCode::BAD_REQUEST);
    assert_eq!(body, empty);

    let (status, _) = send(&app, Method::PATCH, "/tasks/1", Some(r#"{"title": " "}"#)).await;
    assert_eq!(status, StatusCode::BAD_REQUEST);

    let (status, body) = send(
        &app,
        Method::PATCH,
        "/tasks/9999",
        Some(r#"{"done": true}"#),
    )
    .await;
    assert_eq!(status, StatusCode::NOT_FOUND);
    assert_eq!(body, json!({ "error": "할 일을 찾을 수 없습니다" }));
}

#[tokio::test]
async fn put_is_not_allowed() {
    let app = lesson11::app();
    let request = Request::put("/tasks/1")
        .header(header::CONTENT_TYPE, "application/json")
        .body(Body::from(r#"{"title": "x"}"#))
        .unwrap();
    let response = app.oneshot(request).await.unwrap();
    assert_eq!(response.status(), StatusCode::METHOD_NOT_ALLOWED);
    let allow = response.headers()[header::ALLOW].to_str().unwrap();
    assert!(allow.contains("PATCH") && allow.contains("DELETE"));
}

#[tokio::test]
async fn patch_updates_only_the_sent_fields() {
    let app = lesson11::app();
    send(&app, Method::POST, "/tasks", Some(r#"{"title": "운동"}"#)).await;

    let (status, task) = send(&app, Method::PATCH, "/tasks/1", Some(r#"{"done": true}"#)).await;
    assert_eq!(status, StatusCode::OK);
    assert_eq!(task["title"], "운동");
    assert_eq!(task["done"], true);

    let (_, task) = send(
        &app,
        Method::PATCH,
        "/tasks/1",
        Some(r#"{"title": "  아침 운동  "}"#),
    )
    .await;
    assert_eq!(task["title"], "아침 운동");
    assert_eq!(task["done"], true);

    let (_, task) = send(&app, Method::GET, "/tasks/1", None).await;
    assert_eq!(task["title"], "아침 운동");
    assert_eq!(task["done"], true);
}

#[tokio::test]
async fn delete_then_get_returns_404() {
    let app = lesson11::app();
    send(&app, Method::POST, "/tasks", Some(r#"{"title": "임시"}"#)).await;

    let (status, _) = send(&app, Method::DELETE, "/tasks/1", None).await;
    assert_eq!(status, StatusCode::NO_CONTENT);

    let (status, _) = send(&app, Method::GET, "/tasks/1", None).await;
    assert_eq!(status, StatusCode::NOT_FOUND);

    let (status, _) = send(&app, Method::DELETE, "/tasks/1", None).await;
    assert_eq!(status, StatusCode::NOT_FOUND);
}
