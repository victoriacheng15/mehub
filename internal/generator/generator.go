package generator

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mehub-ssg/internal/config"
	"mehub-ssg/internal/content"
)

type PageData struct {
	Config       *config.SiteConfig
	CurrentYear  int
	Title        string
	Posts        []content.Post
	Post         *content.Post
	Tags         []string
	TagCounts    map[string]int
	Archive      map[int][]content.Post
	ArchiveYears []int
	CurrentPage  int
	TotalPages   int
	PathPrefix   string
}

type SiteGenerator struct {
	Config  *config.SiteConfig
	FuncMap template.FuncMap
}

func New(cfg *config.SiteConfig) *SiteGenerator {
	return &SiteGenerator{
		Config: cfg,
		FuncMap: template.FuncMap{
			"split":             strings.Split,
			"replace":           strings.ReplaceAll,
			"trimSpace":         strings.TrimSpace,
			"stringsHasPrefix":  strings.HasPrefix,
			"stringsTrimPrefix": strings.TrimPrefix,
			"add":               func(a, b int) int { return a + b },
			"sub":               func(a, b int) int { return a - b },
			"safeHTML":          func(s string) template.HTML { return template.HTML(s) },
			"cleanYAMLList": func(data interface{}) []string {
				var input string
				switch v := data.(type) {
				case string:
					input = v
				case []string:
					return v
				default:
					return nil
				}

				lines := strings.Split(strings.TrimSpace(input), "\n")
				var result []string
				for _, line := range lines {
					cleaned := strings.TrimSpace(line)
					cleaned = strings.TrimPrefix(cleaned, "- ")
					cleaned = strings.Trim(cleaned, "\"")
					if cleaned != "" {
						result = append(result, cleaned)
					}
				}
				return result
			},
		},
	}
}

func (g *SiteGenerator) RenderPage(dir, filename, tmplPath string, titlePrefix string, data PageData) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create dir %s: %v", dir, err)
	}

	tmpl, err := template.New("base.html").Funcs(g.FuncMap).ParseFiles("internal/templates/base.html", tmplPath)
	if err != nil {
		log.Fatalf("Failed to parse templates for %s: %v", tmplPath, err)
	}

	outputFile, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		log.Fatalf("Failed to create output file %s: %v", filename, err)
	}
	defer outputFile.Close()

	title := g.Config.Site.Title
	if titlePrefix != "" {
		title = titlePrefix + " | " + title
	}

	data.Config = g.Config
	data.CurrentYear = time.Now().Year()
	data.Title = title

	if err := tmpl.ExecuteTemplate(outputFile, "base.html", data); err != nil {
		log.Fatalf("Failed to execute template for %s: %v", tmplPath, err)
	}
}

func (g *SiteGenerator) GenerateStaticPages(distDir string, data *content.ContentData) {
	g.RenderPage(distDir, "index.html", "internal/templates/index.html", "", PageData{})
	g.RenderPage(distDir, "about.html", "internal/templates/about.html", "About", PageData{})
	g.RenderPage(distDir, "now.html", "internal/templates/now.html", "Now", PageData{})
	g.RenderPage(distDir, "404.html", "internal/templates/404.html", "404 - Not Found", PageData{})
	g.RenderPage(distDir, "tags.html", "internal/templates/tags.html", "Tags", PageData{
		Tags:      data.Tags,
		TagCounts: data.TagCounts,
	})
	g.RenderPage(distDir, "archive.html", "internal/templates/archive.html", "Archive", PageData{
		Archive:      data.PostsByYear,
		ArchiveYears: data.ArchiveYears,
	})
}

func (g *SiteGenerator) GenerateBlogPagination(distDir string, data *content.ContentData, pageSize int) {
	totalPages := (len(data.Posts) + pageSize - 1) / pageSize
	for i := 0; i < totalPages; i++ {
		startIdx := i * pageSize
		endIdx := startIdx + pageSize
		if endIdx > len(data.Posts) {
			endIdx = len(data.Posts)
		}
		pagePosts := data.Posts[startIdx:endIdx]
		pageNumber := i + 1

		if pageNumber == 1 {
			g.RenderPage(distDir, "blog.html", "internal/templates/blog.html", "Blog", PageData{
				Posts:       pagePosts,
				CurrentPage: pageNumber,
				TotalPages:  totalPages,
			})
		} else {
			pageDir := filepath.Join(distDir, "blog")
			g.RenderPage(pageDir, fmt.Sprintf("%d.html", pageNumber), "internal/templates/blog.html", fmt.Sprintf("Blog - Page %d", pageNumber), PageData{
				Posts:       pagePosts,
				CurrentPage: pageNumber,
				TotalPages:  totalPages,
				PathPrefix:  "../",
			})
		}
	}
}

