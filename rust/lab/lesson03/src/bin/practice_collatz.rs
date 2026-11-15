fn collatz_steps(start: u64) -> u32 {
    let mut n = start;
    let mut steps = 0;
    while n != 1 {
        if n.is_multiple_of(2) {
            n /= 2;
        } else {
            n = 3 * n + 1;
        }
        steps += 1;
    }
    steps
}

fn main() {
    for start in [1, 6, 27] {
        println!("{start}에서 1까지 {}단계", collatz_steps(start));
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn one_needs_no_steps() {
        assert_eq!(collatz_steps(1), 0);
    }

    #[test]
    fn small_examples() {
        assert_eq!(collatz_steps(2), 1);
        assert_eq!(collatz_steps(6), 8);
    }

    #[test]
    fn famous_example() {
        assert_eq!(collatz_steps(27), 111);
    }
}
