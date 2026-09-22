---
title: "Automated Tag Validation for Static Blogs"
description: "A zero-dependency Python audit script integrated into the local build pipeline enforces frontmatter schema consistency and eliminates tag drift across static blog posts."
date: 2026-09-22
tags: ["platform", "python"]
---

## No Visibility Into What Tags Exist

Running a blog with hundreds of posts inevitably leads to silent frontmatter drift. Without automated validation, there is no mechanism to verify tag spelling, enforce allowable taxonomies, or constrain tag counts. Ad-hoc `grep` sweeps are manual and rarely executed proactively. The result is unvalidated conventions that exist only in documentation.

The problem compounds as the taxonomy expands. A post tagged `kubernets` instead of `kubernetes` never fails a standard Markdown build. Similarly, an `engineering-log` entry that omits its mandatory `retrospective` tag renders without warning. Without a validation gate, schema violations persist indefinitely until discovered by accident.

---

## Writing the Audit Script

The solution is a lightweight script at `scripts/audit_tags.py`. It inspects every `.md` file in `blog/`, parses frontmatter metadata, and enforces four constraints:

1. Verifies all tags exist within a defined `KNOWN_TAGS` set.
2. Caps tags at three per post to prevent taxonomy bloat.
3. Enforces that all `engineering-log-*.md` posts carry the `retrospective` tag.
4. Ensures standard posts do not inadvertently carry `retrospective`.

```text
blog/*.md
    ├── KNOWN_TAGS check      → catches unknown tag names
    ├── max 3 tags check      → enforces tag ceiling
    └── engineering-log check
            ├── must carry 'retrospective'
            └── non-logs must NOT carry 'retrospective'
```

The script prints a deduplicated list of active tags alongside pass/fail statuses for each invariant. By relying strictly on the Python standard library, the script introduces zero dependencies and executes in milliseconds. Standard POSIX exit codes (`0` on success, `1` on violation) make it directly composable with local Makefiles and CI pipelines.

---

## What the First Run Found

Running the script across hundreds of existing posts immediately surfaced a latent violation: `scaling-with-systemd-template-units.md` carried the `retrospective` tag despite being a technical tutorial. Because static site generators ignore semantic frontmatter mismatches, this error had remained invisible.

The audit also established an accurate baseline: 21 unique tags were in active use across the corpus. Exposing this consolidated set makes taxonomy consolidation and near-miss detection trivial. Aside from the single miscategorized systemd post, all remaining posts conformed to the schema rules, establishing a clean baseline for future posts.

---

## Gating the Build Pipeline

A standalone audit script provides visibility, but relying on manual execution creates a clear failure mode: if I forget to invoke the script before committing, invalid metadata still slips into the generated site. The pragmatic solution is wiring the audit directly into `make build` and the local development server. Any frontmatter violation halts site generation immediately before templates render.

```text
make build / dev server
    └── python3 scripts/audit_tags.py
            ├── exit 0  → generate static site in dist/
            └── exit 1  → abort build immediately with errors
```

Running the check locally ensures issues surface immediately in my terminal during development. The static site generator never attempts to render pages if frontmatter contains an unknown tag or exceeds the limit. Diagnostic output points directly to the offending file and the exact violation, keeping the feedback loop instant without adding complex build tooling.

---

## Conclusion

Relying on memory to maintain tag consistency across hundreds of posts was never sustainable. Writing a tiny Python script turned vague conventions into deterministic checks that run on every build. If I make a typo or accidentally introduce an extra tag, the build halts immediately. Shifting validation into the pipeline eliminates cognitive overhead, leaving me with one less thing to worry about when sitting down to write.
