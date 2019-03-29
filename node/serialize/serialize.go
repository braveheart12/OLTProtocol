package serialize

import (
	"github.com/tendermint/go-amino"
)

type Channel int

const (
	CLIENT Channel = iota
	PERSISTENT
	NETWORK
	JSON
)

var aminoCodec *amino.Codec
var JSONSzr Serializer

func init() {
	aminoCodec = amino.NewCodec()

	JSONSzr, _ = GetSerializer(JSON)
}

type Serializer interface {
	Serialize(obj interface{}) ([]byte, error)
	Deserialize(d []byte, obj interface{}) error
}

// GetSerializer for a channel
func GetSerializer(channel Channel, args ...interface{}) (Serializer, error) {

	switch channel {

	case CLIENT:
		return &msgpackStrategy{}, nil

	case PERSISTENT:
		return &msgpackStrategy{}, nil

	case NETWORK:
		return &msgpackStrategy{}, nil

	case JSON:
		return &jsonStrategy{}, nil

	default:
		return nil, ErrIncorrectChannel
	}
}


// functions to register types
func RegisterInterface(obj interface{}) {
	aminoCodec.RegisterInterface(obj, &amino.InterfaceOptions{AlwaysDisambiguate: true})
}

func RegisterConcrete(obj interface{}, name string) {
	aminoCodec.RegisterConcrete(obj, name, nil)
}
