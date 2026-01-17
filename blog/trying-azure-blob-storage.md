---
title: "Trying Azure Blob Storage"
description: "How I loaded CSV data from Azure Blob Storage in a school project—with local fallback and no hardcoded secrets."
date: 2025-12-23
tags: ["platform"]
---

## What is Azure Blob Storage?

Azure Blob Storage is a cloud service for storing unstructured data—like images, logs, backups, or **CSV files**. Think of it as a scalable, secure “cloud hard drive” with an API, access control, and built-in redundancy.

Unlike local files, blobs can be accessed by any authorized app—whether it’s running on your laptop, in Azure, or on a server across the world.

### Why I Tried It

Azure Blob Storage was a **requirement** in a recent academic project. Before that, I hadn’t really thought about how or when I’d use it in my own projects—but I was curious to see how it actually works in practice.

The task was simple: load a dataset (`All_Diets.csv`) from the cloud instead of a local file. So I set up a Blob container and wrote a loader that tries cloud first, then falls back to local.

> At the time, I used a connection string via environment variables (like in a `.env` file), which is a common and practical approach for local development or low-risk projects. Later, I learned about **Azure Key Vault**—a more secure option for managing secrets in cloud environments. Neither is “wrong”; it’s about **choosing the right tool for the context**.

---

## The Fallback Loader: Cloud + Local, One Function

Here’s the entry point used by the app:

```py
def load_dataset(filename="All_Diets.csv"):
    # Try Azure Blob Storage first
    if os.getenv("AZURE_STORAGE_CONNECTION_STRING"):
        try:
            from ..blob_storage import read_csv_from_blob
            return read_csv_from_blob(filename)
        except Exception:
            pass  # Silently fall back to local

    # Local fallback
    csv_path = Path(__file__).parent.parent / "datasets" / filename
    return pd.read_csv(csv_path)
```

This means:

- In Azure: app loads data directly from Blob Storage
- On my laptop: it uses the local datasets/ folder
- Same code, no config changes

No crashing. No manual switching. Just works.

### How `read_csv_from_blob` Works

Behind the scenes, a small utility handles the Azure interaction:

```py
def read_csv_from_blob(blob_name: str, container_name: str = "datasets") -> pd.DataFrame:
    blob_service_client = get_blob_service_client()
    blob_client = blob_service_client.get_blob_client(container=container_name, blob=blob_name)
    
    # Download as bytes → load into pandas
    blob_data = blob_client.download_blob().readall()
    return pd.read_csv(io.BytesIO(blob_data))
```

It relies on the `AZURE_STORAGE_CONNECTION_STRING` environment variable—stored in `.env` during development.

---

## Final Thoughts

I hadn’t planned to use Blob Storage before this assignment—and I’m not sure when I’ll need it in a personal project yet. But having hands-on experience with real Azure services is a huge bonus.

Through academic projects like this, I’m getting familiar with core cloud patterns:

- Storing secrets (Key Vault)
- Managing containers (ACR)
- Hosting data (Blob Storage)

Even if I don’t use them all right away, knowing how they work—and how they fit together—makes the cloud feel less abstract.

> ☁️ Sometimes the best learning comes from a requirement you didn’t ask for—but turns out to be useful anyway.

---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
