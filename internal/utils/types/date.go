package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type OnlyDate time.Time

const layout = "2006-01-02"

func (d *OnlyDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		return nil
	}
	t, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	*d = OnlyDate(t)
	return nil
}

func (d OnlyDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(d).Format(layout))), nil
}

func (d *OnlyDate) Scan(value any) error {
	if value == nil {
		*d = OnlyDate(time.Time{})
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*d = OnlyDate(v)
		return nil
	case []byte:
		t, err := time.Parse(layout, string(v))
		if err != nil {
			return err
		}
		*d = OnlyDate(t)
		return nil
	case string:
		t, err := time.Parse(layout, v)
		if err != nil {
			return err
		}
		*d = OnlyDate(t)
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into OnlyDate", value)
	}
}

func (d OnlyDate) Value() (driver.Value, error) {
	return time.Time(d), nil
}

func (d OnlyDate) ToTime() time.Time {
	return time.Time(d)
}
