#!/bin/bash

MODULE_NAME="github.com/go-squad-5/go-init"

# Check if the user provided a new module name as an argument
NEW_MODULE_NAME="$1"

if [ -z "$NEW_MODULE_NAME" ]; then
  echo "Usage: $0 <new-module-name>"
  exit 1
fi

# Check if the user passed another argument as old module name to replace
if [ -z "$2" ]; then
  OLD_MODULE_NAME="$MODULE_NAME"
else
  OLD_MODULE_NAME="$2"
fi

# Change the module name in go.mod
go mod edit -module "$NEW_MODULE_NAME"

# Change the imports in the `main.go` file to match the new module name.
find . -type f -name "*.go" -exec sed -i "s|$OLD_MODULE_NAME|$NEW_MODULE_NAME|g" {} +
