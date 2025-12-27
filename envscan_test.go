package envscan

import "testing"

// func TestFoo(t *testing.T) {
//     t.Setenv("XYZ_URL", "http://example.com")
//     /* do your tests here */
// }

func TestReadEnvironment(t *testing.T) {
	type MockConfig struct {
		MockString      string   `env:"mock_string"`
		MockSliceString []string `env:"mock_slice_string"`
	}
}
