// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package enigma

import (
	"fmt"
	"strings"
)

const (
	// PaddingNONE is a Padding of type NONE.
	PaddingNONE Padding = iota
	// PaddingPKCS is a Padding of type PKCS.
	PaddingPKCS
)

const _PaddingName = "NONEPKCS"

var _PaddingNames = []string{
	_PaddingName[0:4],
	_PaddingName[4:8],
}

// PaddingNames returns a list of possible string values of Padding.
func PaddingNames() []string {
	tmp := make([]string, len(_PaddingNames))
	copy(tmp, _PaddingNames)
	return tmp
}

var _PaddingMap = map[Padding]string{
	PaddingNONE: _PaddingName[0:4],
	PaddingPKCS: _PaddingName[4:8],
}

// String implements the Stringer interface.
func (x Padding) String() string {
	if str, ok := _PaddingMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Padding(%d)", x)
}

var _PaddingValue = map[string]Padding{
	_PaddingName[0:4]:                  PaddingNONE,
	strings.ToLower(_PaddingName[0:4]): PaddingNONE,
	_PaddingName[4:8]:                  PaddingPKCS,
	strings.ToLower(_PaddingName[4:8]): PaddingPKCS,
}

// ParsePadding attempts to convert a string to a Padding.
func ParsePadding(name string) (Padding, error) {
	if x, ok := _PaddingValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _PaddingValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return Padding(0), fmt.Errorf("%s is not a valid Padding, try [%s]", name, strings.Join(_PaddingNames, ", "))
}
