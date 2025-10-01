/*
Copyright © 2025 Farye Nwede <farye@aeekay.com>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/faryeyay/insights/primitives/pkg/benchmark"
	"github.com/faryeyay/insights/primitives/pkg/datastructures/testbench"
)

// benchmarkCmd represents the benchmark command
var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Running benchmark...")
		log.Printf("Creating benchmark tracker...")
		benchmark := benchmark.New()
		log.Printf("Starting benchmark...")
		benchmark.Start()
		log.Printf("Running benchmark...")
		err := testbench.TestBPlus()
		if err != nil {
			log.Fatalf("error running the benchmark for b+ trees: %s", err)
		}
		log.Printf("Stopping benchmark...")
		benchmark.Stop()
		log.Printf("Report: %s", benchmark.Report())
		log.Printf("Benchmark complete.")
	},
}

func init() {
	rootCmd.AddCommand(benchmarkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// benchmarkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// benchmarkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
