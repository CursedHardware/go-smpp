package sms

import (
	"bytes"
	"io"
	"time"

	"github.com/NiceLabs/go-smpp/coding/semioctet"
)

type Time struct{ time.Time }

func (t *Time) ReadFrom(r io.Reader) (n int64, err error) {
	data := make([]byte, 7)
	if _, err = r.Read(data); err != nil {
		return
	}
	blocks := semioctet.DecodeSemi(data)
	t.Time = time.Date(
		2000+int(blocks[0]),
		time.Month(blocks[1]),
		int(blocks[2]),
		int(blocks[3]),
		int(blocks[4]),
		int(blocks[5]),
		0,
		time.FixedZone("", int(blocks[6])*900),
	)
	return
}

func (t Time) WriteTo(w io.Writer) (n int64, err error) {
	_, offset := t.Time.Zone()
	return semioctet.EncodeSemi(
		w,
		t.Time.Year()-2000,
		int(t.Time.Month()),
		t.Time.Day(),
		t.Time.Hour(),
		t.Time.Minute(),
		t.Time.Second(),
		0,
		offset/900,
	)
}

type Duration struct{ time.Duration }

func (d *Duration) ReadFrom(r io.Reader) (n int64, err error) {
	data := make([]byte, 1)
	if _, err = r.Read(data); err != nil {
		return
	}
	switch n := time.Duration(data[0]); {
	case n <= 143:
		d.Duration = 1 + 5*time.Minute*n
	case n <= 167:
		halfDays := 12 * time.Hour
		halfHours := 30 * time.Minute
		d.Duration = (n-143)*halfHours + halfDays
	case n <= 196:
		d.Duration = (n - 166) * 24 * time.Hour
	default:
		d.Duration = (n - 192) * 7 * 24 * time.Hour
	}
	return
}

func (d Duration) WriteTo(w io.Writer) (n int64, err error) {
	var period time.Duration
	if minutes := d.Duration / time.Minute; minutes <= 5 {
		period = 0
	} else if hours := d.Duration / time.Hour; hours <= 12 {
		period = minutes/5 - 1
	} else if hours <= 24 {
		period = (d.Duration-(hours*12))/(time.Minute*30) + 143
	} else if days := hours / 24; days <= 31 {
		period = hours/24 + 166
	} else if weeks := days / 7; weeks <= 62 {
		period = weeks + 192
	} else {
		period = 255
	}
	var buf bytes.Buffer
	buf.WriteByte(byte(period))
	return buf.WriteTo(w)
}

type ValidityPeriod interface{}
