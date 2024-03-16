package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"shodo/internal/input"
	"strings"
)

func ReadFile(filePath string) (*input.Tag, error) {

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var model input.Tag
	if err := json.Unmarshal(bytes, &model); err != nil {
		return nil, err
	}

	return &model, nil
}

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
