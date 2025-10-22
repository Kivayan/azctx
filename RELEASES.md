# Building and Releasing azctx

This project uses **dynamic versioning** based on git tags (managed by `hatch-vcs`):

- Version is automatically determined from git tags
- No manual version updates needed in code
- Development builds include commit hash: `0.1.dev12+gbeaba485b.d20251017`
- Release builds use clean version from tag: `1.0.0`

**Check current version:**

```bash
azctx --version
```

---

## Local Build

Build the standalone executable on your local machine:

```bash
# Ensure dev dependencies are installed
uv sync --dev

# Build the executable
uv run pyinstaller --onefile src/cli.py --name azctx

# Test the executable
.\dist\azctx.exe --help
```

The built executable will be in the `dist/` directory (typically 10-20MB).

**Troubleshooting:**

- **PyInstaller import errors**: Edit `azctx.spec`, add missing modules to `hiddenimports=['yaml', 'pkg_resources.py2_warn']`, rebuild with `uv run pyinstaller azctx.spec`
- **Executable too large** (over 50MB): Edit `azctx.spec`, add exclusions: `excludes=['tkinter', 'matplotlib', 'IPython', 'notebook', 'jupyter']`

---

## GitHub Build with Release

Create a new release with automated builds via GitHub Actions:

```powershell
$release = "0.1.6"
git add -A && git commit -m "Release $release"
git tag $release
git push origin $release
```

GitHub Actions will automatically:

- Build the Windows executable with version from tag
- Verify size is under 50MB
- Package as `azctx-v0.1.6-windows.zip`
- Create a GitHub Release with the zip attached
- Generate release notes from commits

**Version Numbering** ([Semantic Versioning](https://semver.org/)):

- **1.0.0** - Major release (breaking changes)
- **1.1.0** - Minor release (new features, backward compatible)
- **1.0.1** - Patch release (bug fixes)
