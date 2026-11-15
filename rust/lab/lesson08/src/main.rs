use std::collections::HashMap;

#[derive(Debug, Clone)]
struct Task {
    id: u32,
    title: String,
    done: bool,
    minutes: u32,
}

impl Task {
    fn new(id: u32, title: &str, done: bool, minutes: u32) -> Self {
        Task {
            id,
            title: title.to_string(),
            done,
            minutes,
        }
    }
}

fn sample_tasks() -> Vec<Task> {
    vec![
        Task::new(1, "장보기", true, 30),
        Task::new(2, "보고서 작성", false, 90),
        Task::new(3, "운동", false, 45),
        Task::new(4, "코드 리뷰", false, 20),
        Task::new(5, "책 읽기", true, 60),
    ]
}

// ── 1. 클로저: 주변 변수를 붙잡는 이름 없는 함수 ──────────────────

fn demo_closures() {
    let add = |a: i32, b: i32| a + b;
    println!("add(2, 3) = {}", add(2, 3));

    let limit = 40;
    let is_long = |task: &Task| task.minutes > limit;
    for task in sample_tasks() {
        println!("{} → 오래 걸림? {}", task.title, is_long(&task));
    }
}

// ── 2. Fn / FnMut / FnOnce ─────────────────────────────────────

fn apply_twice<F: Fn(i32) -> i32>(f: F, x: i32) -> i32 {
    f(f(x))
}

fn call_three_times<F: FnMut()>(mut f: F) {
    f();
    f();
    f();
}

fn run_once<F: FnOnce() -> String>(f: F) -> String {
    f()
}

fn make_adder(n: i32) -> impl Fn(i32) -> i32 {
    move |x| x + n
}

type Op = Box<dyn Fn(i32) -> i32>;

fn demo_fn_traits() {
    let bonus = 10;
    println!("apply_twice = {}", apply_twice(|n| n + bonus, 1));

    let mut count = 0;
    call_three_times(|| count += 1);
    println!("count = {count}");

    let owner = String::from("소유권을 가져가는 클로저");
    println!("{}", run_once(move || owner));

    let ops: Vec<(&str, Op)> = vec![("+5", Box::new(make_adder(5))), ("×2", Box::new(|x| x * 2))];
    for (label, op) in &ops {
        println!("{label} 적용: {}", op(7));
    }
}

// ── 3. 이터레이터: 어댑터를 이어 붙이고 마지막에 소비한다 ──────────

fn demo_iterators() {
    let tasks = sample_tasks();

    let open_titles: Vec<String> = tasks
        .iter()
        .filter(|task| !task.done)
        .map(|task| task.title.clone())
        .collect();
    println!("남은 일: {open_titles:?}");

    let total: u32 = tasks.iter().map(|task| task.minutes).sum();
    let remaining = tasks
        .iter()
        .filter(|task| !task.done)
        .fold(0, |acc, task| acc + task.minutes);
    println!("전체 {total}분, 남은 시간 {remaining}분");

    if let Some(task) = tasks.iter().max_by_key(|task| task.minutes) {
        println!("가장 오래 걸리는 일: {}", task.title);
    }
    println!(
        "모두 양수? {} / 끝난 일 있음? {} / 90분 이상 위치: {:?}",
        tasks.iter().all(|task| task.minutes > 0),
        tasks.iter().any(|task| task.done),
        tasks.iter().position(|task| task.minutes >= 90),
    );

    for (index, task) in tasks.iter().enumerate().skip(1).take(2) {
        println!("{index}번째: {}", task.title);
    }

    let (done, todo): (Vec<&Task>, Vec<&Task>) = tasks.iter().partition(|task| task.done);
    println!("끝남 {}개 / 남음 {}개", done.len(), todo.len());

    let by_id: HashMap<u32, &Task> = tasks.iter().map(|task| (task.id, task)).collect();
    println!("3번 일: {}", by_id[&3].title);
}

fn demo_lazy() {
    let numbers = [1, 2, 3, 4, 5, 6];

    let pipeline = numbers
        .iter()
        .inspect(|n| println!("검사: {n}"))
        .map(|n| n * 10);
    println!("파이프라인을 만들었지만 아직 실행 전입니다");
    let first: Vec<i32> = pipeline.take(1).collect();
    println!("첫 결과: {first:?}");

    let first_even_square = numbers
        .iter()
        .inspect(|n| println!("검사: {n}"))
        .filter(|n| **n % 2 == 0)
        .map(|n| n * n)
        .next();
    println!("첫 짝수의 제곱: {first_even_square:?}");
}

