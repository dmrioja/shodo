package main

import (
	"fmt"
	"io/fs"
	"os"
	"shodo/internal/input"
	"shodo/internal/output"
	markdown "shodo/internal/report"
	"shodo/internal/utils"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
)

var fileName string = os.Args[2]
var filePath string = fmt.Sprintf("../../data/%s/", fileName)

func main() {

	files := readFiles()
	tags := make([]*output.Tag, 0, len(files))
	for index, file := range files {
		// TODO: deal better with first version
		previousVersion := "v0.1.0"
		if index != 0 {
			previousVersion = getVersion(files[index-1].Name())
		}
		tags = append(tags, processFile(file.Name(), previousVersion))
	}

	generalReport(tags)

	os.Exit(0)
}

func readFiles() (tagFiles []fs.FileInfo) {

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sort.Slice(files, func(i, j int) bool {
		vi, _ := version.NewVersion(getVersion(files[i].Name()))
		vj, _ := version.NewVersion(getVersion(files[j].Name()))

		return vi.LessThan(vj)
	})

	return files
}

func processFile(fileName, previousVersion string) *output.Tag {

	// TODO: improve this later. I don't feel like doing it now
	tagInput, err := utils.ReadFile(filePath + fileName)
	if err != nil {
		fmt.Printf("Error reading input file [%s]: %s", fileName, err)
		os.Exit(1)
	}

	tagReport := generateTagReport(tagInput, previousVersion)

	if err := markdown.WriteReport(tagReport); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return tagReport
}

func generateTagReport(input *input.Tag, previousVersion string) *output.Tag {
	tag := &output.Tag{
		Name:        input.Name,
		Type:        input.Type,
		Hash:        input.Hash,
		LastVersion: previousVersion,
	}

	for _, commit := range input.Commits {
		tag.NewCommit(commit)
	}
	tag.CalculateAggregates()

	return tag
}

func generalReport(tags []*output.Tag) {
	sort.Slice(tags, func(i, j int) bool {
		vi, _ := version.NewVersion(tags[i].Name)
		vj, _ := version.NewVersion(tags[j].Name)
		return vi.LessThan(vj)
	})

	var md []string
	md = append(md, "```mermaid")
	md = append(md, "xychart-beta")
	md = append(md, "title \"New features along time\"")

	xAxis := "x-axis ["
	line := "line ["

	for _, tag := range tags {
		xAxis = fmt.Sprintf("%s%s,", xAxis, tag.Name)
		line = fmt.Sprintf("%s%d,", line, getTotalCommitsPerType(tag, "fix"))
		// line = fmt.Sprintf("%s%f,", line, tag.GetAvgFilesModifiedPer("feat"))
	}

	xAxis = xAxis[:len(xAxis)-1]
	xAxis = fmt.Sprintf("%s]", xAxis)

	line = line[:len(line)-1]
	line = fmt.Sprintf("%s]", line)

	md = append(md, xAxis)
	md = append(md, line)

	md = append(md, "```")

	// Deal with this error better later
	if err := utils.WriteFile(fileName+"_SHODO.md", md); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getTotalCommitsPerType(tag *output.Tag, label string) int64 {
	if commitType, ok := tag.ConventionalCommits[label]; ok {
		return commitType.TotalCommits
	}
	return 0
}

func getVersion(tagFileName string) string {
	v, found := strings.CutSuffix(tagFileName, "_shodo_report.json")

	if !found {
		fmt.Printf("Error extracting version from file %s\n", tagFileName)
		os.Exit(1)
	}

	return v
}
