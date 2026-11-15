use lesson09::text::{slugify, truncate};
use lesson09::{Task, TaskError};

fn main() {
    let titles = [
        String::from("  보고서   작성  "),
        String::new(),
        String::from("Write the Report!"),
        "가".repeat(60),
    ];

    for (index, title) in titles.iter().enumerate() {
        let id = (index + 1) as u32;
        match Task::new(id, title) {
            Ok(mut task) => {
                task.complete();
                println!(
                    "#{} {} (slug: {:?}, 완료: {})",
                    task.id(),
                    task.title(),
                    slugify(task.title()),
                    task.is_done()
                );
            }
            Err(TaskError::EmptyTitle) => println!("#{id} 건너뜀: 제목이 없습니다"),
            Err(error) => println!("#{id} 실패: {error} → {}", truncate(title, 6)),
        }
    }
}
