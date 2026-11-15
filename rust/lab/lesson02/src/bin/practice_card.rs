use std::env;

fn card(name: &str, role: &str) -> String {
    format!("[{role:^8}] {name}님, 환영합니다!")
}

fn main() {
    let name = env::args().nth(1).unwrap_or_else(|| String::from("방문자"));
    println!("{}", card(&name, "Rust"));
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn card_contains_name_and_role() {
        let text = card("지민", "Rust");
        assert!(text.contains("지민"));
        assert!(text.contains("Rust"));
    }

    #[test]
    fn card_centers_role_in_eight_columns() {
        assert_eq!(card("민수", "ab"), "[   ab   ] 민수님, 환영합니다!");
    }
}
