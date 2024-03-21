package input

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
	Hash    string    `json:"hash"`
	Date    string    `json:"date"`
	Message string    `json:"message"`
	Changes []*Change `json:"changes"`
}

type Change struct {
	Insertions int64  `json:"insertions"`
	Deletions  int64  `json:"deletions"`
	Path       string `json:"path"`
}
