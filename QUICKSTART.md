# Quick Start Guide to Publish on Homebrew

Follow these steps to make `dns-switch` available via `brew install` using the Pinaka tap.

## 1. Push to GitHub (pinaka-io organization)

```bash
# Initialize git if not already done
git init
git add .
git commit -m "Initial commit"

# Push to the pinaka-io organization repository
git remote add origin https://github.com/pinaka-io/dns-switch.git
git branch -M main
git push -u origin main
```

## 2. Create First Release

```bash
# Use the helper script
./scripts/prepare_release.sh 1.0.0

# Or manually:
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

## 3. Create GitHub Release

1. Go to: `https://github.com/pinaka-io/dns-switch/releases/new`
2. Select tag: `v1.0.0`
3. Title: `v1.0.0 - Initial Release`
4. Description: Add release notes
5. Click "Publish release"

## 4. Calculate SHA256

```bash
curl -L https://github.com/pinaka-io/dns-switch/archive/refs/tags/v1.0.0.tar.gz | shasum -a 256
```

Copy the hash (first part before the dash).

## 5. Update Formula

Edit `Formula/dns-switch.rb` and update line 8 with the SHA256 hash:

```ruby
sha256 "YOUR_SHA256_HASH_HERE"
```

The URL is already set to the pinaka-io organization.

## 6. Add to Existing Homebrew Tap

```bash
# Clone the existing pinaka-io tap
git clone https://github.com/pinaka-io/homebrew-tap.git
cd homebrew-tap

# Copy the formula
cp ../dns-switch/Formula/dns-switch.rb Formula/

# Commit and push
git add Formula/dns-switch.rb
git commit -m "Add dns-switch formula"
git push origin main
```

## 7. Test Installation

```bash
# Tap the pinaka-io repository
brew tap pinaka-io/tap

# Install
brew install dns-switch

# Test
sudo dns-switch
```

## 8. Users Install With

Users can now install with:

```bash
brew tap pinaka-io/tap
brew install dns-switch
sudo dns-switch
```

Or in one command:

```bash
brew install pinaka-io/tap/dns-switch
```

## Updating to New Versions

When you release a new version (e.g., v1.1.0):

1. **Create new tag:**
   ```bash
   git tag -a v1.1.0 -m "Release version 1.1.0"
   git push origin v1.1.0
   ```

2. **Create GitHub release** at the URL above

3. **Calculate new SHA256:**
   ```bash
   curl -L https://github.com/pinaka-io/dns-switch/archive/refs/tags/v1.1.0.tar.gz | shasum -a 256
   ```

4. **Update formula** in homebrew-tap:
   ```bash
   cd homebrew-tap
   # Edit Formula/dns-switch.rb
   # - Update version in URL: v1.0.0 -> v1.1.0
   # - Update sha256 with new hash
   git add Formula/dns-switch.rb
   git commit -m "Update dns-switch to v1.1.0"
   git push origin main
   ```

5. **Users update with:**
   ```bash
   brew update
   brew upgrade dns-switch
   ```

## Publish to PyPI (Optional)

To also make it available via `pip install`:

1. Create account on [PyPI](https://pypi.org/)
2. Get API token from PyPI account settings
3. Add token to GitHub secrets as `PYPI_TOKEN`
4. Push a tag - GitHub Actions will auto-publish

```bash
git tag v1.0.0
git push origin v1.0.0
# GitHub Actions will build and publish to PyPI
```

## That's It!

Your tool is now available via:

```bash
brew install pinaka-io/tap/dns-switch
```

🎉

For detailed information, see `BREW_SETUP.md`.
