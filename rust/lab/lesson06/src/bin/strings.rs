fn main() {
    let literal: &str = "안녕, Rust";
    let mut owned = String::from(literal);
    owned.push_str("!!");
    owned.push('?');
    println!("{literal} / {owned}");

    println!(
        "바이트 수: {}, 글자 수: {}",
        literal.len(),
        literal.chars().count()
    );

    let hello = "안녕";
    println!("처음 3바이트: {}", &hello[0..3]);
    for c in hello.chars() {
        println!("글자: {c}");
    }

    let s1 = String::from("a");
    let s2 = String::from("b");
    let s3 = s1 + &s2;
    println!("합친 결과: {s3}, s2는 그대로: {s2}");
    let joined = format!("{s3}-{s2}-{}", 3);
    println!("{joined}");

    println!("[{}]", "  공백 제거  ".trim());
    println!("{}", "hello".to_uppercase());
    println!("{:?}", "a,b,c".split(',').collect::<Vec<_>>());
    println!("{}", "rust".contains("us"));
}
