---
title: "[Engineering Log] Safer ID-Based Memory Deletion"
description: "A mid-level developer retrospective on shifting Echo memory deletion from content matching to ID-based operations to ensure AI tools target database rows with no ambiguity."
date: 2026-06-30
tags: ["go", "retrospective", "mcp"]
draft: true
---

## Context

Designing deletion flows is a common task for mid-level developers building AI integrations. The implementation must remain simple but predictable under fuzzy agent behavior. This log documents the transition from basic database patterns to deterministic tool boundaries.

Echo already had the right human workflow: retrieve, confirm, act. The weak point was the target used during the final action. The user confirmed a displayed memory, but the delete operation still depended on reconstructed content and context values.

### Past Logs

- [\\[Engineering Log\\] Persistent AI Memory with Echo](https://victoriacheng15.dev/blog/engineering-log-persistent-ai-memory-with-echo)
- [\\[Engineering Log\\] Measuring Memory Value in Echo](https://victoriacheng15.dev/blog/engineering-log-measuring-memory-value-in-echo)
- [\\[Engineering Log\\] Fixing DuckDB File Locks](https://victoriacheng15.dev/blog/engineering-log-fixing-duckdb-file-locks)

---

## Challenge

The initial design chose natural keys to simplify the database schema. This decision proved to be a mistake when integrating autonomous tool calls. A small query variation broke the deterministic mapping required for destructive operations.

```sql
DELETE FROM memories WHERE content = ? AND context_key = ?;
```

This pattern became problematic when asking an AI tool call to carry exact content strings across a destructive boundary. Content can be long, visually similar, or byte-different because of whitespace. A model may summarize, trim, normalize, or select a nearby memory while trying to help.

Confirmation can create a false sense of validation if the final delete target remains fragile. A single character mismatch breaks the deletion flow or targets the wrong record. The design needed a more deterministic identifier to ensure reliability.

---

## Investigation

The investigation separated two identity needs. Storage benefits from natural-key deduplication because repeated content in the same context should reinforce the existing memory. Deletion needs exact record targeting because the goal is to remove one confirmed row.

That made the database `id` the better destructive target. The user does not need to reason about the number directly. The user confirms the human-readable record, and the tool acts on the machine-readable primary key returned by search.

```text
Natural Key Deletion:
                      "Hello" + "Ctx A"
                              |
                              v
                   +---------------------+
                   | DELETE FROM memories |
                   +---------------------+
                              |
             +----------------+----------------+
             |                |                |
             v                v                v
      +------------+   +------------+   +------------+
      |   Row 1    |   |   Row 2    |   |   Row 3    |
      | ID: 101    |   | ID: 102    |   | ID: 103    |
      | "Hello"    |   | "Hello"    |   | "Hello"    |
      | (Deleted)  |   | (Deleted)  |   | (Deleted)  |
      +------------+   +------------+   +------------+

ID-Based Deletion:
                           ID: 102
                              |
                              v
                   +---------------------+
                   | DELETE FROM memories |
                   +---------------------+
                              |
                              v
                       +------------+
                       |   Row 2    |
                       | ID: 102    |
                       | "Hello"    |
                       | (Deleted)  |
                       +------------+
                         Row 1 & 3
                         Preserved
```

The same precision pattern also applied to updates. Surgical updates should modify a specific record and preserve metadata unless a field is intentionally replaced. That pushed both deletion and update behavior toward ID-based operations.

---

## Solution

The MCP delete tool changed from content and context arguments to a single numeric `id`. The service layer added `DeleteMemoryByID(id int64)`, and the CLI now mirrors the same decision through `echo delete -id <memory-id>`. There is no content-based CLI delete path in the current model.

The execution protocol became explicit to prevent accidental deletions. The workflow relies on clear verification steps before any mutation occurs. The interface enforces three distinct steps for tool interaction:

- Use `search_for_deletion` to retrieve the exact memory and its unique `id`.
- Display that memory to the user for confirmation.
- Call `delete_memory` with the returned `id` only after confirmation.

```text
+-----------------------+
|  search_for_deletion  | ---> Retrieves ID and Content
+-----------------------+
            |
            v
+-----------------------+
|   Human Reviewer      | ---> Confirms details (ID, tags, context)
+-----------------------+
            |
            v (Yes)
+-----------------------+
|     delete_memory     | ---> Deletes exact row by ID
+-----------------------+
```

This kept human review and machine precision in separate roles. The user confirms content, context, type, tags, and metadata. The tool deletes by primary key.

---

## Evolution

The change was small, but it moved the trust boundary. Echo did not just make deletion easier. It made deletion more deterministic.

The useful lesson was that verification workflows can fail at the final step. Confirmation matters, but confirmation alone is not enough if the act step uses fuzzy identifiers. Destructive tools need stable targets.

This is especially important for AI memory. Lost context can change future behavior, and ambiguous deletion creates unnecessary risk. ID-based operations make the system easier to trust because the confirmed record and deleted record are the same row.


