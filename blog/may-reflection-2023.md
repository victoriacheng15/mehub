---
title: May Reflection 2023
description: "Reflected on what I have learned for the May, what is docker and its usage, Harvard CS50, Back of the Envelope Estimation from System Design. Read the full guide to learn."
date: 2023-05-22
tags: ["monthly-log"]
---

## Monthly Reflection

In this post, we will explore **May Reflection 2023**. This month was a bit slow for me, but I managed to brainstorm at least 5+ blog post ideas. I suppose this is not a bad thing. As always, celebrate this win!

Achieved List for the month:

- Published 2 blogs, `What is Nullish Coalescing Operator` and `Exploring JavaScript Prototype Inheritance`
- Read `What is a Url` blog post
- Read Chapter 2 of System Design - `Back of the Envelope Estimation`
- Learned about `Docker`
- Harvard CS50 - the beginning to week 2

---

## Blogs

### Nullish Coalescing Operator

It is an operator that checks for null or undefined values only. Keep in mind about a pitiful, if the value contains true or false in a boolean context, the nullish coalescing operator would not work on this. In this case, the logical OR operator or ternary operator should be used.

```js
const value = null;

const check = value ?? "2nd value";
console.log(check); // 2nd value
// try assign undefined to value, it will return 2nd value again
```

```js
const value = "1st value";

const check = value ?? "2nd value";
console.log(check); // 1st value
```

- [Blog post: What is Nullish Coalescing Operator](https://victoriacheng15.vercel.app/posts/what-is-nullish-coalescing-operator)

### Exploring JavaScript Prototype Inheritance

I deeply dive into the topic of inheritance for object and function and learned the differences between them. You could use either the `new` keyword (for function) or `Object.create()` (for object and function).

- [Blog post: Exploring JavaScript Prototype Inheritance](https://victoriacheng15.vercel.app/posts/exploring-javascript-prototype-inheritance)

---

## What is a URL?

The blog post by the creator of `curl` highlights a funny URL and explains how different systems handle it. However, it doesn't delve into the impact of varying system interpretations of URLs. A talk by Orange Tsai uncovers inconsistencies between different libraries and the resulting security risks. Understanding the parts of a URL, such as a scheme, username, password, host, port, path, query, and fragment, is essential but challenging due to differing interpretations.

Inconsistencies can lead to parsing errors, as exemplified by different behaviors in query/username and port/path parsing. Handling URLs securely, especially when dealing with user input, requires careful consideration and architectural safeguards.

- [What Is a URL by Azeem Bande-Ali](https://azeemba.com/posts/what-is-a-url.html)

---

## System Design

### Load balancing

I discovered this article on load balancing, which artfully employs visuals to illustrate a gradual progression from the simplest concepts to a comprehensive exploration of various load balancing algorithms, including round robin, weighted round robin, and many more.

- [LinkedIn Post - Load Balancing](https://www.linkedin.com/posts/victoriacheng15_coding-programming-softwareengineering-activity-7055203132854702080-RUZD/?utm_source=share&utm_medium=member_desktop)
- [Load Balancing By Sam Rose](https://samwho.dev/load-balancing/)

### Back of the Envelope Estimation

This chapter introduces the concept of back-of-the-envelope (BOTE) estimation, a quick and rough technique used in system design interviews to estimate metrics like storage, traffic, and processing power. It can be likened to forecasting sales or predicting weather conditions.

For example, estimating the storage size needed for daily user posts on a platform like LinkedIn involves considering daily user activity and extrapolating it to future years. Accurate estimations can help companies scale their databases effectively, plan upgrades or adjustments to cloud services, and obtain quotes from different cloud providers. The chapter prompted the reader to engage in calculations and considerations regarding the amount of storage space required in future years.

- [LinkedIn Post - Back of the Envelope Estimation](https://www.linkedin.com/posts/victoriacheng15_coding-programming-softwareengineering-activity-7059558550699716608-NvqI?utm_source=share&utm_medium=member_desktop)

---

## Docker

What is container and container image? It is a lightweight, standalone, executable package of software that includes everything needed to run an application. This package contains underlying OS dependencies, runtime dependencies (e.g. Python runtime), libraries (e.g. SQL Alchemy) and application code.

- [Docker Crash Course for Absolute Beginners [NEW] By TechWorld with Nana](https://www.youtube.com/watch?v=pg19Z8LL06wp)
- [Complete Docker Course - From BEGINNER to PRO! (Learn Containers) By DevOps Directive](https://www.youtube.com/watch?v=RqTEHSBrYFw)

---

## Harvard CS50

I watched from the beginning to the week 2 section of the video.

Summary:

- Creating your first program in C.
- Working with variables, conditionals, and loops.
- Utilizing the Linux command line.-\_+
- Problem-solving approach for computer science problems.
- Understanding the basics of arrays, strings, and command-line arguments.

Links:

- [Harvard CS50 – Full Computer Science University Course](https://www.youtube.com/watch?v=8mAITcNt710&t=22852s)
- [Week 1 C note](https://cs50.harvard.edu/x/2023/notes/1/)
- [Week 2 Arrays note](https://cs50.harvard.edu/x/2023/notes/2/)

---

## Thank you

Thank you for your time and for reading this
