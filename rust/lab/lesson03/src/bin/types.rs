fn main() {
    let a: i32 = -42;
    let b: u8 = 255;
    let big: i64 = 3_000_000_000;
    let ratio = 2.5;
    let ok = true;
    let letter = 'R';
    let hangul = '한';
    println!("{a} {b} {big} {ratio} {ok} {letter} {hangul}");

    println!("u8 최댓값 {}, i32 최솟값 {}", u8::MAX, i32::MIN);
    println!("255u8 + 1 (checked)  = {:?}", b.checked_add(1));
    println!("255u8 + 1 (wrapping) = {}", b.wrapping_add(1));

    let pair: (i32, &str) = (7, "일곱");
    let (number, word) = pair;
    println!("{number} {word} {}", pair.0);

    let days = [31, 28, 31];
    println!("첫 달은 {}일, 달 수는 {}", days[0], days.len());

    let average = f64::from(days[0] + days[1] + days[2]) / 3.0;
    println!("평균 {average:.2}");

    let truncated = 7.9_f64 as i32;
    println!("7.9 as i32 = {truncated}");
}
