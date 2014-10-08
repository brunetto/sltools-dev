package slt

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/brunetto/goutils"
)

var err error

// Verb control the package-wise verbosity.
// Use with:
// if Verb { ...
var Verb bool
var All bool

// Debug activate the package-wise debug verbosity.
// Use with:
// if Verb { ...
var Debug bool

// ConfName is the name of the JSON configuration file.
var ConfName string

// SlToolsCmd is the main command.
var SlToolsCmd = &cobra.Command{
	Use:   "sltools",
	Short: "Tools for StarLab simulation management",
	Long: `SlTools would help in running simulations with StarLab.
It can create the inital conditions if StarLab is compiled and the 
necessary binaries are available.
SlTools can also prepare ICs from the last snapshot and stich the 
output.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Choose a sub-command or type sltools help for help.")
	},
}

// VersionCmd print the sltools version.
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of slt",
	Long:  `All software has versions. This is sltools' one.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("StarLab Tools v0.9")
	},
}

// ReadConfCmd load the JSON configuration file.
var ReadConfCmd = &cobra.Command{
	Use:   "readConf",
	Short: "Read and print the configuration file",
	Long: `Read and print the configuration specify by the -c flag.
	It must be in the form of a JSON file like:

	{
		"Runs": 10,
		"Comb": 18, 
		"Ncm" : 10000,
		"Fpb" : 0.10,
		"W"   : 5,
		"Z"   : 0.20,
		"EndTime" : 500,
		"Machine" : "plx",
		"UserName" : "bziosi00",
		"PName": "IscrC_VMStars" 
	}
	`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := InitVars(ConfName)
		if Verb {
			fmt.Println("Config:")
			conf.Print()
		}
	},
}

var (
	RunICC bool
)

// CreateICsCmd will launch the functions to create the ICs from JSON configuration file.
var CreateICsCmd = &cobra.Command{
	Use:   "createICs",
	Short: "Create ICs from the JSON config file.",
	Long: `Create initial conditions from the JSON config file.
	Use like:
	sltools createICs -c conf21.json -v -C
	sltools createICs -v -C -A # to create folders and ICs for all the config files`,
	Run: func(cmd *cobra.Command, args []string) {
		if All {
			log.Println("Create all ICs following all the .json config files in this folder")
			CreateICsWrap("all", RunICC)
		} else {
			CreateICsWrap(ConfName, RunICC)
		}
	},
}

var (
	inFileName string
)

// Out2ICsCmd creates new ICs from STDOUT to restart the simulation
var Out2ICsCmd = &cobra.Command{
	Use:   "out2ics",
	Short: "Prepare the new ICs from the last STDOUT",
	Long: `StarLab can restart a simulation from the last complete output.
	The out2ics command prepare the new ICs parsing the last STDOUT and writing
	the last complete snapshot to the new input file.
	Use like:
	sltools out2ics -i out-cineca-comb16-NCM10000-fPB005-W5-Z010-run06-rnd00.txt -n 1`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			cssInfo = make (chan map[string]string, 1)
			inFileNameChan = make (chan string, 1)
		)
		go Out2ICs(inFileNameChan, cssInfo)
		inFileNameChan <- inFileName
		close(inFileNameChan)
		<-cssInfo
	},
}

var (
	icsName      string
	randomNumber string
	simTime      string
	machine string
)

