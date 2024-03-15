package markdown

import (
	"fmt"
	"shodo/internal/output"
)

func ChangesSankey(tagReport *output.Tag) (md []string) {
	md = append(md, "<div class=\"changes-wrapper\">")
	md = append(md, "<div class=\"gross-changes\">\n")
	md = append(md, grossChangesSankey(tagReport)...)
	md = append(md, "\n</div>")
	md = append(md, "<div class=\"net-changes\">\n")
	md = append(md, netChangesSankey(tagReport)...)
	md = append(md, "\n</div>")
	md = append(md, "</div>")

	return
}

func netChangesSankey(tagReport *output.Tag) (md []string) {
	md = append(md, "```mermaid")
	md = append(md, "sankey-beta")
	md = append(md, "%% source,target,value")

	for _, commitType := range tagReport.GetCommitsOrderedByNetChanges() {
		md = append(md, fmt.Sprintf("Net changes,%s,%.0f", commitType.Label, commitType.GetNetChanges()))
	}
	md = append(md, "```")

	return md
}

func grossChangesSankey(tagReport *output.Tag) (md []string) {
	md = append(md, "```mermaid")
	md = append(md, "sankey-beta")
	md = append(md, "%% source,target,value")

	for _, commitType := range tagReport.GetCommitsOrderedByNetChanges() {
		md = append(md, fmt.Sprintf("Gross changes,%s,%d", commitType.Label, commitType.GetGrossChanges()))
	}
	md = append(md, "```")

	return md
}

func AvgFilesChangedPerCommitType(tagReport *output.Tag) (md []string) {
	md = append(md, "\n```mermaid")
	md = append(md, "xychart-beta")
	md = append(md, "title \"Avg Files modified per commit type\"")

	// It would be interesting to draw this along time, to see how repository complexity changes over time
	// See property title

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

func Styles() (md []string) {
	md = append(md, "<style>")
	md = append(md, ".changes-wrapper{display:inline-block;}")
	md = append(md, ".gross-changes{display:block;}")
	md = append(md, ".net-changes{display:none;}")
	md = append(md, ".changes-wrapper:hover .gross-changes{display:none;}")
	md = append(md, ".changes-wrapper:hover .net-changes{display:block;}")
	md = append(md, "</style>")

	return
}
