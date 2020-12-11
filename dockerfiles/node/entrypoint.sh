#!/bin/bash

/gev --datadir ./data init ./genesis.json

exec "$@"
