#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

usage() {
    echo "Usage: $(basename "$0") [--dry-run]"
    echo ""
    echo "Prepare and execute a release for getignore."
    echo ""
    echo "This script:"
    echo "  1. Computes the next semantic version from conventional commits"
    echo "  2. Generates CHANGELOG.md using git-cliff"
    echo "  3. Commits the changelog update"
    echo "  4. Creates an annotated git tag"
    echo "  5. Pushes the commit and tag (triggering the CI release workflow)"
    echo ""
    echo "Options:"
    echo "  --dry-run    Show what would be released without making any changes"
    echo "  -h, --help   Show this help message"
    echo ""
    echo "Prerequisites:"
    echo "  - git-cliff: install via 'brew install git-cliff' or 'cargo install git-cliff'"
    echo "  - Clean working tree on the main branch"
}

DRY_RUN=false

while [[ $# -gt 0 ]]; do
    case "$1" in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Error: unknown argument '$1'"
            usage
            exit 1
            ;;
    esac
done

cd "${REPO_ROOT}"

# --- Prerequisites ---

if ! command -v git-cliff &>/dev/null; then
    echo "Error: git-cliff is not installed."
    echo ""
    echo "Install it via one of:"
    echo "  brew install git-cliff"
    echo "  cargo install git-cliff"
    echo "  https://github.com/orhun/git-cliff/releases"
    exit 1
fi

# Check we're on the main branch
CURRENT_BRANCH="$(git branch --show-current)"
if [[ "${CURRENT_BRANCH}" != "main" ]]; then
    echo "Error: must be on the 'main' branch to release (currently on '${CURRENT_BRANCH}')."
    exit 1
fi

# Check for clean working tree
if ! git diff --quiet || ! git diff --cached --quiet; then
    echo "Error: working tree is not clean. Commit or stash changes before releasing."
    exit 1
fi

# Check for untracked files that might matter
if [[ -n "$(git ls-files --others --exclude-standard)" ]]; then
    echo "Warning: there are untracked files in the working tree."
    echo "They will not be included in the release."
    echo ""
fi

# Compute the next version
NEXT_VERSION="$(git-cliff --bumped-version)"
if [[ -z "${NEXT_VERSION}" ]]; then
    echo "Error: git-cliff could not determine the next version."
    echo "Are there any conventional commits since the last tag?"
    exit 1
fi

# Strip v prefix for display
VERSION_NUM="${NEXT_VERSION#v}"

# Check that the tag doesn't already exist
if git rev-parse "${NEXT_VERSION}" &>/dev/null; then
    echo "Error: tag '${NEXT_VERSION}' already exists."
    exit 1
fi

echo "Next release: ${VERSION_NUM}"
echo ""

# --- Dry run ---

if [[ "${DRY_RUN}" == true ]]; then
    echo "=== DRY RUN ==="
    echo ""
    echo "Would release version ${VERSION_NUM}"
    echo ""
    echo "--- Changelog for this release ---"
    git-cliff --bump --unreleased --strip header
    echo ""
    echo "--- Actions that would be taken ---"
    echo "  1. Regenerate CHANGELOG.md"
    echo "  2. Commit: chore(release): release ${VERSION_NUM}"
    echo "  3. Tag: ${NEXT_VERSION}"
    echo "  4. Push commit and tag to origin"
    echo "  5. CI will create a GitHub Release and build artifacts"
    exit 0
fi

# --- Perform release ---

echo "Generating CHANGELOG.md..."
git-cliff --bump -o CHANGELOG.md

echo "Committing changelog..."
git add CHANGELOG.md
git commit -m "chore(release): release ${VERSION_NUM}"

echo "Creating tag ${NEXT_VERSION}..."
git tag -am "Release ${VERSION_NUM}" "${NEXT_VERSION}"

echo "Pushing to origin..."
git push --follow-tags

echo ""
echo "=== Release ${VERSION_NUM} initiated ==="
echo ""
echo "The tag push will trigger the CI release workflow, which will:"
echo "  1. Create a GitHub Release with changelog notes"
echo "  2. Build and upload binaries for all platforms"
echo ""
echo "Monitor progress at: https://github.com/gotgenes/getignore/actions"
