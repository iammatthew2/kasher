#!/bin/bash
set -e

mkdir -p dist

platforms=(
  "darwin amd64"
  "darwin arm64"
  "linux amd64"
  "windows amd64"
)

for platform in "${platforms[@]}"; do
  read -r GOOS GOARCH <<< "$platform"

  ext=""
  [ "$GOOS" = "windows" ] && ext=".exe"

  binname="kasher$ext"
  zipname="kasher_${GOOS}_${GOARCH}.zip"

  echo "Building for $GOOS/$GOARCH..."

  GOOS=$GOOS GOARCH=$GOARCH go build -o "dist/$binname" .

  (cd dist && zip -q "$zipname" "$binname")
  rm "dist/$binname"
done

echo "All kasher binaries zipped and stored in ./dist"
