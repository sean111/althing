package cmd

import (
	"fmt"
	"os"
	"sean111/althing/internal/council"
	"sean111/althing/internal/formatting"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "althing",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		initConfig()
		council.Init()
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		prompt, _ := cmd.Flags().GetString("prompt")
		promptFile, _ := cmd.Flags().GetString("file")

		if prompt == "" {
			if promptFile != "" {
				temp, err := os.ReadFile(promptFile)
				if err != nil {
					panic(err)
				}
				prompt = string(temp)
			} else {
				fmt.Println(formatting.ErrorStyle.Render("No prompt provided"))
				os.Exit(1)
			}
		}
		council.Run(prompt)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("prompt", "p", "", "prompt to use")
	rootCmd.Flags().StringP("file", "f", "", "file to use for prompt")
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME./.config/althing/")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("althing")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

}
