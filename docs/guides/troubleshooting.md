# Troubleshooting

Common issues and how to fix them.

## SSH Connection Failed

**Symptoms:** `❌ SSH connection failed: ...`

**Possible Causes:**
1. **Wrong Credentials:** Check your username and host in `.gitsync.yaml`.
2. **Network Issue:** Can you `ping` the server?
3. **SSH Key:** Ensure your private key is loaded or the path in `ssh_key_path` is correct.
4. **Permissions:** Ensure your user has SSH access on the server.

**Fix:** Try connecting manually using `ssh <user>@<host>`. If that fails, the issue is with your SSH setup, not GitSynq.

## Remote Setup Failed

**Symptoms:** `❌ Remote setup failed: ...`

**Possible Causes:**
1. **Git not installed:** The server MUST have `git` installed and in the `PATH`.
2. **Permissions:** The user must have write access to the `remote_path`.
3. **Invalid Branch:** Ensure the branch specified in `.gitsync.yaml` exists locally.

**Fix:** Run `gitsync --verbose push` to see the exact error output from the server.

## No New Commits to Bundle

**Symptoms:** `⚠️  No new commits or first push. ...`

**Possible Causes:**
1. You haven't made any new commits since your last `gitsync push`.
2. You are trying to push an incremental update but the server doesn't have the base commits yet.

**Fix:** Use `gitsync push --full` to force a full repository synchronization.

## Merge Conflicts

**Symptoms:** `❌ Merge failed: ...`

**Possible Causes:** Changes were made to the same lines of the same files on both your laptop and the server.

**Fix:**
1. GitSynq will leave your local repository in a "merging" state.
2. Resolve the conflicts manually using your IDE or `git mergetool`.
3. Commit the resolution: `git commit -m "Resolve sync conflicts"`.

## Still having trouble?

Please [open an issue](https://github.com/princetheprogrammerbtw/gitsynq/issues) on GitHub and include:
- The command you ran.
- The full output with the `--verbose` flag.
- Your OS and Git version.
