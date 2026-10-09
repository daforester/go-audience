#!/bin/sh
# Prints the Fantasy backend version, taken from the git tags on the current commit:
#
#   commit tagged fantasy-backend/v0.0.1        0.0.1
#   3 commits after that tag                    0.0.1-3-gabc1234
#   uncommitted changes to the backend          ...+dirty
#   no fantasy-backend tag in the history yet   0.0.0-dev.abc1234
#
# Used by CI and for local image builds:
#   docker build --build-arg VERSION="$(apps/fantasy/scripts/version.sh)" -f Dockerfile-fantasy .
set -e
cd "$(dirname "$0")/.."

sha=$(git rev-parse --short=7 HEAD)
if described=$(git describe --tags --match 'fantasy-backend/v[0-9]*' --abbrev=7 2>/dev/null); then
  version=${described#fantasy-backend/v}
else
  version="0.0.0-dev.$sha"
fi

if ! git diff --quiet HEAD -- . 2>/dev/null; then
  version="$version+dirty"
fi

echo "$version"
