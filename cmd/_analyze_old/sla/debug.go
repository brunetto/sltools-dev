package sla

// https://groups.google.com/d/msg/golang-nuts/N2szzAFrcFM/s2WosyIDNP4J

import (
	"fmt"
	"log"
	"runtime"
	"path/filepath"
)

// Whoami print the name of the calling function.
func Whoami(print bool) string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}
	me := runtime.FuncForPC(pc)
	if me == nil {
		return "unnamed"
	}
	if print{
		fmt.Println("###########################################")
		fmt.Print("###\t\t\t")
		log.Println(filepath.Base(me.Name()))
		fmt.Println("###########################################")
	}
	return filepath.Base(me.Name())
}