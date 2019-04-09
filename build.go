// Copyright © 2019 Jan Arens
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type information struct {
	version        string
	commit         string
	date           string
	goVersion      string
	buildMachineOs string
}

var info information

func main() {
	setVersion()
	setCommit()
	setDate()
	setGoVersion()
	setBuildMachineOs()

	build()
}

func setVersion() {
	out, err := exec.Command("git", "describe", "--tags").CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get version from git tag: Error: %s\n", err)
		return
	}

	info.version = strings.Trim(string(out), "\n")

	fmt.Println("Version        :", info.version)
}

func setCommit() {
	out, err := exec.Command("git", "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get commit from git: Error: %s\n", err)
		return
	}

	info.commit = strings.Trim(string(out), "\n")

	fmt.Println("Commit         :", info.commit)
}

func setDate() {
	info.date = strconv.FormatInt(time.Now().Unix(), 10)

	fmt.Println("Date           :", info.date)
}

func setGoVersion() {
	info.goVersion = runtime.Version()

	fmt.Println("Go             :", info.goVersion)
}

func setBuildMachineOs() {
	info.buildMachineOs = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	fmt.Println("Build on       :", info.buildMachineOs)
}

func build() {
	osList := []string{"windows", "linux"}
	archList := []string{"amd64", "386"}

	for _, buildOs := range osList {
		for _, buildArch := range archList {
			binPath := fmt.Sprintf("bin/backup-and-sync_%s_%s", buildOs, buildArch)
			if buildOs == "windows" {
				binPath += ".exe"
			}

			ldflags := fmt.Sprintf(`-X "main.Version=%s" -X "main.Commit=%s" -X "main.Date=%s" -X "main.GoVersion=%s" -X "main.BuildMachineOs=%s"`,
				info.version, info.commit, info.date, info.goVersion, info.buildMachineOs)

			args := []string{"build", "-ldflags", ldflags, "-o", binPath, "github.com/th3noname/backup-and-sync/src"}

			fmt.Printf("Starting go build (GOOS=%s GOARCH=%s) with arguments: %s\n", buildOs, buildArch, strings.Join(args, " "))

			command := exec.Command("go", args...)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			command.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", buildOs), fmt.Sprintf("GOARCH=%s", buildArch))
			err := command.Run()

			if err != nil {
				fmt.Fprintf(os.Stderr, "go build failed: GOOS=%s GOARCH=%s Error: %s\n", buildOs, buildArch, err)
				os.Exit(1)
			}
		}
	}
}
