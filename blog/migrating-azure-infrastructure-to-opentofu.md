---
title: "Migrating Azure Infrastructure to OpenTofu"
description: "Migrating Cover Craft Azure infrastructure to OpenTofu brings git version control to cloud resources, prevents concurrent state drift, and automates CI deployment runs."
date: 2026-10-06
tags: ["terraform", "cloud", "automation"]
draft: true
---

## Configuration Drift and Portal Fragility

The Azure portal was great for initial prototyping. Pointing, clicking, and wiring up services directly in the browser got Cover Craft running without upfront boilerplate. But as the architecture grew, manual adjustments left no commit history behind.

I wanted infrastructure management to have the same versioning guarantees as application code. Moving to OpenTofu treats cloud resources as versioned files in git, where every change produces a clear diff, a pull request review, and a repeatable deployment.

---

## Modular Topology and Pipeline Boundaries

The cloud setup is split into four focused modules under `tofu/modules/`, with the root configuration wiring inputs and outputs together:

- `storage` provisions the primary storage account for deployment containers and asset storage.
- `application_insights` sets up the Log Analytics workspace and telemetry instance.
- `function_app` manages the serverless service plans, function apps, and deployment containers inside a dedicated resource group to avoid mixed-SKU conflicts.
- `app_service` provisions the Linux hosting plan and the Next.js web app, consuming the function key output from the API module.

A few operational invariants keep deployments predictable:

- Resource definitions live in git alongside application code.
- Remote state locks automatically in Azure Blob Storage during updates.
- Infrastructure runs in its own pipeline step before application code builds.
- Sensitive values like database strings and function keys stay in secrets, never hardcoded in configuration files.

```text
  [ git push ]
       |
       v
  [ CI Pipeline ]
       |
       +---> 1. [ infra-apply ]    ---> Provision or update Azure resources
       |
       +---> 2. [ deploy-api ]      ---> Deploy API function package
       |
       +---> 3. [ deploy-frontend ] ---> Deploy Next.js web application
```

---

## Staged Execution and State Locking

The deployment workflow in `.github/workflows/infra-deploy.yml` isolates platform changes into distinct pipeline phases. The `infra-apply` job runs first to apply OpenTofu updates. Both application jobs (`deploy-go-api` and `deploy-frontend`) declare `needs: infra-apply`, guaranteeing that code builds only proceed after the cloud infrastructure provisions cleanly.

Concurrency control relies on OpenTofu's `azurerm` backend configured against a remote `tfstate` container in Azure Blob Storage. OpenTofu acquires an active blob lease on `terraform.tfstate` during execution. If a local apply is run while a GitHub Actions job is in progress, the lease mechanism halts the second process immediately with a lock conflict, protecting the remote state from corruption.

Common resource tags (`project = "cover-craft"`) are passed through root locals to every module. Having consistent tags makes it straightforward to filter operational costs, review telemetry in Application Insights, and clean up test resources.

---

## Conclusion

Storing `tfstate` in Azure Blob Storage and automating provisioning through GitHub Actions gave me confidence in my deployment pipeline. The safety of this setup became clear when I merged two Dependabot version updates back to back. While the first workflow was applying changes, the second run failed immediately due to the active remote state lease lock.

Enforcing one operation at a time prevents race conditions and protects the platform from state drift. The pipeline fails fast on lock conflicts, guaranteeing that infrastructure changes apply cleanly and sequentially.

Find the source code in the [cover-craft](https://github.com/victoriacheng15/cover-craft) repository.
