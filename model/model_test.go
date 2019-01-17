package model

import (
	"testing"

	"github.com/mediacoin-pro/core/common/bin"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	msg := &TestMessageType{"Ловивших педофилов на живца полицейских задержали в Москве"}
	enc := Encode(msg)

	obj, err := Decode(enc)
	dec, ok := obj.(*TestMessageType)

	assert.Equal(t, TestModelType, TypeOf(msg))
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, msg, dec)
}

//----------------------
type TestMessageType struct {
	Msg string
}

const TestModelType = 100500

var _ = RegisterModel(TestModelType, &TestMessageType{})

func (m *TestMessageType) String() string {
	return m.Msg
}

func (d *TestMessageType) Encode() []byte {
	return bin.Encode(d.Msg)
}

func (d *TestMessageType) Decode(data []byte) error {
	return bin.Decode(data, &d.Msg)
}
