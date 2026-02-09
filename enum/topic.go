package enum

import (
	"bytes"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type Topic uint8

const (
	TopicNone Topic = iota
	TopicDailyLife
	TopicFood
	TopicTravel
	TopicWork
	TopicSchool
	TopicShopping
	TopicHealth
	TopicFamily
	TopicWeather
	TopicHobby
	TopicTransportation
	TopicTechnology
	TopicNature
	TopicCulture
	TopicSports
	TopicTimeDate
	TopicFeelings
	TopicLocation
	TopicOther
)

var TopicName = map[Topic]string{
	TopicNone:           "",
	TopicDailyLife:      "daily_life",
	TopicFood:           "food",
	TopicTravel:         "travel",
	TopicWork:           "work",
	TopicSchool:         "school",
	TopicShopping:       "shopping",
	TopicHealth:         "health",
	TopicFamily:         "family",
	TopicWeather:        "weather",
	TopicHobby:          "hobby",
	TopicTransportation: "transportation",
	TopicTechnology:     "technology",
	TopicNature:         "nature",
	TopicCulture:        "culture",
	TopicSports:         "sports",
	TopicTimeDate:       "time",
	TopicFeelings:       "feelings",
	TopicLocation: 		 "location",
	TopicOther:          "other",
}

var TopicValue = func() map[string]Topic {
	res := map[string]Topic{}
	for k, v := range TopicName {
		res[v] = k
	}
	return res
}()

func (e *Topic) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)

	// number
	if len(data) > 0 && data[0] != '"' {
		i, err := strconv.ParseUint(string(data), 10, 8)
		if err != nil {
			return fmt.Errorf("invalid enum number: %s", data)
		}

		v := Topic(i)
		if _, ok := TopicName[v]; !ok {
			return fmt.Errorf("enum '%d' is not registered, must be one of: %v", i, e.EnumDescriptions())
		}

		*e = v
		return nil
	}

	// string
	data = bytes.Trim(data, "\"")
	v, ok := TopicValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not registered, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v
	return nil
}

func (e Topic) MarshalJSON() ([]byte, error) {
	v, ok := TopicName[e]
	if !ok {
		return []byte(`""`), nil
	}

	return []byte(`"` + v + `"`), nil
}

func (e *Topic) UnmarshalText(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := TopicValue[string(data)]

	if !ok {
		return fmt.Errorf("enum '%s' is not registered, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v
	return nil
}

func (e Topic) MarshalText() ([]byte, error) {
	v, ok := TopicName[e]
	if !ok {
		return []byte(`""`), nil
	}

	return []byte(`"` + v + `"`), nil
}

func (e *Topic) EnumDescriptions() []string {
	r := []string{}
	for k := range TopicValue {
		r = append(r, k)
	}
	return r
}

func (e Topic) MarshalBSONValue() (bsontype.Type, []byte, error) {
	v, ok := TopicName[e]
	if !ok {
		return bson.MarshalValue("")
	}
	return bson.MarshalValue(v)
}

func (e *Topic) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var s string
	if err := bson.UnmarshalValue(t, data, &s); err != nil {
		return err
	}

	v, ok := TopicValue[s]
	if !ok {
		return fmt.Errorf("invalid Topic: %s", s)
	}

	*e = v
	return nil
}
