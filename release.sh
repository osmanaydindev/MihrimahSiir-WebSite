#! /bin/bash

if [[ "$2" =~ ^(backend|frontend)$ ]]; then
  TAG_PREFIX=$2
else
  echo "Unknown tag. Candidates: backend|frontend"
  exit 1
fi

# https://stackoverflow.com/a/27332476/7485480
#get highest tag number
VERSION=$(git tag --sort=committerdate | grep -E '^'"$TAG_PREFIX"'-[0-9]' | tail -1)
[ -z "$VERSION" ] && VERSION="${TAG_PREFIX}-0.0.0"

#replace . with space so can split into an array
VERSION_BITS=(${VERSION//./ })

VNUM=(${VERSION_BITS[0]//-/ }) # başta idare vs yazıyor ya onu ayıklıyoruz
VNUM1=${VNUM[@]: -1}
VNUM2=${VERSION_BITS[1]}
VNUM3=${VERSION_BITS[2]}

if [[ $1 == "major" ]]; then
  VNUM1=$((VNUM1 + 1))
  VNUM2=0
  VNUM3=0

elif [[ $1 == "minor" ]]; then
  VNUM2=$((VNUM2 + 1))
  VNUM3=0

elif [[ $1 == 'patch' ]]; then
  VNUM3=$((VNUM3 + 1))

else
  echo "No version type or incorrect type specified, [major, minor, patch]"
  exit 1
fi

NEW_TAG="${TAG_PREFIX}-$VNUM1.$VNUM2.$VNUM3"

read -p "$VERSION -> $NEW_TAG .Are you sure? [Y/n] " -n 1 -r
echo
if [[ ! $REPLY =~ ^[^Yy]$ ]]; then
    git tag "$NEW_TAG"
    git push --tags
fi
