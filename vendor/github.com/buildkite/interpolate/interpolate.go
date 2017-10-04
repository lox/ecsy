package interpolate

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var (
	variablesWithBracketsRegex   = regexp.MustCompile(`([\\\$]?\$\{([^}]+?)})`)
	variablesWithNoBracketsRegex = regexp.MustCompile(`([\\\$]?\$[a-zA-Z0-9_]+)`)
	substringRegexp              = regexp.MustCompile(`\A\s*:\s*(\-?\s*\d+)(?:\s*:\s*(\-?\s*\d+))?\s*\z`)
)

// Interpolate takes a set of environment and interpolates it into the provided string using shell script expansions
func Interpolate(env map[string]string, str string) (string, error) {
	var err error

	// Do a parse and handle ENV variables with the ${} syntax, i.e. ${FOO}
	interpolated := variablesWithBracketsRegex.ReplaceAllStringFunc(str, func(v string) string {
		// if there has been an error previously, skip interpolation
		if err != nil {
			return v
		}
		key, option := extractKeyAndOptionFromVariable(v)

		// Just return the key by itself if it was escaped
		if isPrefixedWithEscapeSequence(v) {
			return key
		}

		err = isValidPosixEnvironmentVariable(v)
		if err != nil {
			return v
		}

		vv, isEnvironmentVariableSet := env[key]

		switch {
		case substringRegexp.MatchString(option):
			// Substring Expansion -- select a substring of a variable with:
			//
			// ${parameter:offset}
			// ${parameter:offset:length}
			//
			// In the first form select a substring of $parameter starting from
			// 0-indexed offset until the end of $parameter. If offset is
			// negative then it is an offset from the end of $parameter instead.
			//
			// In the second form, length is the number of characters from offset
			// to select. If negeative, length is instead an offset from the end
			// of $parameter.
			//
			match := substringRegexp.FindStringSubmatch(option)
			lenvv := int64(len(vv))

			offset, parseErr := strconv.ParseInt(match[1], 10, 0)
			if err != nil {
				err = parseErr
				return v
			}

			// Negative offsets = from end
			if offset < 0 {
				offset = lenvv - (-offset)
			}

			// Still negative = too far from end? Truncate to start.
			if offset < 0 {
				offset = 0
			}

			// Beyond end? Truncate to end.
			if offset > lenvv {
				offset = lenvv
			}

			// Length?
			if len(match) < 3 || match[2] == "" {
				vv = vv[offset:lenvv]
			} else {
				length, parseErr := strconv.ParseInt(match[2], 10, 0)
				if err != nil {
					err = parseErr
					return v
				}

				if length >= 0 {
					// Positive length = from offset
					length = offset + length

					// Too far? Truncate to end.
					if length > lenvv {
						length = lenvv
					}
				} else {
					// Negative length = from end
					length = lenvv - (-length)

					// Too far? Truncate to offset.
					if length < offset {
						length = offset
					}
				}

				vv = vv[offset:length]
			}

		case strings.HasPrefix(option, "?"):
			if vv == "" {
				errorMessage := option[1:]
				if errorMessage == "" {
					errorMessage = "not set"
				}
				err = fmt.Errorf("$%s: %s", key, errorMessage)
			}

		case strings.HasPrefix(option, ":-"):
			if vv == "" {
				vv = option[2:]
			}

		case strings.HasPrefix(option, "-"):
			if !isEnvironmentVariableSet {
				vv = option[1:]
			}

		case option != "":
			err = fmt.Errorf("Invalid option `%s` for environment variable `%s`", option, key)
		}

		return vv
	})
	if err != nil {
		return str, err
	}

	// Another parse but this time target ENV variables without the {}
	// surrounding it, i.e. $FOO. These ones are super simple to replace.
	interpolated = variablesWithNoBracketsRegex.ReplaceAllStringFunc(interpolated, func(v string) string {
		// if there has been an error previously, skip interpolation
		if err != nil {
			return v
		}
		key, _ := extractKeyAndOptionFromVariable(v)

		// Just return the key by itself if it was escaped
		if isPrefixedWithEscapeSequence(v) {
			return key
		}

		err = isValidPosixEnvironmentVariable(v)
		if err != nil {
			return v
		}

		val, _ := env[key]
		return val
	})

	return interpolated, err
}

func isPrefixedWithEscapeSequence(variable string) bool {
	return strings.HasPrefix(variable, "$$") || strings.HasPrefix(variable, "\\$")
}

var validPosixEnvironmentVariablePrefixRegex = regexp.MustCompile(`\A\${1}\{?[a-zA-Z]`)

// Returns true if the variable is a valid POSIX environment variale. It will
// return false if the variable begins with a number, or it starts with two $$
// characters.
func isValidPosixEnvironmentVariable(variable string) error {
	if validPosixEnvironmentVariablePrefixRegex.MatchString(variable) {
		return nil
	} else {
		return fmt.Errorf("Invalid environment variable `%s` - they can only start with a letter", variable)
	}
}

var firstNonEnvironmentVariableCharacterRegex = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// Takes an environment variable, and extracts the variable name and a suffixed
// option.  For example, ${BEST_COMMAND:-lol} will be turned split into
// "BEST_COMMAND" and ":-lol". Regualr environment variables like $FOO will
// return "FOO" as the `key`, and a blank string as the `option`.
func extractKeyAndOptionFromVariable(variable string) (key string, option string) {
	if strings.HasPrefix(variable, "${") {
		// Trim the first 2 characters `${` and the last character `}`
		trimmed := variable[2 : len(variable)-1]

		optionsIndicies := firstNonEnvironmentVariableCharacterRegex.FindStringIndex(trimmed)
		if len(optionsIndicies) > 0 {
			key = trimmed[0:optionsIndicies[0]]
			option = trimmed[optionsIndicies[0]:]
		} else {
			key = trimmed
		}
	} else {
		// Trim the first character `$`
		key = variable[1:]
	}

	return
}

// Windows isn't case sensitive for env
func normalizeKeyName(key string) string {
	if runtime.GOOS == "windows" {
		return strings.ToUpper(key)
	}
	return key
}

// Converts an env slice (like from os.Environ) to a map, with keys normalized to uppercase in windows
func EnvFromSlice(env []string) map[string]string {
	envMap := map[string]string{}
	for _, l := range env {
		parts := strings.SplitN(l, "=", 2)
		if len(parts) == 2 {
			envMap[normalizeKeyName(parts[0])] = parts[1]
		}
	}
	return envMap
}
