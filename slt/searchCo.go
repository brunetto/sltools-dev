package slt

import (
	"log"
	"regexp"
)

var (
	// BH regexp
	BHRegString string = `(\d{1,5}\+*\d*\+*\d*)` + // group(1) is the BH id
		`\s+\S+_to_black_hole_at_time\s*=\s*` +
		`(\d+\.*\d*)` + // group(2) is the BH formation time in units of FIXME
		`\s*` +
		`(\S+)` // group(3) is the time unit
		//`\s*\(old mass\s*=\s*(\d+\.*\d*)\)`									// group(4) is the progenitor mass (probably wrong)
	BHReg *regexp.Regexp = regexp.MustCompile(BHRegString)
	BHRes []string

	// NS regexp
	NSRegString string = `(\d{1,5}\+*\d*\+*\d*)` + // group(1) is the NS id
		`\s+\S+_to_neutron_star_at_time\s*=\s*` +
		`(\d+\.*\d*)` + // group(2) is the NS formation time in units of FIXME
		`\s*` +
		`(\S+)` // group(3) is the time unit
		//`\s*\(old mass\s*=\s*(\d+\.*\d*)\)`									// group(4) is the progenitor mass (probably wrong)
	NSReg *regexp.Regexp = regexp.MustCompile(NSRegString)
	NSRes []string

	// WD regexp
	WDRegString string = `(\d{1,5}\+*\d*\+*\d*)` + // group(1) is the WD id
		`\s+\S+_to_\S+_dwarf_at_time\s*=\s*` +
		`(\d+\.*\d*)` + // group(2) is the WD formation time in units of FIXME
		`\s*` +
		`(\S+)` // group(3) is the time unit
		//`\s*\(old mass\s*=\s*(\d+\.*\d*)\)`									// group(4) is the progenitor mass (probably wrong)
	WDReg *regexp.Regexp = regexp.MustCompile(WDRegString)
	WDRes []string

	// Merger regexp
	MergerRegString string = `binary_evolution:\s*merger within\s*\(` + // group(1) are the two ids
		`(\d{1,5}\+*\d*\+*\d*,\d{1,5}\+*\d*\+*\d*)` +
		`\)\s*triggered by \d{1,5}\+*\d*\s*at time\s*` +
		`(\d+\.*\d*)` // group(2) is the merger time in units of FIXME
	MergerReg *regexp.Regexp = regexp.MustCompile(MergerRegString)
	MergerRes []string

	// Merger time regexp
	MergerTimeRegString string         = `Collision at time =\s*\d+\.*\d*\s*\((\d+\.*\d*)\s*\[Myr\]\)\s*between`
	MergerTimeReg       *regexp.Regexp = regexp.MustCompile(MergerTimeRegString)
	MergerTimeRes       []string

	// Collision check regexp
	CollisionResultRegString string         = `\d+\s*\((\S+);\s*M\s*=\s*(\d+\.*\S*)\s*\[Msun\]`
	CollisionResultReg    *regexp.Regexp = regexp.MustCompile(CollisionResultRegString)
	CollisionResultRes    []string
	
	// Binary ids check regexp
	BinIdRegString string         = `^\s*(U*)\s*(\(\S+,\S+\)):*\s+a\s=\s(\d+\.*\S*)\s+e\s=\s(\d+\.*\S*)\s+P\s=\s(\d+\.*\S*)`
	BinIdReg    *regexp.Regexp = regexp.MustCompile(CollisionResultRegString)
	BinIdRes    []string

)

func Search(snapChan chan *DumbSnapshot) {
	var (
		lineIn, lineOut string
		snap            *DumbSnapshot
		lineChan        = make(chan string, 1)
	)

	// Read lines and send them to functions with goroutines and channels
	// so if a func needs to read other lines it is possible

	// Read snap lines
	go func() {
		for snap = range snapChan {
			for _, lineIn = range snap.Lines {
				lineChan <- line
			}
		}
	}()

	// Riceive lines and search for CO or mergers or binaries
	go func() {
		for lineOut = range lineChan {
			if SearchMerger(lineOut, lineChan) {
				continue
			} else if SearchCO(lineOut) {
				continue 
			} else if SearchBinaries(lineOutline, lineChan) {
				continue
			}
		}
	}()

}

