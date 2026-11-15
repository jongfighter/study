use std::collections::HashMap;

fn word_frequencies(text: &str) -> Vec<(String, usize)> {
    let mut counts: HashMap<String, usize> = HashMap::new();
    for word in text.split_whitespace() {
        *counts.entry(word.to_lowercase()).or_insert(0) += 1;
    }

    let mut result: Vec<(String, usize)> = counts.into_iter().collect();
    result.sort_by(|a, b| b.1.cmp(&a.1).then(a.0.cmp(&b.0)));
    result
}

fn main() {
    let text = "to be or not to be that is the question to";
    for (word, count) in word_frequencies(text) {
        println!("{word}: {count}");
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn counts_and_sorts_by_frequency_then_word() {
        let result = word_frequencies("b a b c a b");
        assert_eq!(
            result,
            vec![
                (String::from("b"), 3),
                (String::from("a"), 2),
                (String::from("c"), 1),
            ]
        );
    }

    #[test]
    fn ignores_case() {
        let result = word_frequencies("Rust rust RUST go");
        assert_eq!(result[0], (String::from("rust"), 3));
    }

    #[test]
    fn empty_text_gives_empty_list() {
        assert!(word_frequencies("   ").is_empty());
    }
}
