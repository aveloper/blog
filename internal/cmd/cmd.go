package cmd

import (
	"github.com/spf13/cobra"
)

// TODO: Improve all texts

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "blog",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

//versionCmd prints the version information
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of blog",
	Long:  `All software has versions. This is Blog's`,
	Run:   version,
}

//installCmd installs the blog
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "TODO",
	Long:  `TODO`,
	Run:   install,
}

//testCmd tests the settings and config to find out any possible errors
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Tests that all the configs are in order",
	Long:  `TODO`,
	Run:   test,
}

//updateCmd updates the blog to the latest version
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the blog",
	Long:  `TODO`,
	Run:   update,
}

//startCmd starts the blog
var startCmd = &cobra.Command{
	Use:    "start",
	Short:  "start the blog",
	Long:   `TODO`,
	Hidden: true,
	Run:    start,
}

//resetCmd to delete the data from database
var resetCmd = &cobra.Command{
	Use:    "reset",
	Short:  "delete the data from database",
	Long:   `TODO`,
	Hidden: true,
	Run:    reset,
}

//addCommands adds all the commands to the root command
func addCommands() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(resetCmd)
}

//addFlags adds all the flags for the commands
func addFlags() {
	installCmd.Flags().BoolP("force", "f", false, "reset config and force install")
	startCmd.Flags().Bool("api-only", false, "Start only the API server")
}
