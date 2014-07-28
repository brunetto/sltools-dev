package main

import (
	"fmt"
	"log"
	"time"
	
	"bitbucket.org/brunetto/sltools/slt"
	"github.com/brunetto/goutils/debug"
)

func main () () {
	defer debug.TimeMe(time.Now())

	if err := slt.PbsLaunch(); err != nil {
		log.Fatal(err)
	}
	
	
	fmt.Print("\x07") // Beep when finish!!:D
}


