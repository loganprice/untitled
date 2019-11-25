package env

import (
	"os"
	"testing"
)

type validTestStruct struct {
	A string `env:"A"`
}

func setConstEnv(desired map[string]string, t *testing.T) {
	t.Helper()
	os.Clearenv()
	for k, v := range desired {
		os.Setenv(k, v)
	}

}
