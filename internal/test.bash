#!/usr/bin/env bash

goreleaser --release-notes ../changelog.md --rm-dist --skip-publish --skip-validate
