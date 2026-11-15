use std::env;

fn main() {
    println!("프로그램에 넘어온 인자를 출력합니다.");
    for arg in env::args().skip(1) {
        println!("인자: {arg}");
    }
}
