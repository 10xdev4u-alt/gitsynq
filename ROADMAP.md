# 🗺️ GitSynq Roadmap & Future Ideas

> **"The best way to predict the future is to invent it."** – Alan Kay

Welcome to the GitSynq Roadmap! This document serves as a wishlist, a vision board, and a "Help Wanted" sign for the community. We want GitSynq to be the **standard** for air-gapped development, and we need your help to get there.

Whether you're a Go wizard, a shell scripting guru, or a documentation enthusiast, there's something here for you.

---

## 🔥 High Priority (The "Next Level")

These features are practical, high-value, and needed immediately.

- [ ] **Strict SSH Host Key Checking**: Currently, we default to `InsecureIgnoreHostKey` for ease of use. We need a `--strict` flag (or config option) to parse `~/.ssh/known_hosts` and enforce security.
- [ ] **Resume-able Transfers**: Bundle transfers can fail on flaky connections. Implement logic to resume uploads/downloads from where they left off (using SFTP seek).
- [ ] **Compression Options**: Add support for `zstd` or `lz4` compression for bundles. Git bundles are already packed, but an extra layer can help on extremely slow satellite/radio links.
- [ ] **Windows Native Credential Manager**: Integrate with Windows Credential Manager for secure password storage if SSH keys aren't used.

## 🎮 UX & Terminal UI (TUI)

Make the tool feel magical.

- [ ] **Interactive Conflict Resolution**: A TUI (using [Bubble Tea](https://github.com/charmbracelet/bubbletea)) that lets users pick "Ours" vs "Theirs" directly in the terminal when a merge conflict occurs.
- [ ] **Multi-Profile Support**: Allow `gitsync` to manage multiple servers/projects easily.
  ```bash
  gitsync push --target production-server
  gitsync push --target staging-lab
  ```
- [ ] **GitSynq Shell**: A command `gitsync shell` that drops you into an SSH session on the remote server, pre-configured with the project directory.

## 🔐 Security & Encryption

For the paranoid and the protected.

- [ ] **Encrypted Bundles**: Support transparent encryption of `.bundle` files using [Age](https://github.com/FiloSottile/age) or PGP. This ensures code is encrypted *at rest* on the transfer medium.
- [ ] **SSH Agent Forwarding**: Allow the server to use the local machine's SSH keys for further connections (chaining).
- [ ] **MFA/Interactive Auth**: Better handling for SSH servers that require 2FA codes or interactive keyboard input.

## 🔌 Pluggable Transports

Why stop at SSH?

- [ ] **USB "Sneakernet" Mode**: Sync to a mounted USB drive instead of a remote server.
  ```bash
  gitsync push --transport usb --path /mnt/secure-drive
  ```
- [ ] **S3/Cloud Buckets**: Use an S3 bucket as the "dead drop" for exchanging bundles.
- [ ] **Peer-to-Peer (P2P)**: (Crazy idea) Direct laptop-to-laptop sync over local WiFi/Bluetooth using something like `libp2p` or `syncthing` protocol logic.

## 🤖 CI/CD & Automation

Make GitSynq a building block.

- [ ] **JSON Output**: Add `--json` flag to all commands for easier parsing by scripts.
- [ ] **Webhooks**: Trigger a local or remote webhook URL after a successful sync (e.g., to notify Slack or start a build).
- [ ] **GitSynq Daemon**: A robust background service for the `watch` command that survives reboots.

---

## 👷 How to Contribute

1.  **Pick a Task**: Find something above that excites you.
2.  **Open an Issue**: "I'm working on [Feature X]". We'll discuss the design.
3.  **Fork & Branch**: `git checkout -b feat/super-cool-feature`.
4.  **Cook**: Write clean, idiomatic Go. Add tests.
5.  **Pull Request**: Send it in!

Let's build the ultimate offline-first dev tool together. 🚀
