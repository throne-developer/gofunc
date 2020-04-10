package main

import (
	"bufio"
	"os"
)

func loadFile(file string) (lines []string, loadErr error) {
	lines = make([]string, 0, 1024)

	f, err := os.Open(file)
	if err != nil {
		loadErr = err
		return
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	for {
		line, err := readLine(reader)
		if err != nil {
			break
		}

		lines = append(lines, string(line))
	}
	return
}

func writeFileWithLines(file string, lines []string) (writeErr error) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		writeErr = err
		return
	}
	defer f.Close()

	for _, line := range lines {
		if _, err := f.WriteString(line + "\n"); err != nil {
			writeErr = err
			return
		}
	}
	return
}

func readLine(r *bufio.Reader) ([]byte, error) {
	line, isprefix, err := r.ReadLine()
	if !isprefix {
		return line, err
	} else {
		content := make([]byte, 0, len(line))
		content = append(content, line...)
		for isprefix && err == nil {
			line, isprefix, err = r.ReadLine()
			content = append(content, line...)
		}
		return content, err
	}
}

func deleteFile(file string) {
	os.Remove(file)
}