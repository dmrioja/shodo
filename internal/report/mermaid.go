package markdown

import (
	"fmt"
	"shodo/internal/output"
	"strings"
)

func ChangesSankey(tagReport *output.Tag) (md []string) {
	md = append(md, "```mermaid")
	md = append(md, "%%{init:{\"sankey\":{\"nodeAlignment\":\"left\"}}}%%")
	md = append(md, "sankey-beta")

	for _, commitType := range tagReport.GetCommitsOrderedByNetChanges() {
		md = append(md, fmt.Sprintf("Gross changes,%s,%d", strings.ToUpper(commitType.Label), commitType.GetGrossChanges()))
		md = append(md, fmt.Sprintf("%s,%s,%.0f", strings.ToUpper(commitType.Label), commitType.Label, commitType.GetNetChanges()))
		if commitType.GetNetChanges() != 0 {
			md = append(md, fmt.Sprintf("%s,%s,%.0f", commitType.Label, "Net changes", commitType.GetNetChanges()))
		}
	}
	md = append(md, "```")

	return md
}

func AvgFilesChangedPerCommitType(tagReport *output.Tag) (md []string) {
	md = append(md, "\n```mermaid")
	md = append(md, "%%{init:{\"themeVariables\":{\"xyChart\":{\"backgroundColor\":\"transparent\"}}}}%%")
	md = append(md, "xychart-beta")
	md = append(md, "title \"Avg Files modified per commit type\"")

	// It would be interesting to draw this along time, to see how repository complexity changes over time
	// See property title

	// avgfiles modified / total files in the repository -> this could be an index to see:
	// Is the repository well structured? How difficult it is to add new changes??

	mdXAxis := "x-axis ["
	bars := "bar ["
	for _, commitType := range tagReport.GetCommitsOrderedByAvgFilesModified() {
		mdXAxis = fmt.Sprintf("%s%s,", mdXAxis, commitType.Label)
		bars = fmt.Sprintf("%s%.0f,", bars, commitType.GetAvgFilesModified())
	}
	mdXAxis = mdXAxis[:len(mdXAxis)-1]
	mdXAxis = fmt.Sprintf("%s]", mdXAxis)

	bars = bars[:len(bars)-1]
	bars = fmt.Sprintf("%s]", bars)

	md = append(md, mdXAxis)
	md = append(md, bars)

	md = append(md, "```")

	return
}
