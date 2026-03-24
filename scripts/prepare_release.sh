#!/bin/bash
# Script to help prepare a Homebrew release

set -e

VERSION=${1:-"1.0.0"}

echo "Preparing release v${VERSION}..."

# Check if git is clean
if [[ -n $(git status -s) ]]; then
    echo "Error: Git working directory is not clean. Commit or stash changes first."
    exit 1
fi

# Create and push tag
echo "Creating git tag v${VERSION}..."
git tag -a "v${VERSION}" -m "Release version ${VERSION}"
git push origin "v${VERSION}"

echo ""
echo "✓ Tag v${VERSION} created and pushed"
echo ""
echo "Next steps:"
echo "1. Create a GitHub release at: https://github.com/pinaka-io/dns-switch/releases/new"
echo "2. Calculate SHA256:"
echo "   curl -L https://github.com/pinaka-io/dns-switch/archive/refs/tags/v${VERSION}.tar.gz | shasum -a 256"
echo "3. Update Formula/dns-switch.rb with the SHA256"
echo "4. Push the formula to: https://github.com/pinaka-io/homebrew-tap"
echo ""
echo "See BREW_SETUP.md for detailed instructions."
