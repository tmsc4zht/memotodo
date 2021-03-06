package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	Version  = "unset"
	Revision = "unset"
)

func lastTodo(paths []string) (string, bool) {
	var lastDate time.Time

	for _, path := range paths {
		name := filepath.Base(path)

		t, err := time.Parse("2006-01-02-todo.md", name)
		if err != nil {
			continue
		}
		if t.After(lastDate) {
			lastDate = t
		}
	}

	if lastDate.IsZero() {
		return "", false
	}

	return lastDate.Format("2006-01-02-todo.md"), true
}

func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return nil
	}

	return ioutil.WriteFile(dst, input, 0644)
}

func run([]string) error {
	var usageFlag bool
	var versionFlag bool

	flag.BoolVar(&usageFlag, "usage", false, "print usage")
	flag.BoolVar(&versionFlag, "version", false, "print version")

	flag.Parse()

	if usageFlag {
		fmt.Println("open todo")
		return nil
	}

	if versionFlag {
		fmt.Printf("%s(%s)", Version, Revision)
		return nil
	}

	if len(flag.Args()) < 0 {
		return fmt.Errorf("argument required")
	}

	switch flag.Arg(0) {
	case "install":
		return cmdInstall()
	case "new":
		return cmdNewTodo()
	}

	return fmt.Errorf("unknown command")
}

func cmdNewTodo() error {
	memodir := os.Getenv("MEMODIR")
	filename := time.Now().Format("2006-01-02-todo.md")
	memopath := filepath.Join(memodir, filename)
	_, err := os.Stat(memopath)
	if err == nil {
		cmd := exec.Command("memo", "edit", filename)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}

	files, err := filepath.Glob(filepath.Join(memodir, "*-todo.md"))
	if err != nil {
		return fmt.Errorf("could not search memodir: %v", err)
	}

	lastTodoFile, ok := lastTodo(files)
	if !ok {
		cmd := exec.Command("memo", "new", "todo")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}

	if err := copyFile(filepath.Join(memodir, lastTodoFile), memopath); err != nil {
		return fmt.Errorf("could not copy from last todo file: %v", err)
	}

	cmd := exec.Command("memo", "edit", filename)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func cmdInstall() error {
	out, err := exec.Command("memo", "config", "--cat").Output()
	if err != nil {
		return fmt.Errorf("cannot exec memo: %f", err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))

	for scanner.Scan() {
		text := scanner.Text()
		if !strings.HasPrefix(text, "pluginsdir = ") {
			continue
		}
		pluginDirPath := strings.TrimPrefix(text, "pluginsdir = ")
		pluginDirPath = pluginDirPath[1 : len(pluginDirPath)-1]

		srcPath := os.Args[0]
		dstPath := filepath.Join(pluginDirPath, filepath.Base(srcPath))

		input, err := ioutil.ReadFile(srcPath)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(dstPath, input, 0755)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("cannot read exec memo output: %f", err)
	}
	return fmt.Errorf("could not find plugins dir path")
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	os.Exit(0)
}
