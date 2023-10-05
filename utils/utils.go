package utils

import (
	"bytes"
	"encoding/json"
	"github.com/oklog/ulid/v2"
	"math/rand"
	"strings"
	"time"
)

func GenerateId() string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	return strings.ToLower(ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String())
}

func UnPack(in interface{}, target interface{}) error {
	var e1 error
	var b []byte
	switch in := in.(type) {
	case []byte:
		b = in
	// Do something.
	default:
		// Do the rest.
		b, e1 = json.Marshal(in)
		if e1 != nil {
			return e1
		}
	}

	buf := bytes.NewBuffer(b)
	enc := json.NewDecoder(buf)
	enc.UseNumber()
	if err := enc.Decode(&target); err != nil {
		return err
	}
	return nil
}
