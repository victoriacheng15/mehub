---
title: "Merged the First Chaos Mesh PR"
description: "A first LFX Chaos Mesh milestone initiates the Ginkgo to Gherkin migration, proving a readable BDD test path in Kubernetes e2e tests with the first merged pull request."
date: 2026-07-14
tags: ["go", "kubernetes", "cncf"]
---

## The First Merge Landed

Acceptance into the LFX Mentorship with `Chaos Mesh` transitioned the `BDD` migration from a plan into active contributions. The task is to help migrate part of the `e2e` tests for `PodChaos` and `NetworkChaos` from `Ginkgo`-style test code into `Gherkin` scenarios backed by `Go` step definitions. The first milestone was small but important: the first pull request for that migration was merged.

This first change focused on `PodChaos`, with `PodKill` as the starting point. The initial phase focused on bootstrapping `Gherkin` and `Godog`, targeting `PodKill` under `PodChaos`. The first pull request establishes the core test setup rather than migrating the entire suite. This proof of concept verifies that a `Gherkin` scenario can drive `Chaos Mesh` experiments through existing `Kubernetes` end-to-end helpers.

---

## Starting Small

This initial pull request defines the `Gherkin` feature and the `Godog` execution pathway. The primary objective is integrating the new test flow with the established framework. Feature files declare high-level behavior, while `Go` step definitions manage cluster setup, selectors, and custom resources by invoking existing `e2e` fixtures.

```text
go test ./e2e-gherkin/...
  |
  v
Godog runner
  |
  v
Gherkin feature file
  |
  v
Go step definitions
  |
  v
Kubernetes cluster + Chaos Mesh
```

The execution runs as a live end-to-end test. The scenario deploys test workloads, applies the chaos resource, and monitors the cluster state for expected failures. `Gherkin` abstraction improves scenario readability while verifying actual system behavior.

---

## What Changed

The change introduces the first `podchaos.feature` file, a `Godog` runner, and the `Go` step definitions. One scenario verifies that a targeted pod terminates when the resource is applied using label selectors. A second scenario validates experiment pausing by verifying that no further pod terminations occur after the chaos is paused.

Although the initial footprint is small, the flow executes a complete test cycle. This setup shows the value of `BDD` by separating readable behavior from the complexity of `Kubernetes` orchestration and resource polling. The testing path can now scale incrementally, adding scenarios without requiring a massive upfront migration.

---

## What I Learned

Maximizing fixture reuse is a key lesson from this initial implementation. The step definitions delegate workload setup and state assertions to existing packages like `pkg/fixture` and `e2e/util`. Leveraging these shared utilities prevents duplication and aligns the `BDD` suite with established testing practices.

Certain trade-offs exist in this initial release. Test setup is confined to the local `TestContext`, and the main `make e2e` target does not invoke the `Gherkin` suite. Addressing these gaps provides a clear objective for the subsequent pull request.

### Related Links

- [Issue - [LFX Mentorship] Refactor PodChaos and NetworkChaos E2E Tests into Gherkin-based BDD Scenarios](https://github.com/chaos-mesh/chaos-mesh/issues/4902)
- [PR - feat(e2e-test): implement gherkin BDD runner and podkill scenarios](https://github.com/chaos-mesh/chaos-mesh/pull/5013)

---

## Conclusion

The merged pull request establishes the foundation for the `Chaos Mesh` migration. Beyond the mentorship goals, this work prompted an investigation into using `Gherkin` `BDD` as AI guardrails within personal development projects. Defining strict behavioral scenarios provides a structured interface that constrains agentic code generation and validates system outputs. This milestone successfully transitions the testing pathway from a plan into functional, merged code.
