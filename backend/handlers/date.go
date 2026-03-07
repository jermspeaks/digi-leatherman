package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DateParseRequest struct {
	Value string `json:"value"`
}

type DateParseResponse struct {
	ISO      string `json:"iso"`
	Unix     int64  `json:"unix"`
	UnixMs   int64  `json:"unixMs"`
	RFC2822  string `json:"rfc2822"`
	Relative string `json:"relative"`
}

type TimezoneRequest struct {
	Value    string `json:"value"`
	FromZone string `json:"fromZone"`
	ToZone   string `json:"toZone"`
}

type TimezoneResponse struct {
	Result   string `json:"result"`
	FromTime string `json:"fromTime"`
	ToTime   string `json:"toTime"`
}

type CountdownRequest struct {
	Value string `json:"value"`
}

type CountdownResponse struct {
	Days    int    `json:"days"`
	Hours   int    `json:"hours"`
	Minutes int    `json:"minutes"`
	Seconds int    `json:"seconds"`
	IsPast  bool   `json:"isPast"`
	Text    string `json:"text"`
}

var dateFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
	"01/02/2006",
	"01/02/2006 15:04:05",
	"02-01-2006",
	"02-01-2006 15:04:05",
	"Jan 2, 2006",
	"Jan 2, 2006 3:04 PM",
	"January 2, 2006",
	"January 2, 2006 3:04 PM",
	"2006/01/02",
	"2006/01/02 15:04:05",
}

func parseDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, fmt.Errorf("empty date value")
	}

	if ts, err := strconv.ParseInt(value, 10, 64); err == nil {
		if ts > 1e12 {
			return time.UnixMilli(ts), nil
		}
		return time.Unix(ts, 0), nil
	}

	if ts, err := strconv.ParseFloat(value, 64); err == nil {
		sec := int64(ts)
		nsec := int64((ts - float64(sec)) * 1e9)
		return time.Unix(sec, nsec), nil
	}

	for _, layout := range dateFormats {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", value)
}

func formatRelative(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	isPast := diff > 0

	if !isPast {
		diff = -diff
	}

	seconds := int(diff.Seconds())
	minutes := int(diff.Minutes())
	hours := int(diff.Hours())
	days := int(diff.Hours() / 24)
	weeks := days / 7
	months := days / 30
	years := days / 365

	var unit string
	var value int

	switch {
	case years >= 1:
		unit = "year"
		value = years
	case months >= 1:
		unit = "month"
		value = months
	case weeks >= 1:
		unit = "week"
		value = weeks
	case days >= 1:
		unit = "day"
		value = days
	case hours >= 1:
		unit = "hour"
		value = hours
	case minutes >= 1:
		unit = "minute"
		value = minutes
	default:
		unit = "second"
		value = seconds
	}

	if value != 1 {
		unit += "s"
	}

	if isPast {
		return fmt.Sprintf("%d %s ago", value, unit)
	}
	return fmt.Sprintf("in %d %s", value, unit)
}

func DateParse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DateParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	t, err := parseDate(req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := DateParseResponse{
		ISO:      t.UTC().Format(time.RFC3339),
		Unix:     t.Unix(),
		UnixMs:   t.UnixMilli(),
		RFC2822:  t.Format(time.RFC1123Z),
		Relative: formatRelative(t),
	}

	writeJSON(w, resp)
}

func DateTimezone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TimezoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	t, err := parseDate(req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fromLoc, err := time.LoadLocation(req.FromZone)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid fromZone: %s", req.FromZone), http.StatusBadRequest)
		return
	}

	toLoc, err := time.LoadLocation(req.ToZone)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid toZone: %s", req.ToZone), http.StatusBadRequest)
		return
	}

	fromTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), fromLoc)
	toTime := fromTime.In(toLoc)

	resp := TimezoneResponse{
		Result:   toTime.Format(time.RFC3339),
		FromTime: fromTime.Format(time.RFC3339),
		ToTime:   toTime.Format(time.RFC3339),
	}

	writeJSON(w, resp)
}

func DateCountdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CountdownRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	t, err := parseDate(req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	diff := t.Sub(now)
	isPast := diff < 0

	if isPast {
		diff = -diff
	}

	totalSeconds := int(math.Abs(diff.Seconds()))
	days := totalSeconds / 86400
	hours := (totalSeconds % 86400) / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	var text string
	if isPast {
		text = fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds ago", days, hours, minutes, seconds)
	} else {
		text = fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds remaining", days, hours, minutes, seconds)
	}

	resp := CountdownResponse{
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
		IsPast:  isPast,
		Text:    text,
	}

	writeJSON(w, resp)
}
