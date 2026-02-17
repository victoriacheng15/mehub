---
title: "Python Secret Double Underscore Tricks"
description: "Learn Python dunder methods by building a Coffee class step-by-step. Master init, str, repr, eq & len with clear examples. Perfect for beginners!"
date: 2025-06-10
tags: ["python"]
---

## **Dunder Methods Explained**

If you've used Python classes, you know `__init__` - the method that creates objects. But Python has many more *dunder* (double underscore) methods that control how objects behave. Let's explore the most useful ones.

---

## Building the Class: Dunder by Dunder

### **`__init__`: The Foundation**

Sets up your object’s initial state

```python
class Coffee:
    def __init__(self, size, price):
        self.size = size
        self.price = price

# Usage
morning_coffee = Coffee("Large", 4.99)
print(f"Object - {morning_coffee}") # this will print the object
print(f"Coffee size: {morning_coffee.size}")  # Output: Coffee size: Large
```

### **`__str__`: The Human-Friendly Display**

Controls what `print()` shows

```python
class Coffee:
    def __init__(self, size, price):
        self.size = size
        self.price = price
    
    def __str__(self):
        return f"Coffee size: {self.size} and price: {self.price}"

# Usage
morning_coffee = Coffee("Large", 4.99)
print(f"Object - {morning_coffee}") # this will print the object
# this will disappeared once added __str__ method
print(f"Coffee size: {morning_coffee.size}")  # Output: Coffee size: Large
print(f"__str__ - {morning_coffee}")  # Coffee size: Large and price: 4.99
```

### **`__repr__`: The Debugger’s Truth**

Shows how to recreate the objec

```python
class Coffee:
    def __init__(self, size, price):
        self.size = size
        self.price = price
    
    def __str__(self):
        return f"{self.size} coffee (${self.price})"

    def __repr__(self):
        return f"Coffee(size='{self.size}', price={self.price})"

# Usage
morning_coffee = Coffee("Large", 4.99)
print(f"Object - {morning_coffee}") # this will print the object
# this will disappeared once added __str__ method
print(f"Coffee size: {morning_coffee.size}")  # Output: Coffee size: Large
print(f"__str__ - {morning_coffee}")  # Coffee size: Large and price: 4.99
print(f"Debugging: {repr(morning_coffee)}") # Debugging: Coffee(size='Large', price=4.99)
```

### **`__eq__`: The Equality Judge**

Defines what `==` means for your object

```python
class Coffee:
    def __init__(self, size, price):
        self.size = size
        self.price = price
    
    def __str__(self):
        return f"{self.size} coffee (${self.price})"

    def __repr__(self):
        return f"Coffee(size='{self.size}', price={self.price})"

    def __eq__(self, other):
        return (self.size == other.size) and (self.price == other.price)

# Usage
morning_coffee = Coffee("Large", 4.99)
print(f"Object - {morning_coffee}") # this will print the object
# this will disappeared once added __str__ method
print(f"Coffee size: {morning_coffee.size}")  # Output: Coffee size: Large
print(f"__str__ - {morning_coffee}")  # Coffee size: Large and price: 4.99
print(f"Debugging: {repr(morning_coffee)}") # Debugging: Coffee(size='Large', price=4.99)

large_coffee = Coffee("Large", 4.99)
print(morning_coffee == large_coffee) # True
small_coffee = Coffee("Small", 3.99)
print(morning_coffee == small_coffee) # False
```

### **`__len__`: The Size Inspector**

Makes objects work with `len()`

```python
class Coffee:
    def __init__(self, size, price):
        self.size = size
        self.price = price
    
    def __str__(self):
        return f"{self.size} coffee (${self.price})"

    def __repr__(self):
        return f"Coffee(size='{self.size}', price={self.price})"

    def __eq__(self, other):
        return (self.size == other.size) and (self.price == other.price)

    def __len__(self):
        return len(self.size)  # Returns character count of size (e.g., "Large" → 5)

# Usage
morning_coffee = Coffee("Large", 4.99)
print(f"Object - {morning_coffee}") # this will print the object
# this will disappeared once added __str__ method
print(f"Coffee size: {morning_coffee.size}")  # Output: Coffee size: Large
print(f"__str__ - {morning_coffee}")  # Coffee size: Large and price: 4.99
print(f"Debugging: {repr(morning_coffee)}") # Debugging: Coffee(size='Large', price=4.99)

large_coffee = Coffee("Large", 4.99)
print(morning_coffee == large_coffee) # True
med_coffee = Coffee("Medium", 3.99)
print(morning_coffee == med_coffee) # False

print(len(morning_coffee)) # 5
print(len(med_coffee)) # 6
```

---

## Recap

We built a `Coffee` class step by step using special Python methods:

- `__init__` - Sets up the object when created
- `__str__` - Controls what prints when you display the object
- `__repr__` - Shows how to recreate the object (helpful for debugging)
- `__eq__` - Defines how to compare two objects
- `__len__` - Lets you check the object's size with len()

These special methods (called "dunder" methods because of their double underscores) make your objects work naturally with Python's built-in functions. Each one adds useful behavior while keeping everything working together smoothly.

---

## Thank you

Thank you for your time and for reading this!