func (g *SiteGenerator) GenerateTagPages(distDir string, data *content.ContentData) {
	tagsDistDir := filepath.Join(distDir, "tags")
	for tag, tagPosts := range data.PostsByTag {
		g.RenderPage(tagsDistDir, tag+".html", "internal/templates/blog.html", "#"+tag, PageData{
			Posts:      tagPosts,
			PathPrefix: "../",
		})
	}
}

func (g *SiteGenerator) GeneratePostPages(distDir string, data *content.ContentData) {
	blogDistDir := filepath.Join(distDir, "blog")
	for _, post := range data.Posts {
		p := post
		g.RenderPage(blogDistDir, post.Slug+".html", "internal/templates/post.html", post.Title, PageData{
			Post:       &p,
			PathPrefix: "../",
		})
	}
}

func (g *SiteGenerator) GenerateSearchIndex(distDir string, data *content.ContentData) {
	type SearchItem struct {
		Title       string   `json:"title"`
		Slug        string   `json:"slug"`
		Description string   `json:"description"`
		Date        string   `json:"date"`
		Tags        []string `json:"tags"`
	}

	var items []SearchItem
	for _, post := range data.Posts {
		items = append(items, SearchItem{
			Title:       post.Title,
			Slug:        post.Slug,
			Description: post.Description,
			Date:        post.Date.Format("January 02, 2006"),
			Tags:        post.Tags,
		})
	}

	jsonData, err := json.Marshal(items)
	if err != nil {
		log.Fatalf("Failed to marshal search index: %v", err)
	}

	if err := os.WriteFile(filepath.Join(distDir, "search-index.json"), jsonData, 0644); err != nil {
		log.Fatalf("Failed to write search-index.json: %v", err)
	}
}

func (g *SiteGenerator) GenerateRSS(distDir string, posts []content.Post) {
	f, err := os.Create(filepath.Join(distDir, "rss.xml"))
	if err != nil {
		log.Fatalf("Failed to create rss.xml: %v", err)
	}
	defer f.Close()

	fmt.Fprint(f, `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
<channel>
  <title>`+g.Config.Site.Title+`</title>
  <link>`+g.Config.Site.URL+`</link>
  <description>`+g.Config.Site.Description+`</description>
  <language>en-us</language>
`)

	for _, post := range posts {
		link := g.Config.Site.URL + "blog/" + post.Slug + ".html"
		fmt.Fprintf(f, `  <item>
    <title>%s</title>
    <link>%s</link>
    <description>%s</description>
    <pubDate>%s</pubDate>
    <guid>%s</guid>
  </item>
`, post.Title, link, post.Description, post.Date.Format(time.RFC1123), link)
	}

	fmt.Fprint(f, `</channel>
</rss>`)
}

func (g *SiteGenerator) GenerateSitemap(distDir string, posts []content.Post) {
	f, err := os.Create(filepath.Join(distDir, "sitemap.xml"))
	if err != nil {
		log.Fatalf("Failed to create sitemap.xml: %v", err)
	}
	defer f.Close()

	fmt.Fprint(f, `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
`)

	// Static Pages
	pages := []string{"", "about.html", "now.html", "blog.html", "tags.html", "archive.html"}
	for _, page := range pages {
		fmt.Fprintf(f, `  <url>
    <loc>%s%s</loc>
    <lastmod>%s</lastmod>
  </url>
`, g.Config.Site.URL, page, time.Now().Format("2006-01-02"))
	}

	// Blog Posts
	for _, post := range posts {
		fmt.Fprintf(f, `  <url>
    <loc>%sblog/%s.html</loc>
    <lastmod>%s</lastmod>
  </url>
`, g.Config.Site.URL, post.Slug, post.Date.Format("2006-01-02"))
	}

	fmt.Fprint(f, `</urlset>`)
}
