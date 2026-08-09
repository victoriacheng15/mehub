---
title: "Kubernetes Placement and Network Controls"
description: "Configure advanced pod placement using node affinity and taints, and secure container traffic using namespace network policies inside local microservice sandboxes."
date: 2026-08-11
tags: ["kubernetes", "platform"]
draft: true
---

## Unconstrained Placement and Exposed Services

Deploying containers on shared nodes without placement rules caused critical database workloads to run alongside heavy batch jobs. This poor scheduling led to high latency and resource contention on the worker nodes. Additionally, default flat networks allowed any pod in the cluster to access exposed database ports. Implementing explicit node steering and network isolation rules resolved these scheduling and security issues.

| Aspect | Node Affinity | Taints & Tolerations | Network Policy |
| --- | --- | --- | --- |
| Focus | Pod placement on nodes | Pod restriction on nodes | Network traffic between pods |
| Analogy | Matching tags | Lock and key | Firewall rule |
| Rule Direction | Pod attracts itself to nodes | Node repels pods without keys | Pod allows or blocks traffic |
| Example Use | Run web pods on fast SSD nodes | Keep database nodes clean from normal pods | Block web pods from reaching database ports |

---

## Pod Placement with Affinity and Taints

Steering pods to specific hosts requires using matching tags (affinity) and locks (taints). Node affinity ensures a pod schedules only on a host carrying a matching label. Conversely, a taint acts as a node lock that repels pods unless they carry a matching toleration key. The list below explains how these two features work together to secure the database workload.

For example, a cluster with three nodes running as virtual machines can separate workloads into `frontend`, `backend`, and `database` tiers:

```bash
# Label nodes by workload type
kubectl label nodes node-a type=frontend
kubectl label nodes node-b type=backend
kubectl label nodes node-c type=database

# Lock Node C to prevent standard pods from scheduling on it
kubectl taint nodes node-c dedicated=database:NoSchedule
```

- Node affinity rules ensure the database pod schedules only on Node C (`type=database`).
- Adding the lock taint to Node C blocks `frontend` and `backend` pods from running on it unless they carry the matching toleration key.
- This setup guarantees that database workloads remain isolated on dedicated database hardware.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: database-pod
  namespace: app-sandbox
  labels:
    app: database
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: type
            operator: In
            values:
            - database
  tolerations:
  - key: "dedicated"
    operator: "Equal"
    value: "database"
    effect: "NoSchedule"
  containers:
  - name: database
    image: redis:alpine
```

```text
                          [Pod: database-pod]
                           ├── Requires tag: type=database
                           └── Carries key: dedicated=database
                                   │
           ┌───────────────────────┼───────────────────────┐
           ▼                       ▼                       ▼
┌──────────────────────┐┌──────────────────────┐┌──────────────────────────────┐
│ Node A (Worker VM)   ││ Node B (Worker VM)   ││ Node C (Worker VM)           │
│ ├── Tag: frontend    ││ ├── Tag: backend     ││ ├── Tag: type=database       │
│ └── (No match)       ││ └── (No match)       ││ └── Lock: dedicated=database │
└──────────────────────┘└──────────────────────┘└──────────────────────────────┘
           │                       │                       │
       (Blocked)               (Blocked)              (Scheduled)
```

---

## Isolating Workloads with Network Policies

A network policy acts like a firewall that isolates selected pods by default. When an ingress policy selects a target pod, the cluster network plugin (such as Calico or Cilium) blocks all incoming traffic to that pod except from explicitly allowed sources. The YAML file below restricts access so only backend pods within the `app-sandbox` namespace can reach the database on its Redis port.

- Non-matching workloads (such as `app=frontend` or untrusted pods) are blocked when attempting to connect to the database.
- The network plugin enforces these packet-filtering rules at the kernel level before traffic reaches the target container.
- This boundary prevents lateral network movement even if adjacent workloads in the cluster are compromised.

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: secure-db
  namespace: app-sandbox
spec:
  podSelector:
    matchLabels:
      app: database
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: backend
    ports:
    - protocol: TCP
      port: 6379
```

```text
  ┌───────────────────────────┐
  │    [Pod: frontend-app]    │ ────────────X (Blocked: Web tier isolation)
  └───────────────────────────┘              \\
                                              \\
  ┌───────────────────────────┐                ▼
  │  [Pod: unauthorized-app]  │ ────────────X ┌───────────────────────────┐
  └───────────────────────────┘  (Blocked)    │    [Pod: database-pod]    │
                                              │   ├── Label: app=database │
  ┌───────────────────────────┐               │   └── Port:  6379 (Redis) │
  │    [Pod: backend-app]     │ ────────────> └─────────────▲─────────────┘
  └───────────────────────────┘    (Allowed: TCP 6379)      │
```

---

## Conclusion

Working through these examples clarified the difference between cluster hardware boundaries and runtime security boundaries:

1. **Tagging and locking nodes** ensures dedicated workloads (like our database) get predictable CPU and memory on dedicated machines.
2. **Restricting network ingress** ensures that physical placement alone is not trusted for security, blocking lateral traffic from the web tier.

Understanding how labels drive both scheduling decisions and firewall rules is the foundation for building resilient, production-ready clusters.
