---
title: "Comparing Multistage and Single Container Builds"
description: "An empirical analysis of build times and image sizes comparing isolated multi-stage container compilation against single-stage builds for Next.js and serverless APIs."
date: 2026-08-25
tags: ["platform", "linux"]
draft: true
---

## Comparing Container Build Strategies

A full-stack monorepo container build pulling compiled workspaces and graphic dependencies bloated the execution image to over 5 GB. Modifying a single frontend file invalidated the build cache, forcing a complete two-minute reinstall and compilation loop. Restructuring the build process using isolated compile stages isolates build-time tooling and preserves caches.

---

## Evaluating the Monolithic Build Pattern

The first trial utilized a single-stage `Dockerfile.single` compiling all workspaces together in the final execution image. This monolithic build pattern was tested to measure the baseline speed and final file size. The build was executed by running `podman build -t cover-craft-single:latest -f Dockerfile.single .` from the root directory. The configuration resulted in several critical observations:

```text
+-------------------------------------------------------+
|                   Base Stage (Node)                   |
| - Installs build-essential, C++ libraries, devDeps    |
| - Compiles Shared Library, Next.js, and API           |
| - Retains all source files, caches, and tooling       |
+-------------------------------------------------------+
                           |
                           v
+-------------------------------------------------------+
|                  Final Image Output                   |
| - Runs PM2 process manager                            |
| - Size: 5.17 GB (including compiler utilities)        |
+-------------------------------------------------------+
```

- A single compilation run was fast for the initial setup, completing the full build in 1:52.79.
- The final image size was extremely large, registering a total footprint of 5.17 GB.
- Build-time compilers like `build-essential` and local `devDependencies` remained inside the final running layer.
- Modifying a single frontend file invalidated the entire cache, forcing a complete two-minute reinstall and recompilation.

---

## Optimizing via Isolated Build Stages

The second trial restructured the compilation using isolated builder stages inside the root `Dockerfile`. Compilation occurred in temporary builders, and only the compiled files were copied to the final `node:24-slim` base. The build was executed by running `podman build -t cover-craft-all-in-one:latest .` from the root directory. This optimization resolved the caching and storage issues through several improvements:

```text
+--------------------------------------------+
|             Base Stage (Node)              |
+--------------------------------------------+
                      |
                      v
+--------------------------------------------+
|            Shared-Builder Stage            |
| - Compiles @cover-craft/shared package     |
+--------------------------------------------+
                      |
        +-------------+-------------+
        |                           |
        v                           v
+---------------+           +---------------+
|   Frontend    |           |      API      |
|    Builder    |           |    Builder    |
| - Next.js     |           | - Canvas      |
+---------------+           +---------------+
        |                           |
        +-------------+-------------+
                      |
                      v
+--------------------------------------------+
|         Runner Stage (Debian Slim)         |
| - Copies only compiled execution assets    |
| - Excludes compiler tools and source code  |
| - Size: 3.26 GB (Optimized footprint)      |
+--------------------------------------------+
```

- The clean build time was slightly slower at 2:20.31 due to repeating dependencies.
- The final image size was reduced to 3.26 GB, saving 1.91 GB of space by excluding compiler utilities.
- Rebuild times on code changes dropped to under 10 seconds because unchanged builder stages were completely cached.

---

## Empirical Comparison Metrics

The local testing was executed on a Linux host using Podman. The base image `node:24-slim` was pre-pulled, but no previous build caches or final images existed prior to execution. The system monitored build durations, processor utilization, and storage outputs to produce the comparison table below:

```text
+------------------------+-----------------+-----------------+--------------+------------+
| Metric                 | Single-Stage    | Multi-Stage     | Difference   | % Change   |
+------------------------+-----------------+-----------------+--------------+------------+
| Final Image Size       | 5.17 GB         | 3.26 GB         | -1.91 GB     | -36.9%     |
| Network Packets (Est)  | 3,446,667       | 2,173,333       | -1,273,334   |            |
| Build Time             | 1:52.79         | 2:20.31         | +27.52s      | +24.4%     |
| Host CPU Usage (Avg)   | 150%            | 166%            | +16%         | +10.7%     |
+------------------------+-----------------+-----------------+--------------+------------+
```

The network packet estimation assumes a standard Maximum Transmission Unit size of 1,500 bytes. This hardware limit determines the maximum data payload carried by a single Ethernet frame. The total packet count illustrates the physical network load required to move each container image.

---

## Technical Insights and Key Takeaways

Analyzing these results highlights the trade-offs of different containerization strategies. Each configuration affects build caching and image size differently. The core takeaways from this evaluation include:

- Single-stage builds are simpler to write but produce bloated containers containing unnecessary build tools.
- Multi-stage builds increase the initial compilation time slightly but yield major storage savings on the local host.
- Cache isolation in multi-stage builds prevents minor code edits from triggering a full rebuild of unrelated workspaces.
- Isolating compilation from execution is required to optimize local container storage and host resource usage.

---

## Conclusion

I started this experiment to build upon what I learned from my previous School Management Flask API project. In that small python experiment, the build speed differences between single-stage and multi-stage builds were minimal. The full-stack `next.js` and `nodejs` monorepo highlighted a clear image size difference between the two configurations. Saving 1.91 GB avoids transmitting approximately 1.3 million network packets during image distribution.

Running this updated experiment clarified how combining a React frontend with a backend service drastically increases the image footprint depending on the imported packages. Compiling all workspaces in a single stage is convenient but results in a bloated final footprint. The pivot to multi-stage builds resolved this bloat by copying only the compiled assets needed for execution without sacrificing build speed.
