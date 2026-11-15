use std::fmt;

use crate::text::normalize_title;

pub const MAX_TITLE_CHARS: usize = 50;

#[derive(Debug, Clone, PartialEq)]
pub struct Task {
    id: u32,
    title: String,
    done: bool,
}

#[derive(Debug, PartialEq, Eq)]
pub enum TaskError {
    EmptyTitle,
    TitleTooLong { max: usize, actual: usize },
}

impl fmt::Display for TaskError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            TaskError::EmptyTitle => write!(f, "제목이 비어 있습니다"),
            TaskError::TitleTooLong { max, actual } => {
                write!(f, "제목이 너무 깁니다 (최대 {max}자, 입력 {actual}자)")
            }
        }
    }
}

impl std::error::Error for TaskError {}

impl Task {
    pub fn new(id: u32, title: &str) -> Result<Task, TaskError> {
        let title = validate_title(title)?;
        Ok(Task {
            id,
            title,
            done: false,
        })
    }

    pub fn id(&self) -> u32 {
        self.id
    }

    pub fn title(&self) -> &str {
        &self.title
    }

    pub fn is_done(&self) -> bool {
        self.done
    }

    pub fn complete(&mut self) {
        self.done = true;
    }

    pub fn toggle(&mut self) {
        if self.done {
            self.reopen();
        } else {
            self.complete();
        }
    }

    pub fn rename(&mut self, title: &str) -> Result<(), TaskError> {
        self.title = validate_title(title)?;
        Ok(())
    }

    pub(crate) fn reopen(&mut self) {
        self.done = false;
    }
}

fn validate_title(title: &str) -> Result<String, TaskError> {
    let title = normalize_title(title);
    if title.is_empty() {
        return Err(TaskError::EmptyTitle);
    }
    let actual = title.chars().count();
    if actual > MAX_TITLE_CHARS {
        return Err(TaskError::TitleTooLong {
            max: MAX_TITLE_CHARS,
            actual,
        });
    }
    Ok(title)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn new_normalizes_title() {
        let task = Task::new(1, "  보고서   작성 ").unwrap();
        assert_eq!(task.title(), "보고서 작성");
        assert!(!task.is_done());
    }

    #[test]
    fn empty_title_is_rejected() {
        assert_eq!(Task::new(1, "   "), Err(TaskError::EmptyTitle));
    }

    #[test]
    fn too_long_title_reports_lengths() {
        let long = "가".repeat(MAX_TITLE_CHARS + 1);
        assert_eq!(
            Task::new(1, &long),
            Err(TaskError::TitleTooLong {
                max: MAX_TITLE_CHARS,
                actual: MAX_TITLE_CHARS + 1,
            })
        );
    }

    #[test]
    fn toggle_flips_done_back_and_forth() -> Result<(), TaskError> {
        let mut task = Task::new(7, "운동")?;
        task.toggle();
        assert!(task.is_done());
        task.toggle();
        assert!(!task.is_done());
        Ok(())
    }

    #[test]
    #[ignore = "예시: 평소에는 건너뛰고 필요할 때 --ignored로 실행합니다"]
    fn slow_example() {
        let titles: Vec<String> = (0..10_000).map(|n| format!("할 일 {n}")).collect();
        assert_eq!(titles.len(), 10_000);
    }
}
