package output

import (
	"math"
	"regexp"
	"shodo/internal/input"
	"sort"
)

var colors map[string]string

func init() {
	colors = make(map[string]string)
	colors["feat"] = "#27ae60"
	colors["fix"] = "#ff3100"
	colors["build"] = "#d6dbdf"
	colors["chore"] = "#808b96"
	colors["ci"] = "#8e44ad"
	colors["docs"] = "#2874a6"
	colors["style"] = "#f682c5"
	colors["perf"] = "#f9e79f"
	colors["refactor"] = "#5dade2"
	colors["test"] = "#76d7c4"
	colors["patch"] = "#f1c40f"
	colors["others"] = "#000000"
}

type Report struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	// Roots []*RootCommit `json:"roots"`
}

// type RootCommit struct {
// 	Hash string `json:"hash"`
// 	Date int64  `json:"date"`
// }

type Tag struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Hash        string `json:"hash"`
	LastVersion string

	TotalConventionalCommits int
	TotalNetChanges          int
	TotalGrossChanges        int

	ConventionalCommits map[string]*ConventionalCommit
}

func (t *Tag) GetNonConventionalCommits() []*Commit {
	if ncc, ok := t.ConventionalCommits["others"]; ok {
		return ncc.Commits
	}
	return nil
}

// TODO: This should be renamed since it also includes non conventional commits
type ConventionalCommit struct {
	Label        string
	TotalCommits int64
	Insertions   int64
	Deletions    int64
	FilesChanged int
	Color        string
	Commits      []*Commit
	// TODO: we could add a list of commits (date and message of each commit could be interesting)
}

type Commit struct {
	Message string
	Date    string
}

// TODO: remove?
func (cc *ConventionalCommit) GetNetChanges() float64 {
	return math.Abs(float64(cc.Insertions - cc.Deletions))
}

// TODO: remove?
func (cc *ConventionalCommit) GetGrossChanges() int64 {
	return cc.Insertions + cc.Deletions
}

// TODO: remove?
func (cc *ConventionalCommit) GetAvgFilesModified() float64 {
	return float64(cc.FilesChanged) / float64(cc.TotalCommits)
}

func (t *Tag) CalculateAggregates() {
	for _, cc := range t.ConventionalCommits {
		t.TotalNetChanges += int(cc.Insertions - cc.Deletions)
		t.TotalGrossChanges += int(cc.Insertions + cc.Deletions)
	}
	// TODO: need to agreggate non-conventional commits
}

func (t *Tag) GetTotalCommits() (total int64) {
	// TODO: We call here a lot, this may be calculated only once
	for _, cmtType := range t.ConventionalCommits {
		total += cmtType.TotalCommits
	}
	return
}

func (t *Tag) NewCommit(newCommit *input.Commit) {
	if cmtType, ok := extractLabel(newCommit); ok {
		t.AddNewCommit(cmtType, newCommit)
	} else {
		t.AddNewCommit("others", newCommit)
		// TODO: refactor this and make it general
		ncc := t.ConventionalCommits["others"]
		ncc.Commits = append(ncc.Commits, &Commit{
			Message: newCommit.Message,
			Date:    newCommit.Date,
		})
	}
}

func extractLabel(commit *input.Commit) (string, bool) {
	// This validation should be more strict
	match := regexp.MustCompile(
		"^(feat|fix|build|chore|ci|docs|style|perf|refactor|test|revert|fixup)",
	).FindStringSubmatch(commit.Message)

	// Not a pretty orthodox validation
	if len(match) == 2 {
		return match[1], true
	}

	// TODO: is this really necessary??
	return "", false
}

func (t *Tag) AddNewCommit(cmtType string, newCommit *input.Commit) {
	commitType, ok := t.ConventionalCommits[cmtType]

	if !ok {
		if t.ConventionalCommits == nil {
			t.ConventionalCommits = make(map[string]*ConventionalCommit)
		}

		commitType = &ConventionalCommit{
			Label: cmtType,
			Color: colors[cmtType],
		}
		t.ConventionalCommits[cmtType] = commitType
	}
	commitType.TotalCommits++

	for _, change := range newCommit.Changes {
		commitType.Insertions += change.Insertions
		commitType.Deletions += change.Deletions
	}
	commitType.FilesChanged += len(newCommit.Changes)
}

// Is this really necessary ?? (Seems like only for markdown)
func (t *Tag) GetCommitsOrderedByQuantity() []*ConventionalCommit {
	return t.getCommitsOrderedBy(
		func(cc *ConventionalCommit) float64 {
			return float64(cc.TotalCommits)
		},
	)
}

func (t *Tag) GetCommitsOrderedByNetChanges() []*ConventionalCommit {
	return t.getCommitsOrderedBy(
		func(cc *ConventionalCommit) float64 {
			return math.Abs(float64(cc.Insertions - cc.Deletions))
		},
	)
}

func (t *Tag) GetCommitsOrderedByAvgFilesModified() []*ConventionalCommit {
	return t.getCommitsOrderedBy(
		func(cc *ConventionalCommit) float64 {
			return float64(cc.FilesChanged) / float64(cc.TotalCommits)
		},
	)
}

func (t *Tag) getCommitsOrderedBy(criteria func(*ConventionalCommit) float64) []*ConventionalCommit {
	s := make([]*ConventionalCommit, 0, len(t.ConventionalCommits))

	for _, c := range t.ConventionalCommits {
		s = append(s, c)
	}

	sort.Slice(s, func(i, j int) bool {
		return criteria(s[i]) > criteria(s[j])
	})

	return s
}
