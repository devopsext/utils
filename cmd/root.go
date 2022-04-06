package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"syscall"

	"github.com/spf13/cobra"
)

var VERSION = "unknown"

func interceptSyscall() {

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-c
		os.Exit(1)
	}()
}

func Execute() {

	rootCmd := &cobra.Command{
		Use:   "utils",
		Short: "Utils",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

			//
		},
		Run: func(cmd *cobra.Command, args []string) {

			//
		},
	}

	//flags := rootCmd.PersistentFlags()
	//flags.StringVar(&stdoutOptions.Format, "stdout-format", stdoutOptions.Format, "Stdout format: json, text, template")

	interceptSyscall()

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(VERSION)
		},
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
