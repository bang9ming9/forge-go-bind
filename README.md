# forge-go-bind

foundry 으로 작업한 컨트랙트를 go 언어로 바인딩하는 도구입니다.<br>
{project} 의 src 하위에 있는 컨트랙트만 필터하여 바인딩 합니다.<br>

## build
```bash
git clone https://github.com/bang9ming9/forge-go-bind.git
cd forge-go-bind
go build -o {PATH} .
```

## Option
```
  -out string
        output file
  -pkg string
        package name for the generated file (default "bindings")
```

## Usage
```bash
./forge-go-bind -out {OUTPUT_FILE} -pkg {PACKAGE_NAME} {FOUNDRY_PROJECT_PATH}
```