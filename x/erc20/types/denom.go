package types

import (
	"fmt"
	"strings"
)

// CreateDenomDescription generates a string with the coin description
func CreateDenomDescription(address string) string {
	return fmt.Sprintf("Merlion coin token representation of %s", address)
}

// CreateDenom generates a string the module name plus the address to avoid conflicts with names staring with a number
func CreateDenom(address string) string {
	return fmt.Sprintf("%s/%s", DenomPrefix, address)
}

// SanitizeERC20Name enforces snake_case and removes all "coin" and "token"
// strings from the ERC20 name.
func SanitizeERC20Name(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " token", "")
	name = strings.ReplaceAll(name, " coin", "")
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}
