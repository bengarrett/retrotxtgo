#!/usr/bin/env bash

echo "Current version tag $(git describe)"

read -p 'New release semantic version tag? (v1.x.x) ' newtag
read -p 'New release comment? ' newcmmt
echo -e "new commit version: \"$newtag\" comment: \"$newcmmt\"\nchangelog:\n$(cat ../changelog.md)\n"
read -p 'confirm? [y/N] ' confirm

case $confirm in
y | yes | ok) ;;
*)
    exit 0
    ;;
esac

# goreleaser can only be used with an active git commit
# it MUST run before go generate
echo "$newtag" >.version

# generate any updated templates before using git
go generate ./...

git status

git add --verbose ../. &&
    git commit -m "$newcmmt" &&
    git tag -a $newtag -m "$newcmmt" &&
    git push origin $newtag &&
    git pull &&
    echo "You can now run ./deploy.bash"

# notes
# to delete a local tag
# git tag --delete tagname
# to delete a remote tag
# git push --delete origin tagname
