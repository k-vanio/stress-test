package main

import (
	"log"

	"github.com/k-vanio/stress-test/internal/stress"
	"github.com/spf13/cobra"
)

func main() {
	s := stress.New()

	rootCmd := &cobra.Command{
		Use:   "stress",
		Short: "performs load testing on a web service",
		Run:   s.Run,
	}

	rootCmd.Flags().Args()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
