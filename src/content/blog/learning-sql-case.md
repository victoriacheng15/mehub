---
title: "Learning SQL: CASE"
description: "Learn how to use the SQL CASE statement for conditional logic in queries. Master categorization, custom labels, and conditional calculations with examples!"
date: 2025-03-11
tags: ["backend"]
---

## What is CASE Claus?

The `CASE` statement in SQL is a powerful tool that allows you to add conditional logic to your queries. Think of it as a way to create "if-then-else" logic directly within your SQL statements. It’s incredibly useful when you want to categorize data, create custom labels, or perform calculations based on specific conditions. For example, you can use CASE to classify customers into different tiers based on their purchase amounts or to assign a status to orders depending on their delivery dates.

## How Does CASE Work?

The `CASE` statement works by evaluating a series of conditions and returning a value when a condition is met. It’s like a decision-making tool for your data. Once a condition is satisfied, the corresponding result is returned, and the rest of the conditions are skipped. If no conditions are met, you can optionally specify a default value using the ELSE clause.

Basic Syntax:

```sql
CASE
    WHEN condition1 THEN result1
    WHEN condition2 THEN result2
    ...
    ELSE default_result
END
```

- **WHEN**: Specifies the condition to evaluate.
- **THEN**: Defines the result if the condition is true.
- **ELSE**: Provides a default result if none of the conditions are met (optional).
- **END**: Marks the end of the CASE statement.

Example:

Imagine you have a table called `orders` with the following columns: `order_id`, `customer_id`, `order_date`, and `order_amount`. You want to categorize orders into "Small," "Medium," and "Large" based on their amounts. Here’s how you can do it:

```sql
SELECT 
    order_id,
    order_amount,
    CASE
        WHEN order_amount < 100 THEN 'Small'
        WHEN order_amount BETWEEN 100 AND 500 THEN 'Medium'
        ELSE 'Large'
    END AS order_size
FROM orders;
```

**What Happens Here?**

- The CASE statement evaluates each `order_amount`.
- If the amount is less than 100, it returns "Small."
- If the amount is between 100 and 500, it returns "Medium."
- For all other amounts, it returns "Large."
- The result will include a new column called `order_size` with the appropriate label for each order.

## Key Points to Remember

- **Order Matters**: The CASE statement evaluates conditions in the order they are written. Once a condition is met, it stops evaluating further.
- **ELSE is Optional**: If you don’t include an ELSE clause and no conditions are met, the result will be `NULL`.
- **Use in SELECT, WHERE, and ORDER BY**: You can use CASE in various parts of a query, such as in the SELECT clause to create new columns, in the WHERE clause to filter data, or in the ORDER BY clause to sort results conditionally.
- **Combine with Aggregate Functions**: CASE can be used with aggregate functions like SUM, COUNT, and AVG to perform conditional calculations.

Exmaple with Aggregate Functions:

Suppose you want to calculate the total sales for small, medium, and large orders separately. Here’s how you can do it:

```sql
SELECT 
    CASE
        WHEN order_amount < 100 THEN 'Small'
        WHEN order_amount BETWEEN 100 AND 500 THEN 'Medium'
        ELSE 'Large'
    END AS order_size,
    SUM(order_amount) AS total_sales
FROM orders
GROUP BY order_size;
```

## Recap

- The CASE statement adds conditional logic to your SQL queries.
- It’s useful for categorizing data, creating custom labels, and performing conditional calculations.
- You can use it in SELECT, WHERE, and ORDER BY clauses, and even with aggregate functions.

By mastering the CASE statement, you can make your SQL queries more dynamic and insightful. Keep practicing, and soon you’ll be using CASE like a pro!

## Resources

[PostgreSQL CASE](https://neon.tech/postgresql/postgresql-tutorial/postgresql-case)

## Thank you

Thank you for your time and for reading this!
