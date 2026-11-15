use lesson09::text::slugify;
use lesson09::{Task, TaskError};

#[test]
fn task_lifecycle_uses_only_public_api() {
    let mut task = Task::new(1, "  Write   the Report ").unwrap();
    assert_eq!(task.title(), "Write the Report");
    assert_eq!(slugify(task.title()), "write-the-report");

    task.complete();
    assert!(task.is_done());
}

#[test]
fn rename_keeps_state_and_validates_the_new_title() {
    let mut task = Task::new(1, "운동").unwrap();
    task.complete();

    task.rename("  아침   운동 ").unwrap();
    assert_eq!(task.title(), "아침 운동");
    assert!(task.is_done());

    assert_eq!(task.rename("   "), Err(TaskError::EmptyTitle));
    assert_eq!(task.title(), "아침 운동");
}

#[test]
fn invalid_titles_return_errors() {
    assert_eq!(Task::new(1, ""), Err(TaskError::EmptyTitle));

    let error = Task::new(2, &"가".repeat(80)).unwrap_err();
    assert_eq!(
        error.to_string(),
        "제목이 너무 깁니다 (최대 50자, 입력 80자)"
    );
}
