#!/usr/bin/env bash

# goreleaser requires access to github using a personal access token.
# https://github.com/settings/tokens
file=~/.tokens/goreleaser.txt
if [ -f $FILE ]; then
    token=$(<$file) && export GITHUB_TOKEN="$token" &&
        goreleaser --release-notes ../docs/changelog.md --rm-dist &&
        exit 0
else
    echo "The Github token file '$file' could not be found."
    exit 1
fi
