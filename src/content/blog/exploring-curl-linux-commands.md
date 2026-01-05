---
title: "Exploring curl - Linux Commands"
description: "Learn how to use the curl command in Linux to download files, interact with APIs, send POST data, handle redirects, and more with simple examples."
date: 2025-07-15
tags: ["platform"]
---

## What is the curl Command in Linux?

The **curl** command in Linux is a versatile tool used to transfer data to or from a server, using a variety of protocols such as HTTP, HTTPS, FTP, and more. Unlike wget, which is mainly focused on downloading files, curl is designed to work with data transfers in both directions and is often used for testing APIs, uploading data, or interacting with web services.

**Key features include:**

- Supports multiple protocols (HTTP, HTTPS, FTP, SMTP, SCP, SFTP, and more)
- Can send custom HTTP headers and data (useful for APIs)
- Supports HTTP methods like GET, POST, PUT, DELETE
- Handles authentication (Basic, Digest, OAuth, etc.)
- Supports cookies and session management
- Ability to upload files or data
- Shows detailed request/response headers for debugging
- Works well in scripts and automation pipelines

### Basic Syntax

```bash
curl [OPTIONS] URL
```

### Example

```bash
curl https://example.com
```

This will fetch the content of the webpage and output it to the terminal.

## Options

| Option         | What It Does                                  | Example                                      |
|----------------|-----------------------------------------------|----------------------------------------------|
| `-o FILE`      | Save output to a file                         | `curl -o backup.zip [URL]`                   |
| `-O`           | Save with original filename                   | `curl -O [URL]`                              |
| `-L`           | Follow redirects                              | `curl -L [URL]`                              |
| `-d DATA`      | Send POST data                                | `curl -d "a=1&b=2" -X POST [URL]`            |
| `-H HEADER`    | Add custom header                             | `curl -H "Content-Type: application/json"`   |
| `-X METHOD`    | Set request method                            | `curl -X POST [URL]`                         |
| `-u USER:PASS` | Use basic auth                                | `curl -u user:pass [URL]`                    |
| `-s`           | Silent mode (no output except result)         | `curl -s [URL]`                              |
| `-I`           | Show only headers                             | `curl -I [URL]`                              |

## Common Use Cases

- Download and save a file with a custom name

```bash
curl -o backup.zip https://example.com/file.zip
```

- Download and save with the original remote file name

```bash
curl -O https://example.com/file.zip
```

- Follow redirects (e.g., from shortened URLs)

```bash
curl -L http://bit.ly/some-short-url
```

- Send POST form data to an API

```bash
curl -d "username=username&password=password" -X POST https://api.example.com/login
```

- Send JSON data with a custom header

```bash
curl -H "Content-Type: application/json" -d '{"name":"Victoria"}' -X POST https://api.example.com/users
```

- Use basic authentication

```bash
curl -u user:pass https://api.example.com/secure-data
```

- Download a file silently (no progress output)

```bash
curl -s -O https://example.com/largefile.iso
```

- Fetch only HTTP headers (e.g., for debugging)

```bash
curl -I https://example.com
```

## Additional Help

To see all available options:

```bash
curl --help
# or
man curl
```

## Recap

The **curl** command is a powerful, flexible tool for transferring data over a variety of protocols. Whether you need to download files, interact with APIs, upload data, or troubleshoot HTTP requests, curl provides the control and options to handle these tasks efficiently.

By mastering curl, you’ll improve your ability to automate web interactions and debug network services from the Linux command line.

Happy curling! 🌐🚀

##

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
