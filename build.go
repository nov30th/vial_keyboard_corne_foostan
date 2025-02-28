//go:build ignore
// +build ignore

// -ldflags "-s -w"
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// 创建 build 目录
	os.MkdirAll("build", 0755)

	// 定义要构建的目标平台
	platforms := []struct {
		os   string
		arch string
		ext  string
		args []string
	}{
		{"windows", "amd64", ".exe", []string{"-ldflags", "-H windowsgui -s -w"}},
		{"windows", "386", "-x86.exe", []string{"-ldflags", "-H windowsgui -s -w"}},
		{"darwin", "amd64", "-amd64", []string{"-ldflags", "-s -w"}},
		{"darwin", "arm64", "-arm64", []string{"-ldflags", "-s -w"}},
		{"linux", "amd64", "", []string{"-ldflags", "-s -w"}},
	}

	for _, platform := range platforms {
		env := append(os.Environ(),
			fmt.Sprintf("GOOS=%s", platform.os),
			fmt.Sprintf("GOARCH=%s", platform.arch),
		)

		outputName := filepath.Join("build", fmt.Sprintf("keyboard-converter%s", platform.ext))
		args := append([]string{"build"}, platform.args...)
		args = append(args, "-o", outputName, ".")

		cmd := exec.Command("go", args...)
		cmd.Env = env
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Printf("构建 %s/%s...\n", platform.os, platform.arch)
		if err := cmd.Run(); err != nil {
			fmt.Printf("构建失败: %v\n", err)
		} else {
			fmt.Printf("构建成功: %s\n", outputName)
		}
	}
}
