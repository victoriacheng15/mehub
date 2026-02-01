---
title: "[Engineering Log] Cron vs. Systemd: A Deep Dive"
description: "A deep dive into Systemd vs Cron: Automating an observability hub with self-healing timers, centralized logging, and dependency management."
date: 2026-02-03
tags: ["retrospective", "linux", "automation"]
draft: true
---

## 1. Context (The "Why")

`systemd` is the standard service manager in modern Linux systems, but I had always relied on `cron` for scheduling tasks. I wanted to understand if `systemd` timers offered enough value to justify the steeper learning curve.

> **The Friction:** `cron` is simple, but it is "fire and forget." It has no concept of whether the job succeeded, failed, or if the network was even online when it tried to run.

---

## 2. The Challenge / Question

I needed to know: **"Is Systemd just extra complexity, or does it solve real problems that Cron can't?"**

- **The Gap:** `cron` requires manual log redirection (`>> /var/log/my.log 2>&1`) and custom retry logic.
- **The Goal:** A scheduling system that treats background scripts as first-class citizens with built-in observability.

---

## 3. Investigation & Trade-offs

I implemented a "Practice Service" to compare the two approaches.

| Feature | Cron | Systemd Timers | Verdict |
| :--- | :--- | :--- | :--- |
| **Scheduling** | Simple syntax (`* * * * *`). | Verbose Unit files (`.timer`). | 🤝 Cron wins on simplicity |
| **Logging** | Manual redirection required. | Automatic (`journalctl`). | 🏆 Systemd wins |
| **Dependencies** | None (runs blindly). | Can wait for Network/DB. | 🏆 Systemd wins |
| **Recovery** | None (fails silently). | Auto-restart (`Restart=on-failure`). | 🏆 Systemd wins |

### Discovery: The "Kickstart" Gotcha

I discovered that `OnUnitActiveSec` (run X minutes after *last run*) differs from Cron's wall-clock scheduling.

- **The Trap:** If the service has *never* run, the timer won't start.
- **The Fix:** You must manually `systemctl start` the service once to "kickstart" the cycle.

---

## 4. The Solution / Insight

The power of Systemd became clear when I used **Template Units** (`@.service`). This allows running multiple instances of the same logic with different configurations.

```ini
# practice@.service
[Unit]
Description=Practice Service for %i

[Service]
Type=oneshot
ExecStart=/path/to/scripts/practice.sh %i
User=server
```

### The Implementation Pattern

To move from "discovery" to "operational," I followed this lifecycle:

- **Deploy & Reload:** Copy unit files to `/etc/systemd/system/` and run reload:

  ```bash
  sudo systemctl daemon-reload
  ```

- **Enable Timer:** Start the timer on boot and run it now:

  ```bash
  sudo systemctl enable --now practice@alpha.timer
  ```

- **Kickstart:** Required to trigger the first `OnUnitActiveSec` cycle:

  ```bash
  sudo systemctl start practice@alpha.service
  ```

- **Monitor:** Check live logs and verify the schedule:

  ```bash
  journalctl -u practice@alpha -f
  systemctl list-timers
  ```

---

## 5. Putting it into Practice: Observability Hub Sync

To move beyond "practice" scripts, I applied this architecture to a real problem: **automating reading analytics ingestion**.

My requirements were strict:

1. **Catch-up Logic:** If my laptop is asleep at 10:00 AM, the job must run immediately upon wake.
2. **Network Awareness:** The job fails if the local API isn't reachable, so dependency management is key.

### The Configuration

I used `Persistent=true` to handle the "laptop problem" (missed runs) and `RandomizedDelaySec` to prevent resource contention on startup.

#### Service Unit (`reading-sync.service`)

```ini
[Unit]
Description=Trigger Reading Analytics Sync
Wants=reading-sync.timer

[Service]
Type=oneshot
User=server
ExecStart=/usr/bin/curl -X POST https://api.example.com/sync
# Standardize logging for journald
StandardOutput=journal
StandardError=journal
```

#### Timer Unit (`reading-sync.timer`)

```ini
[Unit]
Description=Daily Sync for Reading Analytics

[Timer]
# Run every day at 10:00 AM
OnCalendar=*-*-* 10:00:00
# Ensure catch-up if the machine was off
Persistent=true
# Prevent thundering herd (30m window)
RandomizedDelaySec=1800

[Install]
WantedBy=timers.target
```

---

## 6. Outcome & Learnings

- **Result:** I now have a scheduling system that is **self-healing** (auto-restart) and **observable** (centralized logs).
- **Lesson:** "For simple tasks, use Cron. For critical services that need state awareness, use Systemd."

---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
