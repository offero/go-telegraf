package telegraf

// Taken https://raw.githubusercontent.com/influxdata/telegraf/master/metric/escape.go

import (
	"strings"
)

var (
	// escaper is for escaping:
	//   - tag keys
	//   - tag values
	//   - field keys
	// see https://docs.influxdata.com/influxdb/v1.0/write_protocols/line_protocol_tutorial/#special-characters-and-keywords
	escaper = strings.NewReplacer(`,`, `\,`, `"`, `\"`, ` `, `\ `, `=`, `\=`)
	// unEscaper = strings.NewReplacer(`\,`, `,`, `\"`, `"`, `\ `, ` `, `\=`, `=`)

	// nameEscaper is for escaping measurement names only.
	// see https://docs.influxdata.com/influxdb/v1.0/write_protocols/line_protocol_tutorial/#special-characters-and-keywords
	nameEscaper = strings.NewReplacer(`,`, `\,`, ` `, `\ `)
	// nameUnEscaper = strings.NewReplacer(`\,`, `,`, `\ `, ` `)

	// stringFieldEscaper is for escaping string field values only.
	// see https://docs.influxdata.com/influxdb/v1.0/write_protocols/line_protocol_tutorial/#special-characters-and-keywords
	stringFieldEscaper = strings.NewReplacer(`"`, `\"`)
	// stringFieldUnEscaper = strings.NewReplacer(`\"`, `"`)
)

// MetricPart specifies a part of an Influx metric
type MetricPart int

const (
	// Name is the string name component of the metric
	Name MetricPart = iota
	// FieldKey is the string key part of a field
	FieldKey
	// FieldVal is the int, float, bool or string value part of a field
	FieldVal
	// TagKey is the string key part of a tag
	TagKey
	// TagVal is the string value part of a tag
	TagVal
)

// Escape adds necessary escape characters to a specified part of the Influx
// metric. This must be done before submitting to Influx.
func Escape(s string, t MetricPart) string {
	switch t {
	case FieldKey, TagKey, TagVal:
		return escaper.Replace(s)
	case Name:
		return nameEscaper.Replace(s)
	case FieldVal:
		return stringFieldEscaper.Replace(s)
	}
	return s
}
