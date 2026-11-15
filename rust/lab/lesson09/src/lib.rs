//! Task API에서 쓸 도메인 규칙을 모아 둔 작은 라이브러리 크레이트입니다.

pub mod task;
pub mod text;

pub use task::{Task, TaskError};
