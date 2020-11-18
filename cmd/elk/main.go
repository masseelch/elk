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
	cfgFile string

	rootCmd = &cobra.Command{
		Use:     "elk",
		Short:   "A wrapper around facebook/ent to add code generation features",
		Version: "0.0.5",
	}
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "generate code for your defined schemas (client, handler and flutter)",
		Run: func(cmd *cobra.Command, args []string) {
			entCmd.Run(cmd, args)
			handlerCmd.Run(cmd, args)
			flutterCmd.Run(cmd, args)
		},
	}
	entCmd = &cobra.Command{
		Use:   "ent",
		Short: "generate db client code for your defined schemas",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := &gen.Config{}

			if err := viper.UnmarshalKey("ent", cfg); err != nil {
				log.Fatal(err)
			}

			if err := gen.Generate(cfg); err != nil {
				log.Fatal(err)
			}
		},
	}
	handlerCmd = &cobra.Command{
		Use:   "handler",
		Short: "generate api handlers for your defined schemas",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := &gen.Config{}

			if err := viper.UnmarshalKey("handler", cfg); err != nil {
				log.Fatal(err)
			}

			if err := gen.Handler(cfg); err != nil {
				log.Fatal(err)
			}
		},
	}
	flutterCmd = &cobra.Command{
		Use:   "flutter",
		Short: "A brief description of your command",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := &gen.FlutterConfig{}

			if err := viper.UnmarshalKey("flutter", cfg); err != nil {
				log.Fatal(err)
			}

			if err := gen.Flutter(cfg); err != nil {
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
	generateCmd.AddCommand(entCmd, handlerCmd, flutterCmd)

	// Allow setting flags by config file.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "/path/to/config.yaml")

	registerSourceAndTargetFlag(entCmd, handlerCmd, flutterCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func registerSourceAndTargetFlag(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		cmd.Flags().StringP("source", "s", "./ent/schema", "path/to/schema/definitions")
		cmd.Flags().StringP("target", "t", "", "path/to/target/dir (default: dir one above source)")
	}
}
