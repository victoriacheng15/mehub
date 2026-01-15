---
title: "Exploring scp - Linux Commands"
description: "Learn how to use the SCP command in Linux for secure file transfers between local and remote systems. Master key options & practical examples! 🔐📁"
date: 2025-05-06
tags: ["linux"]
---

## **What is the SCP Command in Linux?**

The **SCP (Secure Copy Protocol)** command in Linux is a powerful and secure way to transfer files between your local machine and a remote server, or between two remote servers. It uses **SSH (Secure Shell)** for data transfer, ensuring that your files are encrypted during transit. SCP is especially useful for securely moving sensitive data, automating file transfers in scripts, or managing files on remote servers without additional tools.

Basic Syntax:

```bash
scp [OPTIONS] source_file destination_file
```

Example:

```bash
scp myfile.txt user@[ip-address]:/home/user/
```

This command copies `myfile.txt` from your local machine to the `/home/user/` directory on the remote server.

### Options

| **Option** | **Description** |
| --- | --- |
| `-P` | Specify a custom SSH port (default is 22). |
| `-r` | Copy directories recursively. |
| `-C` | Enable compression for faster transfers. |
| `-v` | Show verbose output for debugging. |
| `-i` | Use a specific SSH private key for authentication. |
| `-l` | Limit bandwidth usage (in Kbit/s). |

## Common Use Cases

### 1. Copying a File from Local to Remote

To transfer a file from your local machine to a remote server:

```bash
scp file.txt user@remote-server:/path/to/destination/
```

Example:

```bash
scp report.pdf john@192.168.1.100:/home/john/documents/
```

### 2. Copying a File from Remote to Local

To download a file from a remote server to your local machine:

```bash
scp user@remote-server:/path/to/file.txt /local/destination/
```

Example:

```bash
scp john@192.168.1.100:/home/john/backup.zip ~/downloads/
```

### 3. Copying a Directory Recursively

To transfer an entire directory (including subdirectories), use the `-r` flag:

```bash
scp -r /local/folder/ user@remote-server:/remote/path/
```

Example:

```bash
scp -r ~/project/ john@192.168.1.100:/home/john/backups/
```

### 4. Using a Custom SSH Port

If the remote server uses a non-default SSH port (e.g., 2222), specify it with `-P`:

```bash
scp -P 2222 file.txt user@remote-server:/destination/
```

### 5. Limiting Bandwidth Usage

To prevent SCP from consuming all available bandwidth, use `-l` (in Kbit/s):

```bash
scp -l 500 largefile.iso user@remote-server:/downloads/
```

This limits the transfer speed to 500 Kbit/s.

### 6. Enabling Compression for Faster Transfers

If you're transferring large files over a slow connection, use `-C` to compress data:

```bash
scp -C backup.tar.gz user@remote-server:/backups/
```

I decided to have a bit of fun and experiment with the -C option to see how much it actually helps. I used two files:

- test.csv (~1GB, uncompressed)
- test.zip (~0.5GB, already compressed)

Transfer Time Comparison

| File Name | Compression (`-C`) | Transfer Time (seconds) |
| --------- | ------------------ | ----------------------- |
| test.csv  | No                 |                    1:10 |
| test.csv  | Yes                |                    3:30 |
| test.zip  | No                 |                    1:11 |
| test.zip  | Yes                |                    1:46 |

Summary
From my test, enabling -C didn’t always make transfers faster:

- Transferring the uncompressed CSV without -C was the fastest.
- Using -C with the CSV actually slowed it down a lot, likely because of the extra work to compress a large file on the fly.
- For the already-compressed ZIP file, -C made a small improvement — but it was still slower than transferring either file without -C.

Note: These results can vary depending on your system and internet connection. This was just a fun test to understand how the -C option behaves in real situations.

### 7. Using SSH Key Authentication

If you authenticate via SSH keys, specify the private key with `-i`:

```bash
scp -i ~/.ssh/id_rsa file.txt user@remote-server:/home/user/
```

## Additional Help

To see all available SCP options, check the manual:

```bash
scp --help
# or
man scp
```

## Recap

The **SCP command** is an essential tool for securely transferring files in Linux. Whether you're moving data between local and remote machines or automating backups, SCP provides a fast, encrypted, and reliable method for file transfers. With its straightforward syntax and powerful options, mastering SCP will make managing remote files much easier.

**Happy copying!** 🚀

## Thank you

Thank you for your time and for reading this!
