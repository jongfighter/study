fn main() {
    let text = "255";
    let value: u8 = text.parse().unwrap();
    println!("값: {value}");

    let next = value + 1;
    println!("1을 더한 값: {next}");
}
