package slt

import (
	"log"
	"regexp"
)

var (
	BHRegString string = `(\d{1,5}\+*\d*\+*\d*)` + // group(1) is the BH id
		`\s+\S+_to_black_hole_at_time\s*=\s*` +
		`(\d+\.*\d*)` + // group(2) is the BH formation time in units of FIXME
		`\s*` +
		`(\S+)` // group(3) is the time unit
		//`\s*\(old mass\s*=\s*(\d+\.*\d*)\)`									// group(4) is the progenitor mass (probably wrong)
	BHReg *regexp.Regexp = regexp.MustCompile(BHRegString)
	BHRes []string

	NSRegString string = `(\d{1,5}\+*\d*\+*\d*)` + // group(1) is the NS id
		`\s+\S+_to_neutron_star_at_time\s*=\s*` +
		`(\d+\.*\d*)` + // group(2) is the NS formation time in units of FIXME
		`\s*` +
		`(\S+)` // group(3) is the time unit
		//`\s*\(old mass\s*=\s*(\d+\.*\d*)\)`									// group(4) is the progenitor mass (probably wrong)
	NSReg *regexp.Regexp = regexp.MustCompile(NSRegString)
	NSRes []string

	WDRegString string = `(\d{1,5}\+*\d*\+*\d*)` + // group(1) is the WD id
		`\s+\S+_to_\S+_dwarf_at_time\s*=\s*` +
		`(\d+\.*\d*)` + // group(2) is the WD formation time in units of FIXME
		`\s*` +
		`(\S+)` // group(3) is the time unit
		//`\s*\(old mass\s*=\s*(\d+\.*\d*)\)`									// group(4) is the progenitor mass (probably wrong)
	WDReg *regexp.Regexp = regexp.MustCompile(WDRegString)
	WDRes []string

	MergerRegString string = `binary_evolution:\s*merger within\s*\(` + // group(1) are the two ids
		`(\d{1,5}\+*\d*\+*\d*,\d{1,5}\+*\d*\+*\d*)` + 
		`\)\s*triggered by \d{1,5}\+*\d*\s*at time\s*` + 
		`(\d+\.*\d*)` // group(2) is the merger time in units of FIXME
	MergerReg       *regexp.Regexp = regexp.MustCompile(MergerRegString)
	MergerRes       []string
	
	MergerTimeRegString string = `Collision at time =\s*\d+\.*\d*\s*\((\d+\.*\d*)\s*\[Myr\]\)\s*between` 
	MergerTimeReg       *regexp.Regexp = regexp.MustCompile(MergerTimeRegString)
	MergerTimeRes       []string
	
	CollisionResultString string = `\d+\s*\((\S+);\s*M\s*=\s*(\d+\.*\S*)\s*\[Msun\]` 
	CollisionResultReg       *regexp.Regexp = regexp.MustCompile(CollisionResultRegString)
	CollisionResultRes       []string
)


func Search (snapChan chan *DumbSnapshot) () {
	var (
		line string
		snap *DumbSnapshot
		lineChan = make(chan string, 1)
	)

	go SearchCO(lineChan)
	go SearchMerger(lineChan)
	go SearchBinaries(lineChan)
	
	for snap  = range snapChan {
		for _, line = snap.Lines {
			lineChan <- line
		}
	}
	
	
}



func SearchCO () () {
	

}

func SearchMerger () () {

}

func SearchBinaries () () {
	
	
	
	
}










