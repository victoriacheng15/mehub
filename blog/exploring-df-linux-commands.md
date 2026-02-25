---
title: "Exploring df - Linux Commands"
description: "Learn how to use the df command in Linux to check disk space usage, view file system types, and monitor mounted volumes with practical examples and flags. Read more to learn."
date: 2025-07-22
tags: ["linux"]
---

## What is the df Command in Linux?

The **df** command in Linux is used to display information about disk space usage on file systems. It shows how much space is used, available, and where it's mounted. It's a quick and essential tool for monitoring storage on local drives, external disks, and mounted network volumes.

**Key features include:**

- Shows total, used, and available disk space
- Lists mount points for all file systems
- Supports human-readable format for easy reading
- Can filter output by file system type or specific paths

### Basic Syntax

```bash
df [OPTIONS] [FILESYSTEM or PATH]
```

### Example

```bash
df
```

This will display disk usage for all mounted file systems.

---

## Options

| Option        | What It Does                                  | Example                                     |
|---------------|-----------------------------------------------|---------------------------------------------|
| `-h`          | Human-readable sizes (KB, MB, GB)             | `df -h`                                     |
| `-T`          | Show file system type                         | `df -T`                                     |
| `-a`          | Include pseudo and hidden file systems        | `df -a`                                     |
| `-t TYPE`     | Show only file systems of a certain type      | `df -t ext4`                                |
| `-x TYPE`     | Exclude file systems of a certain type        | `df -x tmpfs`                               |

---

## Common Use Cases

- View disk usage in human-readable format

```bash
df -h
```

- Check disk usage of a specific directory

```bash
df -h /home
```

- Show file system types along with disk usage

```bash
df -T
```

- Exclude temporary file systems like `tmpfs`

```bash
df -h -x tmpfs
```

- Display all file systems including dummy mounts

```bash
df -a
```

- Filter to only show ext4 file systems

```bash
df -t ext4
```

---

## Additional Help

To see all available options:

```bash
df --help
# or
man df
```

---

## Recap

The **df** command is a simple yet powerful way to monitor disk space on your Linux system. Whether you're tracking usage on internal drives or mounted volumes, `df` gives you the insight you need with just a few keystrokes.

Keep your storage in check! 🗂️💾

---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
