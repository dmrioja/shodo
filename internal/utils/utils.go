package utils

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func WriteFile(fileName string, markdown []string) error {
	file, err := os.Create(fmt.Sprintf("output/%s", fileName))
	if err != nil {
		return fmt.Errorf("ERROR creating output/%s file: %s", fileName, err)
	}
	defer file.Close()

	if _, err := file.Write([]byte(strings.Join(markdown, lineFeed()))); err != nil {
		return fmt.Errorf("ERROR writing into output/%s file: %s", fileName, err)
	}

	return nil
}

// TODO: check this better
func lineFeed() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