// CreateStartScriptsCmd create start scripts: kiraLaunch and PBSlaunch
var CreateStartScriptsCmd = &cobra.Command{
	Use:   "css",
	Short: "Prepare the new ICs from all the last STDOUTs",
	Long: `StarLab can restart a simulation from the last complete output.
	The createStartScripts write the necessary start scripts to start a 
	simulation from the ICs.
	Use like:
	sltools createStartScripts -i ics-cineca-comb18-NCM10000-fPB020-W5-Z010-run01-rnd00.txt -c conf.json
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			done = make (chan struct{})
			cssInfo = make (chan map[string]string, 1)
			pbsLaunchChannel = make(chan string, 100)
		)
		if machine == "" {
			if ConfName == "" {
				log.Fatal("You must provide a machine name or a valid config file")
			} else {
				conf := InitVars(ConfName)
				machine = conf.Machine
			}
		}
		// Consumes pbs file names
		go func (pbsLaunchChannel chan string) {
			for _ = range pbsLaunchChannel {
			}
		} (pbsLaunchChannel)
		go CreateStartScripts(cssInfo, machine, pbsLaunchChannel, done)
		
		if All {
			runs, runMap, mapErr := FindLastRound("*-comb*-NCM*-fPB*-W*-Z*-run*-rnd*.txt")
			log.Println("Selected to create start scripts for all the runs in the folder")
			log.Println("Found: ")
			for _, run := range runs {
				if mapErr != nil && len(runMap[run]["ics"]) == 0 {
					continue
				}
				fmt.Printf("%v\n", runMap[run]["ics"][len(runMap[run]["ics"])-1])
			}
			fmt.Println()
			// Fill the channel with the last round of each run
			for _, run := range runs {
				if mapErr != nil && len(runMap[run]["ics"]) == 0 {
					continue
				}
				cssInfo <- map[string]string{
					"remainingTime": simTime,
					"randomSeed": "",
					"newICsFileName": runMap[run]["ics"][len(runMap[run]["ics"])-1],
				}
			}

		} else {
			cssInfo <- map[string]string{
					"remainingTime": simTime,
					"randomSeed": randomNumber,
					"newICsFileName": icsName,
			}
		}
		close(cssInfo)
		<- done
		close(done)
	},
}

// Out2ICsCmd + CreateStartScriptsCmd
var ContinueCmd = &cobra.Command{
	Use:   "continue",
	Short: "Prepare the new ICs from all the last STDOUTs",
	Long: `StarLab can restart a simulation from the last complete output.
	The continue command prepare the new ICs parsing all the last STDOUTs and writing
	the last complete snapshot to the new input file. It also write the necessary 
	start scripts.
	Use like:
	sltools continue -o out-cineca-comb19-NCM10000-fPB005-W9-Z010-run08-rnd01.txt`,
	Run: func(cmd *cobra.Command, args []string) {
		if machine == "" {
			if ConfName != "" {
				conf := InitVars(ConfName)
				machine = conf.Machine
			} else {
				log.Fatal("I need to know the machine name by CLI flag or conf file.")
			}
		}
		if All {
			inFileName = "all"
		}
		Continue(inFileName, machine)
	},
}

var (
	OnlyOut  bool
	OnlyErr  bool
)

// StichOutputCmd stiches STDOUT and STDERR from different round of the same simulation
// (if you restarded your simulation). Can be run serially or in parallel on all the
// file in the folder
var StichOutputCmd = &cobra.Command{
	Use:   "stichOutput",
	Short: "Stich output, only for one simulation or for all in the folder",
	Long: `Stich STDOUT and STDERR from different round of the same simulation 
	(if you restarded your simulation). Can be run serially or in parallel on all the
	file in the folder.
	You just need to select one of the files to stich or the --all flag to stich 
	all the files in the folder accordingly to their names.
	Use like:
	sltools stichOutput -c conf19.json -i out-cineca-comb19-NCM10000-fPB005-W9-Z010-run09-rnd00.txt
	sltools stichOutput -c conf19.json -A # to stich all the outputs in the folder`,
	Run: func(cmd *cobra.Command, args []string) {
		if All {
			log.Println("Stich all!")
			StichThemAll(inFileName)
		} else {
			StichOutputSingle(inFileName)
		}
	},
}

// ***
var CacCmd = &cobra.Command{
	Use:   "cac",
	Short: "...",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		CAC()
	},
}

// ***
var endOfSimMyrString string
var CheckEndCmd = &cobra.Command{
	Use:   "checkEnd",
	Short: "...",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		var endOfSimMyr float64 = 100
		if inFileName == "" || endOfSimMyrString == "" {
		log.Fatal("Provide a STDOUT file and a time in Myr to try to find the final timestep")
	} else {
		if endOfSimMyr, err = strconv.ParseFloat(os.Args[2], 64); err != nil {
			log.Fatal(err)
		}
	}
	CheckEnd (inFileName, endOfSimMyr)
	},
}


// ***
var CheckSnapshotCmd = &cobra.Command{
	Use:   "checkSnapshot",
	Short: "...",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if inFileName == "" {
			log.Fatal("Provide a STDOUT from which to check")
		}
		CheckSnapshot(inFileName)
	},
}

// ***
var CheckStatusCmd = &cobra.Command{
	Use:   "checkStatus",
	Short: "...",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		CheckStatus()
	},
}

// ***
var ComOrbitCmd = &cobra.Command{
	Use:   "comorbit",
	Short: "...",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if inFileName == "" {
			log.Fatal("Provide a STDOUT from which to extract the center of mass coordinates for the orbit")
		}
		ComOrbit(inFileName)
	},
}

