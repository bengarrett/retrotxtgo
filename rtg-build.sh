#!/usr/bin/env sh

# source: https://blog.alexellis.io/inject-build-time-vars-golang/
export RTG_COMMIT=$(git rev-list --abbrev-commit -1 HEAD) && \
export RTG_HEADCNT=$(git rev-list --count HEAD) && \
export RTG_DATE=$(date --utc +%H:%M:%S/%Y-%m-%d) && \
go build -ldflags "-X github.com/bengarrett/retrotxtgo/cmd.GoBuildGitCommit=$RTG_COMMIT \
-X github.com/bengarrett/retrotxtgo/cmd.GoBuildGitCount=$RTG_HEADCNT \
-X github.com/bengarrett/retrotxtgo/cmd.GoBuildDate=$RTG_DATE" && \
./retrotxtgo version