---
title: "[Engineering Log] Building Echo Dual Interface"
description: "Turning Echo from an MCP-only memory server into a dual-interface local tool with agent workflows, CLI inspection, shared services, and direct maintenance paths for humans."
date: 2026-07-07
tags: ["go", "retrospective", "mcp"]
---

## Context

Echo began as an MCP memory server for AI agents. The first version let agents store, recall, and search durable project context through JSON-RPC over stdio. That solved the agent workflow, but it created a new operational artifact: a local memory store that humans also needed to inspect.

Once a system stores durable knowledge, the database is no longer just an implementation detail. Someone needs to correct records, audit search results, rebuild indexes, sync analytics, and recover from bad writes. Asking an agent to mediate every maintenance task made the system harder to trust.

### Past Logs

- [\[Engineering Log\] Persistent AI Memory with Echo](https://victoriacheng15.dev/blog/engineering-log-persistent-ai-memory-with-echo)
- [\[Engineering Log\] Measuring Memory Value in Echo](https://victoriacheng15.dev/blog/engineering-log-measuring-memory-value-in-echo)
- [\[Engineering Log\] Fixing DuckDB File Locks](https://victoriacheng15.dev/blog/engineering-log-fixing-duckdb-file-locks)
- [\[Engineering Log\] Safer ID-Based Memory Deletion](https://victoriacheng15.dev/blog/engineering-log-safer-id-based-memory-deletion)

---

## Challenge

An MCP-only interface fit the first use case, but it had practical limits. Manual inspection required either raw SQL or a conversational round trip. Maintenance tasks like FTS5 rebuilds and analytics sync were awkward as chat operations.

The challenge was adding human control without creating a second application. A separate CLI with copied logic would create two policy surfaces. Echo needed two interfaces, but one service layer.

---

## Investigation

The architecture moved toward a single binary with smart routing. If no subcommand is provided and stdin is not a terminal, Echo assumes an MCP host is launching it and starts the stdio server. If stdin is a terminal, it prints human-facing usage.

Known subcommands route to the CLI. That lets one installed `echo` binary support agent workflows and human workflows without separate deployment steps. The routing rule also keeps MCP stdout protocol-clean by sending logs and usage behavior to the right path.

The important design point was avoiding duplicate business logic:

- MCP maps JSON tool arguments into service calls.
- CLI maps flags into service calls.
- `internal/service` owns validation, storage behavior, search behavior, telemetry hooks, and destructive operations.

That kept protocol code at the edge and preserved one memory policy model.

---

## Solution

The CLI became a terminal curation tool on top of the same services used by MCP. It added `store`, `recall`, `search`, `delete`, `maintain`, and `help`. The commands gave humans direct access to the memory system without bypassing validation.

The output layer made the tool usable in more than one mode. `table` supports terminal review, `json` supports scripts and structured inspection, and `csv` supports spreadsheet-style auditing. That turned the CLI from a command wrapper into an operator interface.

The shared service layer kept behavior consistent. Content limits, `context_key` validation, `entry_type` rules, UPSERT reinforcement, hybrid search, ID-based deletion, FTS5 rebuilds, and analytics sync all stayed behind the same core logic.

---

## Evolution

The dual-interface work changed Echo's identity. It stopped being only an MCP server and became a local memory system with an agent interface and an operator interface. That distinction mattered because persistent AI memory became something humans needed to maintain directly.

The useful lesson was that agent-facing systems still need human affordances. A memory server that only an agent can operate is hard to debug, hard to trust, and hard to recover. Durable knowledge needs a control plane.

The CLI is not a fallback. It is the human side of the same system. MCP gives agents access to memory during work, and CLI gives humans access to the memory after it becomes operational state.