var selectedSnapshot string
// ***
var CutSimCmd = &cobra.Command {
	Use:   "cutsim",
	Short: `Shorten a give snapshot to a certain timestep
	Because I don't now how perverted names you gave to your files, 
	you need to fix the STDOUT and STDERR by your own.
	You can do this by running 
	
	cutsim out --inFile <STDOUT file> --cut <snapshot where to cut>
	cutsim err --inFile <STDERR file> --cut <snapshot where to cut>
	
	The old STDERR will be saved as STDERR.bck, check it and then delete it.
	It is YOUR responsible to provide the same snapshot name to the two subcommands
	AND I suggest you to cut the simulation few timestep before it stalled.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Choose a sub-command or type restartFromHere help for help.")
	},	
}

var stdOutCutCmd = &cobra.Command {
	Use:   "out",
	Short: "cut STDOUT",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		CutStdOut(inFileName, selectedSnapshot)
	},	
}

var stdErrCutCmd = &cobra.Command {
	Use:   "err",
	Short: "cut STDERR",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		CutStdErr(inFileName, selectedSnapshot)
	},	
}	



var (
	noGPU, tf, noBinaries bool
	icsFileName string
	intTime string
// 	randomNumber string // already present
)
var KiraWrapCmd = &cobra.Command{
	Use:   "kiraWrap",
	Short: "Wrapper for the kira integrator",
	Long: `Wrap the kira integrator providing
	environment monitoring.
	The "no-GPU" flag allow you to run the non GPU version 
	of kira if you installed kira-no-GPU in $HOME/bin/.
	Run with:
	
	kiraWrap (--no-GPU)
	
	You can also specify you want our modify version with Allen-Santillan 
	tidal field provided that you have that version of kira, named kiraTF in your
	~/bin/ folder. Run with 
	
	kiraWrap -f.`,
	Run: func(cmd *cobra.Command, args []string) {
		if icsFileName == "" || intTime == "" {
			log.Fatal("Provide an ICs file and the integration time.")
		}
		kiraWrap(icsFileName, intTime, randomNumber, noGPU)
	},
}

// ***
var PbsLaunchCmd = &cobra.Command{
	Use:   "pbsLaunch",
	Short: "...",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := PbsLaunch(); err != nil {
			log.Fatal(err)
		}
	},
}

// ***
var ReLaunchCmd = &cobra.Command{
	Use:   "relaunch",
	Short: "...",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Clean folder
		SimClean()
		
		if !goutils.Exists("complete") {
			// Check and continue
			CAC()
			
			// Submit: already included in CAC
	// 		if err := slt.PbsLaunch(); err != nil {
	// 			log.Fatal(err)
	// 		}
		} else {
			log.Println("'complete' file found, assume simulations are complete.")
		}
	},
}


var RestartFromHereCmd = &cobra.Command {
	Use:   "restartFromHere",
	Short: "Prepare a pp3-stalled simulation to be restarted",
	Long: `Too often StarLab stalled while integrating a binary,
	this tool let you easily restart a stalled simulation.
	Because I don't now how perverted names you gave to your files, 
	you need to fix the STDOUT and STDERR by your own.
	You can do this by running 
	
	restartFromHere out --inFile <STDOUT file> --cut <snapshot where to cut>
	restartFromHere err --inFile <STDERR file> --cut <snapshot where to cut>
	
	The old STDERR will be saved as STDERR.bck, check it and then delete it.
	It is YOUR responsible to provide the same snapshot name to the two subcommands
	AND I suggest you to cut the simulation few timestep before it stalled.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Choose a sub-command or type restartFromHere help for help.")
	},	
}

var stdOutRestartCmd = &cobra.Command {
	Use:   "out",
	Short: "Prepare a pp3-stalled stdout to restart the simulation",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		RestartStdOut(inFileName, selectedSnapshot)
	},	
}

var stdErrRestartCmd = &cobra.Command {
	Use:   "err",
	Short: "Prepare a pp3-stalled stderr so that it is synced with the stdout",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		RestartStdErr(inFileName, selectedSnapshot)
	},	
}	

// ***
var SimCleanCmd = &cobra.Command{
	Use:   "simClean",
	Short: "...",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		SimClean()
	},
}





