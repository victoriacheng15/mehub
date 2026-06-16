---
title: "[Engineering Log] Persistent AI Memory with Echo"
description: "Building an MCP server with Go and SQLite to give AI agents scoped, local-first memory, reduce repeated context setup, and keep project decisions available in CLI sessions."
date: 2026-04-14
tags: ["go", "mcp", "retrospective"]
---

## Context

Curiosity about Model Context Protocol design led to this project. The goal was to understand how context is managed between a local environment and a large language model. I adopted a design-first methodology that prioritized the data model before writing any connectivity logic.

This methodical approach ensured that the underlying system was robust from the start. Once the data structure was solid, I set up the specific methods the server needed to function.

---

## Challenge

Constant repetition of project details to CLI agents became tiresome. Every new session required explaining architectural preferences and specific technical constraints repeatedly. This lack of persistence slowed down the workflow and led to inconsistencies in AI-generated code.

The amnesiac agent problem meant that the model started every session with a blank slate. Even with sophisticated system prompts, the context window remained limited and ephemeral. This repetition was not just a minor annoyance but a significant drain on developer velocity.

---

## Investigation

Starting with the schema exposed the first useful constraint: memory retrieval had to be scoped before it could be fast. A simpler single-bucket design would have been easier to build, but it risked mixing unrelated project facts across sessions. That trade-off pushed the implementation toward explicit context keys before the MCP handlers were wired up.

This exploratory phase involved the following key technical decisions:

- **Language choice**: Go was selected for its performance and native support for JSON-RPC over stdio.
- **Storage engine**: SQLite provided an ideal balance of performance and simplicity for local-first tools.
- **Concurrency**: Write-Ahead Logging was enabled to handle reliable simultaneous operations for multiple tools.
- **Search capability**: The built-in FTS5 module offered deterministic keyword-based search for technical snippets.
- **Maintenance**: The zero-configuration nature of SQLite eliminated the operational overhead of a full database server.

These choices kept the implementation simple while ensuring a portable and zero-configuration memory layer.

---

## Solution

A dual-key scoping strategy ensures relevance to the current task. Project Scope (project:name) prevents the agent from confusing requirements between different projects. Global Scope (global) provides universal truths across all workspaces.

The implementation relies on a `context_key` column and a composite index for fast retrieval. This index prioritizes records based on their importance score, ensuring that the most critical context is always available to the agent.

```sql
CREATE TABLE IF NOT EXISTS memories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL CHECK (length(content) > 0 AND length(content) <= 8192),
    context_key TEXT NOT NULL CHECK (length(context_key) > 0),
    entry_type TEXT DEFAULT 'directive' CHECK (entry_type IN ('directive', 'artifact', 'fact')),
    importance_score INTEGER DEFAULT 1 CHECK (importance_score BETWEEN 0 AND 10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    source TEXT DEFAULT 'mcp' CHECK (source IN ('mcp', 'cli')),
    is_active BOOLEAN DEFAULT 1,
    tags TEXT CHECK (tags IS NULL OR json_valid(tags)),
    UNIQUE(content, context_key)
);

CREATE INDEX IF NOT EXISTS idx_context_relevance ON memories(context_key, importance_score DESC);
CREATE INDEX IF NOT EXISTS idx_content ON memories(content);
CREATE INDEX IF NOT EXISTS idx_is_active ON memories(is_active);
```

---

## Evolution

Echo transforms the AI agent into a collaborative partner that understands historical context. Agents check Echo for patterns when scaffolding modules or fixing bugs. This persistence is the missing piece of the AI engineering journey.

By building tools that remember past decisions, the workflow no longer depends on repeating the same setup notes at the start of each session. Persistent context gives the agent a stable reference point for project rules, architectural choices, and recurring constraints. The result is less context rebuilding and more consistent help across CLI sessions.
