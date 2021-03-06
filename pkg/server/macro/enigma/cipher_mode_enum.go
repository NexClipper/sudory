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
	// CipherModeNONE is a CipherMode of type NONE.
	CipherModeNONE CipherMode = iota
	// CipherModeCBC is a CipherMode of type CBC.
	CipherModeCBC
	// CipherModeGCM is a CipherMode of type GCM.
	CipherModeGCM
)

const _CipherModeName = "NONECBCGCM"

var _CipherModeNames = []string{
	_CipherModeName[0:4],
	_CipherModeName[4:7],
	_CipherModeName[7:10],
}

// CipherModeNames returns a list of possible string values of CipherMode.
func CipherModeNames() []string {
	tmp := make([]string, len(_CipherModeNames))
	copy(tmp, _CipherModeNames)
	return tmp
}

var _CipherModeMap = map[CipherMode]string{
	CipherModeNONE: _CipherModeName[0:4],
	CipherModeCBC:  _CipherModeName[4:7],
	CipherModeGCM:  _CipherModeName[7:10],
}

// String implements the Stringer interface.
func (x CipherMode) String() string {
	if str, ok := _CipherModeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("CipherMode(%d)", x)
}

var _CipherModeValue = map[string]CipherMode{
	_CipherModeName[0:4]:                   CipherModeNONE,
	strings.ToLower(_CipherModeName[0:4]):  CipherModeNONE,
	_CipherModeName[4:7]:                   CipherModeCBC,
	strings.ToLower(_CipherModeName[4:7]):  CipherModeCBC,
	_CipherModeName[7:10]:                  CipherModeGCM,
	strings.ToLower(_CipherModeName[7:10]): CipherModeGCM,
}

// ParseCipherMode attempts to convert a string to a CipherMode.
func ParseCipherMode(name string) (CipherMode, error) {
	if x, ok := _CipherModeValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _CipherModeValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return CipherMode(0), fmt.Errorf("%s is not a valid CipherMode, try [%s]", name, strings.Join(_CipherModeNames, ", "))
}
