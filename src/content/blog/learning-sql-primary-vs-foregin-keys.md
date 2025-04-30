---
title: "Learning SQL: PRIMARY VS FOREGIN KEYS"
description: "Understand the difference between PRIMARY KEY and FOREIGN KEY in SQL, how they enforce data integrity, and link database tables with practical examples."
date: 2025-05-13
tags: [sql]
draft: true
---

## What Are PRIMARY and FOREIGN Keys?

In SQL, **PRIMARY KEY** and **FOREIGN KEY** are essential concepts for defining relationships between tables and ensuring data integrity.

- A **PRIMARY KEY** is a column (or set of columns) that uniquely identifies each row in a table.
- A **FOREIGN KEY** is a column that creates a link between two tables by referencing the PRIMARY KEY of another table.

Think of PRIMARY KEY as a unique identifier (like a Social Security Number), while a FOREIGN KEY establishes a relationship (like a customer ID in an orders table that links back to the customers table).

## How Do PRIMARY and FOREIGN Keys Work?

### PRIMARY KEY: Ensuring Uniqueness

A PRIMARY KEY ensures that:

- Each row has a **unique** identifier.
- The column(s) cannot contain **NULL** values.
- A table can have **only one** PRIMARY KEY.

**Basic Syntax for PRIMARY KEY:**

```sql
CREATE TABLE customers (
    customer_id INT PRIMARY KEY,
    customer_name VARCHAR(100),
    email VARCHAR(100)
);
```

Here, `customer_id` is the PRIMARY KEY, meaning no two customers can have the same ID.

### FOREIGN KEY: Linking Tables

A FOREIGN KEY ensures:

- A relationship between two tables (parent and child).
- The value in the FOREIGN KEY column must exist in the referenced PRIMARY KEY column (or be NULL, if allowed).

**Basic Syntax for FOREIGN KEY:**

```sql
CREATE TABLE orders (
    order_id INT PRIMARY KEY,
    customer_id INT,
    order_date DATE,
    amount DECIMAL(10,2),
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id)
);
```

Here, `customer_id` in the `orders` table references `customer_id` in the `customers` table, ensuring every order is linked to a valid customer.

### Example: How PRIMARY and FOREIGN Keys Work Together

Let’s say we have two tables:

1. **`customers` Table (Parent Table)**
    
    
    | **customer_id (PK)** | **customer_name** | **email** |
    | --- | --- | --- |
    | 1 | John Doe | [john@example.com](https://mailto:john@example.com/) |
    | 2 | Jane Smith | [jane@example.com](https://mailto:jane@example.com/) |
2. **`orders` Table (Child Table)**
    
    
    | **order_id (PK)** | **customer_id (FK)** | **order_date** | **amount** |
    | --- | --- | --- | --- |
    | 101 | 1 | 2023-01-15 | 150.00 |
    | 102 | 2 | 2023-01-16 | 200.00 |

**What Happens Here?**

- The `customer_id` in `orders` must match an existing `customer_id` in the `customers` table.
- If you try to insert an order with `customer_id = 3` (which doesn’t exist in `customers`), SQL will reject it to maintain **referential integrity**.

## **Key Points to Remember**

✅ **PRIMARY KEY**

- Must be **unique** and **non-NULL**.
- Only **one per table**.
- Can be a single column or a combination (composite key).

✅ **FOREIGN KEY**

- Ensures **data consistency** by linking to a PRIMARY KEY.
- Can reference a table in the **same or a different database**.
- Prevents **orphaned records** (e.g., an order without a valid customer).

✅ **Common Use Cases**

- **One-to-Many Relationships**: One customer can have many orders.
- **Many-to-Many Relationships**: Requires a **junction table** with FOREIGN KEYS.

## Recap

| **Feature** | **PRIMARY KEY** | **FOREIGN KEY** |
| --- | --- | --- |
| **Purpose** | Uniquely identifies a row | Links to a PRIMARY KEY in another table |
| **NULL Values** | Not allowed | Allowed (if not set to NOT NULL) |
| **Uniqueness** | Must be unique | Can have duplicates (unless constrained) |
| **Number per Table** | Only one | Multiple allowed |

By understanding **PRIMARY and FOREIGN KEYS**, you can design efficient, well-structured databases that maintain data integrity.

## Resources

[PostgreSQL Primary Key](https://neon.tech/postgresql/postgresql-tutorial/postgresql-primary-key)

[PostgreSQL Foreign Key](https://neon.tech/postgresql/postgresql-tutorial/postgresql-foreign-key)

## Thank you!

Thank you for your time and for reading this!