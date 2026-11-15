use std::fmt;

use crate::Task;

#[derive(Debug, Default, Clone, PartialEq)]
pub struct TaskQuery {
    pub done: Option<bool>,
    pub search: Option<String>,
    pub limit: Option<usize>,
    pub offset: usize,
}

#[derive(Debug, PartialEq, Eq)]
pub enum QueryError {
    UnknownKey(String),
    InvalidValue { key: String, value: String },
}

impl fmt::Display for QueryError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            QueryError::UnknownKey(key) => write!(f, "알 수 없는 쿼리 키: {key}"),
            QueryError::InvalidValue { key, value } => {
                write!(f, "{key}의 값이 올바르지 않습니다: {value:?}")
            }
        }
    }
}

impl std::error::Error for QueryError {}

impl TaskQuery {
    /// `done=true&limit=2&offset=1&q=책` 같은 쿼리 문자열을 해석합니다.
    pub fn parse(query: &str) -> Result<TaskQuery, QueryError> {
        let mut parsed = TaskQuery::default();
        for pair in query.split('&').filter(|pair| !pair.is_empty()) {
            let (key, value) = pair.split_once('=').unwrap_or((pair, ""));
            let invalid = || QueryError::InvalidValue {
                key: key.to_string(),
                value: value.to_string(),
            };
            match key {
                "done" => parsed.done = Some(value.parse().map_err(|_| invalid())?),
                "q" => parsed.search = Some(value.to_string()),
                "limit" => parsed.limit = Some(value.parse().map_err(|_| invalid())?),
                "offset" => parsed.offset = value.parse().map_err(|_| invalid())?,
                other => return Err(QueryError::UnknownKey(other.to_string())),
            }
        }
        Ok(parsed)
    }
}

#[derive(Debug, PartialEq)]
pub struct Page<T> {
    pub items: Vec<T>,
    pub total: usize,
}

/// 필터와 검색을 먼저 적용한 뒤, 남은 개수(total)를 세고 나서 offset/limit으로 자릅니다.
pub fn run_query(tasks: &[Task], query: &TaskQuery) -> Page<Task> {
    let matched: Vec<&Task> = tasks
        .iter()
        .filter(|task| query.done.is_none_or(|done| task.done == done))
        .filter(|task| {
            query
                .search
                .as_ref()
                .is_none_or(|word| task.title.contains(word.as_str()))
        })
        .collect();
    let total = matched.len();
    let items = matched
        .into_iter()
        .skip(query.offset)
        .take(query.limit.unwrap_or(usize::MAX))
        .cloned()
        .collect();
    Page { items, total }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn sample() -> Vec<Task> {
        let titles = ["책 읽기", "운동", "책 정리", "장보기", "보고서"];
        titles
            .iter()
            .enumerate()
            .map(|(index, title)| Task {
                id: index as u64 + 1,
                title: title.to_string(),
                done: index % 2 == 0,
            })
            .collect()
    }

    #[test]
    fn parse_reads_every_key() {
        let query = TaskQuery::parse("done=true&limit=2&offset=1&q=책").unwrap();
        assert_eq!(
            query,
            TaskQuery {
                done: Some(true),
                search: Some("책".to_string()),
                limit: Some(2),
                offset: 1,
            }
        );
    }

    #[test]
    fn parse_reports_bad_input() {
        assert_eq!(
            TaskQuery::parse("limit=abc"),
            Err(QueryError::InvalidValue {
                key: "limit".to_string(),
                value: "abc".to_string(),
            })
        );
        assert_eq!(
            TaskQuery::parse("sort=id"),
            Err(QueryError::UnknownKey("sort".to_string()))
        );
        assert_eq!(TaskQuery::parse("").unwrap(), TaskQuery::default());
    }

    #[test]
    fn run_query_filters_then_pages() {
        let tasks = sample();

        let done_only = run_query(&tasks, &TaskQuery::parse("done=true").unwrap());
        assert_eq!(done_only.total, 3);

        let second_page = run_query(&tasks, &TaskQuery::parse("limit=2&offset=2").unwrap());
        assert_eq!(second_page.total, 5);
        let ids: Vec<u64> = second_page.items.iter().map(|task| task.id).collect();
        assert_eq!(ids, vec![3, 4]);

        let books = run_query(&tasks, &TaskQuery::parse("q=책&limit=1").unwrap());
        assert_eq!(books.total, 2);
        assert_eq!(books.items[0].title, "책 읽기");
    }
}
