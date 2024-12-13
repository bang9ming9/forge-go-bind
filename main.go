package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"os"
	"path/filepath"
	"strings"
)

// flags
var (
	outFlag     = flag.String("out", "", "output file")
	packageFlag = flag.String("pkg", "bindings", "package name for the generated file")
)

func main() {

	// 1. flag 입력값을 확인한다.
	pkg, outDir, rootDir, err := checkAndGetArgs()
	if err != nil {
		fmt.Println("fail to check args:", err)
		os.Exit(1)
	}

	// 2. dir/src 에 있는 .sol 파일의 이름을 캐싱한다.
	srcContracts := make([]string, 0)
	err = filepath.Walk(filepath.Join(rootDir, "src"), func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".sol" {
			srcContracts = append(srcContracts, strings.TrimSuffix(filepath.Base(path), ".sol"))
		}
		return nil
	})
	if err != nil {
		fmt.Println("fail to read directory", err, "\n", rootDir)
		os.Exit(1)
	}
	// 2-1. .sol 파일이 없으면 종료한다.
	if len(srcContracts) == 0 {
		fmt.Println("no contracts found")
		os.Exit(1)
	}

	// 3. dir/out 에서 srcContracts 의 .sol 의 디렉토리를 순환하며 .json 파일을 읽는다.
	bindMap := make(map[string]Bind)
	for _, path := range srcContracts {
		readDir := filepath.Join(rootDir, "out", path+".sol")
		err := filepath.Walk(readDir, func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".json" {
				if bytes, err := os.ReadFile(path); err != nil {
					fmt.Println("fail to read file:", err, "\n", path)
				} else {
					b := Bind{}
					if err := json.Unmarshal(bytes, &b); err != nil {
						fmt.Println("fail to unmarshal:", err, "\n", path)
					} else {
						name := strings.TrimSuffix(filepath.Base(path), ".json")
						bindMap[name] = b
					}
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println("fail to read directory", err, "\n", readDir)
			os.Exit(1)
		}
	}

	// 4. go-ethereum 의 Bind 함수를 호출한다.
	for name, data := range bindMap {
		bytes, err := bind.Bind(data.Args(pkg, name))
		if err != nil {
			fmt.Println("fail to bind contracts:", err, "\n", name)
			os.Exit(1)
		}

		// 4-1. 결과값을 go 파일로 저장한다.
		err = os.WriteFile(filepath.Join(outDir, name)+".go", []byte(bytes), 0644)
		if err != nil {
			fmt.Println("fail to write file:", err, "\n", name)
			os.Exit(1)
		}
	}
	fmt.Println("success")
	os.Exit(0)
}
