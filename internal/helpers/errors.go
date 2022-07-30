package helpers

import (
	"fmt"
	"log"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panic(fmt.Sprintf("%s: %s", msg, err))
		return
	}
}
