---
title: "[Engineering Log] Why Not Index Everything?"
description: "An investigation into B-Trees, the library analogy, and why we don't index every column."
date: 2026-01-13
tags: ["retrospective", "backend"]
---

## 1. Context (The "Why")

I intuitively knew SQL indexes were fast, but I didn't understand the underlying mechanics. I treated them like magic: add an index, query gets faster. I wanted to look under the hood to understand the *cost* of that speed.

> **The Analogy:** Imagine a library with 1,000,000 books arranged randomly. Finding one book requires a **Full Table Scan (O(N))**. An index is the card catalog that points you to the exact shelf **(O(log N))**.

## 2. The Challenge / Question

If indexes provide such a significant performance boost (O(log N) vs O(N)), the obvious question arises: **Why don't we just index every single column in every table?**

## 3. Investigation & Trade-offs

I researched how Postgres and MySQL actually implement these catalogs using **B-Trees (Balanced Trees)**.

| Concept | Gain (Pros) | Cost (Cons) | Verdict |
| :--- | :--- | :--- | :--- |
| **B-Tree Indexing** | Near-instant reads; minimizes disk I/O 'hops'. | Write Penalty & Disk Space. | ✅ Use Strategically |

### Discovery: The Write Penalty

I realized that databases don't just "list" keys. They maintain a complex, balanced tree structure.

* **The Cost:** Every `INSERT`, `UPDATE`, or `DELETE` triggers a tree update.
* **The Friction:** The true friction occurs when a B-Tree node (a 'page') is full. To insert a new key, the database must perform a **'page split'**:
  * It creates a **new child page**.
  * It moves half the data from the old page to the new one.
  * It updates the **parent page** to point to both child pages.

## 4. The Solution / Insight

The "Magic" is actually a deliberate architectural trade-off. We are trading **Write Speed** and **Disk Space** for **Read Speed**.

```sql
-- This speeds up SELECT...WHERE title = '...'
CREATE INDEX idx_books_title ON books(title);

-- But it SLOWS down this:
INSERT INTO books (title, author) VALUES ('New Book', 'Author');
-- Because the DB must now also update idx_books_title.
```

## 5. Outcome & Learnings

* **Result:** A shift from "magic" to intentional engineering. I now evaluate the read/write ratio of a table before adding an index.
* **Lesson:** "Everything in engineering is a trade-off." (A core theme from *Software Engineering at Google*).

---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
