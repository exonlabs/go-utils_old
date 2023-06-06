package sqlite

import (
	"fmt"
	"os"

	"github.com/exonlabs/go-utils/pkg/console"
)

func InteractiveConfig(defaults KwArgs) (KwArgs, error) {
	options := KwArgs{}

	con := console.NewConsole()

	database, _ := defaults["database"].(string)
	if input, err := con.GetValue(
		"Enter database path", &console.InputOpts{
			Required: true, Default: database}); err != nil {
		return nil, err
	} else {
		options["database"] = input
	}

	extargs, _ := defaults["extargs"].(string)
	if input, err := con.GetValue(
		"Enter connection extra-args", &console.InputOpts{
			Default: extargs}); err != nil {
		return nil, err
	} else {
		options["extargs"] = input
	}

	return options, nil
}

func InteractiveSetup(options KwArgs) error {
	database, _ := options["database"].(string)
	if len(database) == 0 {
		return fmt.Errorf("invalid database configuration")
	}

	if _, err := os.Stat(database); os.IsNotExist(err) {
		file, err := os.OpenFile(
			database, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
		if err != nil {
			return err
		}
		file.Close()
	}

	return nil
}