// ── 4. Iterator trait를 직접 구현하기 ───────────────────────────

struct Countdown {
    remaining: u32,
}

impl Iterator for Countdown {
    type Item = u32;

    fn next(&mut self) -> Option<u32> {
        if self.remaining == 0 {
            None
        } else {
            let current = self.remaining;
            self.remaining -= 1;
            Some(current)
        }
    }
}

fn parse_minutes(inputs: &[&str]) -> Result<Vec<u32>, std::num::ParseIntError> {
    inputs.iter().map(|text| text.parse::<u32>()).collect()
}

fn demo_custom_iterator() {
    let odd: Vec<u32> = Countdown { remaining: 5 }.filter(|n| n % 2 == 1).collect();
    println!("5부터 거꾸로 센 홀수: {odd:?}");

    let sum: u32 = Countdown { remaining: 4 }.sum();
    println!("4+3+2+1 = {sum}");

    println!("{:?}", parse_minutes(&["30", "45"]));
    println!("{:?}", parse_minutes(&["30", "abc", "45"]));
}

// ── 5. 연습 문제 정답 ───────────────────────────────────────────

struct Fibonacci {
    current: u64,
    next: u64,
}

impl Fibonacci {
    fn new() -> Self {
        Fibonacci {
            current: 0,
            next: 1,
        }
    }
}

impl Iterator for Fibonacci {
    type Item = u64;

    fn next(&mut self) -> Option<u64> {
        let value = self.current;
        self.current = self.next;
        self.next += value;
        Some(value)
    }
}

fn open_titles_by_minutes(tasks: &[Task]) -> Vec<String> {
    let mut open: Vec<&Task> = tasks.iter().filter(|task| !task.done).collect();
    open.sort_by_key(|task| task.minutes);
    open.iter().map(|task| task.title.clone()).collect()
}

fn make_counter() -> impl FnMut() -> u32 {
    let mut count = 0;
    move || {
        count += 1;
        count
    }
}

fn demo_practice() {
    let firsts: Vec<u64> = Fibonacci::new().take(10).collect();
    println!("피보나치 앞 10개: {firsts:?}");
    let even_sum: u64 = Fibonacci::new().filter(|n| n % 2 == 0).take(5).sum();
    println!("짝수 항 5개의 합: {even_sum}");

    println!("짧은 순서: {:?}", open_titles_by_minutes(&sample_tasks()));

    let mut counter = make_counter();
    println!("카운터: {} {} {}", counter(), counter(), counter());
}

fn main() {
    demo_closures();
    demo_fn_traits();
    demo_iterators();
    demo_lazy();
    demo_custom_iterator();
    demo_practice();
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn countdown_counts_down_to_one() {
        let values: Vec<u32> = Countdown { remaining: 3 }.collect();
        assert_eq!(values, vec![3, 2, 1]);
    }

    #[test]
    fn parse_minutes_collects_result() {
        assert_eq!(parse_minutes(&["1", "2"]), Ok(vec![1, 2]));
        assert!(parse_minutes(&["1", "x"]).is_err());
    }

    #[test]
    fn fibonacci_starts_with_zero_and_one() {
        let firsts: Vec<u64> = Fibonacci::new().take(10).collect();
        assert_eq!(firsts, vec![0, 1, 1, 2, 3, 5, 8, 13, 21, 34]);
    }

    #[test]
    fn fibonacci_even_sum_of_first_five() {
        let sum: u64 = Fibonacci::new().filter(|n| n % 2 == 0).take(5).sum();
        assert_eq!(sum, 188);
    }

    #[test]
    fn open_titles_are_sorted_by_minutes() {
        let titles = open_titles_by_minutes(&sample_tasks());
        assert_eq!(titles, vec!["코드 리뷰", "운동", "보고서 작성"]);
    }

    #[test]
    fn counter_remembers_its_state() {
        let mut counter = make_counter();
        assert_eq!(counter(), 1);
        assert_eq!(counter(), 2);
        let mut other = make_counter();
        assert_eq!(other(), 1);
    }
}
