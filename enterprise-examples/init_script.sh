#!/bin/bash

# Name of your Go module
MODULE_NAME="go-azure-servicebus"

# Create root project directory
mkdir -p $MODULE_NAME/{cmd/{orderapi,queueworker,emailsubscriber,inventorysubscriber},internal/{messaging,models}}

# Optional: create main.go files with stubs
for dir in orderapi queueworker emailsubscriber inventorysubscriber; do
  cat <<EOF > $MODULE_NAME/cmd/$dir/main.go
package main

import "fmt"

func main() {
    fmt.Println("$dir microservice running...")
}
EOF
done

# Create placeholder for models
cat <<EOF > $MODULE_NAME/internal/models/order.go
package models

type Order struct {
    ID     string \`json:"id"\`
    Item   string \`json:"item"\`
    Amount int    \`json:"amount"\`
}
EOF

# Initialize go.mod (optional — must run inside the folder)
echo "Done. Now run:"
echo "cd $MODULE_NAME && go mod init github.com/yourusername/$MODULE_NAME"
