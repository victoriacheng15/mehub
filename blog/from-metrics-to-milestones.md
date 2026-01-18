---
title: "From Metrics to Milestones"
description: "Visualizing a project’s journey: a product-first landing page and a timeline generated from a structured YAML history of milestones."
date: 2026-01-20
tags: ["growth", "go"]
draft: true
---

*This is Part 3 of my journey building a Personal Reading Analytics Dashboard. Catch up on [Part 1 - From PIi to Cloud Automation](https://victoriacheng15.vercel.app/blog/from-pi-to-cloud-automation) and [Part 2 - From Links to Reading Insights](link-to-part-2).*

---

## Visualizing the Journey

In software engineering, we often treat documentation as an afterthought. We build the feature, ship the code, and move on.

However, as this project grew, I realized that the code only told half the story. The real value was in the **evolution**—the series of architectural pivots, refactors, and lessons learned along the way. When sharing a project, it's easy to focus only on the end result, but for this dashboard, I felt that the most interesting part was the **process**.

I wanted to visualize this journey from a local Python script to a serverless Go platform by making these internal milestones visible:

- **Iterative Growth:** How the system transitioned from local scripts to a fully automated cloud pipeline.
- **Product Context:** The "Origin Story" and the design principles that guided the development.
- **Technical Wins:** Specific engineering milestones, such as adopting `asyncio` for faster extraction or building a custom Go metrics engine.

---

## The Landing Page: The "Why" Before the "What"

For the `index.html` (Landing Page), I focused on the "Product Pitch."

Engineers often jump straight to "Tech Stack" (I used Go! I used Docker!). But a Staff Engineer asks, "Why does this exist?"

I structured the landing page to answer that first:

- **Origin Story:** Why I needed this tool (too many tabs, scattered reading lists).
- **Design Principles:** Zero infrastructure, cost-effectiveness, and automation.
- **Call to Action:** Direct links to the live analytics and the code.

By separating the "Pitch" (Landing) from the "Data" (Analytics), I created a better user experience for two different audiences:

- **Recruiters** get the high-level summary and impact.
- **Engineers** get the deep-dive metrics and architecture.

---

## Visualizing Growth: The Evolution Page

The **Evolution Page** (`evolution.html`) renders the project's history into a vertical timeline.

This visualization is powerful because it proves **consistency**. It shows that this wasn't a weekend hackathon project; it’s a maintained system that has survived multiple refactors over months. It highlights the transition from "Junior" tasks (writing a script) to "Senior" tasks (architecting observability and testing frameworks).

But I didn't want to hardcode this HTML manually. I wanted the history to be treated as code.

---

## The Solution: Documentation as Data

Instead of editing HTML files for every update, I decided to treat the project's history as structured data.

I created a schema file, `evolution.yml`, that acts as the single source of truth for the project's lifecycle.

```yaml
events:
  - date: "2024-02"
    title: "Article Collection Begins"
    description: |
      - "Automated article collection from technical blogs to Google Sheets using a local Python script."
      - "Utilized Python generators to efficiently stream extracted data."

  - date: "2025-01"
    title: "Fully Automated Daily Collection"
    description: |
      - "Migrated from manual runs to a scheduled, automated workflow using GitHub Actions."
      - "Added structured logging and secure credential handling."
```

This approach allows me to "update" the portfolio by simply appending a new entry to the YAML file. The frontend automatically rebuilds to reflect the new state.

---

## Engineering the Generator

I extended my existing Go-based **Dashboard Generator** to handle this new data source.

The pipeline now looks like this:

1. **Ingest:** Read `metrics.json` (Quantitative Data) AND `evolution.yml` (Qualitative Data).
2. **Process:** Parse the YAML into a strict Go Struct (`EvolutionData`).
3. **Render:** Hydrate distinct HTML templates (`analytics.html`, `evolution.html`, `index.html`).

```go
type Milestone struct {
    Date             string   `yaml:"date"`
    Title            string   `yaml:"title"`
    Description      string   `yaml:"description"` 
    DescriptionLines []string `yaml:"-"` // Processed for template rendering
}
```

---

## Conclusion: The Repository as a Product

We often think of "Product" as the thing the user touches. But for a software engineer, the **Repository** is also a product. The users are your teammates, your future self, and potential employers.

By adding a polished Landing Page and an Evolution Timeline, I transformed this repo from a "code bucket" into a **technical narrative**.

It’s no longer just about *what* I built. It’s about *how* I grew along with it.

---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
