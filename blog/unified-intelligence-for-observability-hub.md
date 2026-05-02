---
title: "Unified Intelligence for Observability Hub"
description: "Transitioning to Agentic Infrastructure via the Model Context Protocol unified gateway, providing specialized tools for automated RCA across telemetry, pods, network flows, and gateway capability state."
date: 2026-05-05
tags: ["observability", "mcp", "platform"]
draft: true
---

## Solving the Cognitive Telemetry Gap

Investigating the design of a Model Context Protocol server for the Observability Hub reveals complex architectural challenges. The core curiosity lies in how a unified protocol can bridge the semantic gap between fragmented telemetry sources and autonomous agents. Traditional static dashboards often lead to severe cognitive overload, causing human operators to miss critical failure patterns at scale. This exploration focuses on transforming the hub from a passive data store into an active, protocol-driven intelligence layer that supports autonomous reasoning.

---

## The Architectural Evolution

The journey to a unified interface was executed in distinct architectural phases to ensure stability. Each phase established a new layer of visibility while maintaining backward compatibility with existing telemetry. The evolution is documented through the following milestones:

- **Phase 1 (ADR 017)**: Established `mcp-telemetry` to provide a semantic bridge to the LGTM stack.
- **Phase 2 (ADR 018)**: Expanded into `mcp-pods` for deep visibility into the K3s cluster lifecycle.
- **Phase 3**: Transitioned away from fragmented host-level visibility (`mcp-hub`) to reduce security surface and protocol overhead.
- **Phase 4 (ADR 026)**: Consolidated the entire fleet into the `mcp-obs-hub` unified gateway.

---

## The Unified Tool Catalog

The current platform state is represented by a single entry point exposing specialized tools. These tools allow agents to reason across metrics, logs, traces, infrastructure events, network flows, and gateway capability state in a single loop. This consolidation eliminates the need for context switching between disparate providers. The catalog is categorized into four primary domains to provide comprehensive system coverage:

- **Telemetry**
  - `query_metrics`: Execute PromQL queries against Prometheus for metrics analysis.
  - `query_logs`: Execute LogQL queries against Loki for high-speed log analysis.
  - `query_traces`: Retrieve distributed traces from Tempo using trace IDs or specific queries.
  - `investigate_incident`: Correlate metrics, logs, and traces to assist agents in producing on-demand incident reports.
- **Infrastructure**
  - `inspect_pods`: List all pods in a namespace with a concise status summary.
  - `describe_pod`: Get detailed configuration and status for a specific pod instance.
  - `list_pod_events`: List all lifecycle events associated with a specific pod instance.
  - `get_pod_logs`: Retrieve stdout and stderr logs from a specific pod or container.
  - `delete_pod`: Terminate a specific pod to trigger a restart or manual remediation.
- **Network**
  - `observe_network_flows`: Query real-time network flows from Hubble Relay for deep traffic visibility.
- **Diagnostics**
  - `mcp_capabilities`: Report registered MCP domains, available tools, and skipped providers.

---

## Resilience via Soft-Fail Patterns

Operating in a hybrid environment requires a defensive approach to service initialization and tool registration. The unified gateway employs a specific pattern that gracefully skips tool registration if a provider like K3s or Thanos is unreachable. Startup records each domain as available or skipped, including skipped reasons, and exposes the same state through `mcp_capabilities`.

This ensures that localized infrastructure failures never blind the agent to the remaining functional parts of the system. Modular security remains a priority, with providers isolated behind `internal/mcp/providers` and `internal/mcp/tools` while sharing a single gateway.

---

## Architectural Significance

The repository leverages a library-first architecture to streamline the integration of new diagnostic tools. This design pattern ensures that core logic remains decoupled from the protocol transport layer. It allows for the seamless unification of diverse tools into a single MCP server rather than managing three fragmented instances. The primary benefits of this consolidation include:

- **Contextual Correlation**: Agents can now correlate network flow anomalies with pod lifecycle events in one reasoning loop.
- **Reduced Protocol Overhead**: A single `stdio` transport reduces the resource footprint on the control plane significantly.
- **Modular Security**: Providers remain isolated behind `internal/mcp/providers` and `internal/mcp/tools` to maintain clean boundaries while sharing a unified gateway.

---

## Conclusion

The transition to the MCP Era reveals how specialized servers bridge the critical gap between human operators and AI agents. This unified protocol creates a shared semantic language that translates complex telemetry into actionable intelligence. Agents now have the depth required to perform investigations that previously overwhelmed human cognitive limits. This foundation sets the stage for the upcoming transition to autonomous resource management.
