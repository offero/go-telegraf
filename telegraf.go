package telegraf

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"
)

type Metric struct {
	name   string
	tags   map[string]string
	fields map[string]interface{}
	t      time.Time
}

func fieldValueToString(val interface{}) string {
	switch val := val.(type) {
	case int8, int16, int32, int64, int:
		return fmt.Sprintf("%di", val)
	case float32, float64:
		return fmt.Sprintf("%f", val)
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

	n := len(m.tags)
	i := 0
	for k, v := range m.tags {
		i++
		buffer.WriteString(Escape(k, TagKey))
		buffer.WriteString("=")
		buffer.WriteString(Escape(v, TagVal))
		if i < n {
			buffer.WriteString(",")
		}
	}

	buffer.WriteString(" ")

	n = len(m.tags)
	i = 0
	for k, v := range m.fields {
		i++
		buffer.WriteString(Escape(k, FieldKey))
		buffer.WriteString("=")
		buffer.WriteString(fieldValueToString(v))
		if i < n {
			buffer.WriteString(",")
		}
	}

	buffer.WriteString(" ")

	buffer.WriteString(strconv.FormatInt(m.t.UnixNano(), 10))

	buffer.WriteString("\n")

	return buffer.Bytes()
}

func NewMetric(name string, tags map[string]string, fields map[string]interface{}) Metric {
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
	body := m.Serialize()
}

func (c *Client) Close() error {
	c.conn.Close()
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