// Init commands and attach flags
func InitCommands() {

	SlToolsCmd.AddCommand(VersionCmd)
	SlToolsCmd.AddCommand(CacCmd)
	SlToolsCmd.AddCommand(CheckEndCmd)
	SlToolsCmd.AddCommand(CheckSnapshotCmd)
	SlToolsCmd.AddCommand(CheckStatusCmd)
	SlToolsCmd.AddCommand(ComOrbitCmd)
	SlToolsCmd.AddCommand(CutSimCmd)
	SlToolsCmd.AddCommand(KiraWrapCmd)
	SlToolsCmd.AddCommand(ComOrbitCmd)
	SlToolsCmd.AddCommand(Out2ICsCmd)
	SlToolsCmd.AddCommand(PbsLaunchCmd)
	SlToolsCmd.AddCommand(ReLaunchCmd)
	SlToolsCmd.AddCommand(RestartFromHereCmd)
	SlToolsCmd.AddCommand(SimCleanCmd)
	SlToolsCmd.AddCommand(ComOrbitCmd)
// 	SlToolsCmd.AddCommand(SLRecompileCmd)
	
	CutSimCmd.AddCommand(stdOutCutCmd)
	CutSimCmd.AddCommand(stdErrCutCmd)
	
	CheckSnapshotCmd.Flags().StringVarP(&inFileName, "inFile", "i", "", "STDOUT to check")
	ComOrbitCmd.Flags().StringVarP(&inFileName, "inFile", "i", "", "STDOUT from which to extract the center of mass coordinates for the orbit")
	CheckEndCmd.Flags().StringVarP(&inFileName, "inFile", "i", "", "STDOUT from which to try to find the final timestep")
	CheckEndCmd.Flags().StringVarP(&endOfSimMyrString, "endOfSimMyr", "e", "", "Time in Myr to try to find the final timestep")
	CutSimCmd.PersistentFlags().StringVarP(&inFileName, "inFile", "i", "", "Name of the input file")
	CutSimCmd.PersistentFlags().StringVarP(&selectedSnapshot, "cut", "c", "", "At which timestep stop")
	
	KiraWrapCmd.PersistentFlags().BoolVarP(&noGPU, "no-GPU", "n", false, "Run without GPU support if kira-no-GPU installed in $HOME/bin/.")
	KiraWrapCmd.PersistentFlags().BoolVarP(&tf, "tf", "f", false, "Run TF version of kira (debug strings).")
	KiraWrapCmd.PersistentFlags().BoolVarP(&noBinaries, "no-binaries", "b", false, "Switch off binary evolution.")
	KiraWrapCmd.PersistentFlags().StringVarP(&icsFileName, "ics", "i", "", "ICs file to start with.")
	KiraWrapCmd.PersistentFlags().StringVarP(&intTime, "time", "t", "", "Number of timestep to integrate before stop the simulation.")
	KiraWrapCmd.PersistentFlags().StringVarP(&randomNumber, "random", "s", "", "Random number.")
	
	Out2ICsCmd.PersistentFlags().StringVarP(&inFileName, "inFile", "i", "", "Name of the STDOUT file to parse")
	
	RestartFromHereCmd.AddCommand(stdOutRestartCmd)
	RestartFromHereCmd.AddCommand(stdErrRestartCmd)
	
	RestartFromHereCmd.PersistentFlags().StringVarP(&inFileName, "inFile", "i", "", "Name of the input file")
	RestartFromHereCmd.PersistentFlags().StringVarP(&selectedSnapshot, "cut", "c", "", "At which timestep stop")
	
	SlToolsCmd.PersistentFlags().BoolVarP(&Verb, "verb", "v", false, "Verbose and persistent output")
	SlToolsCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Debug output")
	SlToolsCmd.PersistentFlags().StringVarP(&ConfName, "confName", "c", "", "Name of the JSON config file")
	SlToolsCmd.PersistentFlags().BoolVarP(&All, "all", "A", false, "Run command on all the relevant files in the local folder")

	SlToolsCmd.AddCommand(ReadConfCmd)

	SlToolsCmd.AddCommand(CreateICsCmd)
	CreateICsCmd.Flags().BoolVarP(&RunICC, "runIcc", "C", false, "Run the creation of the ICs instad of only create scripts")

	SlToolsCmd.AddCommand(ContinueCmd)
	ContinueCmd.Flags().StringVarP(&inFileName, "stdOut", "o", "", "Last STDOUT to be used as input")
	ContinueCmd.Flags().StringVarP(&machine, "machine", "m", "", "Machine where to run")

	SlToolsCmd.AddCommand(Out2ICsCmd)
	Out2ICsCmd.Flags().StringVarP(&inFileName, "stdOut", "o", "", "Last STDOUT to be used as input")

	SlToolsCmd.AddCommand(CreateStartScriptsCmd)
	CreateStartScriptsCmd.Flags().StringVarP(&icsName, "icsName", "i", "", "ICs file name")
	CreateStartScriptsCmd.Flags().StringVarP(&simTime, "simTime", "t", "500", "Remaining simulation time provided by the out2ics command")
	CreateStartScriptsCmd.Flags().StringVarP(&randomNumber, "random", "r", "", "Init random seed provided by the out2ics command")
	CreateStartScriptsCmd.Flags().StringVarP(&machine, "machine", "m", "", "Machine where to run")

	SlToolsCmd.AddCommand(StichOutputCmd)
	StichOutputCmd.Flags().StringVarP(&inFileName, "inFile", "i", "", "STDOUT or STDERR name to find what to stich")
	StichOutputCmd.Flags().BoolVarP(&OnlyOut, "onlyOut", "O", false, "Only stich STDOUTs")
	StichOutputCmd.Flags().BoolVarP(&OnlyErr, "onlyErr", "E", false, "Only stich STDERRs")
}
