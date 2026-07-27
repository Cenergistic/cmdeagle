package flags

import (
	"testing"

	"github.com/Cenergistic/cmdeagle/types"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateDefinition(t *testing.T) {
	t.Run("accepts a valid flag", func(t *testing.T) {
		err := ValidateDefinition(&types.FlagDefinition{Name: "verbose", Type: "boolean", Shorthand: "v"})
		assert.NoError(t, err)
	})

	t.Run("rejects an unknown type", func(t *testing.T) {
		err := ValidateDefinition(&types.FlagDefinition{Name: "mode", Type: "enumeration"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown type")
	})

	t.Run("rejects a multi-character shorthand", func(t *testing.T) {
		err := ValidateDefinition(&types.FlagDefinition{Name: "verbose", Type: "boolean", Shorthand: "vv"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "single character")
	})

	t.Run("rejects a default that cannot be coerced", func(t *testing.T) {
		err := ValidateDefinition(&types.FlagDefinition{Name: "count", Type: "int", Default: "not-a-number"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid default")
	})

	t.Run("accepts a numeric default provided as a float", func(t *testing.T) {
		err := ValidateDefinition(&types.FlagDefinition{Name: "count", Type: "int", Default: 3.0})
		assert.NoError(t, err)
	})
}

func TestGetFlagTypeReturnsErrorForUnknownType(t *testing.T) {
	_, err := GetFlagType("enumeration")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown flag type")
}

func TestCreateFlagsStoreRejectsBadConfigWithoutPanicking(t *testing.T) {
	cmd := &cobra.Command{Use: "t"}
	commandDef := &types.CommandDefinition{
		Flags: []types.FlagDefinition{{Name: "count", Type: "int", Default: "nope"}},
	}
	_, err := CreateFlagsStore(cmd, commandDef)
	require.Error(t, err)
}

func TestInterpolateEmitsEnvReferences(t *testing.T) {
	cmd := &cobra.Command{Use: "t"}
	commandDef := &types.CommandDefinition{
		Flags: []types.FlagDefinition{{Name: "verbose", Type: "boolean"}},
	}
	store, err := CreateFlagsStore(cmd, commandDef)
	require.NoError(t, err)

	// A hostile value must not be spliced into the script; the placeholder is
	// replaced with a reference to the environment variable that carries it.
	got := store.Interpolate(`echo "{{flags.verbose}}"`)
	assert.Equal(t, `echo "${FLAGS_VERBOSE}"`, got)
}
