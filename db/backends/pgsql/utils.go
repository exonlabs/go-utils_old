package pgsql

import "github.com/exonlabs/go-utils/console"

func InteractiveConfig(defaults KwArgs) (KwArgs, error) {
	options := KwArgs{}

	con := console.NewConsole()

	database, _ := defaults["database"].(string)
	if input, err := con.GetValue(
		"Enter database name", &console.InputOpts{
			Required: true, Default: database}); err != nil {
		return nil, err
	} else {
		options["database"] = input
	}

	host, _ := defaults["host"].(string)
	if input, err := con.GetValue(
		"Enter database host", &console.InputOpts{
			Required: true, Default: host}); err != nil {
		return nil, err
	} else {
		options["host"] = input
	}

	port, _ := defaults["port"].(int)
	if input, err := con.GetNumberWithLimits(
		"Enter database port", 0, 65536, &console.InputOpts{
			Default: port}); err != nil {
		return nil, err
	} else {
		options["port"] = input
	}

	username, _ := defaults["username"].(string)
	if input, err := con.GetValue(
		"Enter database username", &console.InputOpts{
			Default: username}); err != nil {
		return nil, err
	} else {
		options["username"] = input
	}

	password, _ := defaults["password"].(string)
	if input, err := con.GetValue(
		"Enter database password", &console.InputOpts{
			Default: password}); err != nil {
		return nil, err
	} else {
		options["password"] = input
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
	return nil
}
