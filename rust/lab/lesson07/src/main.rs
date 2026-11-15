use std::fmt::{self, Display};

// ── 1. 제네릭: 타입을 나중에 정한다 ─────────────────────────────

fn largest<T: PartialOrd>(items: &[T]) -> Option<&T> {
    let mut iter = items.iter();
    let mut max = iter.next()?;
    for item in iter {
        if item > max {
            max = item;
        }
    }
    Some(max)
}

#[derive(Debug)]
struct Page<T> {
    items: Vec<T>,
    total: usize,
}

impl<T> Page<T> {
    fn new(items: Vec<T>, total: usize) -> Self {
        Page { items, total }
    }

    fn has_more(&self) -> bool {
        self.items.len() < self.total
    }
}

impl<T: Display> Page<T> {
    fn describe(&self) -> String {
        let mut parts = Vec::new();
        for item in &self.items {
            parts.push(item.to_string());
        }
        format!("[{}] 총 {}건", parts.join(", "), self.total)
    }
}

fn demo_generics() {
    println!("{:?}", largest(&[3, 9, 4]));
    println!("{:?}", largest(&[1.5, 0.2]));
    println!("{:?}", largest(&["rust", "kotlin", "python"]));
    let empty: [i32; 0] = [];
    println!("{:?}", largest(&empty));

    let page = Page::new(vec![10, 20, 30], 5);
    println!("{} / 더 있나? {}", page.describe(), page.has_more());
}

// ── 2. trait: 같은 행동을 약속한다 ──────────────────────────────

trait Notifier {
    fn name(&self) -> String;
    fn send(&self, to: &str, message: &str) -> String;

    fn broadcast(&self, recipients: &[&str], message: &str) -> Vec<String> {
        let mut results = Vec::new();
        for to in recipients {
            results.push(self.send(to, message));
        }
        results
    }
}

struct EmailNotifier {
    domain: String,
}

struct SlackNotifier {
    channel: String,
}

impl Notifier for EmailNotifier {
    fn name(&self) -> String {
        String::from("email")
    }

    fn send(&self, to: &str, message: &str) -> String {
        format!("[email] {to}@{} <- {message}", self.domain)
    }
}

impl Notifier for SlackNotifier {
    fn name(&self) -> String {
        String::from("slack")
    }

    fn send(&self, to: &str, message: &str) -> String {
        format!("[slack #{}] @{to} <- {message}", self.channel)
    }
}

fn demo_notifiers() {
    let email = EmailNotifier {
        domain: String::from("example.com"),
    };
    let slack = SlackNotifier {
        channel: String::from("dev"),
    };

    println!("{}", email.send("alice", "배포 완료"));
    for line in slack.broadcast(&["alice", "bob"], "점검 시작") {
        println!("{line}");
    }
}

// ── 3. trait bound: 제네릭에 "할 수 있는 일"을 요구한다 ───────────

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
enum Priority {
    Low,
    Normal,
    High,
}

impl Display for Priority {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        let label = match self {
            Priority::Low => "낮음",
            Priority::Normal => "보통",
            Priority::High => "높음",
        };
        write!(f, "{label}")
    }
}

fn notify_one<N: Notifier>(notifier: &N, to: &str, message: &str) -> String {
    notifier.send(to, message)
}

fn announce(notifier: &impl Notifier, message: &str) -> Vec<String> {
    notifier.broadcast(&["alice", "bob"], message)
}

fn describe_largest<T>(items: &[T]) -> String
where
    T: PartialOrd + Display,
{
    match largest(items) {
        Some(max) => format!("가장 큰 값: {max}"),
        None => String::from("비어 있음"),
    }
}

fn demo_bounds() {
    let email = EmailNotifier {
        domain: String::from("example.com"),
    };
    let slack = SlackNotifier {
        channel: String::from("dev"),
    };

    println!("{}", notify_one(&email, "carol", "회의 시작"));
    for line in announce(&slack, "점검 종료") {
        println!("{line}");
    }

    let priorities = [Priority::Low, Priority::High, Priority::Normal];
    println!("{}", describe_largest(&priorities));
    println!("{}", describe_largest(&[2.5, 8.0, 1.0]));
}

// ── 4. dyn Trait: 실행 중에 구현체를 고른다 ──────────────────────

fn send_everywhere(notifiers: &[Box<dyn Notifier>], to: &str, message: &str) {
    for notifier in notifiers {
        println!("{}: {}", notifier.name(), notifier.send(to, message));
    }
}

fn demo_dyn() {
    let notifiers: Vec<Box<dyn Notifier>> = vec![
        Box::new(EmailNotifier {
            domain: String::from("example.com"),
        }),
        Box::new(SlackNotifier {
            channel: String::from("ops"),
        }),
    ];
    send_everywhere(&notifiers, "carol", "장애 복구");
}

// ── 5. 라이프타임: 참조가 얼마나 사는지 컴파일러에게 알린다 ────────

fn longest<'a>(first: &'a str, second: &'a str) -> &'a str {
    if first.len() >= second.len() {
        first
    } else {
        second
    }
}

