---
title: "[Engineering Log] The Migration Strangler: Promtail to Alloy"
description: "Migrating from EOL Promtail to Grafana Alloy via a Strangler Fig pattern, bridging Kubernetes and Docker logs to ensure zero visibility loss during the platform transition."
date: 2026-03-10
tags: ["retrospective", "platform", "kubernetes"]
draft: true
---

## Context

Promtail is reaching **End-of-Life (EOL)**. To stay current with the Grafana ecosystem and move toward a unified telemetry agent (logs, metrics, and traces in one binary), I needed to migrate to **Grafana Alloy**.

---

## Challenge

* **The Pain Point:** Promtail's YAML-based pipelines were reaching their limits for complex relabeling, and staying on EOL software introduces significant security and technical debt.
* **The Question:** How do I migrate the "telemetry brain" of a live system without losing log visibility during the transition?

---

## Investigation

* **Concept: The Strangler Fig Pattern.** Instead of a "Big Bang" replacement, I deployed Alloy alongside Promtail.
* **Discovery:** Alloy uses **River** (an HCL-based language), which is programmable and component-based, unlike Promtail's static YAML.

| Feature | Promtail (Legacy) | Grafana Alloy (New) |
| :--- | :--- | :--- |
| **Config** | Static YAML | Programmable River/HCL |
| **Scope** | Logs only | Logs, Metrics, Traces |
| **Deployment** | Standalone Container | kubernetes-native DaemonSet |

---

## Solution

I deployed Alloy as a kubernetes `DaemonSet` to scrape the host's `/var/log/journal`. The critical insight was the **two-stage cutover**:

1. **Phase 1 (Hybrid):** Alloy (in Kubernetes) was configured to push to the **Docker-Loki** instance. This allowed me to compare Alloy's logs directly against Promtail's logs in the same Grafana instance to verify parity.
2. **Phase 2 (Native):** Once Loki was migrated to Kubernetes, I simply updated the endpoint URL.

```hcl
// Alloy's filtered pipeline: Only keeping what matters
loki.source.journal "systemd" {
  forward_to    = [loki.process.journal_pipeline.receiver]
  relabel_rules = loki.relabel.journal_relabel.rules
  // ...
}

loki.write "local_loki" {
  endpoint {
    // Initial: "http://host.docker.internal:3100/loki/api/v1/push"
    // Final (Kubernetes):
    url = "http://loki.observability.svc.cluster.local:3100/loki/api/v1/push"
  }
}
```

This "Strangler Fig" approach meant I never had a "blind spot" in my logs during the migration.

---

## Outcome

* **Result:** Successfully migrated host-level logs with zero downtime.
* **Lesson:** Components > Static Config. Being able to "wire" telemetry outputs to multiple destinations (Loki + Debug) made validation trivial.

---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
