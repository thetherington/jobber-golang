package utils

import (
	"time"

	"google.golang.org/genproto/googleapis/type/datetime"
)

// help function to create a string pointer from a function return
func Ptr(s string) *string { return &s }

// help function to create a int pointer from a function return
func PtrI(i int) *int { return &i }

// help function to create a int pointer from a function return
func PtrF64(f float64) *float64 { return &f }

// helper function to create a time pointer from a function return
func ToTimePtr(t time.Time) *time.Time { return &t }

// helper function to convert gRPC DateTime to Go time.time
func ToTime(dt *datetime.DateTime) *time.Time {
	if dt == nil {
		now := time.Now()

		dt = &datetime.DateTime{
			Year:    int32(now.Year()),
			Month:   int32(now.Month()),
			Day:     int32(now.Day()),
			Hours:   int32(now.Hour()),
			Minutes: int32(now.Minute()),
			Seconds: int32(now.Second()),
			Nanos:   int32(now.Nanosecond()),
		}

	}

	res := time.Date(int(dt.Year), time.Month(dt.Month), int(dt.Day),
		int(dt.Hours), int(dt.Minutes), int(dt.Seconds), int(dt.Nanos),
		time.UTC)

	return &res
}

// helper function to convert gRPC DateTime to Go time.time
func ToTimeOrNil(dt *datetime.DateTime) *time.Time {
	if dt == nil {
		return nil
	}

	res := time.Date(int(dt.Year), time.Month(dt.Month), int(dt.Day),
		int(dt.Hours), int(dt.Minutes), int(dt.Seconds), int(dt.Nanos),
		time.UTC)

	return &res
}

// helper function to convert Go time.time to gRPC DateTime
func ToDateTime(t *time.Time) *datetime.DateTime {
	if t == nil {
		now := time.Now()

		dt := &datetime.DateTime{
			Year:    int32(now.Year()),
			Month:   int32(now.Month()),
			Day:     int32(now.Day()),
			Hours:   int32(now.Hour()),
			Minutes: int32(now.Minute()),
			Seconds: int32(now.Second()),
			Nanos:   int32(now.Nanosecond()),
		}

		return dt
	}

	res := &datetime.DateTime{
		Year:    int32(t.Year()),
		Month:   int32(t.Month()),
		Day:     int32(t.Day()),
		Hours:   int32(t.Hour()),
		Minutes: int32(t.Minute()),
		Seconds: int32(t.Second()),
		Nanos:   int32(t.Nanosecond()),
	}

	return res
}

// helper function to convert Go time.time to gRPC DateTime
func ToDateTimeOrNil(t *time.Time) *datetime.DateTime {
	if t == nil {
		return nil
	}

	res := &datetime.DateTime{
		Year:    int32(t.Year()),
		Month:   int32(t.Month()),
		Day:     int32(t.Day()),
		Hours:   int32(t.Hour()),
		Minutes: int32(t.Minute()),
		Seconds: int32(t.Second()),
		Nanos:   int32(t.Nanosecond()),
	}

	return res
}

// helper function to get current time as gRPC DateTime
func CurrentDatetime() *datetime.DateTime {
	now := time.Now()

	return &datetime.DateTime{
		Year:       int32(now.Year()),
		Month:      int32(now.Month()),
		Day:        int32(now.Day()),
		Hours:      int32(now.Hour()),
		Minutes:    int32(now.Minute()),
		Seconds:    int32(now.Second()),
		Nanos:      int32(now.Nanosecond()),
		TimeOffset: &datetime.DateTime_UtcOffset{},
	}
}
