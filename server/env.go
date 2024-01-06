package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type environmentVars map[string]string

func loadEnvFile(path string) environmentVars {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	var lines []string
	for {
		line, _, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		lines = append(lines, string(line))
	}
	var res = make(environmentVars)
	for _, line := range lines {
		var keyVal = []string{}
		keyVal = strings.Split(line, "=")
		if len(keyVal) != 2 {
			msg := fmt.Sprintf("%s (invalid key=val pair)", line)
			panic(msg)
		}
		res[keyVal[0]] = keyVal[1]
	}
	return res
}
