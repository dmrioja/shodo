package markdown

import (
	"fmt"
	"shodo/internal/output"
	"shodo/internal/utils"
	"strings"
)

func WriteReport(tagReport *output.Tag) error {

	var md []string
	md = append(md, fmt.Sprintf("# %s ... %s", tagReport.LastVersion, tagReport.Name))

	md = append(md, writeConventionalCommitDetail(tagReport)...)
	// TODO: manage errors better
	return utils.WriteFile(fmt.Sprintf("%s_SHODO.md", tagReport.Name), md)
}

func writeConventionalCommitDetail(tagReport *output.Tag) []string {

	// nolint: prealloc
	var md, mdConfig, mdPieChart []string
	md = append(md, "``` mermaid")
	mdConfig = append(mdConfig, "%%{init:{'theme':'base','themeVariables':{'pieTitleTextColor': '#fff','pieLegendTextColor': '#fff',")

	// TODO: I'm assuming ConventionalCommits map is never empty, which is not a clever assumption
	mdPieChart = append(mdPieChart, fmt.Sprintf("pie showData title Number of commits: %d", tagReport.GetTotalCommits()))

	commits := tagReport.GetCommitsOrderedByQuantity()
	var mdConfigString string
	for index, commitType := range commits {
		mdConfigString = fmt.Sprintf("%s'pie%d':'%s',", mdConfigString, index+1, commitType.Color)
		mdPieChart = append(mdPieChart, fmt.Sprintf("\"%s\": %d", commitType.Label, commitType.TotalCommits))
	}

	mdConfig = append(mdConfig, strings.TrimSuffix(mdConfigString, ","))
	mdConfig = append(mdConfig, "}}}%%")

	md = append(md, mdConfig...)
	md = append(md, mdPieChart...)
	md = append(md, "```")

	// Non conventional commits warning
	// TODO: It would be interesting if we could filter a specific author (eg. CI bot)
	ncc := tagReport.GetNonConventionalCommits()
	if len(ncc) > 0 {
		md = append(md, "> **WARNING**")
		// TODO: take care of singular and plural
		md = append(md, fmt.Sprintf(
			"> <details><summary>%d non conventional commits found:</summary><ul>", len(ncc)))

		for _, commit := range ncc {
			md = append(md, fmt.Sprintf("> <li>%s</li>", commit.Message))
		}

		md = append(md, "> </ul></details>")
	}

	md = append(md, ChangesSankey(tagReport)...)
	md = append(md, AvgFilesChangedPerCommitType(tagReport)...)
	md = append(md, Styles()...)
	//md = append(md, AvgFilesChangedPerCommitType(tagReport)...)

	//md = append(md, fmt.Sprintf())

	return md
}