func SearchMerger(line string, lineChan chan string) (bool) {
	if MergerRes = MergerReg.FindStringSubmatch(line); MergerRes == nil {
		return false
	}
	for line = range lineChan {
		if MergerTimeRes = MergerTimeReg.FindStringSubmatch(line); MergerTimeRes != nil { 
			// FIXME: fill
		} else if CollisionResultRes = CollisionResultReg.FindStringSubmatch(line); CollisionResultRes != nil { 
			// FIXME: fill
		}
	}
	mergerId = strings.Replace(MergerRes[1], ",", "+")
	mergerTime = MergerTimeRes[1]
	collisionResult = CollisionResultRes[1]
	collisionMass = CollisionResultRes[2]
	//FIXME check if collision result is CO in the CO list
	// qui in pratica cerchiamo di capire di che tipo sono i due oggetti coinvolti e 
	// 	di conseguenza che tipo di merger e`:
	// 	* BH + BH
	// 	* NS + NS
	// 	* WD +WD
	// 	* BH + NS
	// 	* BH + WD
	// 	* NS + WD
	// 	* BH + star
	// 	* NS + star
	// 	* WD + star
	// ----
	// Send merger to merger list/channel/whatever
	// send result to appropriate destination (star, WD, NS, BH)
}


func SearchCO(line string) (bool) {
	if BHRes = BHReg.FindStringSubmatch(line); BHRes != nil {
		// Send object to appropriate destination
		return true
	} else if NSRes = NSReg.FindStringSubmatch(line); NSRes != nil {
		// Send object to appropriate destination
		return true
	} else if WDRes = WDReg.FindStringSubmatch(line); WDRes != nil {
		// Send object to appropriate destination
		return true
	} else {
		found = false
	}
}


func SearchBinaries(line string, lineChan chan string) () {
	
	if strings.Contains(line, "Binaries/multiples:") {
		// FIXME hardflag: H
	} else if strings.Contains(line, "ound nn pairs:") {
		// FIXME hardflag: S
	} else {
		// No binary found
		return false
	}
	
	// If we are here, a binary was found
	// Start a loop to find binary data
	for line = range lineChan {
		if strings.Contains(line, "Total binary energy") || strings.Contains(line, "user_diag:"){ 
			// End of binary parameters section
			break
		}
		if BinIdRes = BinIdReg.FindStringSubmatch(line); BinIdRes != nil {
// 			container["system_time"] = qui posso usare il timestep dello snapshot
// 			container["phys_time"] = container["system_time"] * float(self.time_unit[0])
// 			#print "Found binary in binaries section ", bin_ids.group(2)
// 			container["ids"] = bin_ids.group(2).translate(None, "()").split(",")
// 			container["hardflag"] = self.hardflag
// 			container["sma"] = float(bin_ids.group(3)) * self.r_unit[0]
// 			container["ecc"] = bin_ids.group(4)
// 			container["period"] = float(bin_ids.group(5)) * self.time_unit[0]
// 			bin_ids = None# reset ids variable
		} else if strings.Contains(line, "masses") {
// 			mass_line = line.split()[1:-3]
// 						for mass in mass_line:
// 							container["masses"].append(float(mass)*self.mass_unit[0])	
// 						self.container_list.append(container)
// 						# Reinit container
// 						container = {
// 						"system_time": None,
// 						"phys_time": None,
// 						"ids": "--",
// 						"objects": [],
// 						"masses": [],
// 						"sma": "--",
// 						"period": "--",
// 						"ecc": "--"#,
		}
		
	}
	
	return true
}
