#[tokio::main]
async fn main() {
    let port: u16 = std::env::var("PORT")
        .ok()
        .and_then(|value| value.parse().ok())
        .unwrap_or(3000);

    let listener = tokio::net::TcpListener::bind(("127.0.0.1", port))
        .await
        .unwrap_or_else(|error| panic!("{port}번 포트를 열 수 없습니다: {error}"));
    println!("Task API가 http://127.0.0.1:{port} 에서 실행 중입니다. 종료하려면 Ctrl+C");

    axum::serve(listener, lesson11::app())
        .await
        .expect("서버 실행 중 오류가 발생했습니다");
}
