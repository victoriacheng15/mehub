---
title: "Scaling with Systemd Template Units"
description: "Learn to scale Linux automation using Systemd Template Units. Discover how the '@' symbol and gitops-sync timers simplify service management. Read the full guide to learn."
date: 2026-02-10
tags: ["retrospective", "linux"]
---

## The Discovery: Deciphering the '@' Symbol

While setting up the GitOps sync for the **Observability Hub**, I leveraged AI to generate a systemd service file. The output it provided included a curious symbol: `@`. The file was named `gitops-sync@.service`.

Initially, I wasn't sure what the `@` signified. After some research, I learned that this is a **Template Unit**. It allows you to pass dynamic arguments—like a folder name—directly into the service.

In Phase 1 of my GitOps setup, I decided to test this with just one repository to understand the mechanism. I quickly realized how powerful this is: instead of creating a new service file for every new repository or folder I want to sync, I can use a single template and scale horizontally just by changing the instance name.

---

## The Blueprint: How Template Units Work

Systemd template units, denoted by the `@` symbol, function similarly to a class in programming. This applies to both **Services** (the logic) and **Timers** (the schedule).

* **Template Service (`gitops-sync@.service`):** The Worker. It defines *what* to do (run the script). It accepts the instance name (e.g., `frontend`) as an argument `%i` to execute logic specific to that context.
* **Template Timer (`gitops-sync@.timer`):** The Scheduler. It defines *when* to execute. By default, a timer looks for a service with the same name.
* **Instance Relationship:** When `gitops-sync@frontend.timer` fires, systemd automatically **activates** `gitops-sync@frontend.service`. The instance name (`frontend`) is implicitly shared because the timer starts the service of the same instance name.

This discovery transformed my approach from "making it work" for one repo to architecting a system that can handle dozens of sync targets with zero additional configuration files.

---

## Implementation: Scaling with Timers

While the service defines *how* to sync, the timer defines *when*. The real power of this pattern emerged when I coupled the template service with a template **timer**.

First, the **Template Service** (`/etc/systemd/system/gitops-sync@.service`) defines the execution logic:

```ini
[Unit]
Description=GitOps Sync for %i
After=network.target

[Service]
Type=oneshot
# The %i specifier is replaced by the instance name at runtime
ExecStart=path-to-your-folder/gitops_sync.sh %i
User=server
```

Then, the **Template Timer** (`/etc/systemd/system/gitops-sync@.timer`) drives the schedule:

```ini
[Unit]
Description=Trigger GitOps Sync for %i
After=network.target

[Timer]
OnBootSec=5min
OnUnitActiveSec=15min
Unit=gitops-sync@%i.service

[Install]
WantedBy=timers.target
```

Notice the `Unit=gitops-sync@%i.service` line in the timer. This explicitly binds the timer instance (e.g., `gitops-sync@frontend.timer`) to the corresponding service instance (`gitops-sync@frontend.service`).

This allows me to treat the timer as the primary control interface. I don't just "run" the sync; I "schedule" it per folder.

---

## Scaling the Fleet: Management and Operations

The real value of this approach isn't just saving a few lines of config—it's streamlined operations and granular control.

### Management

To add a new folder to the sync schedule, I simply enable the timer for that specific instance:

```bash
# Sync the 'frontend' folder
systemctl enable --now gitops-sync@frontend.timer

# Sync the 'backend' folder
systemctl enable --now gitops-sync@backend.timer
```

### Operations & Debugging

Because they are distinct instances, their logs are isolated. I don't have to grep through a giant shared log file to find errors specific to the backend.

To see the full logs for the `frontend` service:

```bash
journalctl -u gitops-sync@frontend.service
```

To see just the recent activity (e.g., last 20 lines) for debugging:

```bash
journalctl -u gitops-sync@frontend.service -n 20
```

This ensures that if the `frontend` sync fails, it doesn't pollute the status or logs of the `backend` sync. We have achieved **isolated failure domains** with zero extra infrastructure overhead.

To see *all* sync activity across the platform (using patterns):

```bash
journalctl -u "gitops-sync@*" -n 20
```

---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
