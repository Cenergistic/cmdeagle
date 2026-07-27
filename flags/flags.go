package flags

import (
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/migsc/cmdeagle/types"

	"github.com/spf13/cast"
	"github.com/spf13/pflag"
)

//go:embed *
var PackageFS embed.FS

type FlagTypeDef struct {
	Bind func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any
}

// GetFlagType returns the binding definition for a flag type name. It returns an
// error (rather than panicking) when the type is not recognised so that a bad
// config surfaces as a friendly message instead of a stack trace.
func GetFlagType(name string) (FlagTypeDef, error) {
	flagType, ok := flagTypes[name]
	if !ok {
		return FlagTypeDef{}, fmt.Errorf("unknown flag type %q (valid types: %s)", name, strings.Join(ValidTypeNames(), ", "))
	}

	return flagType, nil
}

// IsValidType reports whether name is a recognised flag type.
func IsValidType(name string) bool {
	_, ok := flagTypes[name]
	return ok
}

// ValidTypeNames lists the recognised flag type names, sorted for stable output.
func ValidTypeNames() []string {
	names := make([]string, 0, len(flagTypes))
	for name := range flagTypes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ValidateDefinition checks a single flag definition for the mistakes that would
// otherwise blow up flag binding: an unknown type, a multi-character shorthand,
// or a default value that cannot be coerced to the declared type.
func ValidateDefinition(flagDef *types.FlagDefinition) error {
	if flagDef.Name == "" {
		return fmt.Errorf("flag is missing a name")
	}
	if !IsValidType(flagDef.Type) {
		return fmt.Errorf("flag %q has unknown type %q (valid types: %s)", flagDef.Name, flagDef.Type, strings.Join(ValidTypeNames(), ", "))
	}
	if len(flagDef.Shorthand) > 1 {
		return fmt.Errorf("flag %q has shorthand %q, but shorthands must be a single character", flagDef.Name, flagDef.Shorthand)
	}
	if flagDef.Default != nil {
		if _, err := coerceDefault(flagDef.Type, flagDef.Default); err != nil {
			return fmt.Errorf("flag %q has an invalid default: %w", flagDef.Name, err)
		}
	}
	return nil
}

// coerceDefault converts a YAML-decoded default value to the concrete type the
// flag expects. It never panics; an unconvertible value is returned as an error.
func coerceDefault(flagType string, raw any) (any, error) {
	switch flagType {
	case "string":
		return cast.ToStringE(raw)
	case "boolean", "bool":
		return cast.ToBoolE(raw)
	case "number", "float64":
		return cast.ToFloat64E(raw)
	case "float32":
		return cast.ToFloat32E(raw)
	case "int64":
		return cast.ToInt64E(raw)
	case "int32":
		return cast.ToInt32E(raw)
	case "int16":
		return cast.ToInt16E(raw)
	case "int8":
		return cast.ToInt8E(raw)
	case "int":
		return cast.ToIntE(raw)
	case "uint":
		return cast.ToUintE(raw)
	case "uint64":
		return cast.ToUint64E(raw)
	case "uint32":
		return cast.ToUint32E(raw)
	case "uint16":
		return cast.ToUint16E(raw)
	case "uint8":
		return cast.ToUint8E(raw)
	default:
		return nil, fmt.Errorf("unknown flag type %q", flagType)
	}
}

var flagTypes = map[string]FlagTypeDef{
	"string": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var strVal string
			var flagVal any = &strVal
			defaultVal := ""
			if flagDef.Default != nil {
				defaultVal = cast.ToString(flagDef.Default)
			}
			flagSet.StringVarP(&strVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"boolean": {Bind: bindBool},
	"bool":    {Bind: bindBool},
	"number": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal float64
			var flagVal any = &numVal
			defaultVal := cast.ToFloat64(flagDef.Default)
			flagSet.Float64VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"float64": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal float64
			var flagVal any = &numVal
			defaultVal := cast.ToFloat64(flagDef.Default)
			flagSet.Float64VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"float32": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal float32
			var flagVal any = &numVal
			defaultVal := cast.ToFloat32(flagDef.Default)
			flagSet.Float32VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"int64": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal int64
			var flagVal any = &numVal
			defaultVal := cast.ToInt64(flagDef.Default)
			flagSet.Int64VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"int32": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal int32
			var flagVal any = &numVal
			defaultVal := cast.ToInt32(flagDef.Default)
			flagSet.Int32VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"int16": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal int16
			var flagVal any = &numVal
			defaultVal := cast.ToInt16(flagDef.Default)
			flagSet.Int16VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"int8": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal int8
			var flagVal any = &numVal
			defaultVal := cast.ToInt8(flagDef.Default)
			flagSet.Int8VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"int": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal int
			var flagVal any = &numVal
			defaultVal := cast.ToInt(flagDef.Default)
			flagSet.IntVarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"uint": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal uint
			var flagVal any = &numVal
			defaultVal := cast.ToUint(flagDef.Default)
			flagSet.UintVarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"uint64": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal uint64
			var flagVal any = &numVal
			defaultVal := cast.ToUint64(flagDef.Default)
			flagSet.Uint64VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"uint32": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal uint32
			var flagVal any = &numVal
			defaultVal := cast.ToUint32(flagDef.Default)
			flagSet.Uint32VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"uint16": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal uint16
			var flagVal any = &numVal
			defaultVal := cast.ToUint16(flagDef.Default)
			flagSet.Uint16VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
	"uint8": {
		Bind: func(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
			var numVal uint8
			var flagVal any = &numVal
			defaultVal := cast.ToUint8(flagDef.Default)
			flagSet.Uint8VarP(&numVal, flagDef.Name, flagDef.Shorthand, defaultVal, flagDef.Description)
			return &flagVal
		},
	},
}

func bindBool(val *any, flagSet *pflag.FlagSet, flagDef *types.FlagDefinition) *any {
	var boolVal bool
	var flagVal any = &boolVal
	defaultVal := cast.ToBool(flagDef.Default)
	description := flagDef.Description
	if !strings.HasSuffix(description, ")") {
		description += " (accepts: true/false, t/f, 1/0, yes/no, y/n)"
	}

	if flagDef.Shorthand != "" {
		flagSet.BoolVarP(&boolVal, flagDef.Name, flagDef.Shorthand, defaultVal, description)
	} else {
		flagSet.BoolVar(&boolVal, flagDef.Name, defaultVal, description)
	}

	flag := flagSet.Lookup(flagDef.Name)
	if flag != nil {
		flag.NoOptDefVal = "true"
		// Add custom boolean value parsing
		flag.Value = &boolValue{
			value: &boolVal,
		}
	}
	return &flagVal
}

// boolValue implements pflag.Value interface
type boolValue struct {
	value *bool
}

func (b *boolValue) Set(val string) error {
	val = strings.ToLower(strings.TrimSpace(val))
	switch val {
	case "true", "t", "1", "yes", "y":
		*b.value = true
	case "false", "f", "0", "no", "n":
		*b.value = false
	default:
		return fmt.Errorf("invalid boolean value '%s'. Accepted values: true/false, t/f, 1/0, yes/no, y/n", val)
	}
	return nil
}

func (b *boolValue) String() string {
	if b.value == nil {
		return "false"
	}
	return fmt.Sprintf("%v", *b.value)
}

func (b *boolValue) Type() string {
	return "bool"
}
