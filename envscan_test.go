package envscan

import (
	"errors"
	"testing"
)

func TestReadEnvironmentSuccess(t *testing.T) {
	type MockConfig struct {
		MockString      string   `env:"MOCK_STRING"`
		MockSliceString []string `env:"MOCK_SLICE_STRING"`
	}

	defaultEnvData := map[string]string{
		"MOCK_SLICE_STRING": "create,test",
	}

	expectedScannedObj := &MockConfig{
		MockString:      "yupi",
		MockSliceString: []string{"create", "test"},
	}

	scanningObj := &MockConfig{}

	t.Setenv("MOCK_STRING", "yupi")

	err := ReadEnvironment(scanningObj, defaultEnvData)
	if err != nil {
		t.Errorf("failed to scan env, err: %v", err)
	}

	if expectedScannedObj.MockString != scanningObj.MockString {
		t.Errorf("expected: %s, got: %s", expectedScannedObj.MockString, scanningObj.MockString)
	}

	if len(expectedScannedObj.MockSliceString) != len(scanningObj.MockSliceString) {
		t.Errorf("expected len of slice string attr: %d, got: %d", len(expectedScannedObj.MockSliceString), len(scanningObj.MockSliceString))
	}

	for i := range len(expectedScannedObj.MockSliceString) {
		if expectedScannedObj.MockSliceString[i] != scanningObj.MockSliceString[i] {
			t.Errorf("expected elem of mock slice string attr: %s, got: %s", expectedScannedObj.MockSliceString[i], scanningObj.MockSliceString[i])
		}
	}
}

func TestReadEnvironmentErrNilPointer(t *testing.T) {
	err := ReadEnvironment(nil, make(map[string]string))
	if err == nil {
		t.Errorf("expected %v, got nil", ErrNilPointerDeference)
	}

	if !errors.Is(err, ErrNilPointerDeference) {
		t.Errorf("expected %v, got %v", ErrNilPointerDeference, err)
	}
}

func TestReadEnvironmentErrNotPtr(t *testing.T) {
	type MockConfig struct {
		MockString      string   `env:"MOCK_STRING"`
		MockSliceString []string `env:"MOCK_SLICE_STRING"`
	}

	scanningObj := MockConfig{}

	err := ReadEnvironment(scanningObj, make(map[string]string))
	if !errors.Is(err, ErrVMustBePtr) {
		t.Errorf("expected %v, got %v", ErrVMustBePtr, err)
	}
}
