package config

import (
	"fmt"

	"github.com/Cenergistic/cmdeagle/args"
	"github.com/Cenergistic/cmdeagle/flags"
	"github.com/Cenergistic/cmdeagle/types"
)

// ValidateConfig checks a parsed configuration for the mistakes that would
// otherwise surface as a runtime panic or a silently-ignored setting: unknown
// arg/flag types, multi-character shorthands, and defaults that cannot be
// coerced to their declared type. It is meant to run at build time so authors
// get a clear message before a binary is ever produced.
func ValidateConfig(cfg *types.CmdeagleConfig) error {
	if cfg == nil {
		return fmt.Errorf("configuration is empty")
	}
	if cfg.Name == "" {
		return fmt.Errorf("configuration is missing a `name`")
	}

	// Root-level args and flags.
	if err := validateArgsAndFlags(cfg.Name, cfg.Args, cfg.Flags); err != nil {
		return err
	}

	return validateCommands(cfg.Commands)
}

func validateCommands(commands []types.CommandDefinition) error {
	for i := range commands {
		cmd := &commands[i]
		if cmd.Name == "" {
			return fmt.Errorf("a command is missing a `name`")
		}
		if err := validateArgsAndFlags(cmd.Name, cmd.Args, cmd.Flags); err != nil {
			return err
		}
		if err := validateCommands(cmd.Commands); err != nil {
			return err
		}
	}
	return nil
}

func validateArgsAndFlags(commandName string, argDefs []types.ArgDefinition, flagDefs []types.FlagDefinition) error {
	for i := range argDefs {
		if err := args.ValidateDefinition(&argDefs[i]); err != nil {
			return fmt.Errorf("command %q: %w", commandName, err)
		}
	}
	for i := range flagDefs {
		if err := flags.ValidateDefinition(&flagDefs[i]); err != nil {
			return fmt.Errorf("command %q: %w", commandName, err)
		}
	}
	return nil
}
