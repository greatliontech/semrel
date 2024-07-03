#!/bin/sh

set -e
rm -rf completions
mkdir completions

if [ -z "$1" ]; then
	echo "Usage: $0 <binary-name>"
	exit 1
fi

for sh in bash zsh fish; do
	go run main.go completion "$sh" >"completions/${1}.${sh}"
done
