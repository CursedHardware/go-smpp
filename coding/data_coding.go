package coding

import (
	"fmt"

	"github.com/VoiceGateway/go-smpp/coding/gsm7bit"
	. "golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/unicode"
)

// DataCoding see SMPP v5, section 4.7.7 (123p)

type DataCoding byte

func (c DataCoding) GoString() string {
	return c.String()
}

func (c DataCoding) String() string {
	return fmt.Sprintf("%08b", byte(c))
}

func (c DataCoding) MessageWaitingInfo() (coding DataCoding, active bool, kind int) {
	kind = -1
	coding = NoCoding
	switch c >> 4 & 0b1111 {
	case 0b1100:
	case 0b1101:
		coding = GSM7BitCoding
	case 0b1110:
		coding = UCS2Coding
	default:
		return
	}
	active = c>>3 == 1
	kind = int(c & 0b11)
	return
}

func (c DataCoding) MessageClass() (coding DataCoding, class int) {
	class = int(c & 0b11)
	coding = GSM7BitCoding
	if c>>4&0b1111 != 0b1111 {
		coding = NoCoding
		class = -1
	} else if c>>2&0b1 == 1 {
		coding = UCS2Coding
	}
	return
}

func (c DataCoding) Encoding() Encoding {
	if coding, _, kind := c.MessageWaitingInfo(); kind != -1 {
		return encodingMap[coding]
	} else if coding, class := c.MessageClass(); class != -1 {
		return encodingMap[coding]
	}
	return encodingMap[c]
}

func (c DataCoding) Splitter() Splitter {
	if coding, _, kind := c.MessageWaitingInfo(); kind != -1 {
		return splitterMap[coding]
	} else if coding, class := c.MessageClass(); class != -1 {
		return splitterMap[coding]
	}
	return splitterMap[c]
}

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

var encodingMap = map[DataCoding]Encoding{
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
