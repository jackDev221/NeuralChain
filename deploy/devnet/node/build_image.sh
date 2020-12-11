#!/bin/sh
echo "------ Buiding Gev-Node docker image ------"
BASEDIR=$(dirname "$0")
git clone --single-branch -b develop --single-branch git@github.com:Evrynetlabs/evrynet-node.git "$BASEDIR"/project

cp "$BASEDIR"/genesis.json "$BASEDIR"/project/deploy/devnet/node/genesis.json
cp ./dockerfiles/node/token "$BASEDIR"/project/dockerfiles/node/token
cp ./dockerfiles/node/Dockerfile "$BASEDIR"/project/dockerfiles/node/Dockerfile

pushd "$BASEDIR"/project
docker build -f dockerfiles/node/Dockerfile -t kybernetwork/evrynet-node:1.0.1-dev .
popd

rm -rf "$BASEDIR"/project
