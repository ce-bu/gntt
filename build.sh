#!/bin/bash

pushd $(dirname $0) > /dev/null
S=$(pwd)
popd > /dev/null

B="$S/build"

sudo rm -rf "$B" && mkdir -p "$B"
(cd "$S" && autoreconf -i)
cd $B
$S/configure --prefix=$B/install
make install
