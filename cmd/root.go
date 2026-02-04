// Package cmd 命令行工具
//
//	@update 2024-08-11 01:57:31
package cmd

import (
	"os"

	"github.com/hcd233/aris-mem-api/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Aris Memory API",
	Long:  `Aris Memory API`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logger.Logger().Error("[Command] failed to execute command", zap.Error(err))
		os.Exit(1)
	}
}
