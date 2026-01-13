# Git Bundles Explained

Git bundles are a powerful but often overlooked feature of Git. They allow you to package a set of commits and references into a single binary file.

## Why use bundles?

In most cases, Git uses a network protocol (like SSH or HTTPS) to transfer data between repositories. However, in an air-gapped environment, you don't have a direct network connection.

A Git bundle acts as an "offline remote". You can think of it as a `.git` directory compressed into a single file.

## Key Properties

1. **Self-Contained:** A bundle contains all the objects (commits, trees, blobs) needed to reconstruct the references it includes.
2. **Transferable:** Since it's just a file, you can move it via USB drive, SCP, email, or any other file transfer method.
3. **Verifiable:** Git can verify that a bundle is compatible with your current repository before you try to pull from it.
4. **Efficient:** Just like `git fetch`, bundles can be incremental. You only need to package the commits that the other side doesn't have yet.

## Common Commands

### Create a full bundle

```bash
git bundle create myproject.bundle --all
```

### Create an incremental bundle

If you know the other side has everything up to `v1.0`, you can package only the new commits:

```bash
git bundle create myproject-update.bundle v1.0..main
```

### Verify a bundle

Before merging, you can check if it's valid:

```bash
git bundle verify myproject.bundle
```

### List contents

See what branches and tags are inside:

```bash
git bundle list-heads myproject.bundle
```

### Fetch from a bundle

```bash
git fetch myproject.bundle main:local-main
```

## How GitSynq uses them

GitSynq automates the creation, transfer, and merging of these bundles so you don't have to remember the complex syntax or manually SCP files.
