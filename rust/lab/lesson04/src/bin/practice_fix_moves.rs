fn byte_len(s: &str) -> usize {
    s.len()
}

fn append_step(log: &mut String, step: &str) {
    log.push_str(" -> ");
    log.push_str(step);
}

fn main() {
    let title = String::from("소유권 연습");
    let n = byte_len(&title);
    println!("{title}의 길이는 {n}바이트");

    let mut log = String::from("시작");
    let first = log.clone();
    append_step(&mut log, "진행");
    println!("{first} / {log}");
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn byte_len_counts_bytes() {
        let s = String::from("가나");
        assert_eq!(byte_len(&s), 6);
        assert_eq!(s, "가나");
    }

    #[test]
    fn append_step_adds_arrow() {
        let mut log = String::from("시작");
        append_step(&mut log, "진행");
        assert_eq!(log, "시작 -> 진행");
    }
}
