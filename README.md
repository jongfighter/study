# 공부 모음

Spring·Kotlin 경험자를 위한 한국어 HTML 학습 코스 모음입니다. 읽으면서 그대로 따라 할 수 있고, 레슨마다 VS Code 단축키와 설정을 하나씩 익힙니다.

| 코스 | 내용 | 열기 |
|---|---|---|
| Go 실전 입문 | 문법부터 net/http Task REST API까지 12개 레슨 | [`go/index.html`](go/index.html) |
| Rust 실전 입문 | cargo, 소유권부터 axum Task REST API까지 13개 레슨 | [`rust/index.html`](rust/index.html) |

## 폰이나 브라우저에서 보기

**https://jongfighter.github.io/study/** 를 열면 됩니다. GitHub Pages로 배포되어 있어 폰 브라우저에서도 그대로 읽을 수 있습니다.

- 이 저장소 화면(github.com)에서 `.html` 파일을 누르면 렌더링되지 않고 소스 코드가 보입니다. 위 주소로 들어가세요.
- 체크한 진도는 보고 있는 브라우저에만 저장되므로 PC와 폰의 진도는 따로 관리됩니다.
- `lab/` 안의 `.go`·`.rs` 같은 코드 파일은 위 주소에서는 다운로드로 처리될 수 있습니다. 코드는 이 저장소 화면(github.com)에서 파일을 눌러 보는 편이 읽기 좋습니다.

## PC에서 따라 하기

1. 이 저장소를 내려받습니다.
2. 원문 자료를 아래 명령으로 내려받아 각 코스 폴더 안에 둡니다(레슨이 이 폴더를 가리킵니다).
3. VS Code에서 `go` 또는 `rust` 폴더를 **폴더 열기**로 열고 `index.html`의 안내를 따릅니다.

```bash
git clone --depth 1 https://github.com/golang/tour go/tour-of-go
git clone --depth 1 https://github.com/mmcgrana/gobyexample go/go-by-example
git clone --depth 1 https://github.com/quii/learn-go-with-tests go/learn-go-with-tests
git clone --depth 1 https://github.com/teivah/100-go-mistakes go/100-go-mistakes

git clone --depth 1 https://github.com/rust-lang/book rust/the-rust-book
git clone --depth 1 https://github.com/rust-lang/rust-by-example rust/rust-by-example
git clone --depth 1 https://github.com/rust-lang/rustlings rust/rustlings
git clone --depth 1 https://github.com/tokio-rs/axum rust/axum
```

`go/tour-of-go`는 `golang/tour`(공식 미러 저장소)입니다. 원문 자료는 각 저장소의 라이선스를 따르며 이 저장소에는 포함하지 않습니다.

### 이 코스가 기준으로 삼은 자료 버전

위 명령은 항상 최신 커밋을 받습니다. 레슨이 가리키는 파일 경로나 제목은 아래 커밋 기준이므로, 최신 자료와 다르면 해당 커밋을 기준으로 확인하세요(`git -C <폴더> fetch --depth 1 origin <커밋>` 후 `git -C <폴더> checkout FETCH_HEAD`).

| 폴더 | 커밋 | 커밋 날짜 |
|---|---|---|
| `go/tour-of-go` | `11c9ad7eadf5e916a1137ad4dcc113f307f27874` | 2025-10-24 |
| `go/go-by-example` | `5d9b1f14cef8efbe5de247553c1acb9826134c1c` | 2026-10-31 |
| `go/learn-go-with-tests` | `592e9f747e9e851af40eb0a5cb0ac4fe86ef18a6` | 2026-10-26 |
| `go/100-go-mistakes` | `c13d0b48075f2218b4d63116705005f88bfdc293` | 2026-04-21 |
| `rust/the-rust-book` | `2771df165f092a4fbb9f46b5b2e7450f19010adb` | 2026-10-28 |
| `rust/rust-by-example` | `dcc14f912bf8790792a2df082eddc82ab2539229` | 2026-09-29 |
| `rust/rustlings` | `f7b9b3c5882a1a43275130b95142d94f847526f5` | 2026-10-26 |
| `rust/axum` | `2449f1600c6626a9b8b7dbcc3500bf64995968e1` | 2026-11-15 |

참고: `go/go-by-example`의 `go.mod`는 Go 1.27.1을 요구합니다. 설치된 Go가 더 낮으면 그 폴더 안에서는 `go` 명령을 실행하지 말고 읽기만 하세요.
