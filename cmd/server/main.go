package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"mehub/internal"
)

const port = ":8080"

// liveReloadSnippet is injected before </body> in every HTML response.
// It opens an SSE connection to /dev-reload. When the connection drops
// (because the server is restarting), it polls until the new server is up,
// then reloads the page.
const liveReloadSnippet = `<script>
(() => {
	function connect() {
		const es = new EventSource('/dev-reload');
		es.onerror = () => {
			es.close();
			const interval = setInterval(async () => {
				try {
					await fetch('/dev-reload');
					clearInterval(interval);
					location.reload();
				} catch (e) {}
			}, 200);
		};
	}
	connect();
})();
</script>`

func main() {
	if err := build(); err != nil {
		log.Fatalf("initial build failed: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/dev-reload", sseHandler)
	mux.HandleFunc("/", serveHandler)

	log.Printf("dev server → http://localhost%s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

// sseHandler holds an SSE stream open. The browser reload is triggered
// when the connection drops during a server restart.
func sseHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	fmt.Fprintf(w, "data: connected\n\n")
	flusher.Flush()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			fmt.Fprintf(w, "event: heartbeat\ndata: \n\n")
			flusher.Flush()
		}
	}
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

	body := string(content) + liveReloadSnippet
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
