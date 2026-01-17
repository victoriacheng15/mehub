---
title: "Python Orderability Dunder Methods Revealed"
description: "Master Python comparison operators (<, >, <=, >=) with dunder methods. Learn how to make objects sortable using a CoffeeShop example with clear code."
date: 2025-07-01
tags: ["python"]
---

## Mastering Orderability Dunders

In Python, orderability refers to an object's ability to be compared using operators like **`<`**, **`>`**, **`<=`**, and **`>=`**. These operators can be customized for your classes using special dunder methods, just like we did with **`__eq__`** for equality comparisons.

## The 4 Methods That Enable Comparisons

- **`<` Operator: `__lt__(self, other)`**

**What it does**: Teaches Python what "less than" means for your objects

```python
# Sort coffee by price (cheapest first)
def __lt__(self, other):
    return self.price < other.price
```

- **`<=` Operator: `__le__(self, other)`**

**What it does**: Defines "less than or equal to"

```python
def __le__(self, other):
    return self.price <= other.price
```

- **`>` Operator: `__gt__(self, other)`**

What it does: The opposite of `<`

```python
def __gt__(self, other):
    return self.price > other.price  # Explicit is clearer than auto-generated
```

- **`>=` Operator: `__ge__(self, other)`**

**What it does**: Defines "greater than or equal to"

```python
def __ge__(self, other):
    return self.price >= other.price
```

## **Coffee Shop Example**

```python
# coffee_shop.py
class CoffeeShop:
    """A class that demonstrates Python's orderability methods"""
    
    def __init__(self, name, daily_customers, daily_profit):
        self.name = name
        self.daily_customers = daily_customers
        self.daily_profit = daily_profit
    
    # 1. Less Than (<) - The MOST important for sorting
    def __lt__(self, other):
        print(f"🔍 Comparing {self.name} < {other.name}")
        # First compare customers, then profit
        return (self.daily_customers, self.daily_profit) < (other.daily_customers, other.daily_profit)
    
    # 2. Less Than or Equal (<=)
    def __le__(self, other):
        print(f"🔍 Comparing {self.name} <= {other.name}")
        return self < other or self == other
    
    # 3. Greater Than (>)
    def __gt__(self, other):
        print(f"🔍 Comparing {self.name} > {other.name}")
        return not (self <= other)  # Reuses __le__
    
    # 4. Greater Than or Equal (>=)
    def __ge__(self, other):
        print(f"🔍 Comparing {self.name} >= {other.name}")
        return not (self < other)  # Reuses __lt__
    
    # 5. Equality (==) - Required for <= and >=
    def __eq__(self, other):
        return (self.daily_customers == other.daily_customers and 
                self.daily_profit == other.daily_profit)
    
    def __str__(self):
        return f"{self.name} (Customers: {self.daily_customers}, Profit: ${self.daily_profit})"

# Let's see it in action!
if __name__ == "__main__":
    print("=== Creating Coffee Shops ===")
    downtown = CoffeeShop("Downtown Beans", 150, 1200.50)
    suburban = CoffeeShop("Suburban Sips", 200, 1500.75)
    mall = CoffeeShop("Mall Mugs", 200, 1400.00)
    
    print("\n=== Direct Comparisons ===")
    print("Downtown < Suburban:", downtown < suburban)  # Uses __lt__
    print("Mall <= Suburban:", mall <= suburban)       # Uses __le__
    print("Suburban > Mall:", suburban > mall)         # Uses __gt__
    print("Mall >= Downtown:", mall >= downtown)       # Uses __ge__
    
    print("\n=== Sorting Magic ===")
    shops = [downtown, suburban, mall]
    print("Default Sort (uses __lt__):")
    for shop in sorted(shops):
        print("-", shop)
    
    print("\nReverse Sort (still uses __lt__):")
    for shop in sorted(shops, reverse=True):
        print("-", shop)
```

## Result

```python
=== Creating Coffee Shops ===

=== Direct Comparisons ===
🔍 Comparing Downtown Beans < Suburban Sips
Downtown < Suburban: True
🔍 Comparing Mall Mugs <= Suburban Sips
🔍 Comparing Mall Mugs < Suburban Sips
Mall <= Suburban: True
🔍 Comparing Suburban Sips > Mall Mugs
🔍 Comparing Suburban Sips <= Mall Mugs
🔍 Comparing Suburban Sips < Mall Mugs
Suburban > Mall: True
🔍 Comparing Mall Mugs >= Downtown Beans
🔍 Comparing Mall Mugs < Downtown Beans
Mall >= Downtown: True

=== Sorting Magic ===
Default Sort (uses __lt__):
🔍 Comparing Suburban Sips < Downtown Beans
🔍 Comparing Mall Mugs < Suburban Sips
🔍 Comparing Mall Mugs < Suburban Sips
🔍 Comparing Mall Mugs < Downtown Beans
- Downtown Beans (Customers: 150, Profit: $1200.5)
- Mall Mugs (Customers: 200, Profit: $1400.0)
- Suburban Sips (Customers: 200, Profit: $1500.75)

Reverse Sort (still uses __lt__):
🔍 Comparing Suburban Sips < Mall Mugs
🔍 Comparing Downtown Beans < Suburban Sips
🔍 Comparing Downtown Beans < Suburban Sips
🔍 Comparing Downtown Beans < Mall Mugs
- Suburban Sips (Customers: 200, Profit: $1500.75)
- Mall Mugs (Customers: 200, Profit: $1400.0)
- Downtown Beans (Customers: 150, Profit: $1200.5)
```

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
