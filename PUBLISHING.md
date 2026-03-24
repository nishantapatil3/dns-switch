# Publishing to PyPI

This guide explains how to publish dns-switch to PyPI for pipx/pip installation.

## Prerequisites

1. PyPI account: https://pypi.org/account/register/
2. PyPI API token: https://pypi.org/manage/account/token/
3. Python build tools:
   ```bash
   pip install build twine
   ```

## One-Time Setup

### 1. Create PyPI Account
- Go to https://pypi.org/account/register/
- Verify your email

### 2. Generate API Token
- Go to https://pypi.org/manage/account/token/
- Create token with scope: "Entire account"
- Save the token (starts with `pypi-`)

### 3. Configure GitHub Secret (for automated releases)
- Go to: https://github.com/pinaka-io/dns-switch/settings/secrets/actions
- Click "New repository secret"
- Name: `PYPI_TOKEN`
- Value: Your PyPI token
- Click "Add secret"

## Manual Publishing

### 1. Update Version

Edit `pyproject.toml` and `setup.py`:
```python
version = "1.0.1"  # Increment version
```

### 2. Build Package

```bash
# Clean previous builds
rm -rf dist/ build/ *.egg-info

# Build
python -m build
```

This creates:
- `dist/dns-switch-1.0.1.tar.gz` (source)
- `dist/dns_switch-1.0.1-py3-none-any.whl` (wheel)

### 3. Test Upload (Optional)

Test on TestPyPI first:
```bash
# Upload to test.pypi.org
twine upload --repository testpypi dist/*

# Test installation
pipx install --index-url https://test.pypi.org/simple/ dns-switch
```

### 4. Upload to PyPI

```bash
twine upload dist/*
# Enter username: __token__
# Enter password: [your PyPI token]
```

### 5. Test Installation

```bash
pipx install dns-switch
dns-switch --version
```

## Automated Publishing (Recommended)

We have GitHub Actions set up to auto-publish on git tags.

### Create a Release

```bash
# Update version in pyproject.toml and setup.py first
git add pyproject.toml setup.py
git commit -m "Bump version to 1.0.1"
git push

# Create and push tag
git tag v1.0.1
git push origin v1.0.1
```

GitHub Actions will automatically:
1. Build the package
2. Create a GitHub release
3. Publish to PyPI

### Check Release Status

- GitHub Actions: https://github.com/pinaka-io/dns-switch/actions
- PyPI releases: https://pypi.org/project/dns-switch/#history

## Version Guidelines

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR** (1.0.0 → 2.0.0): Breaking changes
- **MINOR** (1.0.0 → 1.1.0): New features, backwards compatible
- **PATCH** (1.0.0 → 1.0.1): Bug fixes

## Troubleshooting

### "File already exists"
You can't re-upload the same version. Increment version number.

### "Invalid or non-existent authentication"
Check your PyPI token is correct and has the right scope.

### "Package name already taken"
The package name might be taken. Choose a different name in `pyproject.toml` and `setup.py`.

## First Release Checklist

- [ ] PyPI account created
- [ ] API token generated
- [ ] Token added to GitHub secrets as `PYPI_TOKEN`
- [ ] Version set to 1.0.0
- [ ] Package built successfully
- [ ] Tested locally with `pipx install -e .`
- [ ] Create git tag v1.0.0
- [ ] Push tag to trigger GitHub Actions
- [ ] Verify on PyPI: https://pypi.org/project/dns-switch/
- [ ] Test installation: `pipx install dns-switch`

## Resources

- [Python Packaging Guide](https://packaging.python.org/)
- [PyPI Help](https://pypi.org/help/)
- [Twine Documentation](https://twine.readthedocs.io/)
