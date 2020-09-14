#! /bin/bash

NEW_VERSION=$1

## Check if there are any uncommitted changes.
if ! [ -d .git ]
then
    echo "Error: Your are not in a git repo."
    exit 1
fi

## Check if there are any uncommitted changes.
if ! [ -z "$(git status --porcelain 2&1>/dev/null)" ]
then
    echo "Error: There are uncommitted changes in your git repository."
    exit 1
fi

## Determine the last release tag.
OLD_VERSION=`git tag | tail -1 | cut -d'v' -f2`
if test "$?" != "0"
then
    echo "Error: No previous release tag found."
    exit 1
fi

## Check if version input is set.
if [ -z $NEW_VERSION ]
then
    echo -e "Usage:\n\trelease.sh VERSION"
    echo -e "Example:\n\t./release.sh 1.5"
    echo -e "\nCurrent version is '$OLD_VERSION'.\n"
    exit 1
fi

## Check if goreleaser is installed.
if ! hash -d goreleaser 2>/dev/null
then
    echo "Error: \"goreleaser\" tool is installed in your system."
    exit 1
fi

## Determine the we already made this release.
if ! [ -z "$(git tag | grep v$NEW_VERSION)" ]
then
    echo "Error: $NEW_VERSION already exists."
    exit 1
fi

## Update release string across the repo.
sed -i 's/v'$OLD_VERSION'/v'$NEW_VERSION'/g' README.md
sed -i 's/_'$OLD_VERSION'_/_'$NEW_VERSION'_/g' README.md
sed -i 's/v'$OLD_VERSION'/v'$NEW_VERSION'/g' Dockerfile.release
sed -i 's/_'$OLD_VERSION'_/_'$NEW_VERSION'_/g' Dockerfile.release
sed -i 's/= "'$OLD_VERSION'"/= "'$NEW_VERSION'"/g' cmd/globals.go
git commit -am  "Releasing version v"$NEW_VERSION
git push

## Tag the release.
git tag -a v$NEW_VERSION -m "Releasing version v"$NEW_VERSION
git push origin v$NEW_VERSION

# Use goreleaser to cross-build for all platforms and make a release.
sudo goreleaser
