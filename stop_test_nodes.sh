#!/bin/sh
echo "--- Stop all test nodes ..."
# shellcheck disable=SC2046
kill -INT $(pgrep -f ./gev)
for i in 1 2 3 4
do
  lsof -ti:2200"$i" | xargs kill
done
