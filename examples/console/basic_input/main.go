package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/exonlabs/go-utils/pkg/console"
)

// represent results
func repr(value any, err error) {
	if err != nil {
		reprError(err)
		return
	}
	val := fmt.Sprint(value)
	if len(val) == 0 {
		val = "<empty>"
	}
	fmt.Printf("  * Value: %v\n", val)
}

// represent error or exit
func reprError(err error) {
	if strings.Contains(err.Error(), "EOF") {
		fmt.Print("\n--exit--\n\n")
		os.Exit(0)
	}
	fmt.Printf("Error: %v\n", err)
}

func main() {
	con := console.NewConsole()

	fmt.Println()
	repr(con.GetValue("Enter required string",
		&console.InputOpts{Required: true}))

	fmt.Println()
	repr(con.GetValue("Enter required string with default",
		&console.InputOpts{Required: true, Default: "default val"}))

	fmt.Println()
	repr(con.GetValue("Enter optional string",
		&console.InputOpts{}))

	fmt.Println()
	repr(con.GetValue("Enter optional string with default",
		&console.InputOpts{Default: "default val"}))

	fmt.Println()
	repr(con.GetValue(
		"[input validation] Enter email (user@domain)",
		&console.InputOpts{Required: true,
			Regex: "^[a-zA-Z0-9_.-]+@[a-zA-Z0-9.-]+$"}))

	fmt.Println()
	repr(con.GetValue(
		"[input validation] Enter IPv4 (x.x.x.x)",
		&console.InputOpts{Required: true,
			Regex: "^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}" +
				"(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"}))

	fmt.Println()
	repr(con.GetValue("Test input hidden text",
		&console.InputOpts{Required: true, Hidden: true}))

	fmt.Println()
	if res, err := con.GetValue(
		"[input with confirm] Enter password",
		&console.InputOpts{Required: true, Hidden: true}); err != nil {
		reprError(err)
	} else {
		if err := con.ConfirmValue("Confirm password", res,
			&console.InputOpts{Required: true, Hidden: true}); err != nil {
			reprError(err)
		} else {
			fmt.Printf("  * Password: %v\n", res)
		}
	}

	fmt.Println()
	repr(con.GetNumber("Enter required number",
		&console.InputOpts{Required: true}))

	fmt.Println()
	repr(con.GetNumber("Enter optional number",
		&console.InputOpts{Default: 1234}))

	fmt.Println()
	repr(con.GetNumberWithLimits(
		"Enter required limited number", -100, 100,
		&console.InputOpts{Required: true}))

	fmt.Println()
	repr(con.GetNumberWithLimits(
		"Enter required limited number with default", -100, 100,
		&console.InputOpts{Required: true, Default: 0}))

	fmt.Println()
	repr(con.GetDecimal("Enter required decimal", 4,
		&console.InputOpts{Required: true}))

	fmt.Println()
	repr(con.GetDecimal("Enter optional decimal", 4,
		&console.InputOpts{Default: 123.456}))

	fmt.Println()
	repr(con.GetDecimalWithLimits(
		"Enter required limited decimal", 4, -10.55, 10.88,
		&console.InputOpts{Required: true}))

	fmt.Println()
	repr(con.GetDecimalWithLimits(
		"Enter required limited number with default", 4, -10.55, 10.88,
		&console.InputOpts{Required: true, Default: 0}))

	fmt.Println()
	repr(con.SelectValue("Select from list",
		[]string{"val1", "val2", "val3"},
		&console.InputOpts{Required: true}))

	fmt.Println()
	repr(con.SelectValue("Select from list with default",
		[]string{"val1", "val2", "val3"},
		&console.InputOpts{Required: true, Default: "val2"}))

	fmt.Println()
	repr(con.SelectYesNo("Select Yes/No",
		&console.InputOpts{Required: true}))

	fmt.Println()
	repr(con.SelectYesNo("Select Yes/No with default",
		&console.InputOpts{Required: true, Default: "n"}))

	fmt.Println()
}
