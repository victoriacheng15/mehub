---
title: "January Reflection 2026"
description: "Reflecting on a shift from chatting with AI to refining context files. Learning to treat instructions as code and using the RFC process for better agent collaboration."
date: 2026-01-27
tags: ["growth"]
---

## The Shift: From "Chatting" to "Refining"

---

## What is January Reflection 2026?

In this post, we will explore **January Reflection 2026**. This month, I solidified a major shift in how I interact with AI agents like Gemini. Initially, I treated them like search engines—ask a question, get an answer. But as I tackled more complex tasks (like my Observability Hub), I realized that vague prompts produce "average" boilerplate.

To get useful output, I realized I needed to provide better input. I shifted my focus to refining the "context file"—treating my instructions as code that needs to be debugged and optimized.

---

## The Protocol: Context as Code

I found that the efficiency of an AI agent is linearly correlated to the precision of the instructions provided. I stopped treating it like magic and started treating it like code that needs debugging.

### 1. The Context File (`GEMINI.md`)

I started with simple prompts like "Help me write a LinkedIn post." The result was generic and often ignored my formatting rules.

So, I treated the instruction file (`.gemini/GEMINI.md`) as a config file that I iteratively patched:

- **Bug:** Agent used standard markdown bullet points. Since LinkedIn doesn't support markdown, these look like plain text dashes.
- **Patch:** Added explicit constraint: "Use emojis for lists to create visual structure without markdown."
- **Bug:** Agent hallucinated content.
- **Patch:** Added protocol: "Step 1: Fetch source context via `gh cli`."

### 2. The "RFC" Process (Plan Before Code)

I forced a workflow where the agent must **Plan** before it **Executes**.

- **Bad:** "Write a script to do X."
- **Good:** "Draft a plan to achieve X. Wait for my approval."

The AI doesn't always generate the plan I have in mind. By reviewing the plan first, I can refine the architecture and align the agent's "mental model" with mine. It’s an iterative loop of **Plan -> Review -> Refine** until I'm satisfied. This ensures the final code is exactly what I wanted, not just a generic guess.

---

## The Result

By focusing on these two critical areas—**Specific Rules** and **Iterative Planning**—I've significantly reduced the time spent correcting AI output.

Even with a solid plan, the AI can sometimes derail mid-execution. I've learned to be proactive: if I see the output drifting off-course, I "smash" the ESC button (or whatever the stop command is for your agent) to stop it immediately. There’s no point in letting it finish a wrong path.

The AI is no longer a "magic box" I hope will guess correctly; it’s a tool executing against a precise spec. Getting the most out of it isn't about finding the perfect prompt; it's about systematically refining the context and constraints, just like iterating on a piece of software.

---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
