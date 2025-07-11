package mapping

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// toPgText преобразует string в pgtype.Text
func toPgText(val string) pgtype.Text {
	return pgtype.Text{
		String: val,
		Valid:  val != "",
	}
}

func toPgUUID(s string) pgtype.UUID {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{
		Bytes: parsed,
		Valid: true,
	}
}

// toPgTimestamp преобразует time.Time в pgtype.Timestamp
func toPgTimestamp(t time.Time) pgtype.Timestamptz {
	if t.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}

// toPgBool преобразует bool в pgtype.Bool
func toPgBool(val bool) pgtype.Bool {
	return pgtype.Bool{
		Bool:  val,
		Valid: true,
	}
}

// toPgInt16 преобразует int16 в pgtype.Int2
func toPgInt16(val int16) pgtype.Int2 {
	return pgtype.Int2{
		Int16: val,
		Valid: true,
	}
}

// toPgInt32 преобразует int32 в pgtype.Int4
func toPgInt32(val int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: val,
		Valid: true,
	}
}

// toPgInt64 преобразует int64 в pgtype.Int8
func toPgInt64(val int64) pgtype.Int8 {
	return pgtype.Int8{
		Int64: val,
		Valid: true,
	}
}
