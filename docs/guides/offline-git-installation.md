# 🆘 Help! My Server Doesn't Have Git!

So, you're ready to use GitSynq, but you hit a wall:

```text
bash: git: command not found
```

And since the server has no internet, you can't just run `apt install git` or `yum install git`. 😫

**Don't panic.** This is the "Chicken and Egg" problem of air-gapped systems. Here is how you fix it.

---

## 🏆 Method 1: The Static Binary (Easiest)

The simplest solution is to bring a "portable" version of Git that doesn't need to be installed as root and doesn't depend on complex system libraries.

1.  **On your laptop**, search for a "static git binary" for Linux.
    *   A popular source is the [git-static](https://github.com/jigish/git-static/releases) repo (unofficial but works).
    *   Or check standard kernel.org binaries.
2.  **Download** the binary (e.g., `git-static-amd64`).
3.  **Transfer it** to your server:
    ```bash
    scp git-static-amd64 user@server:~/git
    ```
4.  **On the Server**, make it executable and add it to your path:
    ```bash
    chmod +x ~/git
    # Add this to your .bashrc
    export PATH=$HOME:$PATH
    ```

---

## 📦 Method 2: The Package "Sneakernet"

If you have `sudo` access but no internet, you can download the package files (`.deb` or `.rpm`) on your laptop and carry them over.

### For Ubuntu / Debian

1.  **On your Laptop:**
    You need to download `git` and its dependencies.
    ```bash
    # Create a folder for the files
    mkdir git-debs && cd git-debs
    
    # Download git and key dependencies (this list might vary by OS version)
    apt-get download git git-man liberror-perl
    ```
2.  **Transfer:** SCP the `git-debs` folder to the server.
3.  **On the Server:**
    ```bash
    sudo dpkg -i git-debs/*.deb
    ```

### For CentOS / RHEL / Fedora

1.  **On your Laptop:**
    ```bash
    mkdir git-rpms && cd git-rpms
    
    # Use yumdownloader to get the RPMs
    yumdownloader --resolve git
    ```
2.  **Transfer:** SCP the `git-rpms` folder to the server.
3.  **On the Server:**
    ```bash
    sudo rpm -ivh git-rpms/*.rpm
    ```

---

## 🛠️ Method 3: Compile from Source (The "Pro" Way)

If your server has a C compiler (like `gcc`)—which most research/HPC clusters do—you can compile Git yourself. This is great because you don't need `sudo`.

1.  **On your Laptop:**
    Download the latest source tarball from [kernel.org](https://mirrors.edge.kernel.org/pub/software/scm/git/). Look for `git-2.xx.x.tar.gz`.
2.  **Transfer:** SCP the `.tar.gz` file to the server.
3.  **On the Server:**
    ```bash
    # Extract
    tar -zxf git-2.xx.x.tar.gz
    cd git-2.xx.x
    
    # Configure to install in your home directory (no sudo needed!)
    ./configure --prefix=$HOME/local
    
    # Build and install
    make && make install
    
    # Add to PATH (add to .bashrc to make permanent)
    export PATH=$HOME/local/bin:$PATH
    ```

---

## 🚀 Verification

Once you've done one of the above, run this on the server:

```bash
git --version
```

If you see version text, **you are ready to use GitSynq!**
