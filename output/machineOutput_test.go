package output

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)


/*
        Tests for GetMachineOutputJson
*/
func TestGetMachineOutputJsonSuccess(t *testing.T) {
  output := MachineOutput{
    ExitCode: 0,
    StatusMessage: "Successfully created cluster",
  }

  expectJsonString := `
    {
      "exitCode": 0,
      "statusMessage": "Successfully created cluster"
    }
  `

  actualJsonString, exitCode := output.GetMachineOutputJson()

  assert.Equal(t, exitCode, 0)
  assert.Equal(t, expectJsonString, actualJsonString)
}

func TestGetMachineOutputJsonError(t *testing.T) {
  output := MachineOutput{}

  // Use reflection and unsafe to inject an invalid map key type
	v := reflect.ValueOf(&output).Elem().FieldByName("DetailedMachineStatus")
	ptr := unsafe.Pointer(v.UnsafeAddr())
	v = reflect.NewAt(v.Type(), ptr).Elem()
	v.Set(reflect.ValueOf(map[int]string{
		123: "bad key type",
	})) // map[int]string is invalid for JSON because keys aren't strings

  expectedJsonString := "{\"exitCode\": 100, \"errorMessage\": \"Error unmarshaling machine output\"}"

  actualJsonString, exitCode := output.GetMachineOutputJson()

  assert.Equal(t, exitCode, 100)
  assert.Equal(t, expectedJsonString, actualJsonString)

}
