package args

import (
	"testing"

	"github.com/Cenergistic/cmdeagle/types"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateDefinition(t *testing.T) {
	t.Run("accepts a valid argument", func(t *testing.T) {
		assert.NoError(t, ValidateDefinition(&types.ArgDefinition{Name: "name", Type: "string"}))
	})

	t.Run("defaults an empty type to string", func(t *testing.T) {
		assert.NoError(t, ValidateDefinition(&types.ArgDefinition{Name: "name"}))
	})

	t.Run("rejects an unknown type", func(t *testing.T) {
		err := ValidateDefinition(&types.ArgDefinition{Name: "when", Type: "chronotype"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown type")
	})
}

func TestGetArgTypeReturnsErrorForUnknownType(t *testing.T) {
	_, err := GetArgType("chronotype")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown argument type")
}

func TestCreateArgsStoreRejectsUnknownTypeWithoutPanicking(t *testing.T) {
	cmd := &cobra.Command{Use: "t"}
	defs := []types.ArgDefinition{{Name: "when", Type: "chronotype"}}
	_, err := CreateArgsStore(cmd, &defs, []string{"today"})
	require.Error(t, err)
}

func TestInterpolateEmitsEnvReferences(t *testing.T) {
	cmd := &cobra.Command{Use: "t"}
	defs := []types.ArgDefinition{{Name: "name", Type: "string"}}
	store, err := CreateArgsStore(cmd, &defs, []string{`x"; rm -rf /; echo "`})
	require.NoError(t, err)

	// The hostile value is NOT present in the resulting script; it is referenced
	// through the environment instead.
	got := store.Interpolate(`echo "Hello {{args.name}}"`)
	assert.Equal(t, `echo "Hello ${ARGS_NAME}"`, got)
	assert.NotContains(t, got, "rm -rf")
}
