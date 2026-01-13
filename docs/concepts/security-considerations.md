# Security Considerations

GitSynq aims to be secure by default, leveraging industry-standard protocols.

## 1. SSH Transport

All communication between your laptop and the server happens over **SSH**.

- **Authentication:** GitSynq supports SSH keys (RSA, ED25519, etc.) and integrates with `ssh-agent`.
- **Encryption:** All data in transit is encrypted by the SSH protocol.
- **Integrity:** SSH provides cryptographic integrity checks.

## 2. Host Key Verification

Currently, the default client uses `ssh.InsecureIgnoreHostKey()` for ease of use in LAN environments.

**⚠️ Production Warning:** If you are using GitSynq over a public network (via a bastion or VPN), we recommend configuring strict host key checking. We plan to add support for `known_hosts` validation in a future release.

## 3. Git Bundles

Git bundles are binary files. They contain the same data as a standard `git fetch` operation.

- **No Execution:** Transferring a bundle does not execute code on either side.
- **Validation:** Git validates the bundle's integrity before merging.

## 4. Local Data

- **Config:** The `.gitsync.yaml` file contains server addresses and usernames. It does **not** store passwords.
- **Bundles:** Temporary bundles are stored in `.gitsync-bundles/`. These files are automatically added to `.gitignore` to prevent them from being pushed to public repositories.

## 5. Least Privilege

GitSynq only needs standard user access on the remote server.

- No `sudo` or `root` privileges required.
- Only needs permission to write to the `remote_path` and run `git` commands.
