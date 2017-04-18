package telegraf

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type Tag struct {
	key   string
	value string
}

type tagsLexi []Tag

func (s tagsLexi) Len() int      { return len(s) }
func (s tagsLexi) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s tagsLexi) Less(i, j int) bool {
	return bytes.Compare([]byte(s[i].key), []byte(s[j].key)) < 0
}

type Field struct {
	key   string
	value interface{}
}

type fieldsLexi []Field

func (s fieldsLexi) Len() int      { return len(s) }
func (s fieldsLexi) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s fieldsLexi) Less(i, j int) bool {
	return bytes.Compare([]byte(s[i].key), []byte(s[j].key)) < 0
}

type Metric struct {
	name   string
	tags   []Tag
	fields []Field
	t      time.Time
}

func fieldValueToString(val interface{}) string {
	switch val := val.(type) {
	case int8, int16, int32, int64, int:
		return fmt.Sprintf("%di", val)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		switch val {
		case true:
			return "T"
		case false:
			return "F"
		}
	}
	s := Escape(fmt.Sprintf("%s", val), FieldVal)
	return fmt.Sprintf("\"%s\"", s)
}

func (m *Metric) Serialize() []byte {
	var buffer bytes.Buffer

	buffer.WriteString(Escape(m.name, Name))
	buffer.WriteString(",")

	sort.Sort(tagsLexi(m.tags))
	for i, tag := range m.tags {
		buffer.WriteString(Escape(tag.key, TagKey))
		buffer.WriteString("=")
		buffer.WriteString(Escape(tag.value, TagVal))
		if i < len(m.tags)-1 {
			buffer.WriteString(",")
		}
	}

	buffer.WriteString(" ")

	sort.Sort(fieldsLexi(m.fields))
	for i, field := range m.fields {
		buffer.WriteString(Escape(field.key, FieldKey))
		buffer.WriteString("=")
		buffer.WriteString(fieldValueToString(field.value))
		if i < len(m.fields)-1 {
			buffer.WriteString(",")
		}
	}

	buffer.WriteString(" ")

	buffer.WriteString(strconv.FormatInt(m.t.UnixNano(), 10))

	buffer.WriteString("\n")

	return buffer.Bytes()
}

func NewMetric(name string, tags []Tag, fields []Field) Metric {
	m := Metric{name, tags, fields, time.Now()}
	return m
}

func (m *Metric) SetTime(t time.Time) {
	m.t = t
}

type Client struct {
	conn *net.UDPConn
}

func (c *Client) Send(m Metric) error {
	// body := m.Serialize()
	return nil
}

func (c *Client) Close() error {
	c.conn.Close()
	return nil
}

func NewClient(uri string) (*Client, error) {
	parsedUrl, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("Error parsing UDP url [%s]: %s", uri, err)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", parsedUrl.Host)
	if err != nil {
		return nil, fmt.Errorf("Error resolving UDP Address [%s]: %s",
			parsedUrl.Host, err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, fmt.Errorf("Error dialing UDP address [%s]: %s",
			udpAddr.String(), err)
	}

	return &Client{conn}, nil
}
