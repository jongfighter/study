# Rust 실전 입문 — 시작하기

Spring·Kotlin 경험자를 위한 Rust 학습 코스입니다. 레슨은 모두 한국어 HTML 문서이고, 문서를 읽으면서 VS Code에서 그대로 따라 합니다.

## 시작 순서 (약 5분)

1. **폴더 열기**: VS Code 메뉴 `File` > `Open Folder...`(`Ctrl+K Ctrl+O`)로 이 폴더(`D:\study\rust`)를 엽니다. 파일 하나만 여는 `Open File...`이 아닙니다.
2. **신뢰하기**: 제한 모드(Restricted Mode) 배너가 보이면 `Manage`에서 이 폴더를 신뢰합니다. 신뢰하지 않으면 rust-analyzer가 동작하지 않습니다. "부모 폴더 전체를 신뢰"하는 체크 상자는 켜지 마세요.
3. **권장 확장 설치**: `Ctrl+Shift+X`로 확장 보기를 열고 `@recommended`를 검색해 `rust-lang.rust-analyzer`, `tamasfe.even-better-toml`, `ms-vscode.cpptools`, `humao.rest-client`를 설치합니다(`.vscode/extensions.json`에 적혀 있습니다).
4. **문서 열기**: 탐색기에서 `index.html`을 연 상태로 `Ctrl+K R`을 누르거나 파일을 우클릭해 `Reveal in File Explorer`를 고른 뒤, 파일 탐색기에서 더블 클릭하면 브라우저로 열립니다.
5. **첫 레슨**: `lessons/00-vscode-start.html`(VS Code 첫걸음)부터 읽습니다. 다음은 `lessons/01-environment.html`입니다.

## 폴더 구성

```
D:\study\rust
├─ index.html        코스 홈
├─ lessons\          레슨 HTML (00~12)와 VS Code 치트시트
├─ lab\              레슨 실습 코드 (Cargo 워크스페이스, 레슨마다 크레이트 lessonNN)
├─ .vscode\          VS Code 설정 (권장 확장, 설정, 작업, 디버그)
└─ axum, rustlings, the-rust-book, rust-by-example
                     클론한 읽기 자료 (수정하지 않습니다. GitHub에서 저장소를 새로 받았다면
                     이 폴더들은 없으니 저장소 맨 위 README.md의 git clone 명령으로 받습니다)
```

## 알아 둘 점

- **기본 터미널이 MSYS2 bash이면** `cargo`, `rustc` 명령을 못 찾을 수 있습니다. MSYS2 bash는 Windows의 PATH를 이어받지 않는 것이 기본 동작이라, 사용자 프로필(`--login -i`로 시작하는 `MSYS2 UCRT64`)과 같은 조건에서 두 명령이 "명령을 찾을 수 없음"으로 나오는 것을 확인했습니다. 터미널 오른쪽 위 `+` 옆 화살표에서 PowerShell 프로필로 새 터미널을 열면 됩니다. bash를 계속 쓰려면 사용자 설정의 해당 프로필 `env`에 `"MSYS2_PATH_TYPE": "inherit"`를 추가합니다(적용 후 `cargo`가 찾아지는 것을 확인했습니다).
- **rust-analyzer 범위**: `.vscode/settings.json`의 `rust-analyzer.linkedProjects`가 `lab/Cargo.toml` 하나만 읽도록 고정합니다. 지정하지 않으면 폴더 안의 `Cargo.toml`을 전부 찾아 읽으므로 클론한 `axum` 같은 큰 프로젝트 때문에 느려집니다(임시 복제 구조에서 자동 탐색 시 워크스페이스 3개, 설정 후 1개만 로드되는 것을 확인했습니다). `rust-analyzer.files.exclude`와 `files.watcherExclude`로 클론 폴더와 `target`을 분석·감시에서 뺐습니다.
- **디버거 선택 근거**: Windows(MSVC 도구 체인)에서는 `ms-vscode.cpptools`(디버거 종류 `cppvsdbg`)를 씁니다. VS Code 공식 Rust 문서가 Windows는 Microsoft C++ 확장, macOS·Linux는 CodeLLDB를 안내하고, rust-analyzer의 디버그 기능도 Windows에서 이 확장을 `cppvsdbg` 설정으로 자동 구성하기 때문입니다. Rust 파일에 중단점을 걸려면 `debug.allowBreakpointsEverywhere`를 켜야 하며 이 폴더 설정에 이미 켜 두었습니다.
- **VS Code 설정**: `tasks.json`(`cargo: check all`, `cargo: test all`, `cargo: clippy all`, `cargo: run task api`와 디버그 준비용 `cargo: build task api`, `cargo: build current lesson`), `launch.json`(`Rust: 현재 레슨 디버그`, `Rust: Task API 디버그`)이 들어 있습니다. 각 파일의 `//` 주석에 설명을 달아 두었습니다.
- 모든 단축키와 설정 키는 `lessons/vscode-cheatsheet.html`에 정리되어 있습니다.
