package uid

import "database/sql/driver"

// Value implements driver.Valuer.
func (u UUID) Value() (driver.Value, error) { return u.String(), nil }

// Scan implements sql.Scanner. Returns ErrInvalid for malformed input or SQL NULL; use *UUID for nullable columns.
func (u *UUID) Scan(src any) error {
	switch v := src.(type) {
	case string:
		if id, ok := Parse(v); ok {
			*u = id
			return nil
		}
	case []byte:
		if id, ok := Parse(string(v)); ok {
			*u = id
			return nil
		}
	}
	return ErrInvalid
}
