---
title: "Learning SQL: UNION and UNION ALL"
description: "Learn the difference between SQL UNION and UNION ALL—how to combine query results, remove duplicates, and optimize performance. Essential for data analysis! 🚀"
date: 2025-04-22
tags: ["backend"]
---

## What are UNION and UNION ALL?

The `UNION` and `UNION ALL` operators in SQL are used to combine the results of two or more `SELECT` queries into a single result set. They allow you to merge data from different tables or queries, making them essential for reporting and data analysis.

- **`UNION`** combines results and removes duplicate rows.
- **`UNION ALL`** combines results but keeps all rows, including duplicates.

Think of them as tools for stacking datasets vertically—like appending one list to another.

---

## How Do UNION and UNION ALL Work?

Both operators require that the queries being combined have the **same number of columns** and **compatible data types**. The key difference is in how they handle duplicates:

- **`UNION`** performs a **distinct operation**, eliminating duplicate rows.
- **`UNION ALL`** **does not remove duplicates**, making it faster since it skips the deduplication step.

### Basic Syntax

```sql
-- Using UNION (removes duplicates)
SELECT column1, column2 FROM table1
UNION
SELECT column1, column2 FROM table2;

-- Using UNION ALL (keeps duplicates)
SELECT column1, column2 FROM table1
UNION ALL
SELECT column1, column2 FROM table2;
```

---

## Example: Combining Customer Data

Imagine you have two tables:

- `customers_east` (stores customers from the East region)
- `customers_west` (stores customers from the West region)

You want to create a **single list of all customers**.

### Using UNION (No Duplicates)

```sql
SELECT customer_id, customer_name FROM customers_east
UNION
SELECT customer_id, customer_name FROM customers_west;
```

**What Happens?** If a customer exists in **both tables**, only **one copy** appears in the result.

### Using UNION ALL (All Rows, Including Duplicates)

```sql
SELECT customer_id, customer_name FROM customers_east
UNION ALL
SELECT customer_id, customer_name FROM customers_west;
```

**What Happens?** If a customer exists in both tables, **both copies** appear in the result.

---

## When to Use UNION vs. UNION ALL?

| **Scenario** | **Use** | **Reason** |
| --- | --- | --- |
| Need unique records only | `UNION` | Removes duplicates |
| Want all records (faster) | `UNION ALL` | No duplicate check, better performance |
| Combining similar datasets | `UNION ALL` | If duplicates are acceptable or unlikely |

---

## Key Points to Remember

1. **Column Matching:** Queries must have the same number of columns with compatible data types.
2. **Performance:** `UNION ALL` is faster than `UNION` because it doesn’t remove duplicates.
3. **Ordering:** If you need sorted results, add `ORDER BY` at the **end of the last query**.

    ```sql
    SELECT name FROM employees
    UNION ALL
    SELECT name FROM contractors
    ORDER BY name;
    ```

4. **Use in Complex Queries:** You can combine `UNION` with `WHERE`, `GROUP BY`, and other clauses.

---

## Recap

- **`UNION`** merges results and removes duplicates.
- **`UNION ALL`** merges results and keeps duplicates.
- Both require matching columns and data types.
- Use `UNION ALL` when duplicates don’t matter for better performance.

---

## Resources

[PostgreSQL UNION](https://neon.tech/postgresql/postgresql-tutorial/postgresql-union)

---

## Thank you

Thank you for your time and for reading this!
