package dglogger

import "reflect"

var SecretStringType reflect.Type

func init() {
	SecretStringType = GetInterfaceType[SecretString](nil)
}

type SecretString interface {
	Secret() string
}

func GetInterfaceType[T any](iface *T) reflect.Type {
	ifaceType := reflect.TypeOf(iface).Elem()
	return ifaceType
}
