use axum::Json;
use axum::extract::rejection::JsonRejection;
use axum::extract::{Path, State};
use axum::http::{StatusCode, header};
use serde_json::{Value, json};

use crate::error::ApiError;
use crate::model::{CreateTask, Task, UpdateTask};
use crate::store::TaskStore;

const MAX_TITLE_CHARS: usize = 100;

pub async fn health() -> Json<Value> {
    Json(json!({ "status": "ok" }))
}

pub async fn list_tasks(State(store): State<TaskStore>) -> Json<Vec<Task>> {
    Json(store.list())
}

pub async fn create_task(
    State(store): State<TaskStore>,
    payload: Result<Json<CreateTask>, JsonRejection>,
) -> Result<(StatusCode, [(header::HeaderName, String); 1], Json<Task>), ApiError> {
    let Json(input) = payload?;
    let title = validate_title(&input.title.unwrap_or_default())?;
    let task = store.create(title);
    let location = [(header::LOCATION, format!("/tasks/{}", task.id))];
    Ok((StatusCode::CREATED, location, Json(task)))
}

pub async fn get_task(
    State(store): State<TaskStore>,
    Path(id): Path<String>,
) -> Result<Json<Task>, ApiError> {
    let id = parse_id(&id)?;
    store.get(id).map(Json).ok_or(ApiError::NotFound)
}

pub async fn update_task(
    State(store): State<TaskStore>,
    Path(id): Path<String>,
    payload: Result<Json<UpdateTask>, JsonRejection>,
) -> Result<Json<Task>, ApiError> {
    let id = parse_id(&id)?;
    let Json(mut changes) = payload?;
    if changes.title.is_none() && changes.done.is_none() {
        return Err(ApiError::Validation(
            "수정할 필드(title, done)가 없습니다".to_string(),
        ));
    }
    if let Some(title) = &changes.title {
        changes.title = Some(validate_title(title)?);
    }
    store
        .update(id, changes)
        .map(Json)
        .ok_or(ApiError::NotFound)
}

pub async fn delete_task(
    State(store): State<TaskStore>,
    Path(id): Path<String>,
) -> Result<StatusCode, ApiError> {
    let id = parse_id(&id)?;
    if store.delete(id) {
        Ok(StatusCode::NO_CONTENT)
    } else {
        Err(ApiError::NotFound)
    }
}

fn parse_id(text: &str) -> Result<u64, ApiError> {
    text.parse::<i64>()
        .ok()
        .filter(|id| *id >= 1)
        .map(|id| id as u64)
        .ok_or_else(|| ApiError::Validation("id는 1 이상의 정수여야 합니다".to_string()))
}

fn validate_title(title: &str) -> Result<String, ApiError> {
    let title = title.trim();
    if title.is_empty() {
        return Err(ApiError::Validation(
            "title은 비어 있을 수 없습니다".to_string(),
        ));
    }
    if title.chars().count() > MAX_TITLE_CHARS {
        return Err(ApiError::Validation(format!(
            "title은 {MAX_TITLE_CHARS}자 이하여야 합니다"
        )));
    }
    Ok(title.to_string())
}
