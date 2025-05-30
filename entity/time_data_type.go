package entity

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const timeLayout = "15:04:05"

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.Format(timeLayout))), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	parsed, err := time.Parse(timeLayout, s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

func (t *Time) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		t.Time = v
		return nil
	case []byte:
		parsed, err := time.Parse(timeLayout, string(v))
		if err != nil {
			return err
		}
		t.Time = parsed
		return nil
	case string:
		parsed, err := time.Parse(timeLayout, v)
		if err != nil {
			return err
		}
		t.Time = parsed
		return nil
	}
	return fmt.Errorf("cannot convert %T to Time", value)
}

func (t Time) Value() (driver.Value, error) {
	return t.Format(timeLayout), nil
}
func (Time) GormDataType() string {
	return "time"
}
func (Time) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "TIME"
}
