---
title: "Automating My Local Dev Env"
description: "Containerizing a static site workflow with Podman, Go, and live reloading automates manual rebuild loops and eliminates repetitive local development environment friction."
date: 2026-10-20
tags: ["platform", "go", "automation"]
draft: true
---

## Containerizing the Build Toolchain

Running my static site generator locally worked well at first, but the manual feedback loop quickly became tedious. Every time I changed a Tailwind CSS class, I had to stop and rerun the build script to preview changes. Packaging the entire toolchain into a container eliminated this manual friction and gave me a consistent local development pipeline.

I chose a multi-stage `Dockerfile.dev` built on Alpine Linux to keep the image lightweight while encapsulating the toolchain:

- `Go 1.26` compiles and runs the static site generator.
- `Tailwind CLI` compiles atomic utility classes on file changes.
- `Air` watches source files and handles live process recompilation.
- `goimports` formats Go source files during the build cycle.

Because the standalone Tailwind CLI binary targets glibc environments, running it on Alpine's musl libc caused dynamic linking failures. I added `gcompat`, `libgcc`, and `libstdc++` to resolve the missing symbols, allowing the container to compile templates and styles automatically without host dependencies.

---

## Running a Custom Dev Server

Because my static site generator compiles pages and exits immediately, keeping the container alive required a persistent server process. I wrote a small development server in `cmd/server/main.go` to handle this: it executes the generation pipeline on startup, triggers Tailwind compilation, and hosts the generated output directory.

To trigger automatic browser refreshes, the dev server streams Server-Sent Events (SSE). The browser maintains a persistent `EventSource` connection to the `/dev-reload` endpoint. My initial search-and-replace script failed to inject the reload snippet because the HTML minifier stripped closing `</body>` tags. I bypassed this cleanly by appending the script directly to the end of the HTML payload, ensuring it runs regardless of minification.

Writing a custom HTTP handler also gave me full control over how files are read from the `dist/` directory. Reading files directly from disk on every request prevents file-locking collisions with the `air` rebuild process and guarantees the browser always receives the freshest build output:

```text
  [ Source File Change ]
            |
            v
  [ Air Restarts Server ]
       |             |
       v             v
  (Re-run SSG)   (Recompile Tailwind)
       |
       v
  [ Server Drops SSE Connection ]
       |
       v
  [ Browser Reconnects & Reloads Instantly ]
```

---

## Managing the Container Lifecycle

Typing long Podman command flags to mount volumes and map ports creates unnecessary friction during daily development. I encapsulated the container lifecycle inside a `Makefile` to establish a simple, repeatable interface.

Starting my local environment now only requires running `make dev-run` on the host machine. This command launches the Podman container, mounts the workspace with SELinux `:Z` volume flags, and starts the file watcher with live log streaming. When I stop the container, runtime state cleans up automatically without leaving orphaned processes.

Hiding complex container flags behind make targets ensures that security policies and mount options remain version-controlled. It makes the workspace portable and easy to run without having to look up command arguments every morning.

---

## Conclusion

I put off automating this for weeks because typing `make build` felt like a minor inconvenience. Taking the time to orchestrate Podman, Air, and live reloading eliminated that friction completely and made working on this site fun again.

Find the source code in the [mehub](https://github.com/victoriacheng15/mehub) repository.
