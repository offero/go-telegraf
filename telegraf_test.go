package telegraf

import (
	"bytes"
	"strconv"
	"testing"
	"time"
)

func TestMetric(t *testing.T) {
	tags := make([]Tag, 0)
	// Add them out of lexicographic order
	tags = append(tags, Tag{"b", "tagbval"})
	tags = append(tags, Tag{"a", "tagaval"})

	fields := make([]Field, 0)
	// Add them out of lexicographic order
	fields = append(fields, Field{"s", "string field val"})
	fields = append(fields, Field{"d", 123})
	fields = append(fields, Field{"f", 456.456})
	m := NewMetric("metric1", tags, fields)
	t1 := time.Now()
	m.SetTime(t1)

	output := m.Serialize()

	expectedOutput := "metric1,a=tagaval,b=tagbval"
	expectedOutput += " "
	expectedOutput += "d=123i,f=456.456,s=\"string field val\""
	expectedOutput += " "
	expectedOutput += strconv.FormatInt(t1.UnixNano(), 10)
	expectedOutput += "\n"

	comp := bytes.Compare(output, []byte(expectedOutput))
	if comp != 0 {
		t.Errorf("Output does not match (%d):\n%s\n%s", comp, output, expectedOutput)
	}
}
