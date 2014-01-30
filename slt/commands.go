package slt

import (
	"github.com/spf13/cobra"
	"fmt"
)

// Package-wise verbosity
// use with:
// if Verb { ...
var Verb bool 

var SlToolsCmd = &cobra.Command{
	Use:   "sltools",
	Short: "Tools for StarLab simulation management",
	Long: `...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Choose a sub-command or type sltools help for help.")
	},
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of slt",
	Long:  `All software has versions. This is sltools'`,
	Run: func(cmd *cobra.Command, args []string) {
	fmt.Println("StarLab Tools v0.2")
	},
}

var (
	inFileName string
	fileN string
)

var Out2ICCmd = &cobra.Command{
	Use:   "continue",
	Short: "Prepare the new ICs from the last STDOUT",
	Long:  `StarLab can restart a simulation from the last complete output.
	The continue command prepare the new ICs parsing the last STDOUT and writing
	the last complete snapshot to the new input file.`,
	Run: func(cmd *cobra.Command, args []string) {
		Out2IC(inFileName, fileN)
	},
}

var CreateScriptCmd = &cobra.Command{
	Use:   "createScripts",
	Short: "Prepare the new ICs from all the last STDOUTs",
	Long:  `StarLab can restart a simulation from the last complete output.
	The continue command prepare the new ICs parsing all the last STDOUTs and writing
	the last complete snapshot to the new input file.`,
	Run: func(cmd *cobra.Command, args []string) {
		CreateScripts()
	},
}

var InstallSLCmd = &cobra.Command{
	Use:   "installSL",
	Short: "Download and install SL",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		InstallSL()
	},
}

var DryInstallSLCmd = &cobra.Command{
	Use:   "dryInstallSL",
	Short: "Only install SL",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		DryInstallSL()
	},
}

var DownloadSLCmd = &cobra.Command{
	Use:   "downloadSL",
	Short: "Only download SL",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		DownloadSL()
	},
}

var (
	inFileTmpl string
)

var StichOutputCmd = &cobra.Command{
	Use:   "downloadSL",
	Short: "Only download SL",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		StichOutput (inFileTmpl)
	},
}

func InitCommands() () {

	SlToolsCmd.AddCommand(VersionCmd)
	SlToolsCmd.Flags().BoolVarP(&Verb, "verb", "v", false, "Verbose and persistent output")
	
	SlToolsCmd.AddCommand(CreateScriptCmd)
	
	SlToolsCmd.AddCommand(Out2ICCmd)
	Out2ICCmd.Flags().StringVarP(&inFileName, "inputFile", "i", "", "Last STDOUT to be used as input")
	Out2ICCmd.Flags().StringVarP(&fileN, "fileN", "n", "", "Number to be attached to the new IC file name")
	
	SlToolsCmd.AddCommand(InstallSLCmd)
	
	SlToolsCmd.AddCommand(DryInstallSLCmd)
	
	SlToolsCmd.AddCommand(DownloadSLCmd)
	
	SlToolsCmd.AddCommand(StichOutputCmd)
	StichOutputCmd.Flags().StringVarP(&inFileTmpl, "inputFileTmpl", "i", "", 
			"STDOUT template name (the STDOUT name without the extention and the )")
	
}

