---
title: "Kubernetes Quotas and Limits"
description: "Prevent cluster noisy neighbor issues by configuring resource quotas and container limit ranges to enforce resource boundary defaults in shared developer environments."
date: 2026-07-28
tags: ["kubernetes", "platform"]
draft: true
---

## The Noisy Neighbor Problem in Shared Clusters

Deploying containers in a shared node without memory limits caused the host operating system to run out of memory. The kernel `out-of-memory killer` terminated critical runtime services, stopping all adjacent workloads. This resource exhaustion highlighted the danger of running unconstrained pods on shared compute nodes. Implementing logical `namespaces` and strict limit controls resolved these stability failures.

---

## Restricting Aggregates with Resource Quotas

Applying a namespace-level `resource quota` served as the first step to fix the crashes. The YAML manifest below implements this quota by setting hard limits on CPU, memory, and total pod counts. Restricting the namespace to a maximum of three pods prevents over-allocation on the shared compute node. The subsequent ASCII diagram illustrates how the `API server` blocks creation requests when resource or pod quotas are exceeded.

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: dev-quota
  namespace: dev-team
spec:
  hard:
    pods: "3"
    requests.cpu: "1"
    requests.memory: "1Gi"
    limits.cpu: "2"
    limits.memory: "2Gi"
```

```text
  [Request for 4th Pod]
           │
           ▼
  ┌─────────────────┐
  │   API Server    │ <───[ResourceQuota: dev-quota (Limit: 3 Pods)]
  └────────┬────────┘
           │
     (Quota Exceeded)
           │
           ▼
   [Request Blocked]
```

To illustrate this, a namespace `ResourceQuota` acts as a fixed resource boundary that all workloads share. In a theoretical scenario with a limit of one CPU, deploying two pods splits this budget into two allocations of `500m` CPU each. If the first pod contains three containers and the second contains two, the allocations are subdivided differently to fit the pod limits. The block diagram below maps this hierarchical resource division.

```text
  ┌──────────────────────────────────────────────────────────┐
  │ Namespace CPU Budget: 1000m (1 CPU)                      │
  ├────────────────────────────┬─────────────────────────────┤
  │ Pod 1 CPU Request: 500m    │ Pod 2 CPU Request: 500m     │
  ├──────┬──────┬──────────────┼──────────────┬──────────────┤
  │ C1   │ C2   │ C3           │ C4           │ C5           │
  │ 166m │ 166m │ 166m         │ 250m         │ 250m         │
  └──────┴──────┴──────────────┴──────────────┴──────────────┘
```

---

## Enforcing Safety Defaults with Limit Ranges

However, restricting namespaces still allowed single containers to run without default resources. Adding `limit ranges` solves this by injecting default memory and CPU requests into each individual container. For example, if two pods are deployed with three and two empty containers respectively, these defaults apply to each container separately. The YAML configuration below defines these container-level defaults, while the subsequent flowchart maps the injection process for both pods.

```yaml
apiVersion: v1
kind: LimitRange
metadata:
  name: dev-limit-range
  namespace: dev-team
spec:
  limits:
  - type: Container
    default:
      cpu: "500m"
      memory: "512Mi"
    defaultRequest:
      cpu: "250m"
      memory: "256Mi"
    max:
      cpu: "1"
      memory: "1Gi"
    min:
      cpu: "50m"
      memory: "64Mi"
```

```text
           [Pod 1 (3 Containers)]        [Pod 2 (2 Containers)]
                (Empty Specs)                 (Empty Specs)
                      │                             │
                      └──────────────┬──────────────┘
                                     ▼
                          ┌─────────────────────┐
                          │     API Server      │ <───[LimitRange: dev-limit-range]
                          └──────────┬──────────┘
                                     │
                             (Injects Defaults)
                                     │
              ┌──────────────────────┴──────────────────────┐
              ▼                                             ▼
  ┌─────────────────────────────────┐           ┌─────────────────────────────────┐
  │ Pod 1 with:                     │           │ Pod 2 with:                     │
  │  ├── Container 1 (500m / 512Mi) │           │  ├── Container 4 (500m / 512Mi) │
  │  ├── Container 2 (500m / 512Mi) │           │  └── Container 5 (500m / 512Mi) │
  │  └── Container 3 (500m / 512Mi) │           │                                 │
  └─────────────────────────────────┘           └─────────────────────────────────┘
```

---

## Comparing Resource Quotas and Limit Ranges

To illustrate the differences, resource quotas and limit ranges operate at different levels of granularity. While resource quotas control the aggregate namespace capacity, limit ranges regulate individual container behaviors. The table below outlines the core differences in scope, protection, and rejection timing.

| Aspect | ResourceQuota | LimitRange |
| --- | --- | --- |
| Scope | Namespace-level (aggregates all Pods) | Container / Pod-level (individual workloads) |
| Protects against | A namespace consuming more than its allocated share of resources | A container running without resource limits |
| Typical Rules | "This namespace cannot exceed 4 CPUs or 8Gi of RAM in total" | "Every container must have between 50m and 1 CPU, and defaults to 100m" |
| Rejection Timing | Workload creation is blocked once the aggregate quota is saturated | Workload creation is blocked instantly if a single container violates boundaries |

---

## Conclusion

Before this exploration, I was aware of resource limits but remained unclear on the difference between namespace and pod-level scopes. Discovering how to apply `quotas` at the namespace level and `defaults` at the container level resolved that knowledge gap for me. I now observe how namespace-level `quotas` block excessive aggregate usage while container-level `defaults` protect individual pods from running unconstrained. This dual-layered control establishes a reliable foundation for managing shared compute resources.
