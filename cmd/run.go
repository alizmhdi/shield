package main

import (
	"fmt"

	"github.com/alizmhdi/shield/config"
	"github.com/alizmhdi/shield/internal/core"
	"github.com/alizmhdi/shield/internal/k8s"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the shield tool",
	RunE:  run,
}

func init() {
	runCmd.Flags().StringP("config", "c", "config.yaml", "Path to config file")
	rootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) error {
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("could not get config flag: %w", err)
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	k8sClient, err := k8s.NewClient()
	if err != nil {
		return fmt.Errorf("failed to connect to Kubernetes: %w", err)
	}
	annotator := core.NewAnnotator(k8sClient, cfg)
	if err := annotator.ApplyAnnotations(); err != nil {
		return fmt.Errorf("failed to apply annotations: %w", err)
	}
	fmt.Println("All annotations applied successfully.")
	return nil
}

