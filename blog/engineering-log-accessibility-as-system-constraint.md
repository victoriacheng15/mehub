---
title: "[Engineering Log] Accessibility as System Constraint"
description: "Making Cover Craft accessibility part of the image generation contract by enforcing WCAG contrast rules, shared validation, safe random colors, and stable analytics payloads."
date: 2026-08-18
tags: ["retrospective", "platform", "frontend"]
---

## Context

Creating cover images with bad color combinations turns a successful technical render into an unreadable file. This means accessibility is not just a cosmetic feature, but a real system constraint. Adding a validator behind the user controls helped resolve this gap.

A single validation check in the browser was a good start, but it was too easy to bypass. The same contrast math had to run at several different stages of the workflow. The system needed a shared rule that worked in multiple places:

- Form validation before the user generates an image.
- Backend enforcement before the `API` renders a `PNG`.
- Random color selection on the frontend before a palette becomes selectable.
- Analytics payloads that report contrast outcomes consistently.

---

## Challenge

The design required contrast checks on both the frontend and the backend during image generation. The challenge was implementing this duplicate validation without causing the rules to drift apart over time. The system needed a single, shared source of truth to keep the preview and the final API output identical.

Implementing a random color generator on the frontend was meant to be a fun feature that makes it easy for users to pick colors on the spot. The challenge was ensuring this client-side selector only suggests combinations that pass contrast checks. The system had to support quick, creative options without offering unreadable text designs.

---

## Investigation

The first step was using standard web accessibility guidelines to define mathematical limits for colors. The generator needed a hard boundary to reject unreadable pairs before starting any image work. This choice turned a vague design preference into a concrete system requirement.

Trying to write duplicate validation code on both ends of the stack was a wrong turn because it caused immediate maintenance headaches. The same color rules had to be updated in multiple files whenever contrast limits were adjusted. Moving the formula into a shared package kept both environments perfectly aligned from the start.

Analytics tracking exposed another headache during local testing. Dashboard charts expected a stable data structure to render categories correctly. If the API skipped sending empty categories, the frontend charts broke and dropped entire lines.

Designing tests required a clear list of system boundaries to verify. The validation suite had to run checks from the initial form input down to the generated PNG output. This resulted in a simple verification list:

- Does the frontend block unreadable choices?
- Does the backend reject the same unreadable choices?
- Does randomization only return usable color pairs?
- Does analytics preserve expected `AA` and `AAA` categories?

The investigation also drew a line between `UI` warnings and backend enforcement. Helping the user select readable colors is a frontend task, but blocking bad renders is a backend responsibility. This separation kept the codebase clean instead of burying logic inside `React` components.

---

## Solution

The final solution extracted the contrast calculations into a shared utility file. The frontend imports this file to run live UI checks, and the backend uses it to block bad API requests. This shared code ensures the generator never outputs a broken image.

Using a single validator keeps the entire flow aligned. The diagram shows how both sides of the application import the same core rule. This layout stops validation logic from drifting apart over time:

```text
                  [ Shared Contrast Rule ]
                 /                        \\
                /                          \\
               v                            v
    [ UI Form Validation ]        [ API Image Generation ]
               |                            |
               v                            v
     ( User Feedback )              ( Generated PNG )
```

The random color feature on the frontend was also refactored to use the same shared logic. The client-side generator filters out bad color pairs before they can ever reach the user selection menu. This keeps the suggestions creative without offering unreadable combinations.

The analytics backend was updated to normalize empty categories. Database queries now fill in missing values with zero counts before sending the payload. This simple database step keeps the frontend chart lines from crashing.

Writing the validation code once simplified the system in three main places. The frontend and backend both reuse the same package. This shared approach cut down duplication and made test writing much faster:

- Shared contrast utilities defined the reusable rule.
- UI validation used the rule to guide the interaction.
- API generation used the rule to protect the exported artifact.

Testing also became much easier to manage. The test suite runs the same validation checks against both the shared package and the `API` endpoints. If contrast rules need to change later, editing the shared file updates the whole application.

---

## Evolution

The main takeaway is implementing requirements on both the frontend and backend as double guards. The frontend is great for telling the user about rules via the user interface. Meanwhile, the backend acts as the final guard to ensure all requirements are met before generating images.

Find the source code in the [cover-craft](https://github.com/victoriacheng15/cover-craft) repository.
