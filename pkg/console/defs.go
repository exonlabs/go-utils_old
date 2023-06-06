package console

// generic types
type KwArgs = map[string]any
type float = float64

const (
	ESC_REDBRT = "\033[31;1m"
	ESC_BRIGHT = "\033[1m"
	ESC_RESET  = "\033[0m"
)

type InputOpts = struct {
	Default  any
	Required bool
	Hidden   bool
	Trials   int
	Regex    string
}

type Validator = func(string) (any, error)
