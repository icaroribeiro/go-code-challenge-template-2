package env_test

import (
	"os"
	"testing"

	envpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/env"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestGetEnvWithDefaultValue() {
	key := ""
	value := ""
	defaultValue := ""

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingEnvValue",
			SetUp: func(t *testing.T) {
				key = "ENV_VAR"
				value = "value"
				os.Setenv(key, value)
			},
			WantError: false,
			TearDown: func(t *testing.T) {
				os.Unsetenv(key)
			},
		},
		{
			Context: "ItShouldSucceedInReturningDefaultEnvValueWhenEnvVariableIsNotFound",
			SetUp: func(t *testing.T) {
				key = "ENV_VAR"
				defaultValue = "defaultValue"
			},
			WantError: false,
			TearDown:  func(t *testing.T) {},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedValue := envpkg.GetEnvWithDefaultValue(key, defaultValue)

			if !tc.WantError {
				_, ok := os.LookupEnv(key)
				if ok {
					assert.Equal(t, value, returnedValue)
				} else {
					assert.Equal(t, defaultValue, returnedValue)
				}
			}

			tc.TearDown(t)
		})
	}
}
