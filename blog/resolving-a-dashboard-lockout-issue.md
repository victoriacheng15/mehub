---
title: "Resolving a Dashboard Lockout Issue"
description: "Resolving a three-year-old dashboard lockout bug caused by expired browser tokens. Learn how route gating and global request checkers prevent infinite error popup loops."
date: 2026-09-01
tags: ["platform", "cncf"]
draft: true
---

## UI Lockout from Polling Loops

Expired tokens in browser storage can trigger infinite error loops in web applications. The background services of the dashboard continue calling the server every ten seconds even when credentials are invalid. When authentication fails, the application receives a series of failed request warnings. Lacking immediate cleanup of the token, the dashboard repeatedly displays error alerts.

Each failed background request triggers a new warning message on the screen. Without automatic cleanup of the expired key, the application keeps sending requests that do not have access. The user must manually clear the stored browser data to recover. The application remains stuck in this loop until the tokens are wiped.

```text
  [ Expired Token ] ──> [ API Requests ]
          ▲                     │
          │                     ▼
  [ Error Alert ] <─── [ 401 Unauthorized ]
```

---

## The String and Number Comparison Trap

A simple comparison mistake in the code allowed invalid tokens to pass the client checks. The program checked the server response code to identify invalid requests. However, the code compared a number directly to a piece of text. Because a number cannot match text, the check always failed to detect the error.

This error allowed invalid tokens to bypass the first dashboard guard. The failure was only detected when the backend server rejected the requests. Fixing the code to compare the same data types resolved this verification bug.

---

## Gating the Main Screen

An investigation of the lockout bug showed that the main web page loaded before the token check completed. Background queries started running before checking if the token was actually empty. Fixing this required blocking the main screen from loading until the token check finished. The dashboard now delays loading any page components until the token is verified.

```text
  [ User Access ]
         │
         ▼
  [ TopContainer ] ──( Token Missing or Loading? )──> YES ──> [ Show Loading / Prompt Token ]
         │                                                            (Queries Blocked)
         ▼ NO
  [ Mount Outlet ] ──> [ Load Dashboard Pages ] ──> [ Run Queries ]
```

A simple status check ensures that the main dashboard only loads when a valid token exists. When the application starts, it reads the saved credentials from browser storage and checks their values. The app blocks the page content from loading until this check is completed. Gating the root layout keeps the screen clean of unauthorized requests.

If the check fails, the application opens the token input popup immediately. This change occurs before any web client requests are created or registered. By stopping the query loop at the root, the application avoids spamming the backend server. Safe designs require clean boundaries between unauthenticated and authenticated pages.

A review enhancement from one of the maintainers polished the final render sequence. The initial revision hid the entire dashboard window when the token input dialog appeared. The maintainer adjusted this to keep the sidebar and header skeleton visible in the background. This layout recovery preserves a consistent frame for the user interface.

---

## Intercepting the Request Storm

The second part of the fix involved setting up a global error checker for web requests. When the server returns an unauthorized error, this checker catches the failure immediately. The program then deletes the invalid token from the browser storage. Finally, the dashboard opens the token input popup without showing error alerts.

By cleaning up the token before opening the popup, background queries stop instantly. This stops new error messages from appearing on the screen. Gated token validation ensures that users without access cannot trigger background tasks.

One of the maintainers also optimized the logic inside the global error checker. The final code analyzes the error type field first to bypass unnecessary comparisons. This prevents redundant processing steps when handling validation failures. The co-authored changes ensure a highly maintainable request handler.

### Related Links

- [PR - fix(ui): resolve token validation and lockout loop](https://github.com/chaos-mesh/chaos-mesh/pull/5046)

---

## Conclusion

Lessons from a three-year-old bug show that old issues can hide in plain sight until triage happens. Finding a type comparison error after years of dashboard updates highlights the value of methodical checking. Resolving these historical edge cases is just as rewarding as building new features. Closing this long-standing issue makes the contribution journey highly satisfying.
