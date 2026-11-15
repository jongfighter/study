fn main() {
    println!("안녕하세요, Rust!");

    let name = "Cargo";
    let day = 2;
    println!("{name}와 함께하는 {day}일차");
    println!("{}일차의 다음 날은 {}일차입니다", day, day + 1);
    println!("{0} + {0} = {1}", day, day * 2);

    let price = 12.3456;
    println!("가격은 {price:.2}원입니다");
    println!("[{:>6}] [{:<6}] [{:^6}]", 42, 42, 42);
    println!("[{:07.2}]", price);

    let tuple = (1, "둘", 3.0);
    println!("{tuple:?}");
    println!("{tuple:#?}");

    eprintln!("이 줄은 표준 오류로 나갑니다");
}
