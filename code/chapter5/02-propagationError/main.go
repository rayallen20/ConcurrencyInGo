package main

import (
	"code/chapter5/02-propagationError/intermediate"
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)
	err := intermediate.RunJob("1")
	if err != nil {
		msg := "There was an unexpected issue; please report this as a bug."
		if _, ok := err.(intermediate.IntermediateErr); ok {
			msg = err.Error()
		}

		handlerError(1, err, msg)
	}
}

func handlerError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v] ", key))
	log.Printf("%#v", err)
	fmt.Printf("[%v] %v\n", key, message)
}
