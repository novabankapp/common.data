package base

import (
	"reflect"

	"github.com/google/uuid"
)

type Entity interface {
}

func FillDefaults[E Entity](entity E) {
	metaValue := reflect.ValueOf(entity).Elem()
	if metaValue.Type() == reflect.TypeOf("") {
		metaValue.FieldByName("ID").SetString(uuid.New().String())
	}
}
