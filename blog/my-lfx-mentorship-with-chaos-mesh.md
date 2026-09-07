---
title: "My LFX Mentorship With Chaos Mesh"
description: "Exploring my LFX Mentorship with Chaos Mesh, refactoring E2E tests into Gherkin BDD scenarios, contributing UI fixes, and growing into an active open-source contributor."
date: 2026-09-02
tags: ["kubernetes", "go", "cncf"]
draft: true
---

## Why Chaos Mesh?

My path into open source was far from linear. While I had experimented with minor contributions in the past, I struggled to find a way to dive deep into a large codebase. Entering a new project can be intimidating without clear initial direction, which is why the Cloud Native Computing Foundation (CNCF) LFX Mentorship program caught my attention. It stood out by offering a structured onboarding path, a clear project scope, and direct weekly syncs with core maintainers, which was exactly the support I needed to get started.

While browsing the list of LFX mentorship projects, [Chaos Mesh](https://github.com/chaos-mesh/chaos-mesh) caught my eye. The concept of deliberately injecting faults into pods to simulate real-world failure scenarios, such as high-traffic events, seemed fascinating, so I decided to apply. After submitting my application, I installed Chaos Mesh on my local machine to play around with it and see how it worked in general. Setting up the dashboard, simulating outages on local pods, and watching them recover made the concepts of chaos engineering click, making me even more excited about the prospect of working on the project.

---

## The Mentorship Experience

Working with mentors was a highly collaborative experience. Regular syncs and code reviews provided a helpful feedback loop for structuring the framework, managing Go modules, and aligning with the project's development standards. Asynchronous discussions on Slack and GitHub kept technical momentum moving across different timezones.

The focus of my mentorship was refactoring the PodChaos and NetworkChaos E2E tests into Gherkin-based BDD (Behavior-Driven Development) scenarios using the `godog` framework. In the legacy test suite, the test intent, fixture setups, custom resource creation, assertions, and cleanups were all mixed together in Go code. This created a high cognitive load for contributors, who had to understand both user-facing chaos behaviors and low-level Kubernetes implementation details simultaneously. Migrating to Gherkin decoupled these concerns, separating the high-level test specifications from the underlying Go execution.

```text
+-------------------------------------------------------------+
| [Layer 1] Specification (Gherkin Feature File)              |
| - Declares human-readable test intent (Given / When / Then) |
| - Example: Scenario: Inject PodKill on Target Pod           |
+-------------------------------------------------------------+
                              |
                              v executes step bindings
+-------------------------------------------------------------+
| [Layer 2] Test Automation (Godog Step Definitions in Go)    |
| - Translates steps into Kubernetes API client calls         |
| - Applies and cleans up Chaos Custom Resources (CRs)        |
| - Probes pod lifecycle and asserts recovery behavior        |
+-------------------------------------------------------------+
                              |
                              v submits CRs & queries state
+-------------------------------------------------------------+
| [Layer 3] System Under Test (Kubernetes Cluster & Engine)   |
| - Kube API Server: Persists Chaos CRs and workload state    |
| - Chaos Controller & Daemon: Reconciles CRs to inject faults|
+-------------------------------------------------------------+
```

During the 12-week program, the project was structured into four key milestones:

* [x] **BDD Framework Setup & CI Integration ([#5028](https://github.com/chaos-mesh/chaos-mesh/pull/5028)):** Integrated the core BDD test runner into the existing CI workflows to run side-by-side with the legacy test suite.
* [x] **PodChaos Migration ([#5013](https://github.com/chaos-mesh/chaos-mesh/pull/5013), [#5059](https://github.com/chaos-mesh/chaos-mesh/pull/5059)):** Ported the core PodChaos E2E test scenarios, such as PodKill actions, into plain-English feature files backed by reusable step definitions.
* [ ] **NetworkChaos Migration:** Porting network-level chaos scenarios (such as partitions, delays, and crossovers) into BDD.
* [ ] **Playbook Documentation:** Authoring a developer guide to ensure future contributors can write new tests by mixing and matching existing steps without writing any Go code.

---

## Beyond the Mentorship

My initial goal was simply to get involved in the CNCF ecosystem through the LFX Mentorship. However, the experience inspired me to turn open-source contribution into a regular habit, not just for Chaos Mesh, but for other open-source projects in the future. Instead of focusing solely on the E2E test migration, I began exploring other areas of the Chaos Mesh monorepo. This allowed me to contribute across the entire stack, including the backend APIs, the React-based UI dashboard, and the container build pipelines.

Some of my subsequent contributions include:

* **Custom Subpath Support for Chaos Dashboard ([#5016](https://github.com/chaos-mesh/chaos-mesh/pull/5016)):** In many deployment environments, operators host the Chaos Dashboard behind reverse proxies under a URL subpath rather than the root domain. Previously, the dashboard UI had hardcoded absolute asset paths, which broke stylesheet and script loading under custom subpaths. I configured Vite's build settings (`base: './'`) and updated the HTML template to load assets dynamically, resolving pathing failures.
* **Dashboard Stability & Authentication Fixes ([#5034](https://github.com/chaos-mesh/chaos-mesh/pull/5034), [#5046](https://github.com/chaos-mesh/chaos-mesh/pull/5046)):** To improve dashboard usability, I worked on several fixes across the UI. I resolved a React runtime crash (TypeError) when loading workflow-type schedule details by adding optional chaining for selector queries. I also debugged a token validation loop in the API interceptors, ensuring expired tokens are correctly cleared to prevent authentication lockout issues.

A complete log of my Chaos Mesh pull requests, issues, and ongoing contributions is available on my [work page](https://victoriacheng15.dev/work).

---

## Looking Ahead

Today, the journey continues. While the migration is still ongoing, my contributions have extended past the official end of the mentorship. The journey doesn't end when the mentorship ends; it is a continuous process of learning, contributing, and collaborating with the community.

Sincere appreciation goes to mentors Zhiqiang ZHOU, Yue Yang, and Shenan Zhang, for their time, guidance, and support throughout the program. Open source thrives when people are willing to help others learn, contribute, and succeed. That supportive environment turns an initial mentorship into a lasting foundation for community participation.
