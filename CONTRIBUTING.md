
# Contributing to flibgolite

This document describes the development workflow, branching rules, commit message conventions, pull request requirements, and the full release process used in this repository.
> *This release workflow applies to all versions starting from v2.2.6 and above.*

The project does not use CI/CD. All builds, tagging, and releases are performed manually.  
Version information is embedded into binaries at compile time, therefore **Git tags must be created before building release artifacts**.

---

# 1. Branching Strategy

All code changes must be made in dedicated branches.  
Allowed branch name prefixes:

- `feature/<name>` — new functionality  
- `fix/<name>` — bug fixes  
- `refactor/<name>` — internal improvements  
- `docs/<name>` — documentation updates  
- `release/<version>` — optional, for preparing large releases  

Examples:

```
feature/optimize-parser
fix/issue-123
docs/update-readme
```

## Creating a branch

### VS Code
1. Click the branch name in the bottom-left corner.
2. Select **Create new branch**.
3. Enter the branch name.

### Terminal
```sh
git checkout -b feature/my-change
```

---

# 2. Commit Message Guidelines

Commit messages must be clear, concise, and meaningful.

Recommended format:

```
<type>: <short description>

[optional longer explanation]
```

Allowed `<type>` values:

- `feat` — new feature  
- `fix` — bug fix  
- `refactor` — code restructuring  
- `docs` — documentation changes  
- `build` — build scripts, Dockerfile, Makefile  
- `chore` — maintenance tasks  

Examples:

```
feat: add new parser for extended syntax
fix: correct handling of empty input
docs: update installation instructions
```

---

# 3. Pull Request Rules

All code changes must go through a Pull Request.

Requirements:

- PR must target the `master` branch.
- PR title must describe the change clearly.
- PR description must explain the purpose of the change.
- Small, focused PRs are preferred.
- Use **Squash and merge** when merging.

### Creating a PR

#### VS Code
1. Open the **Source Control** panel.
2. Push your branch.
3. Click **Create Pull Request** (from GitHub Pull Requests extension).
4. Fill in title and description.
5. Submit.

#### Terminal
```sh
git push -u origin feature/my-change
```
Then open GitHub and create a PR manually.

---

# 4. Direct Commits to `master`

Direct pushes to `master` are allowed **only** for:

- README updates  
- documentation changes  
- comments  
- formatting  
- non-functional changes that do not affect code or build  
- typo fixes  

Direct pushes **must not** modify:

- application code  
- build scripts  
- Dockerfile  
- release logic  
- versioning logic  

---

# 5. Release Policy

Releases follow **Semantic Versioning (SemVer)**:

- **MAJOR** — incompatible changes  
- **MINOR** — new features, backward compatible  
- **PATCH** — bug fixes  

Because version information is embedded into binaries at compile time, the release workflow is:

```
1. Merge all changes into master
2. Create and push Git tag (version)
3. Build binaries (version is now correct)
4. Build Docker image
5. Push Docker image
6. Create GitHub Release and upload binaries
```

---

# 6. Release Workflow (Step-by-Step)

## 6.1 Prepare `master`

### VS Code
- Switch to `master` (bottom-left branch selector)
- Run **Pull** from Source Control menu

### Terminal
```sh
git checkout master
git pull
```

---

## 6.2 Create a version tag (before building!)

### VS Code
- Source Control → `…` menu → **Tags** → **Create Tag**
- Enter version: `vX.Y.Z`
- Push tags when prompted

### Terminal
```sh
git tag -a vX.Y.Z -m "Release vX.Y.Z"
git push origin vX.Y.Z
```

---

## 6.3 Build binaries

Run build script or manual commands.

Example:
```sh
make xbuild
```

---

## 6.4 Build Docker image

```sh
make docker_xbuild
```

---

## 6.5 Push Docker image

```sh
make docker_push
```

---

## 6.6 Publish GitHub Release

1. Open **GitHub → Releases → Draft a new release**
2. Select tag `vX.Y.Z`
3. Upload all binaries
4. Add changelog
5. Publish

---

# 7. Cleanup

### VS Code
- Delete merged branches via GitHub Pull Requests panel
- Switch to `master` and pull latest changes

### Terminal
```sh
git branch -d feature/my-change
git pull
```

---

# 8. Summary

This workflow ensures:

- clean history  
- consistent versioning  
- reproducible releases  
- correct version embedding  
- predictable Docker images  
- safe development process  

Thank you for contributing!
