package console

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

// round float to certain decimal places
func round(f64 float64, decimals int) float64 {
	n := math.Pow(10, float64(decimals))
	return float64(math.Round(f64*n)) / n
}

// validate input string against regex
func validateRegex(input string, regex string) (string, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return "", fmt.Errorf("failed validating, " + err.Error())
	}
	if !re.MatchString(input) {
		return "", fmt.Errorf("invalid input format")
	}
	return input, nil
}

// validate numbers
func validateNumber(input string, vmin *int, vmax *int) (int, error) {
	val, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("invalid number format")
	}
	if (vmin != nil && val < *vmin) ||
		(vmax != nil && val > *vmax) {
		return 0, fmt.Errorf("value out of range")
	}
	return val, nil
}

// validate decimals
func validateDecimal(
	input string, decimals int, vmin *float, vmax *float) (float, error) {
	f64, err := strconv.ParseFloat(input, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid decimal format")
	}
	val := float(round(f64, decimals))
	if (vmin != nil && val < *vmin) ||
		(vmax != nil && val > *vmax) {
		return 0, fmt.Errorf("value out of range")
	}
	return val, nil
}
