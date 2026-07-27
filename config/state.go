package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Cenergistic/cmdeagle/args"
	"github.com/Cenergistic/cmdeagle/envvar"
	"github.com/Cenergistic/cmdeagle/flags"
	"github.com/Cenergistic/cmdeagle/types"
)

type ParamsStateStore struct {
	Args    *args.ArgsStateStore
	Flags   *flags.FlagsStateStore
	Entries map[string]string
}

func CreateParamsStore(argsStore *args.ArgsStateStore, flagsStore *flags.FlagsStateStore) *ParamsStateStore {
	store := &ParamsStateStore{
		Args:    argsStore,
		Flags:   flagsStore,
		Entries: map[string]string{},
	}

	store.Entries["args.json"] = store.Args.ToJSONString()
	store.Entries["flags.json"] = store.Flags.ToJSONString()
	store.Entries["params.json"] = store.ToJSONString()

	return store
}

func (store *ParamsStateStore) Set(key string, value string) {
	store.Entries[key] = value
}

// Interpolate replaces {{<key>}} placeholders (e.g. {{cli.bin_dir}},
// {{args.json}}) with a reference to the corresponding environment variable so
// the value reaches the shell as data rather than as script text. See
// ArgsStateStore.Interpolate for the rationale.
func (store *ParamsStateStore) Interpolate(script string) string {
	for key := range store.Entries {
		placeholder := fmt.Sprintf("{{%s}}", key)
		envRef := "${" + envvar.GetEnvVariableNameFromStateKey(key) + "}"
		script = strings.ReplaceAll(script, placeholder, envRef)
	}

	return script
}

func (store *ParamsStateStore) GetEnvVariables() []types.EnvVar {
	envVars := make([]types.EnvVar, 0)

	envVars = append(envVars, store.Args.GetEnvVariables()...)
	envVars = append(envVars, store.Flags.GetEnvVariables()...)

	for key, val := range store.Entries {
		envVars = append(envVars, types.EnvVar{Name: envvar.GetEnvVariableNameFromStateKey(key), Value: val})
	}

	return envVars
}

func (store *ParamsStateStore) ToJSON() map[string]any {
	return map[string]any{
		"args":  store.Args.ToJSON(),
		"flags": store.Flags.ToJSON(),
	}
}

func (store *ParamsStateStore) ToJSONString() string {
	jsonBytes, err := json.Marshal(store.ToJSON())
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}
