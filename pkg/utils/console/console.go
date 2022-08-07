package console

import (
    "fmt"
    "os"
    "runtime"
    "strings"
    "golang.org/x/term"
)

var (
    redBrightEscape = "\033[31;1m"
    whiteBrightEscape = "\033[37;1m"
    resetEscape = "\033[0m"
)

type OptArgs map[string]any

type Console struct {
    PromptCaret string
}

func NewConsole() *Console {
    return &Console{PromptCaret: ">>"}
}

func (c *Console) getinput(msg string, args OptArgs) string {
    oldstate, err := term.MakeRaw(int(os.Stdin.Fd()))
    if err != nil {
        return "", err
    }
    defer term.Restore(int(os.Stdin.Fd()), oldstate)

    var prompt string
    if runtime.GOOS != "windows" {
        // print BOLD/BRIGHT WHITE text prompt
        prompt = fmt.Sprintf("%v%v %v %v",
            whiteBrightEscape, c.PromptCaret, msg, resetEscape)
    } else {
        prompt = fmt.Sprintf("%v %v ", c.PromptCaret, msg)
    }
    // prompt = fmt.Sprintf("%v %v ", c.PromptCaret, msg)

    th := term.NewTerminal(os.Stdin, prompt)
    // th.SetPrompt(prompt)

    line, err := th.ReadLine()
    if err != nil {
        return "", err
    }

    line = strings.TrimSpace(line)
    if len(line) == 0 {
        if v, ok := opts["default"]; ok {
            line = fmt.Sprint(v)
        }
    }

    return line, nil

}

func (c *Console) GetValue(msg string, args OptArgs) string {

    return c.getinput(msg, opts)
    // fmt.Println(">> user input ? ----", msg)
    // for k, v := range opts {
    //     fmt.Printf("%v %T = %v\n", k, v, v)
    // }
    // // fmt.Printf("%T defval = %v\n", opts.defval, opts.defval)
    // // fmt.Printf("%T hidden = %v\n", opts.hidden, opts.hidden)
    // // fmt.Printf("%T required = %v\n", opts.required, opts.required)
    // // fmt.Printf("%T trials = %v\n", opts.trials, opts.trials)
    // return "none", nil
}

func (c *Console) GetPassword() string {
    fmt.Println(">> user password ? ----")
    return "none", nil
}

// func (c *Console) ConfirmPasswd() {

// }

// func (c *Console) GetNumber() {

// }

// func (c *Console) GetDecimal() {

// }

// func (c *Console) SelectValue() {

// }

// func (c *Console) SelectYesNo() {

// }
