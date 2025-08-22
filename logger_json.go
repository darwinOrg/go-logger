package dglogger

import (
	"reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

var jsonLogger jsoniter.API

func init() {
	jsonLogger = jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
	}.Froze()
	strValEncoder := jsonLogger.EncoderOf(reflect2.TypeOf(""))
	jsonLogger.RegisterExtension(&secretEncoderExtension{
		strValueEncoder: strValEncoder,
	})
}

func Json(value any) ([]byte, error) {
	return jsonLogger.Marshal(value)
}

type secretEncoderExtension struct {
	jsoniter.EncoderExtension
	strValueEncoder jsoniter.ValEncoder
}

type secretValEncoder struct {
	strValueEncoder jsoniter.ValEncoder
	typ             reflect.Type
}

func (sve *secretValEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (sve *secretValEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	newObj := reflect.NewAt(sve.typ, ptr)
	var iface interface{}
	isPointer := false
	if sve.typ.Kind() == reflect.Pointer {
		iface = newObj.Elem().Interface()
		isPointer = true
	} else {
		iface = newObj.Elem().Addr().Interface()
	}
	secret := iface.(SecretString)

	isNilValue := reflect.ValueOf(secret).IsNil()
	if isPointer && isNilValue {
		stream.WriteNil()
		return
	}
	strValue := ""
	if !isNilValue {
		strValue = secret.Secret()
	}
	sve.strValueEncoder.Encode(unsafe.Pointer(&strValue), stream)
}

func (se *secretEncoderExtension) CreateEncoder(typ reflect2.Type) jsoniter.ValEncoder {
	typ1 := typ.Type1()

	if typ1.Implements(SecretStringType) {
		return &secretValEncoder{
			strValueEncoder: se.strValueEncoder,
			typ:             typ1,
		}
	}
	if typ1.Kind() == reflect.Pointer {
		if reflect.PtrTo(typ1).Implements(SecretStringType) {
			return &secretValEncoder{
				strValueEncoder: se.strValueEncoder,
				typ:             typ1,
			}
		}
	}

	return se.EncoderExtension.CreateEncoder(typ)
}
