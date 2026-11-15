use std::time::{Duration, Instant};

use tokio::time::{sleep, timeout};

async fn fetch_user(id: u32) -> String {
    sleep(Duration::from_millis(100)).await;
    format!("사용자#{id}")
}

fn rounded_ms(started: Instant) -> u128 {
    started.elapsed().as_millis() / 100 * 100
}

#[tokio::main]
async fn main() {
    let started = Instant::now();
    let first = fetch_user(1).await;
    let second = fetch_user(2).await;
    println!(
        "순차 실행: {first}, {second} (약 {}ms)",
        rounded_ms(started)
    );

    let started = Instant::now();
    let (third, fourth) = tokio::join!(fetch_user(3), fetch_user(4));
    println!(
        "join! 실행: {third}, {fourth} (약 {}ms)",
        rounded_ms(started)
    );

    let handles: Vec<_> = (5..=7).map(|id| tokio::spawn(fetch_user(id))).collect();
    for handle in handles {
        println!("spawn 결과: {}", handle.await.unwrap());
    }

    match timeout(Duration::from_millis(30), fetch_user(8)).await {
        Ok(user) => println!("제때 도착: {user}"),
        Err(error) => println!("시간 초과: {error}"),
    }
}
