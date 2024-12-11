<img style="width: 250px; text-align: center; display: inline-block; margin: 0 auto;" src="https://github.com/user-attachments/assets/c531f18b-5886-464a-a189-971b39134aee" />


# Browsir 🎩

A simple yet powerful command-line tool to manage multiple browser profiles and shortcuts.

## Features ✨

- Launch different browser profiles with a single command
- Create and manage shortcuts to your favorite websites
- Support for both global (config file) and local shortcuts
- Smart shortcut suggestions when typos occur
- Interactive shortcut creation
- Support for Chrome, Brave and Arc browsers (on macOS)
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

3. Verify the installation:
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

### Configuration 🔧

1. Create or modify `.browsir.yml` in your browsir directory:
   ```yaml
   app_name: browsir
   browser_name: chrome  # can be 'chrome', 'brave', or 'arc'
   profiles:
     - name: personal    # profile name you'll use in commands
       profile_dir: Default  # actual profile directory name
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
- Set your preferred browser (`chrome`, `brave`, or `arc`)
- Define multiple browser profiles with custom names
- Add global shortcuts to frequently visited websites

You can find your Chrome profile directory names by visiting:
- Chrome: `chrome://version`
- Brave: `brave://version`
- Arc: `arc://version`
