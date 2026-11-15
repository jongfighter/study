use std::error::Error;

fn main() -> Result<(), Box<dyn Error>> {
    let good: i32 = "42".parse()?;
    println!("성공: {good}");

    let bad: i32 = "forty-two".parse()?;
    println!("여기는 실행되지 않습니다: {bad}");
    Ok(())
}
