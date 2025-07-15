package vast

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Duration is a VAST duration expressed a hh:mm:ss
type Duration time.Duration

// MarshalText implements the encoding.TextMarshaler interface.
func (dur Duration) MarshalText() ([]byte, error) {
	h := dur / Duration(time.Hour)
	m := dur % Duration(time.Hour) / Duration(time.Minute)
	s := dur % Duration(time.Minute) / Duration(time.Second)
	ms := dur % Duration(time.Second) / Duration(time.Millisecond)
	if ms == 0 {
		return []byte(fmt.Sprintf("%02d:%02d:%02d", h, m, s)), nil
	}
	return []byte(fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (dur *Duration) UnmarshalText(data []byte) (err error) {
	s := string(data)
	s = strings.TrimSpace(s)

	// First chek for bogus data.
	if s == "" || strings.ToLower(s) == "undefined" {
		*dur = 0
		return nil
	}
	parts := strings.SplitN(s, ":", 3)
	if len(parts) > 3 {
		return fmt.Errorf("invalid duration -- too many colons: %s", data)
	}

	// Extract any milliseconds from the last part and remove them from the string.
	finalDuration := Duration(0)
	lastPart := parts[len(parts)-1]
	lastPartPieces := strings.Split(lastPart, ".")
	if len(lastPartPieces) > 2 {
		return fmt.Errorf("invalid duration -- too many periods: %s", data)
	}
	if len(lastPartPieces) == 2 {
		msString := lastPartPieces[1]
		if len(msString) > 3 {
			return fmt.Errorf("invalid duration -- milliseconds too long: %s", data)
		}
		if len(msString) == 0 {
			return fmt.Errorf("invalid duration -- empty milliseconds: %s", data)
		}
		if len(msString) < 3 {
			// Pad with zeros to ensure we have 3 digits.
			for len(msString) < 3 {
				msString += "0"
			}
			lastPartPieces[1] = msString
		}
		parsedMS, err := strconv.ParseInt(msString, 10, 32)
		if err != nil || parsedMS < 0 {
			return fmt.Errorf("invalid duration -- invalid milliseconds: %s", data)
		}
		finalDuration = Duration(parsedMS) * Duration(time.Millisecond)
		parts[len(parts)-1] = lastPartPieces[0]
	}

	multipliers := []time.Duration{time.Second, time.Minute, time.Hour}[0:len(parts)]

	for i, part := range parts {
		parsedValue, err := strconv.ParseInt(part, 10, 32)
		if err != nil || parsedValue < 0 {
			return fmt.Errorf("invalid duration -- invalid time value: %s %s", data, part)
		}
		finalDuration += Duration(parsedValue) * Duration(multipliers[len(multipliers)-1-i])
	}

	*dur = finalDuration

	return nil
}
