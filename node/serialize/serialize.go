package serialize

import (
	"github.com/tendermint/go-amino"
)

type Channel int

const (
	CLIENT Channel = iota
	PERSISTENT
	NETWORK
)

var aminoCodec *amino.Codec
func init() {
	aminoCodec = amino.NewCodec()
}


type Serializer interface {
	Serialize(obj interface{}) ([]byte, error)
	Deserialize(d []byte, obj interface{}) error
}

func GetSerializer(channel Channel, args ...interface{}) (Serializer, error) {

	switch channel {

	case CLIENT:
		return &jsonStrategy{}, nil

	case PERSISTENT:

		return &jsonStrategy{}, nil

	case NETWORK:

		return &jsonStrategy{}, nil
	default:
		return nil, ErrIncorrectChannel
	}
}

func RegisterInterface(obj interface{}) {
	aminoCodec.RegisterInterface(obj, &amino.InterfaceOptions{AlwaysDisambiguate: true})
}

func RegisterConcrete(obj interface{}, name string) {
	aminoCodec.RegisterConcrete(obj, name, nil)
}