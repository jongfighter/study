use lesson05::{Priority, Task, find_task, first_pending};

fn main() {
    let mut tasks = vec![
        Task::new(1, "우유 사기", Priority::Low),
        Task::new(2, "보고서 쓰기", Priority::High),
        Task::new(3, "운동하기", Priority::Medium),
    ];

    for task in &tasks {
        println!("{} (우선순위 {})", task.summary(), task.priority.weight());
    }

    match find_task(&tasks, 2) {
        Some(task) => println!("찾음: {}", task.title),
        None => println!("2번 할 일이 없습니다"),
    }

    if let Some(task) = tasks.first_mut() {
        task.complete();
    }

    match first_pending(&tasks) {
        Some(task) => println!("다음에 할 일: {}", task.title),
        None => println!("모두 끝났습니다"),
    }

    let label = Priority::from_label("urgent");
    println!("urgent 라벨: {label:?}");
}
