---
title: "[Engineering Log] Killing the Static .env with OpenBao"
description: "Eliminating static .env files by centralizing secrets in OpenBao via a reusable pkg/secrets module, decoupling credentials to streamline the Kubernetes migration phase."
date: 2026-03-03
tags: ["retrospective", "platform"]
draft: true
---

## 1. Context (The "Why")

While static `.env` files worked perfectly for my local setup, they limited my exposure to **Centralized Secret Management**. I wanted to move beyond basic file-based configuration to explore how industry-standard tools handle secrets, rotation, and access control in real-world environments.

---

## 2. The Challenge / Question

* **The Pain Point:** Sticking with `.env` files was safe but "boring"; it offered zero opportunity to learn about secret leasing, audit logs, or the operational overhead of a dedicated secret store.
* **The Question:** Can I implement a production-grade secret store that remains truly open-source (avoiding BSL licenses) while gaining hands-on experience with the Vault ecosystem?

---

## 3. Investigation & Trade-offs

| Option / Concept | Pros / Gain | Cons / Cost | Decision / Verdict |
| :--- | :--- | :--- | :--- |
| **HashiCorp Vault** | Industry standard, massive ecosystem. | Business Source License (BSL) limits. | ❌ |
| **OpenBao** | Fork of Vault, MPL 2.0 (True Open Source). | Newer community, smaller footprint. | ✅ |

* **Discovery:** OpenBao maintains API compatibility with Vault, meaning I can use existing Go SDKs without modification.
* **Trade-off:** I traded **Operational Simplicity** (env vars) for **Security & Scalability** (a dedicated secret service).
* **Strategic Win:** Decoupling secrets from the file system early on proved to be a "Force Multiplier" for the kubernetes migration. Since services already knew how to fetch credentials via API, moving them into pods required zero changes to secret-handling logic.

---

## 4. The Solution / Insight

The breakthrough wasn't just installing OpenBao; it was the creation of a centralized `pkg/secrets` module. By standardizing secret retrieval into a reusable internal package, I ensured that any new service (Go-based) could instantly tap into the secure store without reinventing the connection logic.

```go
// The centralized interface in pkg/secrets that standardized our 'Paved Road'
type SecretStore interface {
    GetSecret(ctx context.Context, path string) (map[string]interface{}, error)
}

// OpenBao implementation: abstracting the 'how' so services focus on 'what'
func (b *BaoStore) GetSecret(ctx context.Context, path string) (map[string]interface{}, error) {
    secret, err := b.client.Logical().Read(path)
    // ... logic to handle fallbacks or missing paths
}
```

This "Paved Road" approach meant that when it came time to migrate to Kubernetes, the services were already "cloud-native" in their configuration handling.

---

## 5. Outcome & Learnings

* **Result:** Migrated `postgres` and `mongodb` to dynamic retrieval; no more secrets in the file system.
* **Lesson:** Abstraction layers (interfaces) are the only way to migrate infrastructure without rewriting core business logic.
