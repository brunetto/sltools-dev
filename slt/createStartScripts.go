package slt

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"time"

	"github.com/brunetto/goutils/debug"
)

// CreateStartScripts create the start scripts (kira launch and PBS launch for the ICs).
func CreateStartScripts(cssInfo chan map[string]string, machine string) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		infoMap        map[string]string
		err            error                                                             // common error container
		icsName        string                                                            // name of the ICs
		currentDir     string                                                            // current local directory
		stdOutFile     string                                                            // STDOUT file for the next run
		stdErrFile     string                                                            // STDOUT file for the next run
		shortName      string                                                            // id for the job
		randomString   string                                                            // random number from the STDERR
		simTime        string                                                            // remaining timesteps for the simulation (e.g. 500 - last run)
		queue          string                                                            // name of the queue on wich we will run
		comb, run, rnd string                                                            //combination, run and round number
		baseName       string                                                            // common part of the name without the extension
		walltime       string                                                            // max time we can run on the queue
		kiraString     string                                                            // string to launch kira
		pbsString      string                                                            // string to submit the job to PBS
		baseNameString string         = `(\w{3})-(\S*\).(\S*)`                           // extract baseName parts
		baseNameExp    *regexp.Regexp = regexp.MustCompile(regString)                    // compile the regexp
		baseNameResult []string                                                          // regexp result
		detailsString  string         = `\w{3}-\S*-comb(\S*)-\S*-run(\S*)-rnd(\d*)\.\S*` // extract filename parts in detail
		detailsExp     *regexp.Regexp = regexp.MustCompile(detailsString)                // compile the regexp
		detailsResult  []string                                                          // regexp result
		kiraFile       *os.File                                                          // where to save kiraString
		pbsFile        *os.File                                                          // where to save PBS string
		kiraOutName    string                                                            // kira file name
		pbsOutName     string                                                            // PBS file name
		home           string                                                            // path to home on the cluster
		kiraBinPath    string                                                            // path to kira binaries
		modules        string                                                            // modules we need to load
	)

	if machine == "eurora" {
		modules = "module purge\n" +
			"module load profile/advanced\n" +
			"module load gnu/4.6.3\n" +
			"module load boost/1.53.0--gnu--4.6.3\n" +
			"module load cuda\n\n" +
			"# # # LD_LIBRARY_PATH=$LD_LIBRARY_PATH:" +
			"/cineca/prod/compilers/cuda/5.0.35/none/lib64:" +
			"/cineca/prod/libraries/boost/1.53.0/gnu--4.6.3/lib\n" +
			"# # # LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/eurora/home/userexternal/mmapelli/\n\n" +
			"LD_LIBRARY_PATH=/cineca/prod/compilers/cuda/5.0.35/none/lib64:/cineca/prod/libraries/boost/1.53.0/gnu--4.6.3/lib\n" +
			"export LD_LIBRARY_PATH\n"
		queue = "parallel"
		walltime = "4:00:00"
		home = "/eurora/home/userexternal/bziosi00"
		kiraBinPath = "/eurora/home/userexternal/bziosi00/starlabjune19_2013/usr/bin/kira"
	} else if machine == "plx" {
		modules = "module purge\n" +
			"module load gnu/4.1.2\n" +
			"module load profile/advanced\n" +
			"module load boost/1.41.0--intel--11.1--binary\n" +
			"module load cuda/4.0\n\n" +
			"LD_LIBRARY_PATH=/cineca/prod/compilers/cuda/4.0/none/lib64:" +
			"/cineca/prod/compilers/cuda/4.0/none/lib:/cineca/prod/" +
			"libraries/boost/1.41.0/intel--11.1--binary/lib:/cineca/" +
			"prod/compilers/intel/11.1/binary/lib/intel64\n" +
			"export LD_LIBRARY_PATH\n\n"
		queue = "longpar"
		walltime = "24:00:00"
		home = "/plx/userexternal/bziosi00"
		kiraBinPath = filepath.Join(home, "slpack", "starlab", "usr", "bin", "kira")
	} else {
		log.Fatal("Uknown machine name ", machine)
	}

	for infoMap = range inFileNameChan {

		if baseNameResult = baseNameExp.FindStringSubmatch(icsName); regResult == nil {
			log.Fatal("Can't find commonName in ", icsName)
		}
		if baseNameResult[1] != "ics" {
			log.Fatal("Please specify an ICs file, found ", regResult[1])
		}

		baseName = regResult[2]
		comb = detailsResult[1]
		run = detailsResult[2]
		rnd = detailsResult[3]

		shortName = "r" + comb + "-" + run + "-" + rnd

		if currentDir, err = os.Getwd(); err != nil {
			log.Fatal("Can't find current working folder!!")
		}
		if currentDir, err = filepath.Abs(currentDir); err != nil {
			log.Fatal("Can't find absolute path to current working folder!!")
		}

		stdOutFile = "out-" + baseName + ".txt"
		stdErrFile = "err-" + baseName + ".txt"
		kiraOutName = "kiraLaunch-" + baseName + ".sh"
		pbsOutName = "PBS-" + baseName + ".sh"

		kiraString = "echo $PWD\n" +
			"echo $LD_LIBRARY_PATH\n" +
			"echo $HOSTNAME\n" +
			kiraBinPath + " -t " + simTime + " -d 1 -D 1 -b 1 -f 0 \\\n" +
			" -n 10 -e 0.000 -B " + randomString + " \\\n" +
			"<  " + filepath.Join(currentDir, icsName) + " \\\n" +
			">  " + filepath.Join(currentDir, stdOutFile) + " \\\n" +
			"2> " + filepath.Join(currentDir, stdErrFile) + " \n"

		pbsString = "#!/bin/bash\n" +
			"#PBS -N r" + shortName + "\n" +
			"#PBS -A IscrC_VMStars\n" +
			"#PBS -q " + queue + "\n" +
			"#PBS -l walltime=" + walltime + "\n" +
			"#PBS -l select=1:ncpus=1:ngpus=2\n\n" +
			modules +
			"sh " + filepath.Join(currentDir, kiraOutName) + "\n"

		if kiraFile, err = os.Create(kiraOutName); err != nil {
			log.Fatal(err)
		}
		defer kiraFile.Close()
		fmt.Fprint(kiraFile, kiraString)

		if pbsFile, err = os.Create(pbsOutName); err != nil {
			log.Fatal(err)
		}
		defer pbsFile.Close()
		fmt.Fprint(pbsFile, pbsString)

	}
	// 	done <- struct{}{}
}
