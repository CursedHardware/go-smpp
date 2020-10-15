package sms

import (
	"bufio"
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
		2000+blocks[0],
		time.Month(blocks[1]),
		blocks[2],
		blocks[3],
		blocks[4],
		blocks[5],
		0,
		time.FixedZone("", blocks[6]*900),
	)
	return
}

func (t *Time) WriteTo(w io.Writer) (n int64, err error) {
	_, offset := t.Time.Zone()
	return semioctet.EncodeSemi(
		w,
		t.Time.Year()-2000,
		int(t.Time.Month()),
		t.Time.Day(),
		t.Time.Hour(),
		t.Time.Minute(),
		t.Time.Second(),
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
		n++
		d.Duration = 5 * time.Minute * n
	case n <= 167:
		const halfDays = 12 * time.Hour
		const halfHours = 30 * time.Minute
		d.Duration = (n-143)*halfHours + halfDays
	case n <= 196:
		d.Duration = (n - 166) * 24 * time.Hour
	default:
		d.Duration = (n - 192) * 7 * 24 * time.Hour
	}
	return
}

func (d *Duration) WriteTo(w io.Writer) (n int64, err error) {
	var period time.Duration
	if minutes := d.Duration / time.Minute; minutes <= 5 {
		period = 0
	} else if hours := d.Duration / time.Hour; hours <= 12 {
		period = minutes/5 - 1
	} else if hours <= 24 {
		const halfDays = 12 * time.Hour
		const halfHours = 30 * time.Minute
		period = (d.Duration-halfDays)/halfHours + 143
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

type EnhancedDuration struct {
	time.Duration
	Indicator byte
}

func (d *EnhancedDuration) ReadFrom(r io.Reader) (n int64, err error) {
	buf := bufio.NewReader(r)
	if d.Indicator, err = buf.ReadByte(); err != nil {
		return
	}
	length := 6
	switch d.Indicator & 0b111 {
	case 0b001: // relative
		var duration Duration
		_, err = duration.ReadFrom(buf)
		d.Duration = duration.Duration
		length--
	case 0b010: // relative seconds
		var second byte
		second, err = buf.ReadByte()
		d.Duration = time.Second * time.Duration(second)
		length--
	case 0b011: // relative hh:mm:ss
		data := make([]byte, 3)
		_, err = buf.Read(data)
		semi := semioctet.DecodeSemi(data)
		d.Duration = time.Duration(semi[0])*time.Hour +
			time.Duration(semi[1])*time.Minute +
			time.Duration(semi[2])*time.Second
		length -= len(data)
	}
	if err == nil {
		_, err = buf.Discard(length)
	}
	return
}

func (d *EnhancedDuration) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer
	buf.WriteByte(d.Indicator)
	switch d.Indicator & 0b111 {
	case 0b001: // relative
		_, _ = (&Duration{d.Duration}).WriteTo(&buf)
	case 0b010: // relative seconds
		buf.WriteByte(byte(d.Duration / time.Second))
	case 0b011: // relative hh:mm:ss
		hh, mm, ss := int(d.Hours()), int(d.Minutes()), int(d.Seconds())
		_, _ = semioctet.EncodeSemi(&buf, hh, mm-(hh*60), ss-(mm*60))
	}
	buf.Write(make([]byte, 7-buf.Len()))
	return buf.WriteTo(w)
}
