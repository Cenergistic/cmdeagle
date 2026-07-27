package executable

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-version"
)

// DefaultVersionCommands is a predefined list of common flags and subcommands
var DefaultVersionCommands = [][]string{
	{"-v"}, {"--v"}, {"version"}, {"-version"}, {"--version"},
}

// ExtractSemver parses and validates a semantic version string (e.g., 1.2.3) from command output.
func ExtractSemver(output string) (string, error) {
	// Use go-version to parse and valifdate the semantic version
	v, err := version.NewVersion(output)
	if err != nil {
		return "", fmt.Errorf("no valid semantic version found in output: %s, error: %v", output, err)
	}
	return v.String(), nil
}

// GetVersionWithFallbacks attempts to retrieve a version using common flags and subcommands
func GetVersionWithFallbacks(binaryPath string) (string, error) {
	for _, cmdArgs := range DefaultVersionCommands {
		versionOutput, err := RunCommand(binaryPath, cmdArgs)
		if err != nil {
			continue // Swallow errors and try the next command
		}

		semver, err := ExtractSemver(versionOutput)
		if err == nil {
			return semver, nil // Return the first successfully parsed version
		}
	}

	return "", fmt.Errorf("unable to determine version for binary: %s", binaryPath)
}

// GetVersion retrieves the version of a binary using a user-specified command or falls back to defaults
func GetVersion(binaryPath string, customCommand []string) (string, error) {
	if len(customCommand) > 0 {
		// Use the custom version command if provided
		versionOutput, err := RunCommand(binaryPath, customCommand)
		if err != nil {
			return "", fmt.Errorf("custom version command failed: %v", err)
		}

		semver, err := ExtractSemver(versionOutput)
		if err == nil {
			return semver, nil
		}
		return "", fmt.Errorf("custom version command output could not be parsed")
	}

	// Fallback to default commands
	return GetVersionWithFallbacks(binaryPath)
}

// CheckVersionCompatibility determines whether a found version satisfies a
// declared constraint. It supports an npm-style dialect (`*`, `^`, `~`, `>=`,
// `>`, `<=`, `<`, `A - B` ranges, and an exact version) and delegates all
// version parsing and comparison to hashicorp/go-version so numeric ordering,
// multi-digit segments, and pre-release handling are correct.
func CheckVersionCompatibility(foundVersion, declaredVersion string) (bool, string) {
	// Handle wildcard case first
	if declaredVersion == "*" {
		return true, ""
	}

	// The found version must be a full major.minor.patch. go-version would
	// happily parse "1.2" as "1.2.0", so guard the segment count explicitly to
	// keep the stricter behaviour this package has always promised.
	if !hasThreeSegments(foundVersion) {
		return false, "Invalid found version format"
	}
	found, err := version.NewVersion(foundVersion)
	if err != nil {
		return false, "Invalid found version format"
	}
	foundCore := coreVersion(found)

	declaredVersion = strings.TrimSpace(strings.ToLower(declaredVersion))

	switch {
	case strings.HasPrefix(declaredVersion, "^"):
		// Major version must match exactly, minor and patch can be greater.
		bound, ok := parseBound(strings.TrimPrefix(declaredVersion, "^"))
		if !ok {
			return false, "Invalid declared version format"
		}
		if foundCore.Segments()[0] != bound.Segments()[0] {
			return false, "Major version mismatch"
		}
		return true, ""

	case strings.HasPrefix(declaredVersion, "~"):
		// Major and minor must match exactly, only patch can be greater.
		bound, ok := parseBound(strings.TrimPrefix(declaredVersion, "~"))
		if !ok {
			return false, "Invalid declared version format"
		}
		if foundCore.Segments()[0] != bound.Segments()[0] || foundCore.Segments()[1] != bound.Segments()[1] {
			return false, "Major or minor version mismatch"
		}
		return true, ""

	case strings.HasPrefix(declaredVersion, ">="):
		bound, ok := parseBound(strings.TrimPrefix(declaredVersion, ">="))
		if !ok {
			return false, "Invalid minimum version format"
		}
		if foundCore.LessThan(bound) {
			return false, "Version too low"
		}
		return true, ""

	case strings.HasPrefix(declaredVersion, ">"):
		bound, ok := parseBound(strings.TrimPrefix(declaredVersion, ">"))
		if !ok {
			return false, "Invalid minimum version format"
		}
		if !foundCore.GreaterThan(bound) {
			return false, "Version too low"
		}
		return true, ""

	case strings.HasPrefix(declaredVersion, "<="):
		bound, ok := parseBound(strings.TrimPrefix(declaredVersion, "<="))
		if !ok {
			return false, "Invalid maximum version format"
		}
		if foundCore.GreaterThan(bound) {
			return false, "Version too high"
		}
		return true, ""

	case strings.HasPrefix(declaredVersion, "<"):
		bound, ok := parseBound(strings.TrimPrefix(declaredVersion, "<"))
		if !ok {
			return false, "Invalid maximum version format"
		}
		if !foundCore.LessThan(bound) {
			return false, "Version too high"
		}
		return true, ""

	case strings.Contains(declaredVersion, " - "):
		parts := strings.Split(declaredVersion, " - ")
		if len(parts) != 2 {
			return false, "Invalid version range format"
		}
		min, minOK := parseBound(parts[0])
		max, maxOK := parseBound(parts[1])
		if !minOK || !maxOK {
			return false, "Invalid version range format"
		}
		if foundCore.LessThan(min) {
			return false, "Version below range"
		}
		if foundCore.GreaterThan(max) {
			return false, "Version above range"
		}
		return true, ""

	default:
		// Exact version match
		bound, ok := parseBound(declaredVersion)
		if !ok {
			return false, "Invalid declared version format"
		}
		if !foundCore.Equal(bound) {
			return false, "Version does not match exactly"
		}
		return true, ""
	}
}

// hasThreeSegments reports whether raw has at least major.minor.patch, matching
// the strictness this package has always applied to input versions.
func hasThreeSegments(raw string) bool {
	raw = strings.TrimPrefix(strings.TrimSpace(raw), "v")
	return len(strings.Split(raw, ".")) >= 3
}

// parseBound parses a constraint bound, requiring a full major.minor.patch.
func parseBound(raw string) (*version.Version, bool) {
	raw = strings.TrimSpace(raw)
	if !hasThreeSegments(raw) {
		return nil, false
	}
	v, err := version.NewVersion(raw)
	if err != nil {
		return nil, false
	}
	return coreVersion(v), true
}

// coreVersion strips any pre-release or build metadata, leaving only
// major.minor.patch so comparisons ignore pre-release ordering (e.g. 1.2.3-beta
// is treated as 1.2.3). Segments() is always zero-padded to at least three.
func coreVersion(v *version.Version) *version.Version {
	segs := v.Segments()
	core, _ := version.NewVersion(fmt.Sprintf("%d.%d.%d", segs[0], segs[1], segs[2]))
	return core
}
