---
title: "[Engineering Log] HTML Minification Network Savings"
description: "Analyzing the impact of compile-time HTML minification as a static blog grows to over 240 posts, measuring build speeds, raw byte savings, and network packet reductions."
date: 2026-07-21
tags: ["go", "automation", "retrospective"]
---

## Context

The static site generator originally prioritized layout formatting and readability over build output sizes. There was no initial concern regarding cache eviction rates or network packet limits. As a result, pages were written directly to disk as unminified `html` documents.

A recent review of the repository showed that the `blog/` folder has grown to 246 posts. This growth made me curious about the total byte size and the future outlook if the catalog reaches 500 posts. The accumulated baseline output size had already reached multiple megabytes. Evaluating the build size and compile times became necessary to sustain long-term project portability.

---

## Challenge

The initial assumption was that larger files simply translate to more data passing through the network. While this byte increase was not an immediate bottleneck, it motivated a curiosity-driven deep dive. The core interest was to understand the mechanical connection between raw file size and network packet transmission. This search aimed to connect system build size directly to low-level transport metrics.

---

## Investigation

The investigation focused on measuring the baseline output and analyzing standard transport protocols. Calculations mapped the raw file size directly to network packet transmission limits. These findings showed the operational impact of unminified static pages:

- Baseline `html` size across the site measured 6,001,624 bytes.
- The unminified build required transmitting approximately 4,110 `tcp` packets to clients.

A standard Ethernet interface enforces a strict maximum transmission unit (`mtu`) constraint. The table below details the packet payload allocation and headers overhead under standard network configurations. This structure determines the bytes available for actual data transmission:

| Packet Component | `ipv4` + `tcp` (Bytes) | `ipv6` + `tcp` (Bytes) | Description / Role |
| :--- | :--- | :--- | :--- |
| `ip` header | 20 | 40 | Layer 3 routing and packet delivery information |
| `tcp` header | 20 | 20 | Layer 4 connection control and sequencing |
| Data payload (`mss`) | 1,460 | 1,440 | Actual content segment size transmitted per packet |
| Total `mtu` | 1,500 | 1,500 | Maximum transmission unit capacity |

---

## Solution

The package `github.com/tdewolff/minify/v2` was integrated into the static page generator. The generator initializes the minifier and configures it to handle the `html` media type. During execution, the page rendering function outputs templates to a buffer and minifies the bytes. The minified string is then written to the final output file.

This integration successfully reduced the total `html` output size while maintaining all formatting structures. The table below outlines the comparative metrics before and after incorporating compile-time minification. These statistics illustrate the performance shift across files and compilation times.

| Metric / Asset | Pre-Minification | Post-Minification | Difference | Savings |
| :--- | :--- | :--- | :--- | :--- |
| `html` documents (133 files) | 6,001,624 bytes | 4,890,517 bytes | -1,111,107 bytes | 18.5% |
| `html` documents (Size in MB) | 6.00 MB | 4.89 MB | -1.11 MB | - |
| Network transport (`tcp` packets) | ~4,111 packets | ~3,350 packets | -761 packets | - |
| SSG Compilation Time | ~660 ms | ~614 ms | -46 ms | 7.0% |

The static site generator only executes minification on the generated `html` documents. The stylesheet size is excluded from this optimization as it is already minified by the `Tailwind CSS` CLI. The network transport calculations use `ipv4` packet limits as the baseline. If a client connects using `ipv6`, the savings increase to 772 packets due to the smaller 1,440-byte payload limit.

---

## Evolution

I integrated `html` minification as a proactive experiment after observing the blog grow to 246 posts at the time of writing this post. The initial motivation was mostly curiosity to see how much difference minification would make before and after the change. While a 1.1 MB reduction does not affect edge network performance or cdn distribution on `Cloudflare Pages`, the exercise shows how optimization scales over time. As the post catalog grows, these cumulative transfer reductions and `tcp` packet savings will prevent future overhead and keep local compilation fast.
