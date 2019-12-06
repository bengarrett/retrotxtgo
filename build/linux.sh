#!/usr/bin/env sh

# https://sohlich.github.io/post/go_makefile/ < move to this
# https://www.thepolyglotdeveloper.com/2017/04/cross-compiling-golang-applications-raspberry-pi/
# https://sahilm.com/makefiles-for-golang/
# https://oddcode.daveamit.com/2018/08/16/statically-compile-golang-binary/
# https://jfrog.com/blog/5-best-practices-for-golang-ci-cd/
# https://ops.tips/blog/minimal-golang-makefile/

# source: https://blog.alexellis.io/inject-build-time-vars-golang/
export RTG_COMMIT=$(git rev-list --abbrev-commit -1 HEAD) && \
export RTG_HEADCNT=$(git rev-list --count HEAD) && \
export RTG_DATE=$(date --utc +%H:%M:%S/%Y-%m-%d) && \
go build -o artifacts/retrotxtgo -ldflags "-X github.com/bengarrett/retrotxtgo/cmd.GoBuildGitCommit=$RTG_COMMIT \
-X github.com/bengarrett/retrotxtgo/cmd.GoBuildGitCount=$RTG_HEADCNT \
-X github.com/bengarrett/retrotxtgo/cmd.GoBuildDate=$RTG_DATE" && \
./artifacts/retrotxtgo version