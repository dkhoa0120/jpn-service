package enum

import (
	"bytes"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type JLPTLevel uint8

const (
	JLPTN5 JLPTLevel = iota + 1
	JLPTN4
	JLPTN3
	JLPTN2
	JLPTN1
)

var JLPTLevelName = map[JLPTLevel]string{
	JLPTN5: "N5",
	JLPTN4: "N4",
	JLPTN3: "N3",
	JLPTN2: "N2",
	JLPTN1: "N1",
}

var JLPTLevelValue = func() map[string]JLPTLevel {
	res := map[string]JLPTLevel{}
	for k, v := range JLPTLevelName {
		res[v] = k
	}
	return res
}()

func (e *JLPTLevel) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)

	if len(data) > 0 && data[0] != '"' {
		i, err := strconv.ParseUint(string(data), 10, 8)
		if err != nil {
			return fmt.Errorf("invalid enum number: %s", data)
		}

		v := JLPTLevel(i)
		if _, ok := JLPTLevelName[v]; !ok {
			return fmt.Errorf(
				"enum '%d' is not registered, must be one of: %v",
				i,
				e.EnumDescriptions(),
			)
		}

		*e = v
		return nil
	}

	data = bytes.Trim(data, "\"")
	v, ok := JLPTLevelValue[string(data)]
	if !ok {
		return fmt.Errorf(
			"enum '%s' is not registered, must be one of: %v",
			data,
			e.EnumDescriptions(),
		)
	}

	*e = v
	return nil
}

func (e JLPTLevel) MarshalJSON() ([]byte, error) {
	v, ok := JLPTLevelName[e]
	if !ok {
		return []byte("\"\""), nil
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)

	return buffer.Bytes(), nil
}

func (e *JLPTLevel) UnmarshalText(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := JLPTLevelValue[string(data)]

	if !ok {
		return fmt.Errorf(
			"enum '%s' is not registered, must be one of: %v",
			data,
			e.EnumDescriptions(),
		)
	}

	*e = v
	return nil
}

func (e JLPTLevel) MarshalText() ([]byte, error) {
	v, ok := JLPTLevelName[e]
	if !ok {
		return []byte("\"\""), nil
	}

	return []byte(`"` + v + `"`), nil
}

func (e *JLPTLevel) EnumDescriptions() []string {
	r := []string{}
	for k := range JLPTLevelValue {
		r = append(r, k)
	}
	return r
}

func (e JLPTLevel) MarshalBSONValue() (bsontype.Type, []byte, error) {
	v, ok := JLPTLevelName[e]
	if !ok {
		return bson.MarshalValue("")
	}
	return bson.MarshalValue(v)
}

func (e *JLPTLevel) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var s string
	if err := bson.UnmarshalValue(t, data, &s); err != nil {
		return err
	}

	v, ok := JLPTLevelValue[s]
	if !ok {
		return fmt.Errorf("invalid JLPTLevel: %s", s)
	}

	*e = v
	return nil
}
