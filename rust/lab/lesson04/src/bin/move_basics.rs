fn takes_ownership(s: String) {
    println!("함수 안: {s}");
}

fn makes_copy(n: i32) {
    println!("함수 안: {n}");
}

fn gives_back(s: String) -> String {
    s
}

fn main() {
    let s1 = String::from("hello");
    let s2 = s1;
    println!("s2 = {s2}");

    let s3 = s2.clone();
    println!("s2 = {s2}, s3 = {s3}");

    let n = 5;
    let m = n;
    println!("n = {n}, m = {m}");

    makes_copy(n);
    println!("makes_copy 뒤에도 n = {n}");

    takes_ownership(s3);

    let s4 = gives_back(s2);
    println!("돌려받은 s4 = {s4}");

    println!(
        "String 값 하나의 크기: {}바이트",
        std::mem::size_of::<String>()
    );
}
