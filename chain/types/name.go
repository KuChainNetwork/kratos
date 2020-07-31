package types

import (
	"bytes"
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

const (
	NameStrLenMax = 17

	NameBytesLen         = 16
	NameBytesTypeCodeLen = 1
	NameBytesVersionLen  = 1
	NameBytesLengthLen   = 1
	NameBytesHeaderLen   = NameBytesTypeCodeLen + NameBytesVersionLen + NameBytesLengthLen

	NameStrLengthIdx = 2
)

const (
	currentNameType    uint8 = 1
	currentNameVersion uint8 = 1
)

const (
	nameByteMask  byte = 0xFF
	nameByteMask2      = nameByteMask >> 2
	nameByteMask6      = nameByteMask >> 6
)

const (
	CharValueSeq       byte = 0  // @
	CharValueNil       byte = 63 // 111111
	CharValueDot       byte = 49 // .
	CharValueUnderline byte = 50 // _
)

type Name struct {
	Value []byte `json:"value,omitempty" yaml:"value"`
}

// NewName create Name by string
func NewName(str string) (Name, error) {
	if !VerifyNameString(str) {
		return Name{}, ErrNameStrNoValid
	}
	return parseName(str)
}

// NewNameFromBytes create Name from bytes
func NewNameFromBytes(b []byte) Name {
	res := Name{
		make([]byte, NameStrLenMax),
	}
	copy(res.Value[:], b)
	return res
}

// MustName create Name by string if name is error then panic
func MustName(str string) Name {
	n, err := parseName(str)
	if err != nil {
		panic(err)
	}
	return n
}

// VerifyNameString return if str is valid for name
func VerifyNameString(str string) bool {
	if str == "" {
		return true
	}

	if len(str) > NameStrLenMax {
		return false
	}

	if str[0] == '@' || str[len(str)-1] == '@' {
		return false
	}

	seqCharCount := 0
	for _, c := range str {
		if c == '@' {
			seqCharCount++
		}
		if CharValueNil == char2byte(byte(c)) {
			return false
		}
	}

	return seqCharCount <= 1
}

// parseName parse name from string
func parseName(nameStr string) (Name, error) {
	nameStrLen := len(nameStr)
	if nameStrLen > NameStrLenMax {
		return Name{}, ErrNameParseTooLen
	}

	if nameStr == "" {
		return Name{}, nil
	}

	res := Name{make([]byte, NameStrLenMax)}
	res.Value[0] = currentNameType
	res.Value[1] = currentNameVersion

	res.Value[NameStrLengthIdx] = uint8(nameStrLen)

	var appendLocTyp uint8 = 0
	loc := NameBytesHeaderLen
	for _, c := range nameStr {
		cc := char2byte(byte(c))

		if cc == CharValueNil {
			return Name{}, ErrNameCharError
		}

		switch appendLocTyp {
		case 0:
			res.Value[loc] |= (cc << 2)
		case 1:
			res.Value[loc] |= (cc >> 4)
			res.Value[loc+1] |= (cc << 4)
			loc++
		case 2:
			res.Value[loc] |= (cc >> 2)
			res.Value[loc+1] |= (cc << 6)
		case 3:
			res.Value[loc+1] |= cc
			loc += 2
		}

		appendLocTyp++
		appendLocTyp &= nameByteMask6
	}

	return res, nil
}

// char2byte char to byte
func char2byte(c byte) byte {
	if c >= byte('a') && c <= byte('z') {
		return 1 + (c - 'a')
	}

	if c >= byte('0') && c <= byte('9') {
		return 32 + (c - '0')
	}

	switch c {
	case '@':
		return CharValueSeq
	case '.':
		return CharValueDot
	case '_':
		return CharValueUnderline
	}

	return CharValueNil
}

// byte2char byte convert to char
func byte2char(c byte) byte {
	switch c {
	case CharValueSeq:
		return '@'
	case CharValueDot:
		return '.'
	case CharValueUnderline:
		return '_'
	}

	if c <= 26 && c >= 1 {
		return c + ('a' - 1)
	}

	if c <= 41 && c >= 32 {
		return c + ('0' - 32)
	}

	return '*'
}

// name2StringV1 name to string in version 1
func name2StringV1(n Name) []byte {
	strLen := int(n.Value[NameStrLengthIdx])
	res := make([]byte, 0, strLen+1)

	var appendLocTyp uint8 = 0
	loc := NameBytesHeaderLen
	for i := 0; i < strLen; i++ {
		var charByte byte

		switch appendLocTyp {
		case 0:
			charByte |= (n.Value[loc] >> 2)
		case 1:
			charByte |= (n.Value[loc] << 4)
			charByte |= (n.Value[loc+1] >> 4)
			loc++
		case 2:
			charByte |= (n.Value[loc] << 2)
			charByte |= (n.Value[loc+1] >> 6)
		case 3:
			charByte = n.Value[loc+1]
			loc += 2
		}

		appendLocTyp++
		appendLocTyp &= nameByteMask6

		res = append(res, byte2char(charByte&nameByteMask2))
	}

	return res
}

// String get name readable string
func (n Name) String() string {
	if n.Empty() {
		return ""
	}

	switch n.Value[1] {
	case 1:
		return string(name2StringV1(n))
	default:
		return ""
	}
}

// Eq if name is eq to other
func (n Name) Eq(o Name) bool {
	// Note: When version update, need check
	if (!n.Empty()) && (!o.Empty()) {
		return bytes.Equal(n.Value[NameStrLengthIdx:], o.Value[NameStrLengthIdx:])
	}

	return n.Empty() && o.Empty()
}

// IsNameEq return is r eq l
func IsNameEq(r, l Name) bool {
	return r.Eq(l)
}

// ForEach for each char to iter
func (n Name) Foreach(op func(c byte) bool) {
	strLen := n.Len()
	var appendLocTyp uint8 = 0

	loc := NameBytesHeaderLen
	for i := 0; i < strLen; i++ {
		var charByte byte

		switch appendLocTyp {
		case 0:
			charByte |= (n.Value[loc] >> 2)
		case 1:
			charByte |= (n.Value[loc] << 4)
			charByte |= (n.Value[loc+1] >> 4)
			loc++
		case 2:
			charByte |= (n.Value[loc] << 2)
			charByte |= (n.Value[loc+1] >> 6)
		case 3:
			charByte = n.Value[loc+1]
			loc += 2
		}

		appendLocTyp++
		appendLocTyp &= nameByteMask6

		if !op(byte2char(charByte & nameByteMask2)) {
			break
		}
	}
}

// Len return name string len
func (n Name) Len() int {
	if len(n.Value) <= NameStrLengthIdx {
		return 0
	}
	return int(n.Value[NameStrLengthIdx])
}

// Empty return is name a empty
func (n Name) Empty() bool {
	return n.Len() == 0
}

// MarshalJSON just return a json string
func (n Name) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}

// UnmarshalJSON unmarshal from JSON assuming Bech32 encoding.
func (n *Name) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	nn, err := NewName(s)
	if err != nil {
		return err
	}

	*n = nn
	return nil
}

// MarshalYAML marshals to YAML using Bech32.
func (n Name) MarshalYAML() (interface{}, error) {
	return n.String(), nil
}

// UnmarshalYAML unmarshals from JSON assuming Bech32 encoding.
func (n *Name) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	nn, err := NewName(s)
	if err != nil {
		return err
	}

	*n = nn
	return nil
}

// Bytes returns the raw address bytes.
func (n Name) Bytes() []byte {
	return n.Value[:]
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (n Name) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(n.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", n.Value[:])))
	default:
		s.Write([]byte(fmt.Sprintf("%X", n.Value[:])))
	}
}
