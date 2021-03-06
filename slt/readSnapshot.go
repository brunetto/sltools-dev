package slt

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/brunetto/goutils/debug"

	"github.com/brunetto/goutils/readfile"
)

// DumbSnapshot contains one snapshot without knowing anything about it
type DumbSnapshot struct {
	Timestep     string
	Integrity    bool
	CheckRoot bool
	NestingLevel int
	MaxNesting int
	Lines        []string
}

// WriteSnapshot pick the snapshot line by line and write it to the
// output file
func (snap *DumbSnapshot) WriteSnapshot(nWriter *bufio.Writer) (err error) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	for _, line := range snap.Lines {
		_, err = nWriter.WriteString(line + "\n")
// 		if err = nWriter.Flush(); err != nil {
// 			log.Fatal(err)
// 		}
	}
	if err = nWriter.Flush(); err != nil {
		log.Fatal(err)
	}
	return err
}

// ReadOutSnapshot read one and only one snapshot at a time
func ReadOutSnapshot(nReader *bufio.Reader) (*DumbSnapshot, error) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	var (
		snap       *DumbSnapshot = new(DumbSnapshot)
		line       string
		err        error
		regSysTime = regexp.MustCompile(`system_time\s*=\s*(\d+)`)
		resSysTime []string
		cumulativeNesting int = 0
	)

	// Init snapshot container
	snap.Lines = make([]string, 0)
	snap.Integrity = false
	snap.CheckRoot = false
	snap.NestingLevel = 0
	

	for {
		// Read line by line
		if line, err = readfile.Readln(nReader); err != nil {
			// Mark snapshot as corrupted
			snap.Integrity = false
			return snap, err
		}

		// Add line to the snapshots in memory
		snap.Lines = append(snap.Lines, line)

		// Search for timestep number
		if resSysTime = regSysTime.FindStringSubmatch(line); resSysTime != nil {
			snap.Timestep = resSysTime[1]
		}

		// Check if entering or exiting a particle
		// and update the nesting level
		if strings.Contains(line, "(Particle") {
			snap.NestingLevel++
			cumulativeNesting++
		} else if strings.Contains(line, ")Particle") {
			snap.NestingLevel--
		}
		
		// Doesn't work with ICs because they are without name = root grrrr
		if strings.Contains(line, "name = root") {
			snap.CheckRoot = true
		}

		// Check whether the whole snapshot is in memory
		// (root particle complete) and if true, return
		// We need cumulative nesting in case of a header.
		// Without it the result is always a bad snapshot because 
		// the nesting is 0 and no root is found but only
		// because we are not into the particles section
		if snap.NestingLevel == 0 && cumulativeNesting != 0 {
			if !snap.CheckRoot {
				outFile, err := os.Create("badTimestep.txt")
				defer outFile.Close()
				if err != nil {log.Fatal(err)}
				nWriter := bufio.NewWriter(outFile)
				defer nWriter.Flush()
				snap.WriteSnapshot(nWriter)
				fmt.Println()
				log.Fatal("No root particle in a timestep that seems complete, please check!")
			}
			snap.Integrity = true
			if Verb {
				log.Println("Timestep ", snap.Timestep, " integrity set to: ", snap.Integrity)
			} else {
				fmt.Fprintf(os.Stderr, "\r\tTimestep %v integrity set to: %v", snap.Timestep, snap.Integrity)
			}
			return snap, err
		}
	}
}

// ReadErrSnapshot read one and only one snapshot at a time
func ReadErrSnapshot(nReader *bufio.Reader) (*DumbSnapshot, error) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	var (
		snap       *DumbSnapshot = new(DumbSnapshot)
		line       string
		err        error
		regSysTime = regexp.MustCompile(`^Time = (\d+)`)
		resSysTime []string
		endOfSnap  string = "----------------------------------------"
		// This variables are the idxs to print the last or last 10 lines
		dataStartIdx int = 0
		dataEndIdx   int
	)

	// Init snapshot container
	snap.Lines = make([]string, 0) //FIXME: check if 0 is ok!!!
	snap.Integrity = false
	snap.CheckRoot = false
	snap.Timestep = "-1"

	for {
		// Read line by line
		if line, err = readfile.Readln(nReader); err != nil {
			if err.Error() == "EOF" {
				if Verb {
					fmt.Println()
					log.Println("File reading complete...")
					log.Println("Timestep not complete.")
					log.Println("Last ten lines:")
					dataEndIdx = len(snap.Lines) - 1

					// Check that we have more than 10 lines
					if dataEndIdx > 10 {
						dataStartIdx = dataEndIdx - 10
					}
					for idx, row := range snap.Lines[dataStartIdx:dataEndIdx] {
						fmt.Println(idx, ": ", row)
					}
				}
			} else {
				log.Fatal("Non EOF error while reading ", err)
			}
			// Mark snapshot as corrupted
			snap.Integrity = false
			return snap, err
		}

		// Add line to the snapshots in memory
		snap.Lines = append(snap.Lines, line)

		// Search for timestep number
		if resSysTime = regSysTime.FindStringSubmatch(line); resSysTime != nil {
			snap.Timestep = resSysTime[1]
		}

		// Check if entering or exiting a particle
		// and update the nesting level
		if strings.Contains(line, endOfSnap) {
			snap.Integrity = true
			if Verb {
				log.Println("Timestep ", snap.Timestep, " integrity set to: ", snap.Integrity)
			} else {
				fmt.Fprintf(os.Stderr, "\r\tTimestep %v integrity set to: %v", snap.Timestep, snap.Integrity)
			}
			return snap, err
		}
	}
}
