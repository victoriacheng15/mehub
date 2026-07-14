---
title: "Building a GitOps Homelab Platform Pipeline"
description: "Manual pod updates became a GitOps platform pipeline where CI builds images, Kustomize separates dev and prod, ArgoCD reconciles drift, and short-SHA tags trace rollbacks."
date: 2026-05-12
tags: ["docker", "kubernetes", "platform"]
---
## From Sideloading to a Platform Pipeline

Deploying software to a homelab often starts with simple, one-off commands. While using Makefile commands is manageable, manual sideloading and redeploying pods remain prone to error:

- ⚓ Forgetting to rebuild an image before deployment.
- 🚀 Failing to redeploy a pod after a local build.
- 🔐 Facing state inconsistencies and configuration drift.
- ⚡ Getting lost in debugging rabbit holes before identifying version mismatches.

Implementing a GitOps model using ArgoCD fundamentally changed the deployment perspective. The homelab became less of a manually managed Kubernetes box and more of a DevOps pipeline with platform engineering boundaries: Git defines intent, CI produces artifacts, manifests promote versions, ArgoCD reconciles runtime state, and observability closes the loop. The most significant discovery is realizing that the cluster state is defined by manifests in Git, not by what is currently running in Kubernetes. If a change isn't in the `k3s/` directory, it simply doesn't exist in the runtime.

---

## Discovering the Promotion Gate

A structured Promotion Gate was implemented by separating CI from CD to ensure architectural clarity. This is where the homelab started to resemble a real platform workflow instead of a collection of scripts: developers change source, CI builds immutable artifacts, and deployment only happens when the desired state is promoted through Git. Prior experiences with 'latest' tags in academic projects highlighted the inherent traceability risks:

- **The 'Latest' Pitfall**: While convenient, using 'latest' tags makes identifying the active image version nearly impossible during troubleshooting.
- **CI (GitHub Actions)**: Handles the build logic and pushes Docker images to the GitHub Container Registry (GHCR).
- **CD (ArgoCD)**: Monitors manifests and only triggers a deployment when the YAML version tag is updated.
- **Immutability**: Moving away from generic tags to specific Git SHAs (e.g., sha-7d8708b) ensures every pod runs a traceable commit.

This approach eliminates configuration drift and makes rollbacks as simple as reverting a YAML change in Git. The clarity provided by seeing the cluster state match the Git history is a transformative experience because it creates a small but complete platform contract: build, promote, reconcile, observe, and recover.

---

## Using Secrets and Sync

Transitioning to private images in GHCR required a layer of security through 'imagePullSecrets'. Learning to scope these secrets correctly was a critical lesson in infrastructure safety:

- **Authentication**: Established a secure bridge between the local K3s node and GitHub via namespace-scoped credentials.
- **Event-Driven Sync**: Leveraged a custom webhook and sync script to eliminate standard ArgoCD polling delays.
- **Latency Reduction**: Achieved near-instant reconciliation, reducing deployment feedback from minutes to seconds.
- **Matrix Builds**: Used GitHub Actions' matrix strategy to manage a growing fleet of simulation services from a single blueprint.

These steps ensure that workloads only have access to the credentials required for their specific namespace while maintaining high-speed feedback loops. That balance is what moves the setup closer to platform engineering: the pipeline is fast for developers, but the runtime remains bounded by declarative state, scoped credentials, and observable reconciliation.

---

## The SHA Pivot Discovery

Encountering an 'ErrImagePull' during the deployment of a chaos controller and sensor fleet led to a vital "aha!" moment regarding tag formats:

- **Initial Failure**: Manifests referenced full-length Git SHAs while CI published short-SHA image tags, causing image pull failures during reconciliation.
- **The Pivot**: Transitioning the CI workflow to a seven-character short SHA format.
- **Immediate Success**: Updating the manifest with the short hash triggered ArgoCD to automatically synchronize the cluster state.

This "wrong turn" provided empirical proof that the automation loop, from Registry to Secret to Manifest, was working exactly as intended.

---

## Conclusion

Running the homelab as a GitOps pipeline helped me understand the platform loop by working through it directly:

- **Git as the source of truth**: Kubernetes state is managed from manifests, not manual pod changes.
- **ArgoCD reconciliation**: Pod updates happen through an event-driven reconciliation flow when the desired state changes.
- **Dev and prod separation**: Kustomize makes it possible to manage dev and prod versions without duplicating the whole deployment.
- **Canary validation**: Dev can use a canary tag to check whether the new image and manifest work correctly.
- **Versioned promotion**: Prod can move to a short-SHA image tag after validation, making the deployed version traceable.

At this point, the homelab stopped being a set of manual deployment steps and became a working CI/CD pipeline: CI produced the image, Git stored the desired state, ArgoCD reconciled Kubernetes, and dev-to-prod promotion became explicit through Kustomize and short-SHA tags.

This hands-on process made GitOps feel less like a tool choice and more like a platform engineering workflow: build, validate, promote, reconcile, and observe.
