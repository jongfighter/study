fn main() {
    let mut scores = vec![85, 92, 78];
    scores.push(64);
    println!("점수: {scores:?}, 개수: {}", scores.len());

    println!("첫 점수: {}", scores[0]);
    println!("get(10): {:?}", scores.get(10));

    for score in &scores {
        println!("점수 {score}");
    }
    for score in &mut scores {
        *score += 5;
    }
    println!("5점 가산 후: {scores:?}");

    let total: i32 = scores.iter().sum();
    println!("합계 {total}");

    scores.sort();
    println!("정렬: {scores:?}");
    let last = scores.pop();
    println!("pop: {last:?}, 남은 점수: {scores:?}");

    let names: Vec<String> = Vec::new();
    println!("빈 벡터인가? {}", names.is_empty());
}
