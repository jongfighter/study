use std::sync::{Arc, Mutex, mpsc};
use std::thread;

fn spawn_and_join() {
    let handle = thread::spawn(|| (1..=100).sum::<u32>());
    println!("메인 스레드: 작업이 끝나기를 기다립니다");
    let total = handle.join().unwrap();
    println!("1부터 100까지의 합 = {total}");
}

fn move_into_thread() {
    let mut names = vec![String::from("alice"), String::from("bob")];
    names.push(String::from("carol"));
    let handle = thread::spawn(move || names.len());
    println!("스레드가 센 이름 수 = {}", handle.join().unwrap());
}

fn channel_workers() {
    let (tx, rx) = mpsc::channel();
    for worker in 0..3 {
        let tx = tx.clone();
        thread::spawn(move || {
            let result = (worker + 1) * 100;
            tx.send((worker, result)).unwrap();
        });
    }
    drop(tx);

    let mut results: Vec<(i32, i32)> = rx.iter().collect();
    results.sort();
    println!("채널로 받은 결과 = {results:?}");
}

fn shared_counter() {
    let counter = Arc::new(Mutex::new(0));
    let mut handles = Vec::new();
    for id in 0..4 {
        let counter = Arc::clone(&counter);
        let handle = thread::Builder::new()
            .name(format!("worker-{id}"))
            .spawn(move || {
                for _ in 0..1000 {
                    *counter.lock().unwrap() += 1;
                }
            })
            .unwrap();
        handles.push(handle);
    }
    for handle in handles {
        handle.join().unwrap();
    }
    println!("최종 카운터 = {}", *counter.lock().unwrap());
}

fn scoped_threads() {
    let numbers: Vec<i32> = (1..=6).collect();
    let (left, right) = numbers.split_at(3);
    let (left_sum, right_sum) = thread::scope(|scope| {
        let a = scope.spawn(|| left.iter().sum::<i32>());
        let b = scope.spawn(|| right.iter().sum::<i32>());
        (a.join().unwrap(), b.join().unwrap())
    });
    println!("스코프 스레드 합계 = {left_sum} + {right_sum}");
}

fn main() {
    spawn_and_join();
    move_into_thread();
    channel_workers();
    shared_counter();
    scoped_threads();
}
