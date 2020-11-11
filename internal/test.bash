#!/usr/bin/env bash

goreleaser --release-notes ../docs/changelog.md --rm-dist --skip-publish --skip-validate
