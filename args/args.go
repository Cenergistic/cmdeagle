package args

import (
	"embed"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Cenergistic/cmdeagle/types"

	cast "github.com/spf13/cast"
)

//go:embed *
var PackageFS embed.FS

type ArgTypeDef struct {
	DefaultVal any
	Convert    func(val string) (any, error)
}

var ArgRuleConditionals = []string{"MatchAll", "MatchAny", "MatchNone", "And", "Or", "Nand", "Not"}

func AddArgType(name string, defaultVal any, convert func(val string) (any, error)) {
	if _, ok := argTypes[name]; ok {
		panic(fmt.Sprintf("Arg type `%s` already exists", name))
	}

	argTypes[name] = ArgTypeDef{
		Convert:    convert,
		DefaultVal: defaultVal,
	}
}

// GetArgType returns the type definition for an argument type name. It returns
// an error (rather than panicking) when the type is not recognised so that a bad
// config surfaces as a friendly message instead of a stack trace.
func GetArgType(name string) (ArgTypeDef, error) {
	argType, ok := argTypes[name]
	if !ok {
		return ArgTypeDef{}, fmt.Errorf("unknown argument type %q (valid types: %s)", name, strings.Join(ValidTypeNames(), ", "))
	}

	return argType, nil
}

// IsValidType reports whether name is a recognised argument type.
func IsValidType(name string) bool {
	_, ok := argTypes[name]
	return ok
}

// ValidTypeNames lists the recognised argument type names, sorted for stable output.
func ValidTypeNames() []string {
	names := make([]string, 0, len(argTypes))
	for name := range argTypes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ValidateDefinition checks a single argument definition for an unknown type.
func ValidateDefinition(argDef *types.ArgDefinition) error {
	if argDef.Name == "" {
		return fmt.Errorf("argument is missing a name")
	}
	argType := argDef.Type
	if argType == "" {
		argType = DefaultArgType
	}
	if !IsValidType(argType) {
		return fmt.Errorf("argument %q has unknown type %q (valid types: %s)", argDef.Name, argType, strings.Join(ValidTypeNames(), ", "))
	}
	return nil
}

var argTypes = map[string]ArgTypeDef{
	"string": {
		DefaultVal: "",
		Convert:    func(val string) (any, error) { return cast.ToStringE(val) },
	},
	"date": {
		DefaultVal: time.Now(),
		Convert:    func(val string) (any, error) { return cast.StringToDate(val) },
	},
	"time": {
		DefaultVal: time.Now(),
		Convert:    func(val string) (any, error) { return cast.ToTimeE(val) },
	},
	"duration": {
		DefaultVal: time.Duration(0),
		Convert:    func(val string) (any, error) { return cast.ToDurationE(val) },
	},
	"boolean": {
		DefaultVal: false,
		Convert:    func(val string) (any, error) { return cast.ToBoolE(val) },
	},
	"bool": {
		DefaultVal: false,
		Convert:    func(val string) (any, error) { return cast.ToBoolE(val) },
	},
	"number": {
		DefaultVal: 0.0,
		Convert: func(val string) (any, error) {
			// First try to convert to float64
			if f, err := cast.ToFloat64E(val); err == nil {
				return f, nil
			}
			// If that fails, try converting from int to float64
			if i, err := cast.ToInt64E(val); err == nil {
				return float64(i), nil
			}
			return 0.0, fmt.Errorf("cannot convert %v to number (float64)", val)
		},
	},
	"float64": {
		DefaultVal: 0.0,
		Convert:    func(val string) (any, error) { return cast.ToFloat64E(val) },
	},
	"float32": {
		DefaultVal: 0.0,
		Convert:    func(val string) (any, error) { return cast.ToFloat32E(val) },
	},
	"int64": {
		DefaultVal: int64(0),
		Convert:    func(val string) (any, error) { return cast.ToInt64E(val) },
	},
	"int32": {
		DefaultVal: int32(0),
		Convert:    func(val string) (any, error) { return cast.ToInt32E(val) },
	},
	"int16": {
		DefaultVal: int16(0),
		Convert:    func(val string) (any, error) { return cast.ToInt16E(val) },
	},
	"int8": {
		DefaultVal: int8(0),
		Convert:    func(val string) (any, error) { return cast.ToInt8E(val) },
	},
	"int": {
		DefaultVal: int(0),
		Convert:    func(val string) (any, error) { return cast.ToIntE(val) },
	},
	"uint": {
		DefaultVal: uint(0),
		Convert:    func(val string) (any, error) { return cast.ToUintE(val) },
	},
	"uint64": {
		DefaultVal: uint64(0),
		Convert:    func(val string) (any, error) { return cast.ToUint64E(val) },
	},
	"uint32": {
		DefaultVal: uint32(0),
		Convert:    func(val string) (any, error) { return cast.ToUint32E(val) },
	},
	"uint16": {
		DefaultVal: uint16(0),
		Convert:    func(val string) (any, error) { return cast.ToUint16E(val) },
	},
	"uint8": {
		DefaultVal: uint8(0),
		Convert:    func(val string) (any, error) { return cast.ToUint8E(val) },
	},
}
