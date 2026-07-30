package config

import (
	"testing"

	"github.com/Cenergistic/cmdeagle/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateConfig(t *testing.T) {
	t.Run("accepts a valid config", func(t *testing.T) {
		cfg := &types.CmdeagleConfig{
			Name: "demo",
			Commands: []types.CommandDefinition{
				{
					Name:  "greet",
					Args:  []types.ArgDefinition{{Name: "name", Type: "string"}},
					Flags: []types.FlagDefinition{{Name: "loud", Type: "boolean", Shorthand: "l"}},
				},
			},
		}
		assert.NoError(t, ValidateConfig(cfg))
	})

	t.Run("requires a name", func(t *testing.T) {
		err := ValidateConfig(&types.CmdeagleConfig{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "name")
	})

	t.Run("reports a bad flag in a nested command with context", func(t *testing.T) {
		cfg := &types.CmdeagleConfig{
			Name: "demo",
			Commands: []types.CommandDefinition{
				{
					Name: "parent",
					Commands: []types.CommandDefinition{
						{
							Name:  "child",
							Flags: []types.FlagDefinition{{Name: "count", Type: "int", Default: "nan"}},
						},
					},
				},
			},
		}
		err := ValidateConfig(cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "child")
		assert.Contains(t, err.Error(), "count")
	})
}
