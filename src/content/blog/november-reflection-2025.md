---
title: "November Reflection 2025"
description: "Built a cover image generator with Azure Functions to learn serverless—explored cold starts, CI/CD, and planning-first development."
date: 2025-11-25
tags: ["growth"]
---

## Shipping Serverless

This November, I finished **Cover Craft**, a lightweight cover image generator I built to get hands-on experience with **Azure Functions**. I wanted to understand serverless architecture by using it myself from writing a function to deploying it in the cloud. I also plan to use it to generate cover images for my LinkedIn posts, so it serves a real purpose for me.

The core task is simple: accept text, colors, and dimensions; return a PNG. But how it is built and especially how it is deployed taught me far more than the feature ever could.

Here is what I learned by focusing on the backend and pipeline.

---

## Azure Functions: The Power and the Tradeoffs of Serverless

Working with Azure Functions confirmed both the promise and the reality of serverless. On one hand, there are clear benefits: no VMs to manage, no infrastructure to provision, and automatic scaling that just works. For a task like generating an image on demand, the model fits naturally.

On the other hand, I experienced the most common drawback firsthand: **cold starts**. When the function had not been used for a while, the first request took about **1.5 to 2 seconds** to respond while Azure spun up a new execution environment. For a personal tool like Cover Craft, that delay is fine. In a high-traffic or latency-sensitive application, it would need mitigation.

The biggest win was operational simplicity. I never touched a virtual machine, never configured a load balancer, and never thought about OS updates. All that complexity disappeared.

But new responsibilities appeared in its place. I had to manage:

- Secure CI/CD authentication using an Azure **service principal** stored in **GitHub Secrets**
- Local development with the **Azure Functions Core Tools**, which are essential for testing before deployment

These are not optional extras. They are the practical work of making serverless reliable.

The result was a **fully automated deployment pipeline** that goes from `git push` to a live function in under two minutes with no manual steps.

> **Lesson**: Serverless eliminates infrastructure management, but not engineering responsibility.

---

## DevOps That Runs Only What’s Necessary

The biggest operational insight from this project was not about cloud APIs, it was about **workflow design**. Early on, I knew I did not want every pull request to trigger the full suite of frontend and backend tests, especially since Cover Craft has a clear separation between the API and the frontend. Running everything on every change would waste CI minutes and slow down feedback.

So I designed **path-filtered GitHub Actions workflows from the start**:

- Changes in `api/**`? → Trigger **only** the Azure Functions CI and deployment.
- Changes in `frontend/**`? → Run frontend linting and tests (in parallel).
- Touch a `.md` file? → Just validate markdown formatting.

This was not an optimization I added later, it was baked in to keep things **fast, focused, and frugal**.

My `azure_function.yml` workflow runs **only when it needs to** on merge to `main` with changes in `api/**`, or via manual dispatch. It:

1. Installs dependencies and builds the project  
2. Runs all tests (blocking deployment on failure)  
3. Deploys via `az functionapp deployment source config-zip`  

Because this workflow is defined in code and triggered precisely, **my deployment process is reproducible, auditable, and safe** without unnecessary noise or cost.

> **Lesson**: Good DevOps is not about running more, it is about running *only what matters*, *when it matters*.

---

## Why This Matters

This project mattered because it gave me a real way to learn Azure Functions by doing. I worked with both the Azure CLI and the web portal, which helped me understand how serverless actually works beyond the documentation.

I saw firsthand that even a simple API involves more than just writing a function. You need to set up secure authentication for CI/CD, structure your code so it can be tested, and understand what happens during deployment. These are not abstract ideas. They are the practical details that show up as soon as you try to run your code in the cloud.

And that is exactly what I set out to learn.

---

## Looking Ahead

Now that I have a working image generator, I have started thinking about small ways to extend it. One idea I find interesting is **batch generation**. What if I could upload a CSV file with a list of titles and generate multiple cover images at once? It would be a natural way to test how Azure Functions handles multiple invocations or longer running workloads.

I am also curious whether I could use **Terraform to create the Function App and its resources**, instead of setting everything up manually in the portal. Infrastructure as code feels like the logical next step to make the setup repeatable and version controlled.

But it is encouraging that finishing this small project has already sparked a few concrete ideas for what to try next.

---

## Final Thought

This time I tried something different. Instead of jumping straight into code, I started by planning the backend architecture. I wrote down how I wanted the API to work, defined the data flow, and sketched out the core components before writing a single function. Once the design felt solid, I built the backend and tested it locally. Only after it worked as expected did I deploy it to Azure Functions to verify it behaved the same in the cloud.

With the backend stable, I moved to the frontend. I even sketched the layout using ASCII art to visualize the structure, planned the color theme, and mapped out shared components ahead of time.

This upfront planning felt new to me and effective. Knowing the shape of the system early made the code easier to write, the CI/CD pipeline simpler to configure, and the whole project less chaotic. It turns out that a little design before development can go a long way.

You can check out the [repo](https://github.com/victoriacheng15/cover-craft) on GitHub!

---

## Thank you

Big thanks for reading! You are awesome, and I hope this post helped. Until next time!
