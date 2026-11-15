# Go 실전 입문 — 시작하기

Spring·Kotlin 경험자를 위한 Go 학습 코스입니다. 레슨은 모두 한국어 HTML 문서이고, 문서를 읽으면서 VS Code에서 그대로 따라 합니다.

## 시작 순서 (약 5분)

1. **폴더 열기**: VS Code 메뉴 `File` > `Open Folder...`(`Ctrl+K Ctrl+O`)로 이 폴더(`D:\study\go`)를 엽니다. 파일 하나만 여는 `Open File...`이 아닙니다.
2. **신뢰하기**: 제한 모드(Restricted Mode) 배너가 보이면 `Manage`에서 이 폴더를 신뢰합니다. 신뢰하지 않으면 Go 확장이 동작하지 않습니다. "부모 폴더 전체를 신뢰"하는 체크 상자는 켜지 마세요.
3. **권장 확장 설치**: `Ctrl+Shift+X`로 확장 보기를 열고 `@recommended`를 검색해 `golang.go`, `humao.rest-client`를 설치합니다(`.vscode/extensions.json`에 적혀 있습니다).
4. **문서 열기**: 탐색기에서 `index.html`을 연 상태로 `Ctrl+K R`을 누르거나 파일을 우클릭해 `Reveal in File Explorer`를 고른 뒤, 파일 탐색기에서 더블 클릭하면 브라우저로 열립니다.
5. **첫 레슨**: `lessons/00-vscode-start.html`(VS Code 첫걸음)부터 읽습니다. 다음은 `lessons/01-environment.html`입니다.

## 폴더 구성

```
D:\study\go
├─ index.html        코스 홈
├─ lessons\          레슨 HTML (00~11)와 VS Code 치트시트
├─ lab\              레슨 실습 코드 (Go 모듈 golab)
├─ .vscode\          VS Code 설정 (권장 확장, 설정, 작업, 디버그)
├─ go.work           Go 작업 공간 설정 (lab 모듈 등록)
└─ tour-of-go, go-by-example, learn-go-with-tests, 100-go-mistakes
                     클론한 읽기 자료 (수정하지 않습니다. GitHub에서 저장소를 새로 받았다면
                     이 폴더들은 없으니 저장소 맨 위 README.md의 git clone 명령으로 받습니다)
```

## 알아 둘 점

- **기본 터미널이 MSYS2 bash이면** `go` 명령을 못 찾을 수 있습니다. MSYS2 bash는 Windows의 PATH를 이어받지 않는 것이 기본 동작이라, 사용자 프로필(`--login -i`로 시작하는 `MSYS2 UCRT64`)과 같은 조건에서 `go`가 "명령을 찾을 수 없음"으로 나오는 것을 확인했습니다. 터미널 오른쪽 위 `+` 옆 화살표에서 PowerShell 프로필로 새 터미널을 열면 됩니다. bash를 계속 쓰려면 사용자 설정의 해당 프로필 `env`에 `"MSYS2_PATH_TYPE": "inherit"`를 추가합니다(적용 후 `go`가 찾아지는 것을 확인했습니다).
- **go.work**: `lab` 모듈 하나만 등록되어 있어 gopls가 `lab`을 바로 인식하고, 폴더 어디서든 `go run ./lab/lesson03`처럼 실행할 수 있습니다. 클론한 읽기 자료 폴더 안에서 `go test` 등을 직접 돌려 보려면 `GOWORK=off`가 필요합니다(PowerShell은 `$env:GOWORK = "off"`). 단, `go-by-example/go.mod`는 `go 1.27.1`을 요구하므로 설치된 1.27.0에서는 도구 체인 내려받기가 시도될 수 있습니다. 클론 폴더는 읽기만 하세요.
- **VS Code 설정**: `.vscode/settings.json`(저장 시 서식 정리, 인레이 힌트), `tasks.json`(`go: test all`, `go: vet all`, `go: run task api`), `launch.json`(`Go: 현재 파일 디버그`, `Go: Task API 디버그`)이 들어 있습니다. 각 파일의 `//` 주석에 설명을 달아 두었습니다.
- 모든 단축키와 설정 키는 `lessons/vscode-cheatsheet.html`에 정리되어 있습니다.
