package input

import (
	"encoding/json"
	"io"
	"os"
)

type Repository struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	// Roots []*RootCommit `json:"roots"`
}

// type RootCommit struct {
// 	Hash string `json:"hash"`
// 	Date int64  `json:"date"`
// }

type Tag struct {
	Name    string    `json:"name"`
	Type    string    `json:"type"`
	Hash    string    `json:"hash"`
	Commits []*Commit `json:"commits"`
}

type Commit struct {
	Hash        string    `json:"hash"`
	Author      string    `json:"author"`
	AuthorEmail string    `json:"authorEmail"`
	Date        string    `json:"date"`
	Message     string    `json:"message"`
	Changes     []*Change `json:"changes"`
}

type Change struct {
	Insertions int64  `json:"insertions"`
	Deletions  int64  `json:"deletions"`
	Path       string `json:"path"`
}

func ReadFile(filePath string) (*Tag, error) {

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var model Tag
	if err := json.Unmarshal(bytes, &model); err != nil {
		return nil, err
	}

	return &model, nil
}
