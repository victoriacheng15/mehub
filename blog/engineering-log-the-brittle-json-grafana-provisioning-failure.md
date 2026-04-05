---
title: "[Engineering Log] The Brittle JSON: Grafana Provisioning Failure"
description: "Fixing Grafana dashboard provisioning by moving from brittle embedded JSON in values.yaml to standalone files, resolving 'invalid character' bugs once and for all."
date: 2026-04-07
tags: ["retrospective", "platform"]
draft: true
---

## Context

Originally, I wanted to keep the project structure as lean as possible. To minimize the number of files, I combined the Grafana dashboard configurations directly into the `values.yaml` file as embedded JSON strings. I thought this would make the orchestration easier to manage by keeping everything in one place. This approach initially seemed efficient for a small-scale local lab deployment.

---

## The Challenge

- **The Pain Point**: After a service restart, the dashboards failed to load due to "invalid character" errors.
- **Visual Friction**: Because the JSON was embedded within a large YAML block, it was nearly impossible to see where the escaping had failed.
- **The Question**: Is the convenience of "fewer files" worth the loss of readability and the brittle nature of manual JSON-in-YAML embedding?

The dashboards remained broken as long as the data was trapped inside the scalar block. Finding a single malformed character in a 500-line block created significant visual fatigue during the debugging process. This visual noise acted as a barrier to identifying the actual syntax error. The architectural trade-off had prioritized directory aesthetics over system observability and reliability.

---

## Investigation

- **Discovery**: The "invalid character" errors were caused by **Configuration Fragility**.
- **Structural Failure**: Embedding complex JSON containing `rawSql` with newlines into a YAML block scalar (`|`) is highly error-prone.
- **The Conflict**: Finding a single malformed character inside a 500-line YAML block scalar is a massive friction point compared to linting a standalone JSON file.

The mapping error became systemic when the YAML parser misinterpreted the white space of the SQL queries. This brittle connection between formats meant that even a minor change to a query could break the entire dashboard. Testing these changes required a full deployment cycle because the source was not independently verifiable. The investigation confirmed that the orchestration layer was mangling the source integrity of the JSON data.

---

## Solution

I abandoned the "fewer files" goal in favor of **Source Integrity**. I moved all dashboards into dedicated `.json` files and updated the orchestration to mount them directly. This resolved the parsing issues "in one go" because the source files were no longer being mangled by YAML's block scalar logic. The transition to a directory-based provisioning model ensures that the dashboards are deployed with byte-for-byte accuracy.

```makefile
# The permanent fix implemented in makefiles/kubernetes.mk
k3s-grafana-dashboards:
 kubectl create configmap grafana-dashboards \\
  --from-file=k3s/grafana/dashboards/ \\
  --dry-run=client -o yaml | kubectl apply -f -
```

---

## Evolution

- **Result**: Dashboards now provision perfectly every time without manual escaping.
- **Maturity Milestone**: Decoupling the data (JSON) from the configuration (YAML) removed the entire class of "invalid character" bugs.
- **Lesson**: Don't sacrifice the integrity of your data files for the sake of having a smaller file count.

Standalone files are easier to lint, easier to debug, and much more reliable than embedded strings. This structural shift represents a move toward architectural maturity by prioritizing data robustness over file count. The Hub now benefits from a more stable provisioning pipeline that respects the native format of its components. Future additions to the monitoring layer will follow this decoupled pattern to maintain high-signal observability.
