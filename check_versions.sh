#!/bin/bash

# Checkout the branch from which the pull request is being merged
if [ -n "$GITHUB_HEAD_REF" ]; then
  git fetch origin $GITHUB_HEAD_REF
  git checkout $GITHUB_HEAD_REF
fi

# Fetch all tags
git fetch --tags

VERSION=$(cat "VERSION")
CHANGELOG_VERSION=$(grep -oP '^## \K[vV]?[0-9]+\.[0-9]+\.[0-9]+' "CHANGELOG.md" | head -n 1)
README_VERSION=$(grep -oP 'v[0-9]+\.[0-9]+\.[0-9]+' "README.md" | head -n 1)
LATEST_TAG=$(git tag --sort=-v:refname | head -n 1)

RELEASE_EXISTS=$(curl -s https://api.github.com/repos/grosheth/gysmo/releases | jq -r '.[].tag_name')

# Check if the version matches the format vX.X.X
if ! [[ "$VERSION" =~ ^v[0-9]+(\.[0-9]+)*$ ]]; then
  echo "Version format is incorrect! Expected format: vX.X.X"
  echo "VERSION file: $VERSION"
  exit 1
fi

if [ "$VERSION" != "$CHANGELOG_VERSION" ] || [ "$VERSION" != "$README_VERSION" ] || [ "$VERSION" != "$LATEST_TAG" ] || [[ " $RELEASE_EXISTS " == *" $VERSION "* ]]; then
  echo "Version mismatch detected or release already exists!"
  echo "VERSION file: $VERSION"
  echo "CHANGELOG.md: $CHANGELOG_VERSION"
  echo "README.md: $README_VERSION"
  echo "Latest tag: $LATEST_TAG"
  echo "Existing release: $RELEASE_EXISTS"
  exit 1
fi

echo "All versions match and the release does not already exist: $VERSION"
