fn length(s: &str) -> usize {
    s.len()
}

fn add_world(s: &mut String) {
    s.push_str(", world");
}

fn main() {
    let mut s = String::from("hello");

    let len = length(&s);
    println!("'{s}'의 길이는 {len}");

    add_world(&mut s);
    println!("바뀐 값: {s}");

    let r1 = &s;
    let r2 = &s;
    println!("읽기 참조 두 개: {r1} / {r2}");

    let r3 = &mut s;
    r3.push('!');
    println!("가변 참조: {r3}");
    println!("원래 변수: {s}");
}
