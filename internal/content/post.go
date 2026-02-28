package content

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

type Frontmatter struct {
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	Date        time.Time `yaml:"date"`
	Tags        []string  `yaml:"tags"`
	Draft       bool      `yaml:"draft"`
}

type RelatedPost struct {
	Title string
	Slug  string
}

type Post struct {
	Frontmatter
	Slug         string
	Content      string
	RelatedPosts []RelatedPost
}

type ContentData struct {
	Posts        []Post
	PostsByTag   map[string][]Post
	PostsByYear  map[int][]Post
	Tags         []string
	TagCounts    map[string]int
	ArchiveYears []int
}

var mermaidRegex = regexp.MustCompile(`(?s)<pre><code class="language-mermaid">.*?</code></pre>`)

func ParsePost(path string) (*Post, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parts := strings.SplitN(string(data), "---", 3)
	if len(parts) < 3 {
		return nil, nil
	}

	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(parts[1]), &fm); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
			),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	if err := md.Convert([]byte(parts[2]), &buf); err != nil {
		return nil, err
	}

	content := buf.String()
	content = mermaidRegex.ReplaceAllStringFunc(content, func(m string) string {
		inner := strings.TrimPrefix(m, `<pre><code class="language-mermaid">`)
		inner = strings.TrimSuffix(inner, `</code></pre>`)
		return `<div class="mermaid">` + inner + `</div>`
	})

	slug := strings.TrimSuffix(filepath.Base(path), ".md")

	return &Post{
		Frontmatter: fm,
		Slug:        slug,
		Content:     content,
	}, nil
}

func GetPosts(contentDir string) ([]Post, error) {
	var posts []Post

	files, err := os.ReadDir(contentDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			post, err := ParsePost(filepath.Join(contentDir, file.Name()))
			if err != nil {
				return nil, err
			}
			if post != nil && !post.Draft {
				posts = append(posts, *post)
			}
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts, nil
}

func ProcessPosts(posts []Post) *ContentData {
	data := &ContentData{
		Posts:       posts,
		PostsByTag:  make(map[string][]Post),
		PostsByYear: make(map[int][]Post),
		TagCounts:   make(map[string]int),
	}

	tagMap := make(map[string]bool)
	var archiveYears []int

	for i := range posts {
		// Tags grouping
		for _, tag := range posts[i].Tags {
			tagMap[tag] = true
			data.PostsByTag[tag] = append(data.PostsByTag[tag], posts[i])
			data.TagCounts[tag]++
		}

		// Year grouping
		year := posts[i].Date.Year()
		if len(data.PostsByYear[year]) == 0 {
			archiveYears = append(archiveYears, year)
		}
		data.PostsByYear[year] = append(data.PostsByYear[year], posts[i])
	}

	// Calculate related posts via tag intersection
	for i := range data.Posts {
		tagSet := make(map[string]bool)
		for _, tag := range data.Posts[i].Tags {
			tagSet[tag] = true
		}

		type scoring struct {
			idx   int
			score int
		}
		var scores []scoring

		for j := range data.Posts {
			if i == j {
				continue // Skip self
			}
			score := 0
			for _, tag := range data.Posts[j].Tags {
				if tagSet[tag] {
					score++
				}
			}
			if score > 0 {
				scores = append(scores, scoring{idx: j, score: score})
			}
		}

		// Sort related matches descending by shared tags. Tie-breaker is date (which inherently sorts by array index since GetPosts sorts date desc)
		sort.SliceStable(scores, func(a, b int) bool {
			return scores[a].score > scores[b].score
		})

		// Take top 3 max
		maxRelated := 3
		if len(scores) < maxRelated {
			maxRelated = len(scores)
		}

		for _, s := range scores[:maxRelated] {
			data.Posts[i].RelatedPosts = append(data.Posts[i].RelatedPosts, RelatedPost{
				Title: data.Posts[s.idx].Title,
				Slug:  data.Posts[s.idx].Slug,
			})
		}
	}

	for tag := range tagMap {
		data.Tags = append(data.Tags, tag)
	}
	sort.Slice(data.Tags, func(i, j int) bool {
		if data.TagCounts[data.Tags[i]] != data.TagCounts[data.Tags[j]] {
			return data.TagCounts[data.Tags[i]] > data.TagCounts[data.Tags[j]]
		}
		return data.Tags[i] < data.Tags[j]
	})
	sort.Sort(sort.Reverse(sort.IntSlice(archiveYears)))
	data.ArchiveYears = archiveYears

	return data
}
