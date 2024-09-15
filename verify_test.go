package cli

import (
	"testing"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
)

func TestVerifyReadKeyValues(t *testing.T) {

	logger := sthingsBase.StdOutFileLogger("/tmp/machineShop.log", "2006-01-02 15:04:05", 50, 3, 28)

	tests := []struct {
		templateValues []string
		enableVault    bool
		expectedResult map[string]interface{}
	}{
		{
			templateValues: []string{"username=admin", "password=secret"},
			enableVault:    false,
			expectedResult: map[string]interface{}{
				"username": "admin",
				"password": "secret",
			},
		},
		// {
		// 	templateValues: []string{"key1=value1", "key2=value2"},
		// 	enableVault:    true,
		// 	expectedResult: map[string]interface{}{
		// 		"key1=value1": "key1=value1",
		// 		"key2=value2": "key2=value2",
		// 	},
		// },
	}

	for _, test := range tests {
		result := VerifyReadKeyValues(test.templateValues, logger, test.enableVault)
		if len(result) != len(test.expectedResult) {
			t.Errorf("Expected result length: %d, got: %d", len(test.expectedResult), len(result))
		}
		for key, expectedValue := range test.expectedResult {
			if value, exists := result[key]; !exists || value != expectedValue {
				t.Errorf("For key '%s', expected value: '%v', got: '%v'", key, expectedValue, value)
			}
		}
	}
}
