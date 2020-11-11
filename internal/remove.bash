#!/usr/bin/env bash

echo "Current version tag $(git describe)"

read -p 'Remove release semantic version tag? (1.x.x) ' deltag
read -p 'confirm? [y/N] ' confirm

case $confirm in
y | yes | ok) ;;
*)
    exit 0
    ;;
esac

git status
git tag --delete $deltag
git push --delete origin $deltag
