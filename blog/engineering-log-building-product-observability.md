---
title: "[Engineering Log] Building Product Observability"
description: "Building Cover Craft analytics into a lightweight observability loop that tracks image generation, latency, accessibility, UI usage, and structured backend behavior."
date: 2026-09-15
tags: ["retrospective", "platform", "system-design"]
draft: true
---

## Context

Building a custom observability system for Cover Craft started as a curiosity-driven project. The goal was to understand system behavior and track how users interacted with the image generator. Creating a local feedback loop below the application boundary made these performance and usage patterns visible.

The project needed a way to monitor system health and track users' favorite features. A large tracking platform was too heavy for this stage of development. The system needed a clean way to track key signals directly:

- How often users generated images.
- How long generation took.
- Whether contrast outcomes stayed accessible.
- Which UI paths were actually used.
- Where backend failures appeared.

### Past Logs

- [Accessibility as System Constraint](https://victoriacheng15.dev/blog/engineering-log-accessibility-as-system-constraint)

---

## Challenge

The main issue was payload structure. Metrics are only useful when the browser and the API agree on the event types. If these definitions drift, the dashboard renders charts with inaccurate data.

Basic averages also hid performance problems. A low average latency can conceal slow generation times for larger images. The metrics needed percentiles to expose actual user friction.

Observability tools can quickly become too complex to build and maintain. A small product does not need complex tracing servers or heavy external databases. The design had to provide visibility using minimal resources.

---

## Investigation

The investigation focused on deciding which metrics would best track system health and user preferences. Identifying the most popular features required tracking specific user interactions without logging unnecessary noise. This analysis helped define strict contracts for every tracked event.

The system initially tracked every single user click and API request. That was a wrong turn because it generated cluttered charts and filled the database with meaningless logs. Consolidating the model around a few key actions kept the dashboard clean and informative.

The event model was simplified to focus only on key actions. Reducing the number of event types made database aggregation much cleaner. The dashboard tracked a small set of metrics:

- Generation counts.
- Latency trends.
- Contrast ratios.
- `WCAG` distribution.
- UI usage.

Visualizing performance required more than basic average metrics. Measuring daily percentiles offered a clearer view of slow response times. Tracking these extreme values helped pinpoint actual bottlenecks during generation.

Empty states were also treated as a first-class concern. A chart that omits empty categories looks clean but hides missing data. Standardizing these gaps made the dashboard more predictable.

---

## Solution

The implementation centralized event types and moved logging to the API. Database queries run aggregations to produce the dashboard views. Structured JSON logs provide a search trail without adding heavy infrastructure.

The ingestion architecture processes events through a local pipeline. The diagram shows how analytics flow from the client to the database. This setup keeps ingestion lightweight and predictable:

```text
  [ React Frontend ] ---> [ API Ingestion Endpoint ] ---> [ MongoDB / Aggregation ]
                                 |
                                 v
                        [ JSON Log Files ]
```

The solution focused on building a dashboard to display only the metrics the developer needs. Database queries run `MongoDB` aggregations to fetch generation counts, latency numbers, font selections, and accessibility levels. This dashboard provides a clean interface to inspect system behavior at a glance.

Backend logs use structured `JSON` formatting to simplify local searches. Consistent fields make slow routes and database errors easier to locate. This format ensures logs remain queryable even without a dedicated logging server.

---

## Evolution

The main takeaway was learning how to design a dashboard to be as informative as possible. Showing an overview of user choices makes it easy to understand user behavior at a glance. Putting thought into this presentation helps identify what the product actually needs next.

Now, this custom observability serves as the foundation for the product roadmap. Analyzing these usage patterns directly guides future design decisions and feature improvements. Learning how to translate user data into next steps is the core benefit of this setup. Find the source code in the [cover-craft](https://github.com/victoriacheng15/cover-craft) repository.
