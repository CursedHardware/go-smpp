package coding

import (
	"github.com/NiceLabs/go-smpp/coding/gsm7bit"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/unicode"
)

// DataCoding see SMPP v5, section 4.7.7 (123p)

const (
	GSM7BitCoding   DataCoding = 0b00000000 // GSM 7Bit
	ASCIICoding     DataCoding = 0b00000001 // ASCII
	Latin1Coding    DataCoding = 0b00000011 // ISO-8859-1 (Latin-1)
	ShiftJISCoding  DataCoding = 0b00000101 // Shift-JIS
	CyrillicCoding  DataCoding = 0b00000110 // ISO-8859-5 (Cyrillic)
	HebrewCoding    DataCoding = 0b00000111 // ISO-8859-8 (Hebrew)
	UCS2Coding      DataCoding = 0b00001000 // UCS-2
	ISO2022JPCoding DataCoding = 0b00001010 // ISO-2022-JP
	EUCJPCoding     DataCoding = 0b00001101 // Extended Kanji JIS (X 0212-1990)
	EUCKRCoding     DataCoding = 0b00001110 // KS X 1001 (KS C 5601)
	NoCoding        DataCoding = 0b10111111 // Reserved (Non-specification definition)
)

var encodingMap = map[DataCoding]encoding.Encoding{
	GSM7BitCoding:   gsm7bit.Packed,
	ASCIICoding:     charmap.ISO8859_1,
	Latin1Coding:    charmap.ISO8859_1,
	ShiftJISCoding:  japanese.ShiftJIS,
	CyrillicCoding:  charmap.ISO8859_5,
	HebrewCoding:    charmap.ISO8859_8,
	UCS2Coding:      unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
	ISO2022JPCoding: japanese.ISO2022JP,
	EUCJPCoding:     japanese.EUCJP,
	EUCKRCoding:     korean.EUCKR,
}

var splitterMap = map[DataCoding]Splitter{
	GSM7BitCoding:   _7BitSplitter,
	ASCIICoding:     _1ByteSplitter,
	HebrewCoding:    _1ByteSplitter,
	CyrillicCoding:  _1ByteSplitter,
	Latin1Coding:    _1ByteSplitter,
	ShiftJISCoding:  _MultibyteSplitter,
	ISO2022JPCoding: _MultibyteSplitter,
	EUCJPCoding:     _MultibyteSplitter,
	EUCKRCoding:     _MultibyteSplitter,
	UCS2Coding:      _UTF16Splitter,
}
