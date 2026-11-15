#[derive(Debug)]
enum Command {
    Add(String),
    Done(u32),
    Remove { id: u32, force: bool },
    List,
}

fn describe(command: &Command) -> String {
    match command {
        Command::Add(title) => format!("추가: {title}"),
        Command::Done(id) => format!("완료 처리: #{id}"),
        Command::Remove { id: 0, .. } => String::from("0번은 지울 수 없습니다"),
        Command::Remove { id, force: true } => format!("강제 삭제: #{id}"),
        Command::Remove { id, force: false } => format!("삭제: #{id}"),
        Command::List => String::from("목록 보기"),
    }
}

fn parse_number(text: &str) -> Option<u32> {
    text.parse().ok()
}

fn main() {
    let commands = [
        Command::Add(String::from("우유 사기")),
        Command::Done(3),
        Command::Remove {
            id: 0,
            force: false,
        },
        Command::Remove { id: 4, force: true },
        Command::List,
    ];
    for command in &commands {
        println!("{}", describe(command));
    }

    let maybe = parse_number("42");
    match maybe {
        Some(n) => println!("숫자 {n}"),
        None => println!("숫자가 아님"),
    }

    if let Some(n) = parse_number("abc") {
        println!("파싱 성공: {n}");
    } else {
        println!("abc는 숫자가 아닙니다");
    }

    let Some(n) = parse_number("7") else {
        println!("여기에는 오지 않습니다");
        return;
    };
    println!("let-else로 꺼낸 값: {n}");

    println!("기본값 사용: {}", parse_number("x").unwrap_or(0));
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn describe_covers_each_variant() {
        assert_eq!(describe(&Command::Add(String::from("a"))), "추가: a");
        assert_eq!(describe(&Command::Done(1)), "완료 처리: #1");
        assert_eq!(describe(&Command::List), "목록 보기");
    }

    #[test]
    fn remove_zero_is_protected() {
        let command = Command::Remove { id: 0, force: true };
        assert_eq!(describe(&command), "0번은 지울 수 없습니다");
    }

    #[test]
    fn parse_number_returns_option() {
        assert_eq!(parse_number("12"), Some(12));
        assert_eq!(parse_number("x"), None);
    }
}
