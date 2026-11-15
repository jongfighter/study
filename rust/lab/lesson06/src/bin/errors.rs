use std::error::Error;
use std::fmt;
use std::num::ParseIntError;

#[derive(Debug, PartialEq)]
struct Task {
    id: u32,
    title: String,
    done: bool,
}

#[derive(Debug, PartialEq)]
enum ParseTaskError {
    WrongFieldCount(usize),
    InvalidId(ParseIntError),
    EmptyTitle,
    InvalidDone(String),
}

impl fmt::Display for ParseTaskError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            ParseTaskError::WrongFieldCount(n) => {
                write!(f, "필드는 3개여야 하는데 {n}개입니다")
            }
            ParseTaskError::InvalidId(e) => write!(f, "id가 숫자가 아닙니다: {e}"),
            ParseTaskError::EmptyTitle => write!(f, "제목이 비어 있습니다"),
            ParseTaskError::InvalidDone(text) => {
                write!(f, "done은 true/false여야 하는데 '{text}'입니다")
            }
        }
    }
}

impl Error for ParseTaskError {}

impl From<ParseIntError> for ParseTaskError {
    fn from(e: ParseIntError) -> ParseTaskError {
        ParseTaskError::InvalidId(e)
    }
}

fn parse_task(line: &str) -> Result<Task, ParseTaskError> {
    let fields: Vec<&str> = line.split(',').collect();
    if fields.len() != 3 {
        return Err(ParseTaskError::WrongFieldCount(fields.len()));
    }

    let id = fields[0].trim().parse::<u32>()?;
    let title = fields[1].trim();
    if title.is_empty() {
        return Err(ParseTaskError::EmptyTitle);
    }
    let done = match fields[2].trim() {
        "true" => true,
        "false" => false,
        other => return Err(ParseTaskError::InvalidDone(other.to_string())),
    };

    Ok(Task {
        id,
        title: title.to_string(),
        done,
    })
}

fn parse_all(text: &str) -> Result<Vec<Task>, ParseTaskError> {
    let mut tasks = Vec::new();
    for line in text.lines() {
        if line.trim().is_empty() {
            continue;
        }
        tasks.push(parse_task(line)?);
    }
    Ok(tasks)
}

fn main() -> Result<(), Box<dyn Error>> {
    let input = "1, 우유 사기, false\n2, 보고서 쓰기, true\n";
    let tasks = parse_all(input)?;
    for task in &tasks {
        println!("{task:?}");
    }

    for bad in ["x, 제목, false", "1, , false", "1, 제목", "1, 제목, maybe"] {
        match parse_task(bad) {
            Ok(task) => println!("뜻밖의 성공: {task:?}"),
            Err(e) => println!("실패한 입력 {bad:?} -> {e}"),
        }
    }

    let result = parse_all("1, 가, true\n2, 나, 3");
    println!("전체 파싱 결과: {result:?}");
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parses_valid_line() {
        let task = parse_task("5, 청소, true").unwrap();
        assert_eq!(
            task,
            Task {
                id: 5,
                title: String::from("청소"),
                done: true
            }
        );
    }

    #[test]
    fn reports_wrong_field_count() {
        assert_eq!(
            parse_task("1, 제목"),
            Err(ParseTaskError::WrongFieldCount(2))
        );
    }

    #[test]
    fn converts_parse_int_error_with_question_mark() {
        let err = parse_task("x, 제목, false").unwrap_err();
        assert!(matches!(err, ParseTaskError::InvalidId(_)));
    }

    #[test]
    fn rejects_empty_title_and_bad_done() {
        assert_eq!(parse_task("1, , false"), Err(ParseTaskError::EmptyTitle));
        assert_eq!(
            parse_task("1, a, maybe"),
            Err(ParseTaskError::InvalidDone(String::from("maybe")))
        );
    }

    #[test]
    fn parse_all_stops_at_first_error() {
        assert_eq!(parse_all("1, a, true\n\n2, b, false").unwrap().len(), 2);
        assert!(parse_all("1, a, true\n2, b").is_err());
    }
}
