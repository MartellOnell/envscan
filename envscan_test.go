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

func TestReadEnvironmentErrNotStructPtr(t *testing.T) {
	type MockConfig struct {
		MockString      string   `env:"MOCK_STRING"`
		MockSliceString []string `env:"MOCK_SLICE_STRING"`
	}

	firstScanningWrongObj := MockConfig{}
	secondWrongObj := "wrong"
	secondScanningWrongObj := &secondWrongObj

	err := ReadEnvironment(firstScanningWrongObj, make(map[string]string))
	if !errors.Is(err, ErrVMustBePtr) {
		t.Errorf("expected %v, got %v", ErrVMustBePtr, err)
	}

	err = ReadEnvironment(secondScanningWrongObj, make(map[string]string))
	if !errors.Is(err, ErrVMustBePtr) {
		t.Errorf("expected %v, got %v", ErrVMustBePtr, err)
	}
}

func TestReadEnvironmentStructErrTagMissing(t *testing.T) {
	type MockConfig struct {
		MockString      string
		MockSliceString []string `env:"MOCK_SLICE_STRING"`
	}

	t.Setenv("MOCK_STRING", "some_value")

	scanningWrongObj := &MockConfig{}
	expectedErr := errors.New("Struct field \"MockString\" is missing 'env' tag")

	err := ReadEnvironment(scanningWrongObj, make(map[string]string))
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != expectedErr.Error() {
		t.Errorf("expected %q, got %q", expectedErr.Error(), err.Error())
	}
}

func TestReadEnvironmentErrVarNotSet(t *testing.T) {
	type MockConfig struct {
		MockString      string   `env:"MOCK_STRING"`
		MockSliceString []string `env:"MOCK_SLICE_STRING"`
	}

	t.Setenv("MOCK_SLICE_STRING", "some_value,some_other_value")

	scanningWrongObj := &MockConfig{}
	expectedErr := errors.New("environment variable MOCK_STRING not set")

	err := ReadEnvironment(scanningWrongObj, make(map[string]string))
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != expectedErr.Error() {
		t.Errorf("expected %q, got %q", expectedErr.Error(), err.Error())
	}
}