struct Excerpt<'a> {
    text: &'a str,
}

impl<'a> Excerpt<'a> {
    fn first_line(&self) -> &str {
        self.text.lines().next().unwrap_or("")
    }
}

fn demo_lifetimes() {
    let title_a = String::from("trait와 제네릭");
    let result;
    {
        let title_b = String::from("라이프타임");
        result = longest(title_a.as_str(), title_b.as_str()).to_string();
    }
    println!("더 긴 제목: {result}");

    let note = String::from("첫 줄입니다\n둘째 줄입니다");
    let excerpt = Excerpt { text: &note };
    println!("발췌: {}", excerpt.first_line());
}

// ── 6. 연습 문제 정답 ───────────────────────────────────────────

struct SmsNotifier {
    country_code: String,
}

impl Notifier for SmsNotifier {
    fn name(&self) -> String {
        String::from("sms")
    }

    fn send(&self, to: &str, message: &str) -> String {
        format!("[sms +{}] {to} <- {message}", self.country_code)
    }
}

fn broadcast_all(
    notifiers: &[Box<dyn Notifier>],
    recipients: &[&str],
    message: &str,
) -> Vec<String> {
    let mut results = Vec::new();
    for notifier in notifiers {
        results.extend(notifier.broadcast(recipients, message));
    }
    results
}

#[derive(Debug)]
struct Stack<T> {
    items: Vec<T>,
}

impl<T> Stack<T> {
    fn new() -> Self {
        Stack { items: Vec::new() }
    }

    fn push(&mut self, item: T) {
        self.items.push(item);
    }

    fn pop(&mut self) -> Option<T> {
        self.items.pop()
    }

    fn peek(&self) -> Option<&T> {
        self.items.last()
    }

    fn len(&self) -> usize {
        self.items.len()
    }

    fn is_empty(&self) -> bool {
        self.items.is_empty()
    }
}

fn longest_line(text: &str) -> &str {
    let mut best = "";
    for line in text.lines() {
        if line.len() > best.len() {
            best = line;
        }
    }
    best
}

fn demo_practice() {
    let notifiers: Vec<Box<dyn Notifier>> = vec![
        Box::new(EmailNotifier {
            domain: String::from("example.com"),
        }),
        Box::new(SmsNotifier {
            country_code: String::from("82"),
        }),
    ];
    for line in broadcast_all(&notifiers, &["dave", "erin"], "회의 10분 전") {
        println!("{line}");
    }

    let mut stack = Stack::new();
    stack.push("첫째");
    stack.push("둘째");
    println!("peek = {:?}, len = {}", stack.peek(), stack.len());
    println!("pop = {:?}", stack.pop());
    println!("pop = {:?}", stack.pop());
    println!("pop = {:?}, 비었나? {}", stack.pop(), stack.is_empty());

    println!("가장 긴 줄: {}", longest_line("짧음\n조금 더 긴 줄\n중간"));
}

fn main() {
    demo_generics();
    demo_notifiers();
    demo_bounds();
    demo_dyn();
    demo_lifetimes();
    demo_practice();
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn largest_finds_max_and_handles_empty() {
        assert_eq!(largest(&[3, 9, 4]), Some(&9));
        assert_eq!(largest(&["b", "a"]), Some(&"b"));
        let empty: [i32; 0] = [];
        assert_eq!(largest(&empty), None);
    }

    #[test]
    fn sms_notifier_formats_message() {
        let sms = SmsNotifier {
            country_code: String::from("82"),
        };
        assert_eq!(sms.send("dave", "안녕"), "[sms +82] dave <- 안녕");
    }

    #[test]
    fn broadcast_all_uses_every_notifier() {
        let notifiers: Vec<Box<dyn Notifier>> = vec![
            Box::new(EmailNotifier {
                domain: String::from("example.com"),
            }),
            Box::new(SmsNotifier {
                country_code: String::from("82"),
            }),
        ];
        let lines = broadcast_all(&notifiers, &["a", "b"], "hi");
        assert_eq!(lines.len(), 4);
        assert_eq!(lines[0], "[email] a@example.com <- hi");
        assert_eq!(lines[3], "[sms +82] b <- hi");
    }

    #[test]
    fn stack_works_like_a_stack() {
        let mut stack = Stack::new();
        assert!(stack.is_empty());
        stack.push(1);
        stack.push(2);
        assert_eq!(stack.peek(), Some(&2));
        assert_eq!(stack.pop(), Some(2));
        assert_eq!(stack.len(), 1);
        assert_eq!(stack.pop(), Some(1));
        assert_eq!(stack.pop(), None);
    }

    #[test]
    fn longest_line_picks_longest_and_first_on_tie() {
        assert_eq!(longest_line("aa\nbbbb\ncc"), "bbbb");
        assert_eq!(longest_line("aa\nbb"), "aa");
        assert_eq!(longest_line(""), "");
    }

    #[test]
    fn longest_returns_longer_string() {
        assert_eq!(longest("abc", "de"), "abc");
        assert_eq!(longest("ab", "cd"), "ab");
    }
}
