const MAX_RETRY: u32 = 3;

fn main() {
    let x = 5;
    println!("x = {x}");

    let mut count = 0;
    count += 1;
    println!("count = {count}");

    let spaces = "   ";
    let spaces = spaces.len();
    println!("공백 개수 = {spaces}");

    let x = x * 2;
    {
        let x = x + 100;
        println!("안쪽 블록의 x = {x}");
    }
    println!("바깥 x = {x}");

    println!("최대 재시도 횟수 = {MAX_RETRY}");
}
