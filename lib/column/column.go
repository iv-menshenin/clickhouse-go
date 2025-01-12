package column

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2/lib/binary"
)

type Type string

func (t Type) params() string {
	switch start, end := strings.Index(string(t), "("), strings.LastIndex(string(t), ")"); {
	case len(t) == 0, start <= 0, end <= 0, end < start:
		return ""
	default:
		return string(t[start+1 : end])
	}
}

type Error struct {
	ColumnType string
	Err        error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.ColumnType, e.Err)
}

type ColumnConverterError struct {
	Op       string
	Hint     string
	From, To string
}

func (e *ColumnConverterError) Error() string {
	var hint string
	if len(e.Hint) != 0 {
		hint += ". " + e.Hint
	}
	return fmt.Sprintf("clickhouse [%s]: converting %s to %s is unsupported%s", e.Op, e.From, e.To, hint)
}

type Interface interface {
	Type() Type
	Rows() int
	Row(i int, ptr bool) interface{}
	ScanRow(dest interface{}, row int) error
	Append(v interface{}) (nulls []uint8, err error)
	AppendRow(v interface{}) error
	Decode(decoder *binary.Decoder, rows int) error
	Encode(*binary.Encoder) error
	ScanType() reflect.Type
}

type CustomSerialization interface {
	ReadStatePrefix(*binary.Decoder) error
	WriteStatePrefix(*binary.Encoder) error
}

type UnsupportedColumnType struct {
	t Type
}

func (u *UnsupportedColumnType) Type() Type                          { return u.t }
func (UnsupportedColumnType) Rows() int                              { return 0 }
func (u *UnsupportedColumnType) Row(int, bool) interface{}           { return nil }
func (u *UnsupportedColumnType) ScanRow(interface{}, int) error      { return u }
func (u *UnsupportedColumnType) Append(interface{}) ([]uint8, error) { return nil, u }
func (u *UnsupportedColumnType) AppendRow(interface{}) error         { return u }
func (u *UnsupportedColumnType) Decode(*binary.Decoder, int) error   { return u }
func (u *UnsupportedColumnType) Encode(*binary.Encoder) error        { return u }
func (u *UnsupportedColumnType) ScanType() reflect.Type              { return reflect.TypeOf(nil) }

func (u *UnsupportedColumnType) Error() string {
	return fmt.Sprintf("clickhouse: unsupported column type %q", u.t)
}

var (
	_ error     = (*UnsupportedColumnType)(nil)
	_ Interface = (*UnsupportedColumnType)(nil)
)
