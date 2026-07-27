package envvar

import (
	"fmt"
	"strings"

	"github.com/Cenergistic/cmdeagle/types"
)

type EnvStateStore struct {
	Entries map[string]string
}

func CreateEnvStore() *EnvStateStore {
	store := &EnvStateStore{
		Entries: map[string]string{},
	}

	return store
}

func (store *EnvStateStore) Set(key string, value string) {
	store.Entries[key] = value
}

// Interpolate replaces {{<key>}} placeholders with a reference to the
// corresponding environment variable so build scripts receive values as data
// rather than as script text. See args.ArgsStateStore.Interpolate for the
// rationale.
func (store *EnvStateStore) Interpolate(script string) string {
	for key := range store.Entries {
		placeholder := fmt.Sprintf("{{%s}}", key)
		envRef := "${" + GetEnvVariableNameFromStateKey(key) + "}"
		script = strings.ReplaceAll(script, placeholder, envRef)
	}

	return script
}

func (store *EnvStateStore) GetEnvVariables() []types.EnvVar {
	envVars := make([]types.EnvVar, 0)

	for key, val := range store.Entries {
		envVars = append(envVars, types.EnvVar{Name: GetEnvVariableNameFromStateKey(key), Value: val})
	}

	return envVars
}
