fn add(a: i32, b: i32) -> i32 {
    a + b
}

fn square(n: i32) -> i32 {
    n * n
}

fn greet(name: &str) {
    println!("안녕하세요, {name}님!");
}

fn classify(n: i32) -> String {
    if n < 0 {
        return String::from("음수");
    }
    if n == 0 {
        String::from("영")
    } else {
        String::from("양수")
    }
}

fn main() {
    let y = {
        let a = 3;
        a + 1
    };
    println!("블록의 값 y = {y}");

    println!("add(2, 3) = {}", add(2, 3));
    println!("square(add(1, 2)) = {}", square(add(1, 2)));
    greet("Rust");

    for n in [-5, 0, 8] {
        println!("{n}은(는) {}", classify(n));
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn add_works() {
        assert_eq!(add(2, 3), 5);
    }

    #[test]
    fn square_works() {
        assert_eq!(square(-4), 16);
    }

    #[test]
    fn classify_handles_all_signs() {
        assert_eq!(classify(-1), "음수");
        assert_eq!(classify(0), "영");
        assert_eq!(classify(9), "양수");
    }
}
