---
title: "[Engineering Log] Measuring Memory Value in Echo"
description: "Adding local FinOps and Knowledge ROI analytics to Echo so persistent memory can measure cost, hit rate, carbon estimates, and low-signal context before quality drifts."
date: 2026-06-16
tags: ["go", "retrospective"]
---

## Context

Building a persistent memory system taught me that context utility degrades quickly as low-signal records accumulate across agent sessions. I observed that without an active feedback loop, search latency increased and LLM contexts became bloated with irrelevant historical data. Although SQLite and FTS5 successfully resolved the persistence layer, I realized that durable storage is useless if the system cannot measure whether the retrieved data remains valuable.

I learned that durable storage layers accumulate directives, facts, and artifacts that continuously skew vector and keyword search relevance. Retaining these records indefinitely introduced query noise into my retrieval step, degrading agent performance. This friction forced me to explore active evaluation methods to systematically identify and prune inactive context.

### Past Logs

- [\[Engineering Log\] Persistent AI Memory with Echo](https://victoriacheng15.dev/blog/engineering-log-persistent-ai-memory-with-echo)

---

## Challenge

My main challenge was building this telemetry loop without compromising the local-first, low-latency design of the memory engine. I wanted to observe query utility and estimate operational cost locally, which required an analytical engine that could process telemetry streams directly on-disk. To prevent database locks, I had to design an asynchronous telemetry mechanism that would not block active application queries.

I initially treated telemetry as a simple reporting feature, but I quickly discovered that static dashboards failed to influence memory lifecycle decisions. This wrong turn forced me to pivot, feeding the analytical results directly back into the execution path to automate record pruning. Through this feedback loop, I was able to retain high-value data while systematically decaying low-value records.

---

## Investigation

I split my investigation into three specific signals: financial cost, carbon impact, and knowledge utility. Mapping execution latency to regional API rates allowed me to estimate the synthetic dollar cost of each memory operation. I then converted latency-derived Joules into synthetic carbon estimates to understand the environmental footprint of my local queries. Finally, I tracked the ratio of successful queries to calculate a clear knowledge return on investment.

To keep my resource footprint small, I chose to decouple the transaction database from the analytical database. I implemented an append-only JSONL stream to log events asynchronously, avoiding table-locking issues during memory writes. I then introduced a local DuckDB instance to parse the JSONL stream, allowing me to run complex analytical queries without blocking active memory retrieval.

- `echo.db` stores operational memory, context keys, tags, and FTS5 search data.
- `events.jsonl` keeps telemetry writes cheap and append-friendly.
- `analytics.duckdb` handles analytical queries over memory activity.
- `RateService` applies configurable cost, energy, carbon, and decay values.

Decoupling these paths proved to be the correct choice, maintaining low-latency writes even during high concurrent loads. I accepted that my energy and cost metrics are directional estimations rather than precise physical measurements. However, observing these relative shifts gave me enough signal to compare performance trends and context utility across different sessions.

---

## Solution

I implemented a structured `TelemetryEvent` model to capture latency, hit status, and synthetic energy estimates. To ensure fast local queries were not lost, I recorded latency with sub-millisecond precision. Each event maps the current tool runtime back to its source interface, creating a clear history of context utility.

Using DuckDB allowed me to read the JSONL event log directly into column-oriented tables using native JSON parser tools. I then exposed these analytical tables to the host agent through the MCP `get_analytics` tool. This integration allowed the agent to dynamically inspect its own resource footprint and memory hit rates.

I connected this data stream to the memory engine using a decay algorithm inside `KnowledgeRefiner`. The refiner queries the analytics table to find memory records that consistently yield low hit rates. I programmed the system to decrement the importance score of these rows, automatically archiving them once the score hits zero.

---

## Evolution

This implementation shifted my design from a simple storage engine to a self-optimizing memory system. I realized that telemetry can transform a passive cache into an actively managed context graph. By evaluating context utility continuously, I can protect the agent from quality drift and information noise.

My most valuable lesson was that analytical databases must remain completely independent of the write path. Appending events to a local JSONL log kept my transactional database fast and lock-free. I designed the synchronization tool to process the log out-of-band, updating DuckDB tables only when the agent explicitly requests analytics.

Introducing DuckDB brought unexpected concurrency challenges due to its single-writer file-locking model. I learned that handling concurrent queries from both the CLI and the MCP server requires careful connection management. Resolving these file locks and coordinating connections has become my next focus in the development process.
