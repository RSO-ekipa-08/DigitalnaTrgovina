package pgutil

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// String converts a string to pgtype.Text
func String(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}

// StringPtr converts a string pointer to pgtype.Text
func StringPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{
		String: *s,
		Valid:  true,
	}
}

// FromString converts pgtype.Text to string
func FromString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// UUID converts uuid.UUID to pgtype.UUID
func UUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}

// UUIDFromString converts string to pgtype.UUID
func UUIDFromString(s string) (pgtype.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{Valid: false}, err
	}
	return UUID(id), nil
}

// StringFromUUID converts pgtype.UUID to string
func StringFromUUID(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	return uuid.UUID(id.Bytes).String()
}

// FromUUID converts pgtype.UUID to uuid.UUID
func FromUUID(id pgtype.UUID) uuid.UUID {
	if !id.Valid {
		return uuid.Nil
	}
	return uuid.UUID(id.Bytes)
}

// Float64 converts float64 to pgtype.Numeric
func Float64(f float64) pgtype.Numeric {
	var num pgtype.Numeric
	err := num.Scan(fmt.Sprintf("%f", f))
	if err != nil {
		return pgtype.Numeric{Valid: false}
	}
	return num
}

// Float64Ptr converts float64 pointer to pgtype.Numeric
func Float64Ptr(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{Valid: false}
	}
	return Float64(*f)
}

// FromFloat64 converts pgtype.Numeric to float64
func FromFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, err := n.Float64Value()
	if err != nil || !f.Valid {
		return 0
	}
	return f.Float64
}

// Int32 converts int32 to pgtype.Int4
func Int32(i int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: i,
		Valid: true,
	}
}

// FromInt32 converts pgtype.Int4 to int32
func FromInt32(i pgtype.Int4) int32 {
	if !i.Valid {
		return 0
	}
	return i.Int32
}

// Bool converts bool to pgtype.Bool
func Bool(b bool) pgtype.Bool {
	return pgtype.Bool{
		Bool:  b,
		Valid: true,
	}
}

// FromBool converts pgtype.Bool to bool
func FromBool(b pgtype.Bool) bool {
	if !b.Valid {
		return false
	}
	return b.Bool
}

// StringArray converts string slice to pgtype.Array[string]
func StringArray(arr []string) pgtype.Array[string] {
	if arr == nil {
		return pgtype.Array[string]{Valid: false}
	}
	return pgtype.Array[string]{
		Elements: arr,
		Valid:    true,
	}
}

// FromStringArray converts pgtype.Array[string] to string slice
func FromStringArray(arr pgtype.Array[string]) []string {
	if !arr.Valid {
		return nil
	}
	return arr.Elements
}

// Time converts time.Time to pgtype.Timestamptz
func Time(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}

// TimePtr converts *time.Time to pgtype.Timestamptz
func TimePtr(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{
		Time:  *t,
		Valid: true,
	}
}

// FromTime converts pgtype.Timestamptz to time.Time
func FromTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// NowTime returns current time as pgtype.Timestamptz
func NowTime() pgtype.Timestamptz {
	return Time(time.Now())
}
