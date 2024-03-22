package int64string

import (
	"encoding/json"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

type Int64String int64

func Int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func (m Int64String) MarshalJSON() ([]byte, error) {
	return json.Marshal(Int64ToString(int64(m)))
}

func (m *Int64String) UnmarshalJSON(data []byte) error {
	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err != nil {
		return err
	}

	intValue, err := strconv.ParseInt(stringValue, 10, 64)
	if err != nil {
		return err
	}

	*m = Int64String(intValue)
	return nil
}

func (m Int64String) MarshalBSONValue() (bson.RawValue, error) {
	return bson.RawValue{Type: bson.TypeString, Value: bson.Raw(Int64ToString(int64(m)))}, nil
}
