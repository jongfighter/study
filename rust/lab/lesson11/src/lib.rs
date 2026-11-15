//! 메모리 저장소를 쓰는 Task REST API (axum 0.8).

mod clock;
mod error;
mod handlers;
mod model;
mod store;

use axum::Router;
use axum::routing::get;

pub use error::ApiError;
pub use model::{CreateTask, Task, UpdateTask};
pub use store::TaskStore;

pub fn app() -> Router {
    app_with_store(TaskStore::default())
}

pub fn app_with_store(store: TaskStore) -> Router {
    Router::new()
        .route("/health", get(handlers::health))
        .route(
            "/tasks",
            get(handlers::list_tasks).post(handlers::create_task),
        )
        .route(
            "/tasks/{id}",
            get(handlers::get_task)
                .patch(handlers::update_task)
                .delete(handlers::delete_task),
        )
        .with_state(store)
}
