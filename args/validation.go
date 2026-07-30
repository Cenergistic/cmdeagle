package args

import (
	"fmt"
	"regexp"

	"github.com/charmbracelet/log"
	"github.com/Cenergistic/cmdeagle/params"
	"github.com/Cenergistic/cmdeagle/types"

	"github.com/spf13/cobra"
)

func ValidateArgs(cobraCmd *cobra.Command, argsConfigDef *[]types.ArgDefinition, store *ArgsStateStore) error {
	log.Debug("Validating args", "cobraCmd", cobraCmd, "argsConfigDef", argsConfigDef, "store", store)
	if argsConfigDef == nil || store == nil {
		return nil
	}

	log.Debug("Validating args", "argsConfigDef", argsConfigDef)
	argsConfigDevVal := *argsConfigDef

	if argsConfigDevVal == nil {
		return nil
	}

	for index := range argsConfigDevVal {
		entry := store.GetAt(index)
		log.Debug("Validating args", "index", index, "entry", entry)

		if entry == nil {
			continue
		}

		log.Debug("Validating args", "entry", entry)
		if entry.Err != nil {
			return entry.Err
		}

		// Validate constraints
		if entry.Def != nil {
			constraints := entry.Def.Constraints
			log.Debug("Validating args", "constraints", constraints)
			err := params.ValidateConstraint(&constraints, entry.Val)
			if err != nil {
				return fmt.Errorf("validation failed for argument %s: %v", entry.Def.Name, err)
			}
		}

		// Validate dependencies
		if entry.Def != nil && entry.Def.DependsOn != nil {
			for _, dependency := range entry.Def.DependsOn {
				err := params.ValidateConstraint(dependency.When, store.GetVal(dependency.Name))
				if err != nil {
					return fmt.Errorf("dependency validation failed for argument %s: %v", entry.Def.Name, err)
				}
			}
		}

		// Validate conflicts
		if entry.Def != nil && entry.Def.ConflictsWith != nil {
			for _, conflict := range entry.Def.ConflictsWith {
				conflictVal := store.GetVal(conflict)
				if conflictVal != nil && conflictVal != "" {
					return fmt.Errorf("argument %s conflicts with %s", entry.Def.Name, conflict)
				}
			}
		}

		// Validate pattern
		if entry.Def != nil && entry.Def.Pattern != "" {
			var err error
			var match bool
			var pattern *regexp.Regexp

			pattern, err = regexp.Compile(entry.Def.Pattern)

			if err != nil {
				return fmt.Errorf("invalid pattern for argument %s: %v", entry.Def.Name, err)
			}
			log.Debug("Validating pattern for argument", "pattern", pattern, "value", entry.Val.(string))
			match = pattern.MatchString(entry.Val.(string))
			found := pattern.FindString(entry.Val.(string))

			log.Debug("Validating pattern for argument", "pattern", pattern, "value", entry.Val, "match", match, "err", err, "found", found)

			if !match {
				return fmt.Errorf("pattern validation failed for argument %s: %v", entry.Def.Name, entry.Val)
			}
		}
	}
	return nil
}
