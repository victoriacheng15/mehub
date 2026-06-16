---
title: "Use Azure Key Vault for Secrets"
description: "Securely manage app secrets in Azure using Key Vault—store credentials like GitHub OAuth IDs safely and handle local development with smart fallbacks. Read more to learn."
date: 2025-12-09
tags: ["cloud", "platform", "python"]
---

## What is Key Vault in Azure?

Azure Key Vault is a cloud service that lets you store sensitive stuff—like API keys, passwords, or OAuth credentials—outside of your actual code.

Before learning about Key Vault, I (like many others) used `.env` files with `os.getenv("VARIABLE")` to avoid hardcoding secrets. That works okay locally… but it’s easy to accidentally commit `.env` to Git if you forget to add it to `.gitignore`. One `git add .` and—oops—your secrets are public 😅.

Key Vault helps avoid that by keeping secrets in a centralized, access-controlled place in Azure. Your app fetches them at runtime, so nothing sensitive lives in your repo or config files.

---

## Using Key Vault in an Academic Project

Recently, an academic project required using Azure Key Vault to manage secrets—so I got to try it out firsthand!

For example, instead of putting `GITHUB_OAUTH_CLIENT_ID` and `GITHUB_OAUTH_CLIENT_SECRET` in code (or relying only on `.env`), they were stored in Key Vault as `github-oauth-client-id` and `github-oauth-client-secret`. (Turns out Key Vault doesn’t allow underscores in secret names—only hyphens!)

To make things work both in Azure and on the local machine, I built a small helper with three functions.

### `get_keyvault_client()`

This sets up the connection to Key Vault. It uses `DefaultAzureCredential`, which is pretty smart: in Azure, it uses Managed Identity (no passwords needed!), and on the local machine, it uses my Azure CLI login (after running `az login`).

You just need to set `AZURE_KEYVAULT_URL` once—usually as an app setting in Azure or in your local environment.

```py
def get_keyvault_client():
    keyvault_url = os.getenv("AZURE_KEYVAULT_URL")
    if not keyvault_url:
        raise ValueError("AZURE_KEYVAULT_URL environment variable not set")
    credential = DefaultAzureCredential()
    return SecretClient(vault_url=keyvault_url, credential=credential)
```

### `get_secret(secret_name)`

This fetches a secret, but with a small trick: it automatically converts underscores to hyphens so you can write code like `get_secret("GITHUB_OAUTH_CLIENT_ID")` even though the actual secret name is `github-oauth-client-id`.

It keeps the code clean and avoids having to remember two different naming styles.

```py
def get_secret(secret_name: str) -> str:
    try:
        client = get_keyvault_client()
        # Converts "GITHUB_OAUTH_CLIENT_ID" → "github-oauth-client-id"
        keyvault_secret_name = secret_name.replace("_", "-")
        secret = client.get_secret(keyvault_secret_name)
        return secret.value
    except Exception as e:
        raise Exception(
            f"Failed to retrieve secret '{secret_name}' from Key Vault: {str(e)}"
        )
```

This allows the same logical name to be used in code regardless of the underlying storage format.

### `get_secret_with_fallback()`

This function attempts to retrieve a secret from Key Vault first. If the secret doesn’t exist, Key Vault is unreachable, or authentication fails, it doesn’t crash silently. Instead:

- It prints a clear warning showing which secret name was attempted (e.g., github-oauth-client-id),
- Then checks for a fallback environment variable (like GITHUB_OAUTH_CLIENT_ID),
- If that’s also missing, it tries a provided default,
- Only if all options fail does it raise an error.

```py
def get_secret_with_fallback(
    secret_name: str, env_var_name: str = None, default: str = None
) -> str:
    try:
        return get_secret(secret_name)
    except Exception as e:
        keyvault_secret_name = secret_name.replace("_", "-")
        print(
            f"[KEY VAULT] Warning: Could not retrieve '{keyvault_secret_name}' from Key Vault: {str(e)}"
        )

        if env_var_name:
            env_value = os.getenv(env_var_name)
            if env_value:
                print(
                    f"[KEY VAULT] Using fallback environment variable: {env_var_name}"
                )
                return env_value

        if default is not None:
            print(f"[KEY VAULT] Using default value for '{secret_name}'")
            return default

        raise ValueError(
            f"Could not retrieve secret '{secret_name}' from Key Vault or environment"
        )
```

---

## Final Thoughts

I’m new to Azure Key Vault, but using it made me realize how easy it is to start building more secure habits—even in school projects.

It’s not about being “100% secure” (nothing really is), but about reducing obvious risks: no hardcoded keys, no accidental Git commits, and clearer separation between code and config.

And honestly? Once it’s set up, it just works. Plus, seeing those fallback logs really helps debug what’s going on during development.

If you’re working on a cloud project—even for class—it’s worth giving Key Vault a try. You might be surprised how smooth it can be!

> 🔐 Small step, big difference: keep secrets out of your code.

---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
