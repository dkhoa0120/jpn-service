package enum

import (
	"bytes"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type VocabularyType uint8

const (
	VocabularyTypeNone VocabularyType = iota
	VocabularyTypeKatakana
	VocabularyTypeHiragana
)

var VocabularyTypeName = map[VocabularyType]string{
	VocabularyTypeNone:     "",
	VocabularyTypeKatakana: "ktkn",
	VocabularyTypeHiragana: "hrgn",
}

var VocabularyTypeValue = func() map[string]VocabularyType {
	res := map[string]VocabularyType{}
	for k, v := range VocabularyTypeName {
		res[v] = k
	}
	return res
}()

func (e *VocabularyType) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)

	if len(data) > 0 && data[0] != '"' {
		i, err := strconv.ParseUint(string(data), 10, 8)
		if err != nil {
			return fmt.Errorf("invalid enum number: %s", data)
		}

		v := VocabularyType(i)
		if _, ok := VocabularyTypeName[v]; !ok {
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
	v, ok := VocabularyTypeValue[string(data)]
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

func (e VocabularyType) MarshalJSON() ([]byte, error) {
	v, ok := VocabularyTypeName[e]
	if !ok {
		return []byte("\"\""), nil
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)

	return buffer.Bytes(), nil
}

func (e *VocabularyType) UnmarshalText(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := VocabularyTypeValue[string(data)]

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

func (e VocabularyType) MarshalText() ([]byte, error) {
	v, ok := VocabularyTypeName[e]
	if !ok {
		return []byte("\"\""), nil
	}

	return []byte(`"` + v + `"`), nil
}

func (e *VocabularyType) EnumDescriptions() []string {
	r := []string{}
	for k := range VocabularyTypeValue {
		r = append(r, k)
	}
	return r
}

func (e VocabularyType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	v, ok := VocabularyTypeName[e]
	if !ok {
		return bson.MarshalValue("")
	}
	return bson.MarshalValue(v)
}

func (e *VocabularyType) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var s string
	if err := bson.UnmarshalValue(t, data, &s); err != nil {
		return err
	}

	v, ok := VocabularyTypeValue[s]
	if !ok {
		return fmt.Errorf("invalid VocabularyType: %s", s)
	}

	*e = v
	return nil
}
