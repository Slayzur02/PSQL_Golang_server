package myErrors

import (
	"log"
	"os"
)

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
var errLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

func Check(err error) {
	infoLog.Printf("%v", err)
	errLog.Printf("%v", err)	
}
