//! 캡스톤 확장 과제용 도구 상자입니다. 표준 라이브러리만 사용합니다.

pub mod codec;
pub mod query;
pub mod repo;

pub use query::{Page, QueryError, TaskQuery, run_query};
pub use repo::{FileRepository, InMemoryRepository, TaskRepository};

#[derive(Debug, Clone, PartialEq)]
pub struct Task {
    pub id: u64,
    pub title: String,
    pub done: bool,
}
