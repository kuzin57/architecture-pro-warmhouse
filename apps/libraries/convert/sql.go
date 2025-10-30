package convert

import "database/sql"

func ToNullString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}

	return sql.NullString{String: *value, Valid: true}
}

func FromNullString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}

	return &value.String
}
