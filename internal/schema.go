package internal

import (
	"mehub/internal/contents"
	"mehub/internal/post"
)

type PageData struct {
	Config       *contents.SiteConfig
	CurrentYear  int
	Title        string
	Posts        []post.Post
	Post         *post.Post
	Tags         []string
	TagCounts    map[string]int
	Archive      map[int][]post.Post
	ArchiveYears []int
	CurrentPage  int
	TotalPages   int
	PathPrefix   string
}

type SearchItem struct {
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	Tags        []string `json:"tags"`
}

type BlogItem struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Date        string   `json:"date_published"`
	Tags        []string `json:"skills"`
}

type ProjectItem struct {
	Title            string   `json:"title"`
	ShortDescription string   `json:"short_description"`
	Link             string   `json:"link"`
	Techs            []string `json:"tech_stack"`
}

type BlogRegistry struct {
	TotalPosts int        `json:"total_posts"`
	Posts      []BlogItem `json:"posts"`
}

type ProfileAbout struct {
	Paragraphs []string `json:"paragraphs"`
}

type ProfileRegistry struct {
	URL         string             `json:"url"`
	Title       string             `json:"title"`
	Name        string             `json:"name"`
	Slogan      string             `json:"slogan"`
	Description string             `json:"description"`
	Experience  string             `json:"experience"`
	Status      string             `json:"status"`
	FocusAreas  []string           `json:"focusAreas"`
	About       ProfileAbout       `json:"about"`
	Now         contents.NowConfig `json:"now"`
}

type Manifest struct {
	MCPVersion  string          `json:"mcp_version"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	URL         string          `json:"url"`
	UpdatedAt   string          `json:"updated_at"`
	Profile     ProfileRegistry `json:"profile"`
	Skills      []string        `json:"skills"`
	Projects    []ProjectItem   `json:"projects"`
	Blog        BlogRegistry    `json:"blog"`
}
