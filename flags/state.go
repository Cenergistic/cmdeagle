package flags

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Cenergistic/cmdeagle/envvar"
	"github.com/Cenergistic/cmdeagle/types"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type FlagsStateStore struct {
	// TODO implement global flags
	// TODO implement grouped flag configuration
	cobraCommand *cobra.Command
	pFlagSet     *pflag.FlagSet
	flagDefMap   map[string]*types.FlagDefinition
}

func CreateFlagsStore(cobraCommand *cobra.Command, commandDef *types.CommandDefinition) (*FlagsStateStore, error) {
	// TODO handle persistent flags
	// https://cobra.dev/#persistent-flags

	store := &FlagsStateStore{
		cobraCommand: cobraCommand,
		pFlagSet:     cobraCommand.Flags(),
		flagDefMap:   make(map[string]*types.FlagDefinition),
	}

	for i := range commandDef.Flags {
		flagDef := commandDef.Flags[i]
		log.Debug("\tValidating flag definition", "name", flagDef.Name)
		if err := ValidateDefinition(&flagDef); err != nil {
			return nil, err
		}

		log.Debug("\tGetting flag definition", "name", flagDef.Name)
		flagType, err := GetFlagType(flagDef.Type)
		if err != nil {
			return nil, err
		}

		log.Debug("\tBinding flag", "name", flagDef.Name)
		var flagVal *any
		flagType.Bind(flagVal, store.pFlagSet, &flagDef)

		store.flagDefMap[flagDef.Name] = &flagDef
	}

	return store, nil
}

func (store *FlagsStateStore) Get(key string) *pflag.Flag {
	return store.pFlagSet.Lookup(key)
}

func (store *FlagsStateStore) GetVal(key string) any {
	flag := store.pFlagSet.Lookup(key)
	if flag == nil {
		return nil
	}

	return flag.Value
}

func (store *FlagsStateStore) GetDef(key string) *types.FlagDefinition {
	return store.flagDefMap[key]
}

func (store *FlagsStateStore) VisitAll(fn func(flag *pflag.Flag)) {
	store.pFlagSet.VisitAll(fn)
}

// Interpolate replaces {{flags.<name>}} placeholders with a reference to the
// corresponding environment variable (e.g. ${FLAGS_VERBOSE}) rather than
// splicing the raw value into the script text. The shell treats the contents of
// an expanded variable as data, so a hostile flag value cannot be interpreted
// as shell syntax. The values are supplied via the environment (see
// GetEnvVariables).
func (store *FlagsStateStore) Interpolate(script string) string {
	store.pFlagSet.VisitAll(func(flag *pflag.Flag) {
		placeholder := fmt.Sprintf("{{flags.%s}}", flag.Name)
		envRef := "${FLAGS_" + envvar.GetEnvVariableNameFromStateKey(flag.Name) + "}"
		script = strings.ReplaceAll(script, placeholder, envRef)
	})

	return script
}

func (store *FlagsStateStore) GetEnvVariables() []types.EnvVar {
	envVars := make([]types.EnvVar, 0)

	store.pFlagSet.VisitAll(func(flag *pflag.Flag) {
		envVars = append(envVars, types.EnvVar{Name: "FLAGS_" + envvar.GetEnvVariableNameFromStateKey(flag.Name), Value: fmt.Sprint(flag.Value)})
	})

	return envVars
}

func (store *FlagsStateStore) ToJSON() map[string]any {
	result := make(map[string]any)

	store.pFlagSet.VisitAll(func(flag *pflag.Flag) {
		result[flag.Name] = flag.Value.String()
	})

	return result
}

func (store *FlagsStateStore) ToJSONString() string {
	jsonBytes, err := json.Marshal(store.ToJSON())
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}
