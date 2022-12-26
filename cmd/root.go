package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/cryptodeal/gen-napi/config"
	"github.com/cryptodeal/gen-napi/napi"
	"github.com/spf13/cobra"
)

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "gen-napi",
		Short: "Tool for generating Node_API wrappers from C++ Header Defs",
		Long:  `Gen-Napi generates Node_API by parsing C++ Header files.`,
	}

	rootCmd.PersistentFlags().String("config", "gen_napi.yaml", "config file to load (default is gen_napi.yaml in the current folder)")
	rootCmd.Version = Version() + " " + Target() + " (" + CommitDate() + ") " + Commit()
	rootCmd.PersistentFlags().BoolP("debug", "D", false, "Debug mode (prints debug messages)")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "generate",
		Short: "Generate and write to disk",
		Run:   generate,
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generate(cmd *cobra.Command, args []string) {
	cfgFilepath, err := cmd.Flags().GetString("config")
	if err != nil {
		log.Fatal(err)
	}
	napiConfig := config.ReadFromFilepath(cfgFilepath)
	t := napi.New(&napiConfig)

	err = t.Generate()
	if err != nil {
		log.Fatalf("Gen-NAPI failed: %v", err)
	}
}
