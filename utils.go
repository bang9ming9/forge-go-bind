package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkAndGetArgs() (
	pkg string,
	outDir string,
	rootDir string,
	err error,
) {
	flag.Parse()

	arg0 := flag.Arg(0)
	if arg0 == "" {
		arg0 = "."
	}

	pkg = *packageFlag
	if pkg == "" {
		err = errors.New("package name is required")
	} else if outDir, err = filepath.Abs(*outFlag); err != nil {
		err = fmt.Errorf("failed to get absolute path: %w", err)
	} else if rootDir, err = filepath.Abs(arg0); err != nil {
		err = fmt.Errorf("failed to get absolute path: %w", err)
	} else {
		// 출력 디렉토리가 package 으로 끝나지 않으면 경로에 추가한다.
		if !strings.HasSuffix(outDir, *packageFlag) {
			outDir = filepath.Join(outDir, *packageFlag)
		}
		// 출력 디렉토라가 존재하지 않으면 생성한다.
		if _, err = os.Stat(outDir); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return
			}
			if err = os.MkdirAll(outDir, 0755); err != nil {
				err = fmt.Errorf("failed to create directory: %w", err)
				return
			}
			err = nil
		}

		// 입력된 프로젝트 경로를 확인한다.
		if info, e := os.Stat(rootDir); e != nil {
			err = e
		} else if !info.IsDir() {
			err = errors.New("not a directory")
		} else {
			if strings.HasSuffix(rootDir, "test") ||
				strings.HasSuffix(rootDir, "out") ||
				strings.HasSuffix(rootDir, "src") {
				rootDir = filepath.Join(rootDir, "..")
			}
		}
	}

	return
}
