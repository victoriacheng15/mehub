---
title: "[Engineering Log] Fixing DuckDB File Locks"
description: "Fixing DuckDB file locks in Echo by replacing the startup connection with on-demand database sessions that open, sync, query, and close during each analytics request."
date: 2026-06-23
tags: ["go", "retrospective", "mcp"]
draft: true
---

## Context

Running two commands at the same time caused DuckDB to lock up the database file. Echo uses SQLite for regular app data and DuckDB for analytics. The database file locked up whenever the CLI or server tried to access it at the same time.

Keeping these databases separate made sense because they do different work. SQLite needs to stay open constantly to handle active memory. DuckDB only needs to run when someone requests an analytics report.

### Past Logs

- [\\[Engineering Log\\] Persistent AI Memory with Echo](https://victoriacheng15.dev/blog/engineering-log-persistent-ai-memory-with-echo)
- [\\[Engineering Log\\] Measuring Memory Value in Echo](https://victoriacheng15.dev/blog/engineering-log-measuring-memory-value-in-echo)

---

## Challenge

The initial code started DuckDB as soon as the app booted up. `AnalyticsService` ran immediately, even if analytics were not needed yet. For the MCP server, this meant the database file stayed open as long as the server was running.

This caused issues when running other commands at the same time. A tool call might request analytics while a developer tried to run `maintain -sync` or inspect the database file. These different actions constantly blocked each other.

Opening the file too early caused all of these conflicts. The queries and schemas were actually correct. The real problem was keeping the database connection open for too long.

---

## Investigation

The first bug appeared when loading the rate configuration. The app looked for a relative `configs/rates.yml` file, which failed outside the project root folder. Moving the rate card details into code using `DefaultRateCard` and `NewRateService()` fixed that startup failure.

However, the file locking problem remained. The MCP server still opened the database file before anyone requested analytics. Embedded databases still lock files even if there is no separate database server running.

The solution was to change when the file opens and closes. Instead of keeping a connection open, the app can open the database only when needed. A simple five-step sequence keeps the file lock short.

- Start the MCP server without opening DuckDB.
- Open DuckDB only inside the `get_analytics` handler.
- Sync `events.jsonl` into the analytical store.
- Run the project impact query.
- Close DuckDB before returning from the tool call.

This changed how the app manages the database file. The connection now closes as soon as the query finishes. Other tools and CLI commands can now access the database file without getting locked out.

---

## Solution

The MCP code now takes a directory path instead of an active connection. When `get_analytics` is called, the handler opens `AnalyticsService`, schedules a deferred `Close()`, runs the queries, and returns the data. This keeps the database file open only during the query.

The CLI commands do not have this issue because they exit quickly. A CLI command can open the file, finish its work, and close the file when the command ends. The main issue was the server holding the connection open forever.

A cleanup guard was also added for system shutdowns. The app catches signals like `SIGINT` or `SIGTERM` and closes any active database connections. This acts as a backup to make sure the database is not left in a locked state.

---

## Evolution

The main lesson was that embedded databases still require lifecycle management. Removing a database server does not remove file locks or process conflicts. The app still needs to manage when it connects to the database file.

The final setup keeps DuckDB but limits when the app opens it. SQLite handles all the continuous operational data. DuckDB runs only when analytics are requested.

This change makes the analytics database safe for multiple tools to share. The server, CLI, and test suite can access the same file without locking each other out. The fix was simply opening and closing the database connection at the right time.
