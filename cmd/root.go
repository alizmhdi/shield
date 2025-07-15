package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "shield",
	Short: "Annotate Kubernetes Ingresses with whitelist IPs",
	Run: nil,
}
