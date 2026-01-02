---
title: "Trying Azure Container Registry"
description: "Learn how I used Azure Container Registry to securely store Docker images from a personal project, with GitHub Actions CI and zero hardcoded secrets."
date: 2025-12-16
tags: ["azure", "cloud"]
---

## What is Azure Container Registry?

Azure Container Registry (ACR) is a managed, private Docker registry service in Azure. It lets you store, build, and manage container images and artifacts—securely and at scale—without relying on public registries like Docker Hub.

Think of it like a personal warehouse for your Docker images: only you (or your team) can access it, and it lives right inside your Azure environment.

### I First Encountered ACR in Class

I first encountered Azure Container Registry in a recent academic project. At the time, it was a requirement: “Use ACR to store your container images.” I followed the steps, got it working, and moved on.

But after the assignment ended, I kept wondering:  
> *What’s it actually like to use ACR for something I care about?*

So I decided to try it with a personal project—my `school-management-api`—just to see how the whole flow feels outside of a graded task.

## Why ACR?

I’ve used Docker locally for a while, but I’ve never pushed my images to a remote registry—I just rebuilt them each time I needed to run the container. That works fine for local testing… but I started wondering:

> How do real cloud services like Azure Web App actually pull container images?

That curiosity led me to ACR. I wanted to see firsthand how a private registry in Azure works—and how an Azure service might authenticate and pull an image from it.

ACR offers a secure, managed place to store images inside your Azure subscription—no public exposure, no manual cleanup. And even if I’m not deploying anything yet, understanding this flow feels like a key piece of the cloud puzzle.

## My GitHub Actions Workflow

To automate things, I set up a simple CI pipeline in GitHub Actions:

- On every push to `main`:
  - Log in to Azure using a service principal (stored securely in the `AZURE_CREDENTIALS` GitHub secret)
  - Log in to my ACR instance (`az acr login`)
  - Build the image using a multi-stage `Dockerfile`
  - Tag it with the commit SHA **and** `latest`
  - Push both tags to ACR

```yaml
env:
  CONTAINER_NAME: school-management-api

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v5.0.0

      - name: Log in to Azure
        uses: azure/login@v2.3.0
        with:
          creds: ${{ secrets.AZURE_CREDENTIALS }}

      - name: Log in to Azure Container Registry
        run: |
          az acr login --name schoolAPI

      - name: Build and tag Docker image
        run: |
          docker build -f Dockerfile.multi-stage -t ${{ secrets.ACR_NAME }}.azurecr.io/$CONTAINER_NAME:${{ github.sha }} .
          docker tag ${{ secrets.ACR_NAME }}.azurecr.io/$CONTAINER_NAME:${{ github.sha }} ${{ secrets.ACR_NAME }}.azurecr.io/$CONTAINER_NAME:latest

      - name: Push Docker image to ACR
        run: |
          docker push ${{ secrets.ACR_NAME }}.azurecr.io/$CONTAINER_NAME:${{ github.sha }}
          docker push ${{ secrets.ACR_NAME }}.azurecr.io/$CONTAINER_NAME:latest
```

Notice the use of GitHub Actions secrets like secrets.ACR_NAME and secrets.AZURE_CREDENTIALS—this keeps sensitive info out of the workflow file itself. No hardcoded registry names, no leaked credentials. Just secure, repeatable builds.

This gives me traceable builds (via SHA) and a convenient latest tag for quick local pulls or testing.

Check out this [repo](https://github.com/Azure/login) on GitHub for how to get creds for Azure login.

---

## A Side Discovery: GHCR

While reading about container registries, I also stumbled upon GitHub Container Registry (GHCR). It’s neat—especially if you’re already using GitHub—but I didn’t switch. My goal wasn’t to find the “best” registry, but to understand ACR better on my own terms.

---

## Final Thoughts

Using ACR in class felt like checking a box.
Using it for my own project? That’s when it clicked.

It’s not magic—but once you have the right permissions and a few az commands, the whole process (build → push → pull) becomes reliable, repeatable, and secure by default.

And honestly, that’s a great feeling—knowing my containers live somewhere safe, versioned, and under my control.

> 🐳 Curiosity > compliance. Sometimes the best learning starts after the assignment ends.

 ---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
