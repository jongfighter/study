use std::num::ParseIntError;

fn parse_sum(input: &str) -> Result<i64, ParseIntError> {
    let mut total = 0;
    for piece in input.split(',') {
        let number: i64 = piece.trim().parse()?;
        total += number;
    }
    Ok(total)
}

fn main() {
    for input in ["1, 2, 3", "10,-4, 6", "1, two, 3"] {
        match parse_sum(input) {
            Ok(total) => println!("{input:?} 의 합은 {total}"),
            Err(e) => println!("{input:?} 은(는) 계산할 수 없습니다: {e}"),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn sums_valid_numbers() {
        assert_eq!(parse_sum("1, 2, 3"), Ok(6));
        assert_eq!(parse_sum("10,-4, 6"), Ok(12));
    }

    #[test]
    fn single_number_works() {
        assert_eq!(parse_sum("42"), Ok(42));
    }

    #[test]
    fn invalid_piece_returns_error() {
        assert!(parse_sum("1, two, 3").is_err());
        assert!(parse_sum("").is_err());
    }
}
