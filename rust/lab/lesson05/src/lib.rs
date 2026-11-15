#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Priority {
    Low,
    Medium,
    High,
}

impl Priority {
    pub fn weight(self) -> u8 {
        match self {
            Priority::Low => 1,
            Priority::Medium => 2,
            Priority::High => 3,
        }
    }

    pub fn from_label(label: &str) -> Option<Priority> {
        match label {
            "low" => Some(Priority::Low),
            "medium" => Some(Priority::Medium),
            "high" => Some(Priority::High),
            _ => None,
        }
    }
}

#[derive(Debug, Clone, PartialEq)]
pub struct Task {
    pub id: u32,
    pub title: String,
    pub done: bool,
    pub priority: Priority,
}

impl Task {
    pub fn new(id: u32, title: &str, priority: Priority) -> Task {
        Task {
            id,
            title: title.to_string(),
            done: false,
            priority,
        }
    }

    pub fn complete(&mut self) {
        self.done = true;
    }

    pub fn summary(&self) -> String {
        let mark = if self.done { "x" } else { " " };
        format!("[{mark}] #{} {}", self.id, self.title)
    }
}

pub fn find_task(tasks: &[Task], id: u32) -> Option<&Task> {
    tasks.iter().find(|task| task.id == id)
}

pub fn first_pending(tasks: &[Task]) -> Option<&Task> {
    tasks.iter().find(|task| !task.done)
}

#[cfg(test)]
mod tests {
    use super::*;

    fn sample() -> Vec<Task> {
        let mut done_task = Task::new(1, "우유 사기", Priority::Low);
        done_task.complete();
        vec![
            done_task,
            Task::new(2, "보고서 쓰기", Priority::High),
            Task::new(3, "운동하기", Priority::Medium),
        ]
    }

    #[test]
    fn new_task_starts_pending() {
        let task = Task::new(7, "청소", Priority::Low);
        assert!(!task.done);
        assert_eq!(task.summary(), "[ ] #7 청소");
    }

    #[test]
    fn complete_marks_done() {
        let mut task = Task::new(7, "청소", Priority::Low);
        task.complete();
        assert_eq!(task.summary(), "[x] #7 청소");
    }

    #[test]
    fn find_task_returns_some_or_none() {
        let tasks = sample();
        assert_eq!(
            find_task(&tasks, 2).map(|t| t.title.as_str()),
            Some("보고서 쓰기")
        );
        assert!(find_task(&tasks, 99).is_none());
    }

    #[test]
    fn first_pending_skips_done_tasks() {
        let tasks = sample();
        assert_eq!(first_pending(&tasks).map(|t| t.id), Some(2));
        assert!(first_pending(&tasks[..1]).is_none());
    }

    #[test]
    fn priority_weights_and_labels() {
        assert_eq!(Priority::High.weight(), 3);
        assert_eq!(Priority::from_label("medium"), Some(Priority::Medium));
        assert_eq!(Priority::from_label("urgent"), None);
    }
}
