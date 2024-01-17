/*
Copyright © 2024 zvirgilx
*/
package main

import (
	"github.com/zvirgilx/searxng-go/kernel/cmd"
	_ "github.com/zvirgilx/searxng-go/kernel/internal/engines"
)

func main() {
	cmd.Execute()
}
