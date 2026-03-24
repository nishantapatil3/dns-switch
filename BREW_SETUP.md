# Homebrew Distribution Setup

This guide will help you distribute `dns-switch` via Homebrew.

## Prerequisites

1. GitHub account
2. Repository pushed to GitHub
3. Git version tag created

## Step 1: Prepare the Release

1. **Push your code to GitHub (pinaka-io organization):**
   ```bash
   git remote add origin https://github.com/pinaka-io/dns-switch.git
   git add .
   git commit -m "Initial release"
   git push -u origin main
   ```

2. **Create a version tag:**
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

3. **Create a GitHub release:**
   - Go to: https://github.com/pinaka-io/dns-switch/releases/new
   - Select tag `v1.0.0`
   - Title: `v1.0.0`
   - Description: Add release notes
   - Click "Publish release"

4. **Get the SHA256 of the tarball:**
   ```bash
   curl -L https://github.com/pinaka-io/dns-switch/archive/refs/tags/v1.0.0.tar.gz | shasum -a 256
   ```

   Copy the SHA256 hash and update it in `Formula/dns-switch.rb` (line 8).

## Step 2: Add to Pinaka Homebrew Tap

The tap already exists at: https://github.com/pinaka-io/homebrew-tap

1. **Clone the tap repository:**
   ```bash
   git clone https://github.com/pinaka-io/homebrew-tap.git
   cd homebrew-tap
   ```

2. **Copy the formula:**
   ```bash
   # Copy formula from dns-switch project
   cp /path/to/dns-switch/Formula/dns-switch.rb Formula/

   # Commit and push
   git add Formula/dns-switch.rb
   git commit -m "Add dns-switch formula"
   git push origin main
   ```

## Step 3: Install via Homebrew

Now users can install the tool with:

```bash
# Add the pinaka-io tap
brew tap pinaka-io/tap

# Install dns-switch
brew install dns-switch

# Or in one command
brew install pinaka-io/tap/dns-switch
```

## Step 4: Update the Formula (for future releases)

When you release a new version:

1. Create a new tag:
   ```bash
   git tag -a v1.1.0 -m "Release version 1.1.0"
   git push origin v1.1.0
   ```

2. Update the formula:
   - Update the `url` to the new version
   - Update the `sha256` hash
   - Update the `version` if it's explicitly set

3. Push the updated formula:
   ```bash
   cd homebrew-tap
   git add Formula/dns-switch.rb
   git commit -m "Update dns-switch to v1.1.0"
   git push origin main
   ```

## Alternative: Submit to Homebrew Core

To get into the main Homebrew repository (more visibility):

1. Ensure your formula meets [Homebrew's requirements](https://docs.brew.sh/Acceptable-Formulae)
2. Test your formula thoroughly:
   ```bash
   brew install --build-from-source Formula/dns-switch.rb
   brew test dns-switch
   brew audit --strict dns-switch
   ```
3. Submit a PR to [Homebrew/homebrew-core](https://github.com/Homebrew/homebrew-core)

## Testing Locally

Before publishing, test your formula locally:

```bash
# Install from local formula file
brew install --build-from-source Formula/dns-switch.rb

# Test it works
sudo dns-switch

# Uninstall
brew uninstall dns-switch
```

## Troubleshooting

**Issue: SHA256 mismatch**
- Make sure you're using the correct release tarball URL
- Recalculate the SHA256 from the actual tarball

**Issue: Dependencies not found**
- Ensure all Python dependencies are listed as resources
- Update resource URLs and SHA256s from PyPI

**Issue: Command not found after install**
- Check the entry point in setup.py matches the formula
- Verify the bin wrapper script has correct paths

## Notes

- The formula creates a config at `~/.config/dns-switch/config.yaml` on first run
- Users need `sudo` to run the tool
- The formula uses a Python virtualenv to avoid conflicts
- Update the GitHub URLs in the formula to match your actual repository

## Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Python for Formula Authors](https://docs.brew.sh/Python-for-Formula-Authors)
- [How to Create Homebrew Tap](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)
