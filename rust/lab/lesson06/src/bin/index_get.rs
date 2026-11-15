fn main() {
    let mut scores = vec![85, 92, 78];
    scores.push(64);
    let index = scores.len() + 6;

    match scores.get(index) {
        Some(score) => println!("{index}번째 점수는 {score}"),
        None => println!("{index}번째 점수는 없습니다"),
    }
}
