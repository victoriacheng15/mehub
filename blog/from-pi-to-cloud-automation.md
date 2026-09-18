---
title: "From Pi to Cloud Automation"
description: "Automate content gathering with a Python script and GitHub Actions. No servers, no Pi—just daily updates straight to your Google Sheet. Explore this comprehensive guide."
date: 2025-08-05
tags: ["cloud", "python"]
---

## From Pi to Cloud Automation?

Like many developers, some of my favorite learning projects come from scratching a personal itch. In this case, I wanted a simple system to collect articles I was reading — from sources like freeCodeCamp, Substack, GitHub, and Shopify — so I could see all the links in one place without having to visit each site individually.

That personal itch turned into `articles-extractor`, a Python project that started small and ran locally on my machine. Later, I moved it to a Raspberry Pi to automate scheduled runs. Eventually, I deployed it using GitHub Actions to run reliably in the cloud — enabling the workflow to run anywhere, anytime without depending on my own hardware.

At its core, `articles-extractor` is a Python script that helps **centralize links** from multiple platforms into one place.

The script works by:

- Reads a list of content sources from a **Google Sheet**, which contains two worksheets: one for the source list and another for storing extracted article links
- Fetches web pages from those sources (e.g., blog feeds or update pages)
- Uses BeautifulSoup with source-specific extractor functions to parse the page
- Extracts the essentials: `title`, `link`, and `date`
- Appending that info to the other worksheet in the same Google Sheet, which serves as a clickable reading list

The goal isn’t full content scraping — just a convenient way to view and access fresh articles from all your favorite sources in one place.

It uses:

- `requests` for HTTP calls  
- `BeautifulSoup` for HTML parsing  
- `gspread` for Google Sheets API  
- `python-dotenv` for managing credentials and API keys

---

## 🧪 Version 1: Running Locally

The first version ran on my laptop. Initially, I used it to test and validate how the script extracted article data from different sources. It helped me confirm that the extractor functions were working as expected.

But after a while, the process became **redundant and tedious** — opening a terminal and running the same script manually every day just to get updated links.

So I started thinking:

> How do I automate this locally?

I managed to get it running as a scheduled task on my desktop using a cron-like tool. It worked — but only when my desktop was powered on at the right time.

That limitation made it clear that local automation wasn’t reliable enough for this job. I needed something that didn’t depend on my computer being awake or online.

---

## 🍓 Version 2: Raspberry Pi Automation

To improve automation, I moved the project to a **Raspberry Pi 3**. I used `cron` to schedule the script to run daily, allowing it to execute without manual input.

At first, it felt like a small win: a headless Linux box doing the work for me!

But over time, I ran into multiple problems:

- The Pi would **overheat** because I didn’t have a proper cooling system
- The script would sometimes **fail silently**, and since the logs were stored locally on the Pi, I had to SSH into the device to check what went wrong — which added friction every time something broke

These issues added up, and I started to ask myself:

> How can I run this job reliably without having to worry about it?

I just wanted a solution that was dependable, accessible from anywhere, and easy to monitor. That’s what led me to explore GitHub Actions.

---

## ☁️ Version 3: GitHub Actions

I already knew GitHub Actions was mostly used for CI/CD — like running tests or deployments when code changes.

But then I thought, why not use it for something else? Since it supports cron scheduling, maybe I could have it run my script every day on a schedule. That way, I wouldn’t have to keep my Pi or desktop on all the time.

That simple idea led me to set up a cloud-based workflow where GitHub runs the job daily, so all my article links get updated and ready whenever I want to check them.

### 🛠️ What I Learned

- **Secrets Management:** I securely stored environment variables and API keys using **GitHub Secrets** (`Settings → Secrets → Actions`). This approach keeps sensitive credentials out of the codebase while allowing the script to access them securely at runtime.  
- **`workflow_dispatch`:** I added a `workflow_dispatch` trigger so I could manually run the script from the GitHub UI whenever needed — which is especially helpful for debugging or testing changes.  
- **Logging with Artifacts:** Treating this as a learning opportunity, I explored how **`actions/upload-artifact`** works in GitHub Actions. Rather than SSH’ing into a server to check logs (something I could do but found cumbersome), I modified the script to write logs to a `.txt` file and configured the workflow to upload this log as an artifact. This makes it easy to download and inspect logs directly from the GitHub UI after each run.

Now the whole system runs automatically in the cloud, fully self-contained within the GitHub repo — no local devices, no Pi, and no extra infrastructure to maintain.

It was a simple shift in mindset — from "CI/CD" to "scheduled automation" — but it opened up a whole new use case for GitHub Actions that fits solo projects perfectly.

---

## 🔄 Infrastructure Evolution: A Timeline

| Phase            | What It Looked Like                            | Pain Points Solved                        |
|------------------|--------------------------------------------------|--------------------------------------------|
| Local laptop     | Manual runs via terminal                        | No automation                              |
| Raspberry Pi     | Daily cron job + script                         | Some automation, but unstable hardware     |
| GitHub Actions   | Cloud scheduler + environment-managed secrets   | Fully automated, maintainable from anywhere|

---

## 💡 Reflections

I started this as a fun project — just a simple way to collect articles and save links in one place. I didn’t have any grand plans or long-term goals. But as I used it more, I ran into small annoyances and limitations. Instead of leaving them as-is, I followed each pain point and improved the project step by step.

Over time, it naturally grew into something more polished — something that fits into my daily routine and makes it easier.

Along the way, I ended up learning a lot:

- How to manage configuration with environment variables
- How to use GitHub Actions beyond CI workflows
- How to think about making small tools more reliable and easier to maintain

---

## 📈 What’s Next?

I don’t have a strict roadmap for what’s next — the project is already doing what I need it to do. That said, a few ideas have crossed my mind:

- Expanding support for more **engineering blogs or sources** I follow

For now, I’m happy with the system as it is. I’ll likely continue improving it **organically**, based on real usage and small annoyances I run into over time.

---

## ✅ TL;DR

- **Project**: `articles-extractor` scrapes basic article info and logs it to Google Sheets
- **Tools**: Python, BeautifulSoup, gspread, GitHub Actions
- **Learnings**: Automating with GitHub Actions was a better long-term solution than relying on Raspberry Pi hardware
- **Takeaway**: Start small, iterate often, and use each roadblock as a learning opportunity

---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
