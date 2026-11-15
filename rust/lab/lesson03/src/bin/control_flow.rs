fn main() {
    let number = 7;
    if number % 2 == 0 {
        println!("{number}은(는) 짝수");
    } else {
        println!("{number}은(는) 홀수");
    }

    let label = if number > 5 { "큰 수" } else { "작은 수" };
    println!("{number}은(는) {label}");

    let mut counter = 0;
    let found = loop {
        counter += 1;
        if counter * counter > 50 {
            break counter;
        }
    };
    println!("제곱이 50을 넘는 첫 수: {found}");

    let mut remaining = 3;
    while remaining > 0 {
        println!("카운트다운 {remaining}");
        remaining -= 1;
    }

    let fruits = ["사과", "배", "귤"];
    for fruit in fruits {
        println!("과일: {fruit}");
    }
    for i in 1..=3 {
        println!("{i}번째");
    }
    for (index, fruit) in fruits.iter().enumerate() {
        println!("{index}: {fruit}");
    }
}
