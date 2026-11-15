fn first_word(s: &str) -> &str {
    let bytes = s.as_bytes();
    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[0..i];
        }
    }
    s
}

fn main() {
    let sentence = String::from("hello wide world");
    let word = first_word(&sentence);
    println!("첫 단어: {word}");

    let hello = &sentence[0..5];
    let world = &sentence[11..];
    println!("{hello} / {world}");

    let numbers = [10, 20, 30, 40, 50];
    let middle = &numbers[1..4];
    println!("가운데 조각: {middle:?}, 길이 {}", middle.len());

    println!(
        "&str 참조 하나의 크기: {}바이트",
        std::mem::size_of::<&str>()
    );
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn picks_first_word() {
        assert_eq!(first_word("hello wide world"), "hello");
    }

    #[test]
    fn single_word_is_returned_whole() {
        assert_eq!(first_word("rust"), "rust");
    }

    #[test]
    fn empty_input_gives_empty_word() {
        assert_eq!(first_word(""), "");
    }
}
