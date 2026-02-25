---
title: "Single-Stage or Multi-Stage Docker Builds"
description: "Explore single-stage vs multi-stage Docker builds with a small Flask API, learning differences in image size, build speed, and container optimization. Read more to learn."
date: 2025-11-11
tags: ["platform"]
---

## Understanding Docker Builds

## What is Single-Stage or Multi-Stage Docker Builds?

In this post, we will explore **Single-Stage or Multi-Stage Docker Builds**. While containerizing **the School Management Flask API**, I explored **single-stage and multi-stage Docker builds** to understand their impact on performance, image size, and deployment efficiency. This exploration taught me practical lessons about optimizing images, improving CI/CD pipelines, and separating build-time dependencies from runtime environments.

---

### Single-Stage Builds

A **single-stage build** uses one image for both building and running the application. All dependencies, build tools, and source code are included in the final image.

**Example:**

```docker
FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY . .

RUN chmod +x /app/entrypoint.sh

EXPOSE 5000

ENTRYPOINT ["/app/entrypoint.sh"]
```

**Key Takeaways:**

- Simple and straightforward to set up
- Includes all build and runtime dependencies in one image
- Final image may be larger and include unnecessary build artifacts

**My Experience:**  

- Measured build time: ~1.67s  
- Image contained build tools and cached files, resulting in a larger final image  

---

### Multi-Stage Builds

A **multi-stage build** separates build-time dependencies from the runtime environment. You compile or build your app in a full-featured image, then copy only the necessary artifacts into a smaller, cleaner runtime image.

**Example:**

```docker
# Stage 1 - builder
FROM python:3.12-slim AS builder

WORKDIR /app

COPY requirements.txt /app/

RUN pip install --no-cache-dir -r requirements.txt


# Stage 2 - final image
FROM python:3.12-slim

WORKDIR /app

COPY --from=builder /usr/local/lib/python3.12/site-packages /usr/local/lib/python3.12/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin

COPY . .

RUN chmod +x /app/entrypoint.sh

EXPOSE 5000

ENTRYPOINT ["/app/entrypoint.sh"]
```

**Key Takeaways:**

- Separates build-time dependencies from the runtime image
- Produces a smaller, cleaner final image
- Improves maintainability and keeps runtime focused

**My Experience:**  

- Measured build time: ~1.08s (faster on my system)  
- Multi-stage image excluded build tools and cached files, resulting in a cleaner, focused image  
- Slightly more complex Dockerfile, but worth it for maintainability  

---

### Why It Matters

For small projects like the School Management API, build speed and image size differences may be minimal, and results depend on your system (CPU, RAM, disk speed). However, understanding these approaches still provides benefits for maintainability, security, and best practices.

- Reduces unsed files in production images
- Optimizes CI/CD pipelines and deployment efficiency  
- Encourages better project structure by separating build and runtime concerns  
  
---

## Key Learning

Exploring the differences between single-stage and multi-stage builds was a great opportunity to see how both approaches behave. Even for a small project like the School Management API, multi-stage builds produced smaller, cleaner images, while single-stage builds were simpler to set up. This exercise helped me better understand build optimization and containerization concepts for personal learning.

---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
