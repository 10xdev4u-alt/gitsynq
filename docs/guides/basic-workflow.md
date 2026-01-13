# Basic Workflow Guide

This guide describes the typical day-to-day workflow when using GitSynq.

## 1. Starting your day

Before heading to the air-gapped environment, ensure your local repository is up to date with GitHub:

```bash
git pull origin main
```

## 2. Pushing to the server

Once you're in the LAN and ready to work on the server, push your latest local commits:

```bash
gitsync push
```

GitSynq will create a bundle of new commits since your last sync and transfer it to the server.

## 3. Working on the server

SSH into your server and work as you normally would:

```bash
ssh prince@192.168.1.100
cd ~/projects/my-awesome-project
# Edit files, run tests, etc.
git add .
git commit -m "Implement feature X on server"
```

## 4. Pulling changes back

When you're done working on the server and want to sync back to your laptop (and eventually GitHub):

```bash
# On your laptop
gitsync pull
```

## 5. Pushing to GitHub

After pulling from the server, your local repository has the remote changes. Now you can push them to GitHub:

```bash
git push origin main
```

**Pro Tip:** You can combine steps 4 and 5 by using the `--push` flag:

```bash
gitsync pull --push
```

## Summary

Think of `gitsync push` as `git push` to the server, and `gitsync pull` as `git pull` from the server.
