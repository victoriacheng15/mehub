---
title: "Python Dunder Methods: Arithmetic operators"
description: "Unlock the power of Python's arithmetic operators! Learn how dunder methods like __add__ and __mul__ work under the hood to give you full control over your objects."
date: 2025-07-08
tags: ["python"]
---

## What are Arithmetic Operators in Python?

Arithmetic is one of the first things we learn in Python — `x + y`, `x * y`, `x / y` — but under the hood, Python handles all of these using something more advanced: **dunder methods**. Let’s break it all down.

Python's **arithmetic operators** are symbols used to perform mathematical operations:

- `+` for addition  
- `-` for subtraction  
- `*` for multiplication  
- `/` for division  
- `%` for modulo (remainder)  
- `//` for floor division (rounding down)  
- `**` for exponentiation (power)  
- `@` for matrix multiplication (introduced in Python 3.5)

These are called **binary operators** because they work between two values: `x + y`.

There are also **unary operators** like `+x` or `-x`, which act on just one value.

But these operators are more than just symbols. In Python, they're backed by special methods called **"dunder methods"**, short for "double underscore". These let you override and define what the operators actually *do* for your objects.

---

## What are the Main Arithmetic Dunder Methods?

Python defines a pair of dunder methods for most arithmetic operators:

| Operation | Left-Hand Method | Right-Hand Method | Description         |
|-----------|------------------|-------------------|---------------------|
| `x + y`   | `__add__`        | `__radd__`        | Add or concatenate  |
| `x - y`   | `__sub__`        | `__rsub__`        | Subtract            |
| `x * y`   | `__mul__`        | `__rmul__`        | Multiply            |
| `x / y`   | `__truediv__`    | `__rtruediv__`    | Divide              |
| `x % y`   | `__mod__`        | `__rmod__`        | Modulo              |
| `x // y`  | `__floordiv__`   | `__rfloordiv__`   | Floor divide        |
| `x ** y`  | `__pow__`        | `__rpow__`        | Exponentiate        |
| `x @ y`   | `__matmul__`     | `__rmatmul__`     | Matrix multiply     |

If Python evaluates x + y and x.**add**(y) returns NotImplemented, it will then try y.**radd**(x). This fallback mechanism allows custom classes to handle operations even when the left-hand side doesn't know how which is especially useful when mixing user-defined objects with built-in types.

---

## NumberBox

To see how arithmetic dunder methods work in practice, let’s look at a simple class called `NumberBox`. It wraps a number and overloads Python’s arithmetic operators so we can trace exactly which dunder method is called.

This class lets us experiment with both:

- Binary operators (e.g. x + y, x * y, x ** y)
- Unary operators (e.g. -x, +x, abs(x))

Each method prints its name when called, so you can clearly observe how Python handles different cases — including fallback to the right-hand method (like **radd**) when needed.

```python
class NumberBox:
    def __init__(self, value):
        self.value = value

    def __repr__(self):
        return f"NumberBox({self.value})"

    # Binary arithmetic operators
    def __add__(self, other):
        print("__add__ called")
        if isinstance(other, NumberBox):
            return NumberBox(self.value + other.value)
        return NotImplemented

    def __radd__(self, other):
        print("__radd__ called")
        if isinstance(other, (int, float)):
            return NumberBox(other + self.value)
        return NotImplemented

    def __sub__(self, other):
        print("__sub__ called")
        if isinstance(other, NumberBox):
            return NumberBox(self.value - other.value)
        return NotImplemented

    def __rsub__(self, other):
        print("__rsub__ called")
        if isinstance(other, (int, float)):
            return NumberBox(other - self.value)
        return NotImplemented

    def __mul__(self, other):
        print("__mul__ called")
        if isinstance(other, NumberBox):
            return NumberBox(self.value * other.value)
        return NotImplemented

    def __rmul__(self, other):
        print("__rmul__ called")
        if isinstance(other, (int, float)):
            return NumberBox(other * self.value)
        return NotImplemented

    def __truediv__(self, other):
        print("__truediv__ called")
        if isinstance(other, NumberBox):
            return NumberBox(self.value / other.value)
        return NotImplemented

    def __rtruediv__(self, other):
        print("__rtruediv__ called")
        if isinstance(other, (int, float)):
            return NumberBox(other / self.value)
        return NotImplemented

    def __floordiv__(self, other):
        print("__floordiv__ called")
        if isinstance(other, NumberBox):
            return NumberBox(self.value // other.value)
        return NotImplemented

    def __rfloordiv__(self, other):
        print("__rfloordiv__ called")
        if isinstance(other, (int, float)):
            return NumberBox(other // self.value)
        return NotImplemented

    def __mod__(self, other):
        print("__mod__ called")
        if isinstance(other, NumberBox):
            return NumberBox(self.value % other.value)
        return NotImplemented

    def __rmod__(self, other):
        print("__rmod__ called")
        if isinstance(other, (int, float)):
            return NumberBox(other % self.value)
        return NotImplemented

    def __pow__(self, other):
        print("__pow__ called")
        if isinstance(other, NumberBox):
            return NumberBox(self.value ** other.value)
        return NotImplemented

    def __rpow__(self, other):
        print("__rpow__ called")
        if isinstance(other, (int, float)):
            return NumberBox(other ** self.value)
        return NotImplemented

    def __matmul__(self, other):
        print("__matmul__ called")
        if isinstance(other, NumberBox):
            return NumberBox(self.value * other.value)  # simplified matrix multiply
        return NotImplemented

    def __rmatmul__(self, other):
        print("__rmatmul__ called")
        if isinstance(other, (int, float)):
            return NumberBox(other * self.value)
        return NotImplemented

    # Unary arithmetic operators
    def __neg__(self):
        print("__neg__ called")
        return NumberBox(-self.value)

    def __pos__(self):
        print("__pos__ called")
        return NumberBox(+self.value)

    def __abs__(self):
        print("__abs__ called")
        return NumberBox(abs(self.value))
```

### Operator Overloading in Action

The code below demonstrates how each overloaded operator works with `NumberBox`. Each operation prints the name of the dunder method being called, so you can clearly see which method Python is using under the hood. This is especially helpful for understanding how left-hand vs. right-hand methods behave when mixing `NumberBox` with regular numbers.

```python
a = NumberBox(10)
b = NumberBox(3)

print(a + b)     # __add__
print(5 + a)     # __radd__
print(a - b)     # __sub__
print(5 - a)     # __rsub__
print(a * b)     # __mul__
print(2 * a)     # __rmul__
print(a / b)     # __truediv__
print(20 / a)    # __rtruediv__
print(a // b)    # __floordiv__
print(20 // a)   # __rfloordiv__
print(a % b)     # __mod__
print(20 % a)    # __rmod__
print(a ** b)    # __pow__
print(2 ** a)    # __rpow__
print(a @ b)     # __matmul__
print(2 @ a)     # __rmatmul__

print(-a)        # __neg__
print(+a)        # __pos__
print(abs(a))    # __abs__
```

---

## Recap

- Python arithmetic operators (`+`, `-`, `*`, `/`, etc.) are powered by **dunder methods** like `__add__`, `__sub__`, and so on.
- Each operator has two method forms: a **left-hand version** (e.g., `__add__`) and a **right-hand version** (e.g., `__radd__`).
- If the left-hand method returns `NotImplemented`, Python tries the right-hand method.
- You can **overload** these methods in your own classes to define custom behavior.
- Unary operators like `-x` or `abs(x)` also have corresponding dunder methods (`__neg__`, `__abs__`, etc.).
- A class like `NumberBox` is a great way to practice and visualize how operator overloading works in Python.

---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
