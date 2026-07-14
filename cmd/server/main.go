package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"mehub/internal"
)

const port = ":8080"

// buildID is set once at startup. Because air restarts the process on every
// rebuild, each new server has a different ID. The browser detects the change
// and reloads.
var buildID = strconv.FormatInt(time.Now().UnixMilli(), 10)

// liveReloadSnippet is injected before </body> in every HTML response.
// It polls /dev-reload every second and reloads when the build ID changes.
const liveReloadSnippet = `<script>
(function() {
	var id = null;
	setInterval(function() {
		fetch('/dev-reload')
			.then(function(r) { return r.text(); })
			.then(function(v) {
				if (id === null) { id = v; return; }
				if (v !== id) { location.reload(); }
			})
			.catch(function() {});
	}, 1000);
})();
</script>`

func main() {
	if err := build(); err != nil {
		log.Fatalf("initial build failed: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/dev-reload", devReloadHandler)
	mux.HandleFunc("/", serveHandler)

	log.Printf("dev server → http://localhost%s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

// devReloadHandler returns the current build ID as plain text.
// The browser compares successive responses and reloads on change.
func devReloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Cache-Control", "no-cache")
	fmt.Fprint(w, buildID)
}

// serveHandler serves files from dist/. HTML files have the live-reload script
// injected before </body>. All other files are served directly.
func serveHandler(w http.ResponseWriter, r *http.Request) {
	p := filepath.Join("dist", filepath.Clean(r.URL.Path))

	info, err := os.Stat(p)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if info.IsDir() {
		p = filepath.Join(p, "index.html")
	}

	if strings.HasSuffix(p, ".html") {
		serveHTML(w, r, p)
		return
	}

	http.ServeFile(w, r, p)
}

// serveHTML reads an HTML file, injects the live-reload script, and writes it.
func serveHTML(w http.ResponseWriter, r *http.Request, path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	body := strings.Replace(string(content), "</body>", liveReloadSnippet+"</body>", 1)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, body)
}

// build runs the SSG pipeline then compiles Tailwind CSS.
func build() error {
	start := time.Now()

	count, err := internal.RunPipeline(
		"dist",
		"internal/templates/contents",
		"internal/templates",
		"blog",
		"internal/templates/static",
	)
	if err != nil {
		return fmt.Errorf("ssg pipeline: %w", err)
	}

	if err := runTailwind(); err != nil {
		return fmt.Errorf("tailwind: %w", err)
	}

	log.Printf("✅ built %d posts in %v", count, time.Since(start))
	return nil
}

// runTailwind compiles input.css into dist/styles.css.
func runTailwind() error {
	cmd := exec.Command(
		"tailwindcss",
		"-i", "internal/templates/input.css",
		"-o", "dist/styles.css",
		"--minify",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w\n%s", err, out)
	}
	return nil
}
