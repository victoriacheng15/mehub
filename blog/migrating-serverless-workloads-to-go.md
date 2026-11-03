---
title: "Migrating Serverless Workloads to Go"
description: "Rewriting the Cover Craft serverless API in Go eliminated node_modules packaging bloat, cut deployment sizes by 73 percent, and reduced CI build durations by two minutes."
date: 2026-11-03
tags: ["go", "platform", "cloud"]
---

## The Serverless Node.js Packaging Tax

Building the initial API in TypeScript helped prototype features quickly, but deploying Node.js to serverless environments introduced a heavy packaging penalty: bundling thousands of loose files inside `node_modules` into a zip archive on every release.

The serverless deployment step alone took three minutes just to upload and extract dependency trees. With total CI/CD runs dragging past five minutes, even small bug fixes resulted in sluggish feedback loops:

| Pipeline Component | Baseline Run Time |
| :--- | :--- |
| Serverless API Deployment | ~3m (± 15-30s tolerance) |
| Overall CI/CD Workflow | ~5m 32s |

---

## Compiled Custom Handlers and Contract Isolation

Deploying loose dependency trees to serverless consumption tiers introduces a physical I/O bottleneck: providers must unpack thousands of small files onto ephemeral instances before execution begins. Replacing the Node.js backend with an Azure Functions custom handler eliminated that extraction tax. Compiling with `CGO_ENABLED=0` produces a single standalone Linux binary, stripping out runtime dependency trees, package managers, and dynamic linkers.

The rewrite also separated frontend and backend boundaries without sacrificing type safety:

- The Go binary handles 2D drawing with native libraries, avoiding Cgo compilation overhead.
- `openapi.yaml` serves as the contract to generate TypeScript types and Go structs.
- Cold starts remain bounded below 5ms without paying for provisioned concurrency.
- A centralized `Makefile` orchestrates builds, test suites, and deployment packaging.

```text
  [ Next.js Frontend ] ---> [ OpenAPI Contract ] <--- [ Go API Handler ]
                                                              |
                                                              v
                                                   [ Linux Static Binary ]
                                                              |
                                                   (Sub-5ms Cold Start)
```

---

## Benchmark Results and Pipeline Acceleration

Rather than deprecating Node.js immediately, both backends coexisted in the codebase and ran concurrently through CI/CD across several workflows. Past workflow logs provided a reliable Node.js baseline, but running both in parallel allowed gathering hard comparative data on artifact sizes, file counts, and deployment duration before committing to a permanent cutover:

| Metric | Node.js Backend (Legacy) | Go Backend (Compiled) | Reduction |
| :--- | :--- | :--- | :--- |
| Source Language | TypeScript | Go (Custom Handler) | - |
| Source Code Volume | ~1,100 lines | ~850 lines | ~23% |
| Build Size | ~60MB | 16MB | 73% |
| File Count | ~1,200+ files | ~12 files | 99% |

Shipping a 16MB binary made an immediate impact on pipeline duration. Comparing consecutive workflow runs confirmed that deployment time dropped from three minutes down to approximately one minute:

| Pipeline Component | Baseline (Node.js) | Post-Migration (Go) | Improvement |
| :--- | :--- | :--- | :--- |
| Serverless API Deployment | ~3m | ~1m | ~67% |
| Overall CI/CD Workflow | ~5m 32s | ~3m 30s | ~36% |

---

## Conclusion

Node.js dependency bloat is manageable in long-running containers, but serverless deployment boundaries heavily penalize loose file trees. Running both runtimes concurrently through CI/CD provided the comparative data required to justify the cutover before deprecating the legacy codebase.

The metrics justified the migration: cutting deployment duration from three minutes to one minute eliminated a persistent CI bottleneck. Trading interpreted runtime convenience for a statically compiled binary established tighter contract boundaries, bounded cold starts below 5ms, and permanently reduced iteration cycle time.

Find the source code in the [cover-craft](https://github.com/victoriacheng15/cover-craft) repository.
