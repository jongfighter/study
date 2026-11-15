fn square(x: i32) -> i32 {
    x * x
}

fn main() {
    let numbers = [1, 2, 3];
    let total = square(4);
    println!("{total}");
    for number in numbers {
        println!("{number}");
    }
}
