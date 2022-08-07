package db

import (
	"fmt"
	"regexp"
)


func SqlIdentifier(name string) string {
	if name == "" {
		panic("invalid empty sql identifier")
	}

	match, err := regexp.MatchString("^[a-zA-Z0-9_]+$", name)
	if err != nil {
		panic(fmt.Sprintf("failed validating sql identifier - %v", err))
	} else if match != true {
		panic(fmt.Sprintf("invalid sql identifier [%v]", name))
	}

	return name
}
