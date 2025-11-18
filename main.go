/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/the-yex/gvm/cmd"
	_ "github.com/the-yex/gvm/internal/tui/list"
	_ "github.com/the-yex/gvm/internal/tui/progress"
	_ "github.com/the-yex/gvm/internal/tui/spinner"
	_ "github.com/the-yex/gvm/pkg"
)

func main() {
	cmd.Execute()
}
