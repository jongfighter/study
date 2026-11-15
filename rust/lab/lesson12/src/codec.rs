use std::fmt;

use crate::Task;

#[derive(Debug, PartialEq, Eq)]
pub struct DecodeError(pub String);

impl fmt::Display for DecodeError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "줄을 해석할 수 없습니다: {}", self.0)
    }
}

impl std::error::Error for DecodeError {}

/// Task 한 건을 `id<TAB>done<TAB>title` 한 줄로 바꿉니다.
pub fn encode_task(task: &Task) -> String {
    let title = task.title.replace(['\t', '\n'], " ");
    format!("{}\t{}\t{}", task.id, task.done, title)
}

pub fn decode_task(line: &str) -> Result<Task, DecodeError> {
    let mut parts = line.splitn(3, '\t');
    let bad = || DecodeError(line.to_string());
    let id = parts.next().ok_or_else(bad)?.parse().map_err(|_| bad())?;
    let done = parts.next().ok_or_else(bad)?.parse().map_err(|_| bad())?;
    let title = parts.next().ok_or_else(bad)?.to_string();
    Ok(Task { id, title, done })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn encode_then_decode_round_trips() {
        let task = Task {
            id: 7,
            title: "탭\t포함 제목".to_string(),
            done: true,
        };
        let decoded = decode_task(&encode_task(&task)).unwrap();
        assert_eq!(decoded.id, 7);
        assert!(decoded.done);
        assert_eq!(decoded.title, "탭 포함 제목");
    }

    #[test]
    fn decode_rejects_broken_lines() {
        assert!(decode_task("abc\ttrue\t제목").is_err());
        assert!(decode_task("1\tmaybe\t제목").is_err());
        assert!(decode_task("1\ttrue").is_err());
    }
}
