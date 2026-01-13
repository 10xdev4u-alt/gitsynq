# Your First Sync

This guide walks you through the details of your very first synchronization with GitSynq.

## Step 1: Initialize the Local Repo

Ensure you are in a Git repository.

```bash
gitsync init
```

You will be asked:
1. **Project name:** This will be the folder name on the server.
2. **Server IP:** The address of your air-gapped server.
3. **Username:** Your login on that server.
4. **Remote path:** The parent directory where you want your project to live (e.g., `~/projects`).
5. **SSH Key:** (Optional) Path to your private key.
6. **Main branch:** Usually `main`.

## Step 2: The Initial Push

Since the server doesn't have your repository yet, you must perform a **full push**.

```bash
gitsync push --full
```

**What happens behind the scenes?**
1. GitSynq creates a bundle of your *entire* repository history.
2. It transfers this bundle to the server.
3. It runs `git clone <bundle> <project-name>` on the server.
4. It sets up the branches correctly.

## Step 3: Verify on the Server

SSH into your server to confirm everything is there:

```bash
ssh <user>@<host>
cd <remote_path>/<project_name>
git log --oneline
```

You should see your full commit history!

## Step 4: Make a change on the server

Let's test the pull functionality. Create a file on the server:

```bash
echo "Hello from the server" > server-file.txt
git add server-file.txt
git commit -m "Add file from server"
```

## Step 5: Pull back to Laptop

Back on your laptop, run:

```bash
gitsync pull
```

GitSynq will download the new commit and merge it into your local branch.

## Congratulations! 🎉

You've successfully bridged the air-gap. You can now work seamlessly across both environments.
