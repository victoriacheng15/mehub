---
title: "Exploring wget - Linux Commands"
description: "Linux wget command guide: download files, resume transfers, limit speed & automate tasks with practical examples for efficient downloading."
date: 2025-06-03
tags: [linux]
---

## What is the `wget` Command in Linux?

The `wget` command in Linux is a powerful, non-interactive tool for downloading files from the internet. It supports HTTP, HTTPS, and FTP protocols, and can work in the background, making it ideal for scripting and automated downloads. Key features include:

- Resuming interrupted downloads
- Downloading recursively (entire websites)
- Working with proxies
- Limiting bandwidth usage
- Running without user interaction

Basic Syntax:

```bash
wget [OPTIONS] URL
```

Example:

```bash
wget https://example.com/file.zip
```

## Options

| Option | What It Does | Example |
| --- | --- | --- |
| **`-O FILE`** | Rename downloaded file | **`wget -O backup.zip [URL]`** |
| **`-P DIR`** | Save to directory | **`wget -P ~/downloads [URL]`** |
| **`-c`** | Resume broken download | **`wget -c [URL]`** |
| **`-nc`** | Skip existing files | **`wget -nc [URL]`** |
| **`--limit-rate=SPEED`** | Limit speed | **`wget --limit-rate=2M [URL]`** |
| **`-q`** | Silent mode | **`wget -q [URL]`** |
| **`-b`** | Download in background | **`wget -b [URL]`** |
| **`-i FILE`** | Download URLs from file | **`wget -i urls.txt`** |
| **`-t N`** | Retry attempts (default=20) | **`wget -t 3 [URL]`** |
| **`-T SEC`** | Timeout in seconds | **`wget -T 30 [URL]`** |

## Common Use Case

- Download and rename

```bash
wget -O gnu-keyring.gpg https://ftp.gnu.org/gnu/gnu-keyring.gpg
```

- Save to specific directory

```bash
wget -P ~/Downloads https://www.gutenberg.org/files/1342/1342-0.txt
```

- Resume failed download

```bash
wget -c https://cdn.kernel.org/pub/linux/kernel/v6.x/linux-6.5.11.tar.xz
```

- Download without overwriting

```bash
wget -nc https://www.gutenberg.org/files/11/11-0.txt  # Alice in Wonderland
wget -nc https://www.gutenberg.org/files/11/11-0.txt  # Skips if already downloaded
```

- Limit download speed

```bash
wget --limit-rate=100k https://speed.hetzner.de/100MB.bin
```

- Silent download (for scripts)

```bash
wget -q https://www.gutenberg.org/files/84/84-0.txt  # Frankenstein
```

- Background download

```bash
wget -b https://www.gutenberg.org/files/2701/2701-0.txt  # Moby Dick
tail -f wget-log  # Check progress
```

### Additional Help

To see all available wget options:

```bash
wget --help
# or
man wget
```

## Recap

The `wget` command is an indispensable tool for downloading files from the command line. Whether you need to download single files, mirror entire websites, or automate downloads in scripts, `wget` provides a robust solution with many configuration options.

Happy downloading! 💾🚀

## Thank you!

Thank you for your time and for reading this!