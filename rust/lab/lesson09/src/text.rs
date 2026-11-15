/// 제목의 앞뒤 공백을 지우고, 연속된 공백을 하나로 줄입니다.
///
/// ```
/// use lesson09::text::normalize_title;
///
/// assert_eq!(normalize_title("  장보기   목록  "), "장보기 목록");
/// ```
pub fn normalize_title(title: &str) -> String {
    title.split_whitespace().collect::<Vec<_>>().join(" ")
}

/// 제목을 URL에 쓸 수 있는 소문자 slug로 바꿉니다.
///
/// 영문 소문자·숫자만 남기고 나머지 글자는 `-` 하나로 합칩니다.
/// 한글은 이 예제의 규칙에서 빠지므로 한글만 있는 제목은 빈 문자열이 됩니다.
///
/// ```
/// use lesson09::text::slugify;
///
/// assert_eq!(slugify("Write the  Report!"), "write-the-report");
/// assert_eq!(slugify("보고서"), "");
/// ```
pub fn slugify(title: &str) -> String {
    let mut slug = String::new();
    for ch in title.chars() {
        if is_slug_char(ch) {
            slug.push(ch.to_ascii_lowercase());
        } else if !slug.is_empty() && !slug.ends_with('-') {
            slug.push('-');
        }
    }
    slug.trim_end_matches('-').to_string()
}

/// 제목이 `max_chars` 글자를 넘으면 잘라 내고 `…`를 붙입니다.
///
/// # Panics
///
/// `max_chars`가 0이면 패닉합니다.
///
/// ```
/// use lesson09::text::truncate;
///
/// assert_eq!(truncate("가나다라마", 3), "가나다…");
/// assert_eq!(truncate("가나", 3), "가나");
/// ```
pub fn truncate(title: &str, max_chars: usize) -> String {
    assert!(max_chars > 0, "max_chars는 1 이상이어야 합니다");
    if title.chars().count() <= max_chars {
        return title.to_string();
    }
    let mut cut: String = title.chars().take(max_chars).collect();
    cut.push('…');
    cut
}

/// 제목에 들어 있는 단어(공백으로 구분) 수를 셉니다.
///
/// ```
/// use lesson09::text::word_count;
///
/// assert_eq!(word_count("  장보기   목록 "), 2);
/// assert_eq!(word_count(""), 0);
/// ```
pub fn word_count(title: &str) -> usize {
    title.split_whitespace().count()
}

fn is_slug_char(ch: char) -> bool {
    ch.is_ascii_alphanumeric()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn normalize_collapses_whitespace() {
        assert_eq!(normalize_title("  a   b\t c \n"), "a b c");
    }

    #[test]
    fn slugify_keeps_ascii_words() {
        assert_eq!(slugify("Hello, Rust World!"), "hello-rust-world");
        assert_eq!(slugify("---a---"), "a");
    }

    #[test]
    fn truncate_cuts_by_chars_not_bytes() {
        assert_eq!(truncate("가나다라", 2), "가나…");
    }

    #[test]
    #[should_panic(expected = "max_chars")]
    fn truncate_rejects_zero() {
        truncate("abc", 0);
    }

    #[test]
    fn word_count_ignores_extra_spaces() {
        assert_eq!(word_count("  a   b  c "), 3);
    }
}
