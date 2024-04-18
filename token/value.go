package token

import (
	"encoding"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var _ encoding.BinaryMarshaler = (*Value)(nil)
var _ encoding.BinaryUnmarshaler = (*Value)(nil)

type Value struct {
	UserId    int64
	Token     string
	CreatedAt time.Time
	ExpiredAt time.Time
	Extras    []byte
}

func NewValue(userId int64) *Value {
	return &Value{
		UserId:    userId,
		Token:     uuid.New().String(),
		CreatedAt: time.Now(),
	}
}

// IsExpired checks if the token is expired
func (v *Value) IsExpired() bool {
	return v.ExpiredAt.Before(time.Now())
}

// IsValid checks if the token is valid
func (v *Value) IsValid() bool {
	return !v.IsExpired() && v.Token != "" && v.UserId > 0
}

// Get gets the value of the key
func (v *Value) Get(key string) (value gjson.Result) {
	return gjson.GetBytes(v.Extras, key)
}

// Set sets the value of the key
func (v *Value) Set(key string, value any) (err error) {
	v.Extras, err = sjson.SetBytes(v.Extras, key, value)
	return
}

// MarshalBinary implements the encoding.BinaryMarshaler interface
func (v *Value) MarshalBinary() (data []byte, err error) {
	return json.Marshal(v)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (v *Value) UnmarshalBinary(data []byte) (err error) {
	return json.Unmarshal(data, v)
}
