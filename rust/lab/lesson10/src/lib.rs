use std::thread;
use std::time::Duration;

use tokio::time::sleep;

pub fn parallel_sum(numbers: &[u64], workers: usize) -> u64 {
    assert!(workers > 0, "workers는 1 이상이어야 합니다");
    let chunk_size = numbers.len().div_ceil(workers).max(1);
    thread::scope(|scope| {
        let handles: Vec<_> = numbers
            .chunks(chunk_size)
            .map(|chunk| scope.spawn(move || chunk.iter().sum::<u64>()))
            .collect();
        handles
            .into_iter()
            .map(|handle| handle.join().unwrap())
            .sum()
    })
}

pub async fn fetch_user(id: u32) -> String {
    sleep(Duration::from_millis(50)).await;
    format!("사용자#{id}")
}

pub async fn fetch_all(ids: Vec<u32>) -> Vec<String> {
    let handles: Vec<_> = ids
        .into_iter()
        .map(|id| tokio::spawn(fetch_user(id)))
        .collect();
    let mut users = Vec::new();
    for handle in handles {
        users.push(handle.await.unwrap());
    }
    users
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parallel_sum_matches_sequential_sum() {
        let numbers: Vec<u64> = (1..=1000).collect();
        assert_eq!(parallel_sum(&numbers, 4), 500_500);
        assert_eq!(parallel_sum(&numbers, 7), 500_500);
        assert_eq!(parallel_sum(&[], 3), 0);
    }

    #[test]
    #[should_panic(expected = "workers")]
    fn parallel_sum_rejects_zero_workers() {
        parallel_sum(&[1, 2, 3], 0);
    }

    #[tokio::test]
    async fn fetch_all_keeps_input_order() {
        let users = fetch_all(vec![3, 1, 2]).await;
        assert_eq!(users, vec!["사용자#3", "사용자#1", "사용자#2"]);
    }
}
