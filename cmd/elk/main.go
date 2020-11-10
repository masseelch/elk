package main

import (
	"fmt"
	"github.com/masseelch/elk/pkg/gen"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use:   "elk",
		Short: "A wrapper around facebook/ent to add code generation features",
	}
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "generate code for your defined schemas",
		Run: func(cmd *cobra.Command, args []string) {
			s, err := cmd.Flags().GetString("source")
			if err != nil {
				log.Fatal(err)
			}

			t, err := cmd.Flags().GetString("target")
			if err != nil {
				log.Fatal(err)
			}

			if err := gen.Generate(s, t); err != nil {
				log.Fatal(err)
			}
		},
	}
	handlerCmd = &cobra.Command{
		Use:   "handler",
		Short: "generate api handlers for your defined schemas",
		Run: func(cmd *cobra.Command, args []string) {
			s, err := cmd.Flags().GetString("source")
			if err != nil {
				log.Fatal(err)
			}
			t, err := cmd.Flags().GetString("target")
			if err != nil {
				log.Fatal(err)
			}

			if err := gen.Handler(s, t); err != nil {
				log.Fatal(err)
			}
		},
	}
	flutterCmd = &cobra.Command{
		Use:   "flutter",
		Short: "A brief description of your command",
		Run: func(cmd *cobra.Command, args []string) {
			s, err := cmd.Flags().GetString("source")
			if err != nil {
				log.Fatal(err)
			}

			t, err := cmd.Flags().GetString("target")
			if err != nil {
				log.Fatal(err)
			}

			if err := gen.Flutter(s, t); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(handlerCmd, flutterCmd)

	// Persistent flags for the generate command.
	generateCmd.PersistentFlags().StringP("source", "s", "./ent/schema", "path/to/schema/definitions")
	generateCmd.PersistentFlags().StringP("target", "t", "", "path/to/target/dir (default: dir one above source)")
}

func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
