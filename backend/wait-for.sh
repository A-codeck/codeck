#!/bin/sh
# wait-for.sh

set -e

host="$1"
shift
cmd="$@"

until nc -z $host; do
  >&2 echo "Waiting for $host to be available..."
  sleep 1
done

>&2 echo "$host is up - executing command"
exec $cmd
