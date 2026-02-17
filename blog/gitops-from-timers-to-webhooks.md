---
title: "GitOps: From Timers to Webhooks"
description: "From systemd timers to event-driven webhooks: How I leveraged a Go proxy to eliminate O(N) management toil and achieve instant GitOps syncs for multiple repositories."
date: 2026-02-17
tags: ["go", "automation"]
---

## The Systemd Approach

It started with an idea to leverage systemd's instantiated services (the `@` syntax) to scale my GitOps sync timers. The setup was simple: a systemd timer would trigger every 15 minutes to run a custom git-sync bash script that checked if a repository's local branch needed an update.

I started with two repositories to test the workflow. It worked flawlessly; the polling was reliable, and the overhead was negligible. Emboldened by the initial results, I scaled it up to five repositories.

---

## Hitting the O(N) Wall

That’s when the reality of maintenance set in.

What felt like automation was actually a hidden operational tax. Every time I revised a service configuration or performed system maintenance, I had to manually synchronize the state across my entire "fleet" of timers:

```bash
sudo systemctl daemon-reload
sudo systemctl restart gitops-sync@repo1.timer
sudo systemctl restart gitops-sync@repo2.timer
sudo systemctl restart gitops-sync@repo3.timer
# ... repeat for every single repository
```

This was O(N) management in a world that needed O(1) efficiency. I briefly considered hiding the mess behind a `Makefile`, but that would just be trading terminal commands for the manual chore of maintaining a hardcoded list. I was still managing infrastructure by hand; I just had a better script for it.

I realized I was optimizing the wrong thing. The problem wasn't *how* I was managing the timers; the problem was that I was *polling* in the first place.

---

## "Don't Call Us, We'll Call You"

I needed to invert the control. To do this, I leveraged the Go proxy already running in my [**Observability Hub**](https://github.com/victoriacheng15/observability-hub/tree/main/proxy) repository and integrated it with my existing GitOps sync bash script to ensure my local repositories always stay updated to the latest version.

### The Network Gap

Of course, my server sits safely behind a firewall, invisible to the public internet. To bridge this gap without opening dangerous ports on my router, I used **Tailscale Funnel**. It acts as a secure, public-facing ingress that accepts the webhook from GitHub and tunnels it directly to my local Go application.

It felt a bit like the movie *"Are We There Yet?"* where I was constantly asking "Is there an update?" every 15 minutes. So, I decided to stop asking and let the system tell me instead. I added a secure `api/webhook/gitops` route to the proxy that validates the GitHub signature (HMAC-SHA256) and listens for specific events.

The logic supports both direct pushes and PR merges to `main`:

```go
// ... logic to verify signature ...

// Case 1: Push to main
if eventType == "push" && payload.Ref == "refs/heads/main" {
    shouldTrigger = true
}

// Case 2: PR merged to main
if eventType == "pull_request" && payload.Action == "closed" && payload.PullRequest.Merged {
    shouldTrigger = true
}

if shouldTrigger {
    // Fire and forget: Run sync asynchronously
    go func(repo string) {
        cmd := exec.Command(".../gitops_sync.sh", repo)
        // ... handle output ...
    }(repoName)
}
```

---

## The Payoff

Whew.

The shift changed everything. I no longer have to manage individual systemd services for every new repository I spin up. Honestly, just imagining the alternative (manually restarting twenty different timers every time I tweaked a config) makes me want to scream.

While I still maintain a centralized security allowlist in my sync script, the heavy lifting is gone. Adding a new repository to the pipeline is now as simple as adding a webhook in GitHub and a single entry in the list.

The system is reactive, instant, and scales without the drama. I went from maintaining state to handling events, and that made all the difference. It was a fun way to figure out a better way to manage this, and the end result is much more satisfying.

---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
