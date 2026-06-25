// Package internal defines the core static site generator logic and unified data schemas.
package internal

import (
	"time"
)

// ============================================================================
// Site & Navigation Configuration Schemas
// ============================================================================

// NavItem represents a single hyperlink item in the site navigation headers or footers.
type NavItem struct {
	Href string `yaml:"href"`
	Text string `yaml:"text"`
}

// NavigationConfig holds the lists of NavItem structures for header and footer navigation.
type NavigationConfig struct {
	Header []NavItem `yaml:"header"`
	Footer []NavItem `yaml:"footer"`
}

// Social holds connection link and display metadata for external social networking platforms.
type Social struct {
	Name string `yaml:"name"`
	Href string `yaml:"href"`
	Icon string `yaml:"icon"`
}

// Project outlines portfolio project attributes, linking, and technological classification tags.
type Project struct {
	Title            string   `yaml:"title"`
	ShortDescription string   `yaml:"shortDescription"`
	Link             string   `yaml:"link"`
	Techs            []string `yaml:"techs"`
	Emoji            string   `yaml:"emoji"`
}

// Skill represents a technology tag paired with a matching display icon or CSS utility.
type Skill struct {
	Name string `yaml:"name"`
	Icon string `yaml:"icon"`
}

// NowCategory partitions recent activity or interests under a specific topic heading.
type NowCategory struct {
	Title string   `yaml:"title"`
	Items []string `yaml:"items"`
}

// AboutConfig specifies self-descriptive imagery and multi-paragraph bio details.
type AboutConfig struct {
	Image      string   `yaml:"image"`
	Paragraphs []string `yaml:"paragraphs"`
}

// NowConfig aggregates dynamic recent status information, logs, and categories.
type NowConfig struct {
	LastUpdated string        `yaml:"lastUpdated"`
	Categories  []NowCategory `yaml:"categories"`
}

// SiteMetadata gathers primary context on site-wide branding, status, experience, and biographic data.
type SiteMetadata struct {
	URL        string      `yaml:"url"`
	Title      string      `yaml:"title"`
	Name       string      `yaml:"name"`
	Slogan     string      `yaml:"slogan"`
	Experience string      `yaml:"experience"`
	Status     string      `yaml:"status"`
	FocusAreas []string    `yaml:"focusAreas"`
	About      AboutConfig `yaml:"about"`
	Now        NowConfig   `yaml:"now"`
}

// SiteConfig represents the full composite profile parsed from YAML configurations in the repository.
type SiteConfig struct {
	Site       SiteMetadata     `yaml:"site"`
	Navigation NavigationConfig `yaml:"navigation"`
	Socials    []Social         `yaml:"socials"`
	Projects   []Project        `yaml:"projects"`
	Skills     []Skill          `yaml:"skills"`
}

// ============================================================================
// Blog Post & Content Schemas
// ============================================================================

// Frontmatter represents metadata defined in the YAML header of Markdown blog files.
type Frontmatter struct {
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	Date        time.Time `yaml:"date"`
	Tags        []string  `yaml:"tags"`
	Draft       bool      `yaml:"draft"`
}

// RelatedPost maps target link slugs for displaying behavior-related posts in templates.
type RelatedPost struct {
	Title string
	Slug  string
}

// Post encapsulates a full blog item, linking frontmatter metadata with its converted HTML content body.
type Post struct {
	Frontmatter
	Slug         string
	Content      string
	RelatedPosts []RelatedPost
}

// ContentData bundles loaded blog contents, pre-grouped index tables, and tag analytics.
type ContentData struct {
	Posts        []Post
	PostsByTag   map[string][]Post
	PostsByYear  map[int][]Post
	Tags         []string
	TagCounts    map[string]int
	ArchiveYears []int
}

// ============================================================================
// Page Template Data & API Registry Manifests
// ============================================================================

// PageData represents the context contextually applied to HTML templates during site generation.
type PageData struct {
	Config       *SiteConfig
	CurrentYear  int
	Title        string
	Posts        []Post
	Post         *Post
	Tags         []string
	TagCounts    map[string]int
	Archive      map[int][]Post
	ArchiveYears []int
	CurrentPage  int
	TotalPages   int
	PathPrefix   string
}

// SearchItem maps structure for index searching on the frontend search index payload.
type SearchItem struct {
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	Tags        []string `json:"tags"`
}

// BlogItem maps public API details of individual blog items for the registry manifest.
type BlogItem struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Date        string   `json:"date_published"`
	Tags        []string `json:"skills"`
}

// ProjectItem maps public API details of projects for the registry manifest.
type ProjectItem struct {
	Title            string   `json:"title"`
	ShortDescription string   `json:"short_description"`
	Link             string   `json:"link"`
	Techs            []string `json:"tech_stack"`
}

// BlogRegistry lists posts metadata on the public blog manifest registry.
type BlogRegistry struct {
	TotalPosts int        `json:"total_posts"`
	Posts      []BlogItem `json:"posts"`
}

// ProfileAbout defines paragraphs of bio context on the public manifest.
type ProfileAbout struct {
	Paragraphs []string `json:"paragraphs"`
}

// ProfileRegistry lists professional information and latest status details for AI queries.
type ProfileRegistry struct {
	URL        string       `json:"url"`
	Title      string       `json:"title"`
	Name       string       `json:"name"`
	Slogan     string       `json:"slogan"`
	Experience string       `json:"experience"`
	Status     string       `json:"status"`
	FocusAreas []string     `json:"focusAreas"`
	About      ProfileAbout `json:"about"`
	Now        NowConfig    `json:"now"`
}

// Manifest defines the top-level Model Context Protocol (MCP) compatible structure describing the site context.
type Manifest struct {
	MCPVersion string          `json:"mcp_version"`
	Name       string          `json:"name"`
	URL        string          `json:"url"`
	UpdatedAt  string          `json:"updated_at"`
	Profile    ProfileRegistry `json:"profile"`
	Skills     []string        `json:"skills"`
	Projects   []ProjectItem   `json:"projects"`
	Blog       BlogRegistry    `json:"blog"`
}
