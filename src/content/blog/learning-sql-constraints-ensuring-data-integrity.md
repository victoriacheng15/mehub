---
title: "Learning SQL: Constraints – Ensuring Data Integrity"
description: "Learn essential SQL constraints - NOT NULL, UNIQUE, CHECK, and DEFAULT - to enforce data integrity and maintain database accuracy with practical examples."
date: 2025-05-20
tags: [sql]
draft: true
---

## What Are SQL Constraints?

SQL **constraints** are rules that enforce data integrity by controlling what values can be stored in database columns. They help prevent invalid, inconsistent, or duplicate data. This post covers four fundamental constraints:

- `NOT NULL` – Requires a column to always have a value
- `UNIQUE` – Ensures all values in a column are different
- `CHECK` – Validates data against specific conditions
- `DEFAULT` – Provides a fallback value when none is specifie

## How Do These Constraints Work?

### `NOT NULL` – Preventing Empty Values

Ensures a column cannot contain `NULL` (missing/unknown values).

Example:

```sql
CREATE TABLE customers (  
    customer_id INT,  
    name VARCHAR(100) NOT NULL,  -- Must be provided  
    email VARCHAR(100)  
);  
```

- ✅ **Works:** `INSERT INTO customers VALUES (1, 'John Doe', 'john@example.com');`
- ❌ **Fails:** `INSERT INTO customers (customer_id, email) VALUES (2, 'test@example.com');` *(Missing name)*
- **Key Use:** Critical fields like usernames, product IDs, or order dates.

### `UNIQUE` – Eliminating Duplicates

Guarantees all values in a column (or column combination) are distinct.

Example:

```sql
CREATE TABLE products (  
    product_id INT,  
    sku VARCHAR(50) UNIQUE,  -- No duplicate SKUs allowed  
    price DECIMAL(10,2)  
);  
```

- ✅ **Works:** Different SKUs like "A123" and "B456"
- ❌ **Fails:** Inserting two products with `sku = "A123"`
- **Pro Tip:** Can apply to multiple columns:

### `CHECK` – Custom Data Validation

Restricts values based on logical conditions.

Examples:

```sql
CREATE TABLE employees (  
    employee_id INT,  
    age INT CHECK (age >= 18),           -- Minimum age  
    salary DECIMAL(10,2) CHECK (salary > 0),  -- Positive salary  
    department VARCHAR(50) CHECK (department IN ('Sales', 'IT', 'HR'))  
);  
```

❌ Rejects:

- `age = 16`
- `salary = -5000`
- `department = 'Accounting'` (if not in list)

### `DEFAULT` – Automatic Fallback Values

Provides a predefined value when data is omitted.

Example:

```sql
CREATE TABLE orders (  
    order_id INT,  
    order_date DATE DEFAULT CURRENT_DATE,  -- Uses today's date  
    status VARCHAR(20) DEFAULT 'Pending',  
    priority INT DEFAULT 3  
);  
```

✅ **Result:** If you run:

```sql
INSERT INTO orders (order_id) VALUES (1001);  
```

The order automatically gets:

- `order_date = today`
- `status = 'Pending'`
- `priority = 3`

Common Defaults:

- `DEFAULT 0` for numeric fields
- `DEFAULT FALSE` for boolean flags
- `DEFAULT CURRENT_TIMESTAMP` for audit logs

## Recap: SQL Constraints at a Glance

| **Constraint** | **Purpose** | **Example** |
| --- | --- | --- |
| **NOT NULL** | Blocks empty values | `phone VARCHAR(20) NOT NULL` |
| **UNIQUE** | Prevents duplicates | `email VARCHAR(100) UNIQUE` |
| **CHECK** | Enforces custom rules | `rating INT CHECK (rating BETWEEN 1 AND 5)` |
| **DEFAULT** | Provides backup values | `created_at TIMESTAMP DEFAULT NOW()` |

Best Practices:

1. Use `NOT NULL` for mandatory fields
2. Combine `UNIQUE` + `NOT NULL` for natural keys (like email)
3. Apply `CHECK` constraints for business rules (e.g., positive prices)
4. Set `DEFAULT` values for optional columns with common choices

## Putting It All Together

Example Table:

```sql
CREATE TABLE bank_accounts (  
    account_id INT,  
    account_number VARCHAR(20) UNIQUE NOT NULL,  
    balance DECIMAL(15,2) DEFAULT 0.00 CHECK (balance >= 0),  
    account_type VARCHAR(10) CHECK (account_type IN ('Checking', 'Savings')),  
    opened_date DATE DEFAULT CURRENT_DATE NOT NULL  
);  
```

This ensures:

- ✅ No duplicate/empty account numbers
- ✅ Balances never go negative
- ✅ Only valid account types
- ✅ Automatic date tracking

## Resources

[PostgreSQL CHECK Constraints](https://neon.tech/postgresql/postgresql-tutorial/postgresql-check-constraint)

[PostgreSQL UNIQUE Constraints](https://neon.tech/postgresql/postgresql-tutorial/postgresql-unique-constraint)

[PostgreSQL NOT NULL Constraints](https://neon.tech/postgresql/postgresql-tutorial/postgresql-not-null-constraint)

[PostgreSQL DEFAULT Constraints](https://neon.tech/postgresql/postgresql-tutorial/postgresql-default-value)

## Thank you!

Thank you for your time and for reading this!