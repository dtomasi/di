package utils_test

import (
	"github.com/dtomasi/di/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type ReflectTestStruct struct {
}

func TestGetType(t *testing.T) {
	assert.Equal(t, "ReflectTestStruct", utils.GetType(reflect.TypeOf(ReflectTestStruct{})))
	assert.Equal(t, "*ReflectTestStruct", utils.GetType(reflect.TypeOf(&ReflectTestStruct{})))
}
