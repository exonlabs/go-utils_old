package console

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/term"
)

type Console struct {
	PromptCaret string
}

func NewConsole() *Console {
	return &Console{
		PromptCaret: ">>",
	}
}

// make copy from options and set defaults
func (con *Console) initOptions(opts *InputOpts) *InputOpts {
	var tmpopts InputOpts
	if opts == nil {
		tmpopts = InputOpts{}
	} else {
		tmpopts = *opts
	}

	if tmpopts.Trials <= 0 {
		tmpopts.Trials = 3
	}

	return &tmpopts
}

// print error message
func (con *Console) printError(msg string) {
	msg = fmt.Sprintf(" -- %v", msg)
	// print RED colored error message
	if runtime.GOOS != "windows" {
		msg = ESC_REDBRT + msg + ESC_RESET
	}
	fmt.Println(msg)
}

// read 1 line input from terminal
func (con *Console) readLine(
	msg string, opts *InputOpts) (string, error) {

	// get term status to restore when finish
	oldstate, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldstate)

	prompt := fmt.Sprintf("%v %v ", con.PromptCaret, msg)
	// print BOLD/BRIGHT text prompt
	if runtime.GOOS != "windows" {
		prompt = ESC_BRIGHT + prompt + ESC_RESET
	}

	var line string
	tm := term.NewTerminal(os.Stdin, prompt)
	if opts.Hidden {
		line, err = tm.ReadPassword(prompt)
	} else {
		line, err = tm.ReadLine()
	}
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	if len(line) > 0 {
		return line, nil
	}

	if opts.Default != nil {
		return fmt.Sprint(opts.Default), nil
	}
	return "", nil
}

// get user input with multiple trials
func (con *Console) getInput(
	msg string, opts *InputOpts, vldr Validator) (any, error) {

	msg += ":"
	if opts.Default != nil {
		msg += fmt.Sprintf(" [%v]", opts.Default)
	} else if !opts.Required {
		msg += " []"
	}

	for i := 0; i < opts.Trials; i++ {
		input, err := con.readLine(msg, opts)
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				fmt.Println()
				return nil, err
			}
			con.printError("internal error, " + err.Error())
			continue
		}
		if len(input) == 0 {
			if opts.Required {
				con.printError("input required, please enter value")
				continue
			}
			return opts.Default, nil
		}
		if vldr != nil {
			if val, err := vldr(input); err != nil {
				con.printError(err.Error())
				continue
			} else {
				return val, nil
			}
		}
		return input, nil
	}

	return nil, fmt.Errorf("failed to get valid input")
}

// read general string value
func (con *Console) GetValue(
	msg string, opts *InputOpts) (string, error) {

	opts = con.initOptions(opts)

	var vldr Validator
	if len(opts.Regex) > 0 {
		vldr = func(input string) (any, error) {
			return validateRegex(input, opts.Regex)
		}
	}
	res, err := con.getInput(msg, opts, vldr)

	val, _ := res.(string)
	return val, err
}

func (con *Console) ConfirmValue(
	msg string, val string, opts *InputOpts) error {

	opts = con.initOptions(opts)

	msg += ":"
	for i := 0; i < opts.Trials; i++ {
		input, err := con.readLine(msg, opts)
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				fmt.Println()
				return err
			}
			con.printError("internal error, " + err.Error())
			continue
		}
		if len(input) == 0 {
			con.printError("empty input, please confirm value")
			continue
		}
		if input == val {
			return nil
		}
		con.printError("value not matching, please try again")
	}

	return fmt.Errorf("failed to confirm value")
}

func (con *Console) GetNumber(
	msg string, opts *InputOpts) (int, error) {

	opts = con.initOptions(opts)

	vldr := func(input string) (any, error) {
		return validateNumber(input, nil, nil)
	}
	res, err := con.getInput(msg, opts, vldr)

	val, _ := res.(int)
	return val, err
}

func (con *Console) GetNumberWithLimits(
	msg string, vmin int, vmax int, opts *InputOpts) (int, error) {

	opts = con.initOptions(opts)
	msg += fmt.Sprintf(" (%v...%v)", vmin, vmax)

	vldr := func(input string) (any, error) {
		return validateNumber(input, &vmin, &vmax)
	}
	res, err := con.getInput(msg, opts, vldr)

	val, _ := res.(int)
	return val, err
}

func (con *Console) GetDecimal(
	msg string, decimals int, opts *InputOpts) (float, error) {

	opts = con.initOptions(opts)

	vldr := func(input string) (any, error) {
		return validateDecimal(input, decimals, nil, nil)
	}
	res, err := con.getInput(msg, opts, vldr)

	val, _ := res.(float)
	return val, err
}

func (con *Console) GetDecimalWithLimits(
	msg string, decimals int,
	vmin float, vmax float, opts *InputOpts) (float, error) {

	opts = con.initOptions(opts)
	msg += fmt.Sprintf(" (%v...%v)", vmin, vmax)

	vldr := func(input string) (any, error) {
		return validateDecimal(input, decimals, &vmin, &vmax)
	}
	res, err := con.getInput(msg, opts, vldr)

	val, _ := res.(float)
	return val, err
}

// select string value from certain values
func (con *Console) SelectValue(
	msg string, values []string, opts *InputOpts) (string, error) {

	opts = con.initOptions(opts)
	msg += fmt.Sprintf(" {%v}", strings.Join(values, "|"))

	vldr := func(input string) (any, error) {
		if slices.Contains(values, input) {
			return input, nil
		}
		return "", fmt.Errorf("invalid value")
	}
	res, err := con.getInput(msg, opts, vldr)

	val, _ := res.(string)
	return val, err
}

func (con *Console) SelectYesNo(
	msg string, opts *InputOpts) (bool, error) {

	val, err := con.SelectValue(msg, []string{"y", "n"}, opts)
	if err != nil {
		return false, err
	}
	if val == "y" {
		return true, nil
	}
	return false, nil
}
