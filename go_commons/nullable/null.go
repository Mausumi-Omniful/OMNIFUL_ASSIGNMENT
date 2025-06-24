package nullable

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

func NewNullableString(s string) null.String {
	if len(s) == 0 {
		return null.NewString(s, false)
	}
	return null.NewString(s, true)
}

func NewNullableInt(i int64) null.Int {
	if i == 0 {
		return null.NewInt(i, false)
	}
	return null.NewInt(i, true)
}

func NewNullableFloat(i float64) null.Float {
	if i == 0 {
		return null.NewFloat(i, false)
	}
	return null.NewFloat(i, true)
}

func NewNullableTime(t time.Time) null.Time {
	if t.IsZero() {
		return null.NewTime(t, false)
	}
	return null.NewTime(t, true)
}

func NewNullableBool(b bool) null.Bool {
	return null.NewBool(b, true)
}
