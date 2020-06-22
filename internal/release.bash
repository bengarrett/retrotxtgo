#!/usr/bin/env bash

echo "Current version tag $(git describe)"

read -p 'New release semantic version tag? (1.x.x) ' newtag
read -p 'New release comment? ' newcmmt
echo -e "new commit version: \"$newtag\" comment: \"$newcmmt\"\nchangelog:\n$(cat ../changelog.md)\n"
read -p 'confirm? [y/N] ' confirm

case $confirm in
y | yes | ok) ;;
*)
    exit 0
    ;;
esac

# goreleaser can only be used with an active git commit.

git status
git add ../. &&
    git commit -m "$newcmmt" &&
    git tag -a $newtag -m "$newcmmt" &&
    git push origin $newtag &&
    release

# goreleaser requires access to github using a personal access token.
# https://github.com/settings/tokens

release() {
    file=~/.tokens/goreleaser.txt
    if [ -f $FILE ]; then
        token=$(<$file) && export GITHUB_TOKEN="$token" &&
            goreleaser --release-notes ../changelog.md --rm-dist &&
            exit 0
    else
        echo "The Github token file '$file' could not be found."
        exit 1
    fi
}

# goreleaser --release-notes changelog.md --rm-dist &&
#     go get github.com/bengarrett/retrotxtgo

# notes
# to delete a local tag
# git tag --delete tagname
# to delete a remote tag
# git push --delete origin tagname
