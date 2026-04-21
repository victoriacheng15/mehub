---
title: "[Engineering Log] Swapping Vue for Go"
description: "Rebuilding a Linktree-style Vue project as a Go static site generator to reduce npm maintenance, simplify builds, and learn Go through a small real-world publishing tool."
date: 2026-04-21
tags: ["retrospective", "go"]
---

## Context

BioHub started as a Linktree-style page built with Vue.js. At the time, Vue was a useful way to learn component-driven frontend development through a small real project. The page stayed simple, but the update loop around npm dependencies kept coming back.

The friction showed up in the maintenance loop. A static one-page site still required npm updates, framework packages, and build tooling just to publish a list of links. Vue was not the wrong choice for learning, but it became heavier than the project needed.

That mismatch became harder to ignore over time. Most changes were content edits, not interface changes. The maintenance cost no longer matched the shape of the work.

---

## Challenge

The main problem was not the UI. The page already had a clear visual structure and a small content model. The problem was keeping a JavaScript toolchain alive for a page that did not need runtime reactivity.

- **Objective:** Keep the polished BioHub layout while reducing build and package maintenance.
- **Constraint:** Remove the need for npm dependency updates, frontend build tooling, and client-side JavaScript.
- **Learning Goal:** Use the rewrite as a practical Go exercise instead of a synthetic tutorial.

---

## Investigation

The comparison focused on fit rather than language preference. Vue still offered a strong component model, but BioHub did not need client-side state, routing, or reactive updates. A Go static site generator offered a smaller surface area while still leaving room to practice templates, file generation, and build automation.

| Option | Gain | Cost | Verdict |
| :--- | :--- | :--- | :--- |
| **Vue.js SPA** | Familiar frontend model and reusable components. | npm updates, build tooling, and runtime JavaScript for one static page. | Too much surface area |
| **Go SSG** | Single binary, simpler package management, and static output. | Less frontend interactivity and custom template logic. | Better fit |

### The Humble Pivot: Learning vs. Fit

The useful lesson was that a good learning stack can become the wrong operating stack. Vue helped make frontend patterns more concrete, but the project no longer needed the framework after the learning goal was met. Rebuilding BioHub in Go turned the maintenance problem into a focused systems exercise.

---

## Solution

The implementation moved BioHub from a Vue app to a Go-based static site generator. Structured data became the source of truth, and `html/template` handled the final page rendering at build time. The output became static HTML and CSS instead of a JavaScript application.

- **Before:** Vue components, npm packages, and a frontend build pipeline for a static page.
- **After:** A Go generator, template files, and static assets deployed as plain files.

The build path became easier to reason about. A single Go command produces the generated site without restoring a large dependency tree on every CI run. That change reduced the maintenance surface while keeping the project useful as a Go learning tool.

The rewrite also made the publishing boundary more visible. Input data, templates, and generated output now sit in separate roles. That separation made each future change easier to inspect before deployment.

---

## Evolution

- **Package Management:** The project moved from recurring npm updates to a small Go module surface.
- **Build Simplicity:** The deployment artifact became static HTML and CSS, with no client-side JavaScript required for the page.
- **Learning Outcome:** The rewrite made Go templates, file generation, and static publishing easier to understand through a real project.
- **Lesson Learned:** Tooling fit changes over time, especially when a project moves from learning experiment to maintained utility.
