<div align="center">
  <img src="https://github.com/user-attachments/assets/c531f18b-5886-464a-a189-971b39134aee" alt="Browsir Logo" width="250">
</div>


# Browsir 🎩

A simple yet powerful command-line tool to manage multiple browser profiles and shortcuts.

## Features ✨

- Launch different browser profiles with a single command
- Create and manage shortcuts to your favorite websites
- Support for both global (config file) and local shortcuts
- Smart shortcut suggestions when typos occur
- Interactive shortcut creation
- Support for Firefox, Chrome, Brave and Arc browsers (on macOS)
- Cross-platform: works on macOS, Linux and Windows

## Installation 🚀

### From Source

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/browsir.git
   cd browsir
   ```

2. Build and install:

   ```bash
   make install
   ```

   This will build the binary and install it to `/usr/local/bin`

   It might prompt you for the password. This is because we are trying to write files in locations like `/etc/`.

   You do have the source code tho, so you can either check that everything is nice _or_ you can change the installation folders yourself in the Makefile!

4. Verify the installation:
   ```bash
   browsir --version
   ```

You can also just build without installing:

```bash
make build
./browsir --version
```

## Usage 📖

```bash
# Open browser with profile
browsir [profile] [shortcut | website]

browsir personal mail
browsir personal gmail.com

# Search on google, duckduckgo and bravesearch
# Default search engine is google
browsir [profile] [-se | --search-engine]=[google | brave | duckduckgo] -q=[your query]

browsir personal -q="What's the distance between the moon and the sun"
browsir personal -se=brave -q="Is Brave better for privacy"
```

## Available commands and flags

```bash
-ls, --list-shortcuts, list all shortcuts
-se, --search-engine, set search engine for search
-q, query search engine
-v, --version, check browsir version
-h, --help, help
```

### Configuration 🔧

1. Create or modify `.browsir.yml` in your browsir directory:
   ```yaml
   app_name: browsir
   browser_name: chrome # can be 'firefox', 'chrome', 'brave', or 'arc'
   profiles:
     - name: personal # profile name you'll use in commands
       profile_dir: Default # actual profile directory name
       description: Personal browsing
     - name: work
       profile_dir: Profile 1
       description: Work profile
   shortcuts:
     google: google.com
     github: github.com
     mail: gmail.com
   ```

The configuration file allows you to:

- Set your preferred browser (`firefox`, `chrome`, `brave`, or `arc`)
- Define multiple browser profiles with custom names
- Add global shortcuts to frequently visited websites

You can find your Chrome profile directory names by visiting:

- Chrome: `chrome://version`
- Brave: `brave://version`
- Arc: `arc://version`
- Firefox: `about:profiles`
