package codecs

import "github.com/awillis/bitfan/commons"

type Decoder interface {
	Decode(*interface{}) error
	SetOptions(map[string]interface{}, commons.Logger, string) error
	More() bool
	Buffer() []byte
}
