---
title: "Learning SQL: Subqueries"
description: "Master SQL subqueries & operators (ANY, ALL, EXISTS) with clear examples. Learn nested queries for WHERE, FROM, SELECT clauses—boost your SQL skills today!"
date: 2025-04-08
tags: [sql]
---

## What is a Subquery?

A subquery is simply **a query inside another query,** and also known as **nested queries**. It allows you to use the result of one query as part of another. 

Think of it like asking a question within a question:

*"Find all customers who made purchases above the average order amount."*

To answer this, you’d first need to find the **average order amount** (inner query), then use that number to filter customers (outer query).

Basic Structure:

```sql
SELECT column_name  
FROM table_name  
WHERE column_name OPERATOR (SELECT column_name FROM table_name WHERE condition);  
```

The inner query runs first, and its result is used by the outer query

## Types of Subqueries

Subqueries can be used in different ways:

1. **WHERE Clause Subqueries** – Filter data based on the result of another query.
2. **FROM Clause Subqueries** – Use a subquery as a temporary table.
3. **SELECT Clause Subqueries** – Include a subquery directly in the column list.

### Example 1: Subquery in WHERE Clause

**Problem:** Find all orders greater than the average order amount.

```sql
SELECT order_id, order_amount  
FROM orders  
WHERE order_amount > (SELECT AVG(order_amount) FROM orders);  
```

Here, the inner query (`SELECT AVG(order_amount) FROM orders`) calculates the average, and the outer query uses it to filter orders.

### Example 2: Subquery in FROM Clause

**Problem:** Get the total sales per customer, but only for customers who spent more than $500.

```sql
SELECT customer_id, total_spent  
FROM (  
    SELECT customer_id, SUM(order_amount) AS total_spent  
    FROM orders  
    GROUP BY customer_id  
) AS customer_totals  
WHERE total_spent > 500;  
```

The inner query creates a temporary table (`customer_totals`) with each customer’s total spending, and the outer query filters it.

### Example 3: Subquery in SELECT Clause

**Problem:** Show each order along with how much it differs from the average order.

```sql
SELECT   
    order_id,   
    order_amount,  
    order_amount - (SELECT AVG(order_amount) FROM orders) AS difference_from_avg  
FROM orders;  
```

Here, the subquery calculates the average once and compares each order against it.

## Correlated Subqueries & Special Operators

### Correlated Subqueries

A **correlated subquery** is a special type where the inner query depends on the outer query. It runs once for each row processed by the outer query.

**Example:** Find all customers who placed at least one order above $200.

```sql
SELECT customer_name  
FROM customers c  
WHERE EXISTS (  
    SELECT 1  
    FROM orders o  
    WHERE o.customer_id = c.customer_id  
    AND o.order_amount > 200  
);  
```

### Advanced Subquery Operators: ANY, ALL, EXISTS

These operators supercharge subqueries for precise comparisons:

- **EXISTS Operator**

Checks if a subquery returns **any rows**. Perfect for "has a relationship" checks.

```sql
SELECT customer_name  
FROM customers c  
WHERE EXISTS (  
    SELECT 1 FROM orders o  
    WHERE o.customer_id = c.customer_id  
);
```

Returns customers **only if they have at least one order**.

- **ANY Operator**

Returns `TRUE` if **any value** in the subquery meets the condition.

```sql
-- Find products with ANY order exceeding 100 units  
SELECT product_name  
FROM products  
WHERE product_id = ANY (  
    SELECT product_id FROM order_details  
    WHERE quantity > 100  
);
```

`= ANY` is equivalent to `IN`; `> ANY` means "greater than the smallest value".

- **ALL Operator**

Requires **every value** in the subquery to match.

```sql
-- Find employees older than ALL interns  
SELECT name FROM employees  
WHERE age > ALL (  
    SELECT age FROM interns  
);
```

`> ALL` means "greater than the largest value".

**Quick Comparison**

| **Operator** | **Use Case** | **Example** |
| --- | --- | --- |
| `EXISTS` | Check for relationships | "Has this customer ordered?" |
| `ANY` | Flexible comparisons | "Is this value in the top 50%?" |
| `ALL` | Strict comparisons | "Is this higher than ALL competitors?" |

## Recap

✅ Subqueries help break complex problems into smaller steps.

✅ They can be used in WHERE, FROM, and SELECT clauses.

✅ Correlated subqueries reference the outer query, running row by row.

✅ Always test subqueries separately first to ensure they return the right data.

**Final Tip:** If a subquery feels too complicated, try writing the inner query first, then build the outer query around it. With practice, you’ll use subqueries naturally in your SQL work!

## Resources

[PostgreSQL Subquery](https://neon.tech/postgresql/postgresql-tutorial/postgresql-subquery)

[PostgreSQL Correlated Subquery](https://neon.tech/postgresql/postgresql-tutorial/postgresql-correlated-subquery)

[PostgreSQL ANY Operator](https://neon.tech/postgresql/postgresql-tutorial/postgresql-any)

[PostgreSQL ALL Operator](https://neon.tech/postgresql/postgresql-tutorial/postgresql-all)

[PostgreSQL EXISTS Operator](https://neon.tech/postgresql/postgresql-tutorial/postgresql-exists)

## Thank you!

Thank you for your time and for reading this!