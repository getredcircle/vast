package vast

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDurationMarshaler(t *testing.T) {
	b, err := Duration(0).MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "00:00:00", string(b))
	}
	b, err = Duration(2 * time.Millisecond).MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "00:00:00.002", string(b))
	}
	b, err = Duration(2 * time.Second).MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "00:00:02", string(b))
	}
	b, err = Duration(2 * time.Minute).MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "00:02:00", string(b))
	}
	b, err = Duration(2 * time.Hour).MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "02:00:00", string(b))
	}
}

func TestDurationUnmarshal(t *testing.T) {
	var d Duration
	if assert.NoError(t, d.UnmarshalText([]byte("00:00:00"))) {
		assert.Equal(t, Duration(0), d)
	}
	d = 0
	if assert.NoError(t, d.UnmarshalText([]byte("00:00:02"))) {
		assert.Equal(t, Duration(2*time.Second), d)
	}
	d = 0
	if assert.NoError(t, d.UnmarshalText([]byte(" 00:00:02 "))) {
		assert.Equal(t, Duration(2*time.Second), d)
	}
	d = 0
	if assert.NoError(t, d.UnmarshalText([]byte("00:02:00"))) {
		assert.Equal(t, Duration(2*time.Minute), d)
	}
	d = 0
	if assert.NoError(t, d.UnmarshalText([]byte("02:00:00"))) {
		assert.Equal(t, Duration(2*time.Hour), d)
	}
	d = 0
	if assert.NoError(t, d.UnmarshalText([]byte("00:00:00.123"))) {
		assert.Equal(t, Duration(123*time.Millisecond), d)
	}
	d = 0
	if assert.NoError(t, d.UnmarshalText([]byte("undefined"))) {
		assert.Equal(t, Duration(0), d)
	}
	d = 0
	if assert.NoError(t, d.UnmarshalText([]byte(""))) {
		assert.Equal(t, Duration(0), d)
	}
	if assert.NoError(t, d.UnmarshalText([]byte("00:00:60"))) {
		assert.Equal(t, Duration(60*time.Second), d)
	}
	if assert.NoError(t, d.UnmarshalText([]byte("00:60:00"))) {
		assert.Equal(t, Duration(60*time.Minute), d)
	}
	if assert.NoError(t, d.UnmarshalText([]byte("00:00:70"))) {
		assert.Equal(t, Duration(70*time.Second), d)
	}
	if assert.NoError(t, d.UnmarshalText([]byte("00:00:70.123"))) {
		assert.Equal(t, Duration(70*time.Second)+Duration(time.Millisecond*123), d)
	}
	if assert.NoError(t, d.UnmarshalText([]byte("70.123"))) {
		assert.Equal(t, Duration(70*time.Second)+Duration(time.Millisecond*123), d)
	}
	if assert.NoError(t, d.UnmarshalText([]byte("70"))) {
		assert.Equal(t, Duration(70*time.Second), d)
	}
	if assert.NoError(t, d.UnmarshalText([]byte("1:10.1"))) {
		assert.Equal(t, Duration(1*time.Minute+10*time.Second+100*time.Millisecond), d)
	}
	if assert.NoError(t, d.UnmarshalText([]byte("1:10.100"))) {
		assert.Equal(t, Duration(1*time.Minute+10*time.Second+100*time.Millisecond), d)
	}
	if assert.NoError(t, d.UnmarshalText([]byte("00:1:10.100"))) {
		assert.Equal(t, Duration(1*time.Minute+10*time.Second+100*time.Millisecond), d)
	}
	if assert.NoError(t, d.UnmarshalText([]byte("1:10"))) {
		assert.Equal(t, Duration(1*time.Minute+10*time.Second), d)
	}
	if assert.NoError(t, d.UnmarshalText([]byte("1:10.1"))) {
		assert.Equal(t, Duration(1*time.Minute+10*time.Second+100*time.Millisecond), d)
	}
	if assert.NoError(t, d.UnmarshalText([]byte("1:10.10"))) {
		assert.Equal(t, Duration(1*time.Minute+10*time.Second+100*time.Millisecond), d)
	}

	assert.EqualError(t, d.UnmarshalText([]byte("00:00:00.-1")), "invalid duration -- invalid milliseconds: 00:00:00.-1")
	assert.EqualError(t, d.UnmarshalText([]byte("00:00:00.")), "invalid duration -- empty milliseconds: 00:00:00.")
	assert.EqualError(t, d.UnmarshalText([]byte("00:00:12.34.56")), "invalid duration -- too many periods: 00:00:12.34.56")
	assert.EqualError(t, d.UnmarshalText([]byte("00:00:12.34.")), "invalid duration -- too many periods: 00:00:12.34.")
	assert.EqualError(t, d.UnmarshalText([]byte("00:00:00.1000")), "invalid duration -- milliseconds too long: 00:00:00.1000")
	assert.EqualError(t, d.UnmarshalText([]byte("00h01m")), "invalid duration -- invalid time value: 00h01m 00h01m")
}
