/*
Copyright © 2024 zvirgilx
*/
package cmd

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zvirgilx/searxng-go/kernel/config"
	"github.com/zvirgilx/searxng-go/kernel/internal/complete"
	"github.com/zvirgilx/searxng-go/kernel/internal/engines"
	"github.com/zvirgilx/searxng-go/kernel/internal/engines/traits"
	"github.com/zvirgilx/searxng-go/kernel/internal/result"
)

var loglevel string
var configFile string

var rootCmd = &cobra.Command{
	Use:   "searxng-go",
	Short: "A metasearch engine written in Go",
	Long:  "Searxng-go is a metasearch engine written in Go, inspired by searxng.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&loglevel, "loglevel", "l", "info", "log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "service config file")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	initLog()

	if err := config.InitConfig(configFile); err != nil {
		panic(err)
	}

	complete.InitCompleters(config.Conf.Complete)

	if err := traits.InitTraits(); err != nil {
		panic(err)
	}

	result.InitConfig(config.Conf.Result)

	engines.InitConfiguration(config.Conf.Engines)
}

func initLog() {
	var l slog.Level
	switch strings.ToUpper(loglevel) {
	case "DEBUG":
		l = slog.LevelDebug
	case "INFO":
		l = slog.LevelInfo
	case "WARN":
		l = slog.LevelWarn
	case "ERROR":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: l})
	slog.SetDefault(slog.New(h))
}
