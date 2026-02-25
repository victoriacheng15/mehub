---
title: "Understanding HTTP Methods and Status"
description: "Learn the essentials of HTTP methods (GET, POST, PUT, PATCH, DELETE) and common status codes to build clear, effective APIs. Discover essential tips and practical guides."
date: 2025-08-12
tags: ["platform"]
---

## What is HTTP Methods and Status?

When building or using web APIs, understanding HTTP methods and the status codes that come with them is essential. Whether you’re writing a backend service or integrating with one, these two concepts go hand-in-hand — the method describes *what the client wants to do*, and the status code describes *what actually happened*.

HTTP methods are *verbs* that describe the kind of operation the client wants the server to perform on a resource. For example:

- Want to **read** something? Use `GET`.
- Want to **create** something? Use `POST`.
- Want to **update** something fully? Use `PUT`.
- Want to **update** something partially? Use `PATCH`.
- Want to **delete** something? Use `DELETE`.

When used correctly, they make APIs easier to understand, maintain, and use consistently across teams.

Each method also pairs with different HTTP **status codes**, which are short responses from the server telling you what happened. For example:

- `200 OK` means the request was successful
- `404 Not Found` means the resource doesn’t exist
- `201 Created` means a new resource was successfully added

Let’s walk through each method and its most common status codes.

---

## 🟢 `GET` — Fetch Data

`GET` is used to **retrieve data** from the server. It does not modify anything — just fetches.

Example:

``` json
GET /users/123
```

- This asks the server to return the user with ID `123`.

Common Status Codes:

| Code | Meaning |
|------|---------|
| `200 OK` | Resource found and returned successfully |
| `304 Not Modified` | Resource hasn’t changed (used with caching) |
| `404 Not Found` | Resource does not exist |
| `401 Unauthorized` | User is not authenticated |
| `403 Forbidden` | User doesn’t have access |

---

## 🟡 `POST` — Create a Resource

`POST` is used to **create a new resource** on the server. It usually goes with a JSON body in the request.

Example:

```json
POST /users
Content-Type: application/json

{
  "name": "Alice",
  "email": "alice@example.com"
}
```

- This tells the server to create a new user with the given data.

Common Status Codes:

| Code | Meaning |
|------|---------|
| `201 Created` | New resource was successfully created |
| `400 Bad Request` | Malformed input or missing required fields |
| `409 Conflict` | Duplicate or conflicting data (e.g. email already exists) |
| `422 Unprocessable Entity` | Validation failed (common in RESTful APIs) |

---

## 🔵 `PUT` — Replace a Resource

`PUT` is used to **completely replace** a resource. Think of it as "update the whole thing."

Example:

```json
PUT /users/123
Content-Type: application/json

{
  "name": "Alice Smith",
  "email": "alice.smith@example.com"
}
```

- This replaces the existing user `123` with the new data.

Common Status Codes:

| Code | Meaning |
|------|---------|
| `200 OK` | Resource replaced, and response includes the updated object |
| `204 No Content` | Resource replaced successfully, no response body |
| `400 Bad Request` | Invalid input |
| `404 Not Found` | Resource doesn’t exist to be updated |

---

## 🟣 `PATCH` — Update a Resource Partially

`PATCH` is used to **partially update** a resource — unlike `PUT`, which replaces the whole thing.

Example:

```json
PATCH /users/123
Content-Type: application/json

{
  "name": "Alice S."
}
```

- This only updates the `name` field for user `123`.

Common Status Codes:

| Code | Meaning |
|------|---------|
| `200 OK` | Partial update succeeded, response includes updated resource |
| `204 No Content` | Update succeeded, no body returned |
| `400 Bad Request` | Invalid or incomplete input |
| `404 Not Found` | Target resource not found |

---

## 🔴 `DELETE` — Remove a Resource

`DELETE` is used to **delete a resource** permanently from the server.

Example:

```json
DELETE /users/123
```

- This deletes the user with ID `123`.

Common Status Codes:

| Code | Meaning |
|------|---------|
| `204 No Content` | Successfully deleted, nothing more to say |
| `200 OK` | Successfully deleted, with a message or response body |
| `404 Not Found` | Resource wasn’t found to delete |
| `403 Forbidden` | User is not allowed to delete this resource |

---

## 🎯 Final Thoughts

You don’t need to memorize every HTTP status code — just understand the typical ones that go with each method. Think of methods as **intents** and status codes as **results**. When you design APIs with this mindset, your endpoints become much easier to work with, test, and document.

---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
