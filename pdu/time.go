package pdu

import (
	"fmt"
	"strconv"
	"time"
)

// Time see SMPP v5, section 4.7.23.4 (132p)
type Time struct{ time.Time }

func (t *Time) From(input string) (err error) {
	t.Time = time.Time{}
	if len(input) == 0 {
		return
	}
	parts, symbol := fromTimeString(input)
	if !(symbol == '+' || symbol == '-') {
		err = ErrTimeNotParsed
		return
	}
	t.Time = time.Date(
		int(2000+parts[0]),                    // year
		time.Month(parts[1]),                  // month
		int(parts[2]),                         // day
		int(parts[3]),                         // hour
		int(parts[4]),                         // minute
		int(parts[5]),                         // second
		int(parts[6])*1e8,                     // tenths of second
		time.FixedZone("", int(parts[7]*900)), // timezone offset
	)
	return
}

func (t Time) String() string {
	if t.Time.IsZero() {
		return ""
	}
	_, offset := t.Zone()
	symbol := '+'
	if offset < 0 {
		offset = -offset
		symbol = '-'
	}
	return fmt.Sprintf(
		"%02d%02d%02d%02d%02d%02d%d%02d%c",
		t.Year()-2000,      // year
		int(t.Month()),     // month
		t.Day(),            // day
		t.Hour(),           // hour
		t.Minute(),         // minute
		t.Second(),         // second
		t.Nanosecond()/1e8, // tenths of second
		offset/900,         // offset
		symbol,             // time-zone symbol
	)
}

// Duration see SMPP v5, section 4.7.23.5 (132p)
type Duration struct{ time.Duration }

func (p *Duration) From(input string) (err error) {
	p.Duration = 0
	if len(input) == 0 {
		return
	}
	parts, symbol := fromTimeString(input)
	if symbol != 'R' {
		err = ErrTimeNotParsed
		return
	}
	bases := []time.Duration{
		time.Hour * 8760, time.Hour * 720, time.Hour * 24,
		time.Hour, time.Minute, time.Second, 1e8, 0,
	}
	for i, part := range parts {
		p.Duration += bases[i] * time.Duration(part)
	}
	return
}

func (p Duration) String() string {
	if p.Duration < time.Second {
		return ""
	}
	ts := p.Duration
	parts := []time.Duration{
		time.Hour * 8760, time.Hour * 720, time.Hour * 24,
		time.Hour, time.Minute, time.Second,
	}
	for i, part := range parts {
		parts[i] = ts / part
		ts %= part
	}
	return fmt.Sprintf(
		"%02d%02d%02d%02d%02d%02d%d00R",
		parts[0], parts[1], parts[2],
		parts[3], parts[4], parts[5],
		int(ts.Nanoseconds()/1e8),
	)
}

func fromTimeString(input string) (parts [8]int64, symbol byte) {
	if len(input) != 16 {
		return
	}
	for i := 0; i < 12; i += 2 {
		parts[i/2], _ = strconv.ParseInt(input[i:i+2], 10, 16)
	}
	parts[6], _ = strconv.ParseInt(input[12:13], 10, 16)
	parts[7], _ = strconv.ParseInt(input[13:15], 10, 16)
	symbol = input[15]
	if symbol == '-' {
		parts[7] = -parts[7]
	}
	return
}
