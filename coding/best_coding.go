package coding

import (
	"fmt"
	. "unicode"

	"github.com/NiceLabs/go-smpp/coding/gsm7bit"
	. "golang.org/x/text/encoding"
)

type DataCoding byte

func (c DataCoding) GoString() string   { return c.String() }
func (c DataCoding) String() string     { return fmt.Sprintf("%08b", byte(c)) }
func (c DataCoding) Encoding() Encoding { return encodingMap[c] }
func (c DataCoding) Splitter() Splitter { return splitterMap[c] }

func BestCoding(input string) DataCoding {
	switch {
	case isRangeTable(input, gsm7bit.DefaultAlphabet):
		return GSM7BitCoding
	case isRangeTable(input, _ASCII):
		return ASCIICoding
	case isRangeTable(input, _Latin1):
		return Latin1Coding
	case isRangeTable(input, _Cyrillic):
		return CyrillicCoding
	case isRangeTable(input, _Hebrew):
		return HebrewCoding
	case isRangeTable(input, _Shift_JIS):
		return ShiftJISCoding
	case isRangeTable(input, _EUC_KR):
		return EUCKRCoding
	}
	return UCS2Coding
}

func isRangeTable(input string, table *RangeTable) bool {
	for _, r := range input {
		if !Is(table, r) {
			return false
		}
	}
	return true
}
