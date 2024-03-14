package test

import (
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestCreateUuid(t *testing.T) {

	println(uuid.NewV4().String())
	println(uuid.NewV4().String())
	println(uuid.NewV4().String())
	println(uuid.NewV4().String())
	println(uuid.NewV4().String())
	println(uuid.NewV4().String())
	println(uuid.NewV4().String())
	println(uuid.NewV4().String())
	println(uuid.NewV4().String())
	println(len(uuid.NewV4().String()))
}
