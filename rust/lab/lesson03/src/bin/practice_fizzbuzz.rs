fn fizzbuzz(n: u32) -> String {
    if n.is_multiple_of(15) {
        String::from("FizzBuzz")
    } else if n.is_multiple_of(3) {
        String::from("Fizz")
    } else if n.is_multiple_of(5) {
        String::from("Buzz")
    } else {
        n.to_string()
    }
}

fn main() {
    for n in 1..=15 {
        println!("{}", fizzbuzz(n));
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn plain_numbers_stay_numbers() {
        assert_eq!(fizzbuzz(1), "1");
        assert_eq!(fizzbuzz(7), "7");
    }

    #[test]
    fn multiples_of_three_and_five() {
        assert_eq!(fizzbuzz(9), "Fizz");
        assert_eq!(fizzbuzz(10), "Buzz");
    }

    #[test]
    fn multiples_of_fifteen() {
        assert_eq!(fizzbuzz(15), "FizzBuzz");
        assert_eq!(fizzbuzz(30), "FizzBuzz");
    }
}
