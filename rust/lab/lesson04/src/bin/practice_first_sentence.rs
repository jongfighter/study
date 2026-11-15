fn first_sentence(text: &str) -> &str {
    match text.find('.') {
        Some(end) => &text[..=end],
        None => text,
    }
}

fn main() {
    let memo = String::from("소유권을 배웠다. 이제 빌림을 배운다. 다음은 슬라이스다.");
    let sentence = first_sentence(&memo);
    println!("첫 문장: {sentence}");
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn cuts_at_first_period_including_it() {
        assert_eq!(first_sentence("hello. world."), "hello.");
    }

    #[test]
    fn without_period_returns_everything() {
        assert_eq!(first_sentence("no period here"), "no period here");
    }

    #[test]
    fn works_with_korean_text() {
        assert_eq!(first_sentence("가나다. 라마바."), "가나다.");
    }
}
