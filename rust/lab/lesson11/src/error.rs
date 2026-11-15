use axum::Json;
use axum::extract::rejection::JsonRejection;
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use serde_json::json;

#[derive(Debug)]
pub enum ApiError {
    NotFound,
    Validation(String),
    InvalidBody(JsonRejection),
}

impl IntoResponse for ApiError {
    fn into_response(self) -> Response {
        let (status, message) = match self {
            ApiError::NotFound => (
                StatusCode::NOT_FOUND,
                "할 일을 찾을 수 없습니다".to_string(),
            ),
            ApiError::Validation(message) => (StatusCode::BAD_REQUEST, message),
            ApiError::InvalidBody(rejection) => {
                let status = match rejection {
                    JsonRejection::JsonDataError(_) => StatusCode::BAD_REQUEST,
                    _ => rejection.status(),
                };
                (
                    status,
                    format!("요청 본문을 읽을 수 없습니다: {}", rejection.body_text()),
                )
            }
        };
        (status, Json(json!({ "error": message }))).into_response()
    }
}

impl From<JsonRejection> for ApiError {
    fn from(rejection: JsonRejection) -> Self {
        ApiError::InvalidBody(rejection)
    }
}
