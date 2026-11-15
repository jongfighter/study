use lesson05::{Priority, Task};

fn highest_pending(tasks: &[Task]) -> Option<&Task> {
    let mut best: Option<&Task> = None;
    for task in tasks {
        if task.done {
            continue;
        }
        match best {
            Some(current) if current.priority.weight() >= task.priority.weight() => {}
            _ => best = Some(task),
        }
    }
    best
}

fn main() {
    let tasks = [
        Task::new(1, "우유 사기", Priority::Low),
        Task::new(2, "보고서 쓰기", Priority::High),
        Task::new(3, "운동하기", Priority::Medium),
    ];
    match highest_pending(&tasks) {
        Some(task) => println!("가장 급한 일: {}", task.title),
        None => println!("남은 일이 없습니다"),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn picks_highest_priority_pending_task() {
        let tasks = [
            Task::new(1, "a", Priority::Low),
            Task::new(2, "b", Priority::High),
            Task::new(3, "c", Priority::Medium),
        ];
        assert_eq!(highest_pending(&tasks).map(|t| t.id), Some(2));
    }

    #[test]
    fn skips_done_tasks() {
        let mut urgent = Task::new(1, "a", Priority::High);
        urgent.complete();
        let tasks = [urgent, Task::new(2, "b", Priority::Low)];
        assert_eq!(highest_pending(&tasks).map(|t| t.id), Some(2));
    }

    #[test]
    fn empty_or_all_done_gives_none() {
        assert!(highest_pending(&[]).is_none());
        let mut only = Task::new(1, "a", Priority::Low);
        only.complete();
        assert!(highest_pending(&[only]).is_none());
    }
}
