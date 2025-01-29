package jsonParser

import (
	"fmt"
	"strconv"
)

func ParseJson(in []rune) (any, error) {
	in = skipWhitespaces(in)
	if len(in) == 0 {
		return nil, fmt.Errorf("invalid string")
	} else if in[0] == '{' {
		res, _, err := parseStructure(in)
		return res, err
	} else if in[0] == '[' {
		res, _, err := parseArray(in)
		return res, err
	} else {
		return nil, fmt.Errorf("invalid character: '%c'", in[0])
	}
}

func skipWhitespaces(in []rune) []rune {
	for len(in) > 0 && isWhitespace(in[0]) {
		in = in[1:]
	}
	return in
}

func parseStructure(in []rune) (map[string]any, []rune, error) {
	if in[0] != '{' {
		return nil, in, fmt.Errorf("invalid structure")
	}
	in = skipWhitespaces(in[1:])

	var name string
	var err error
	var value any

	result := make(map[string]any)

	for in[0] != '}' {
		// reading name
		name, in, err = parseString(in)
		if err != nil {
			return nil, in, err
		}

		// ensuring its a property
		in = skipWhitespaces(in)
		if len(in) == 0 {
			return nil, in, fmt.Errorf("input ended unexpectedly")
		}
		if in[0] != ':' {
			return nil, in, fmt.Errorf("invalid propriety")
		}

		// reading value
		value, in, err = parseValue(in[1:])
		if err != nil {
			return nil, in, err
		}

		result[name] = value

		in = skipWhitespaces(in)
		if len(in) == 0 {
			return nil, in, fmt.Errorf("input ended unexpectedly")
		}
		if in[0] == ',' {
			in = skipWhitespaces(in[1:])
		}
	}
	in = skipWhitespaces(in[1:])
	return result, in, nil
}

func parseArray(in []rune) ([]any, []rune, error) {
	if in[0] != '[' {
		return nil, in, fmt.Errorf("invalid array")
	}

	var val any
	var err error
	res := make([]any, 0, 10)
	in = skipWhitespaces(in[1:])

	for in[0] != ']' {
		val, in, err = parseValue(in)
		if err != nil {
			return nil, in, err
		}
		res = append(res, val)
		in = skipWhitespaces(in)
		if len(in) == 0 {
			return nil, in, fmt.Errorf("input ended unexpectedly")
		}
		if in[0] == ',' {
			in = skipWhitespaces(in[1:])
		}
	}

	in = skipWhitespaces(in[1:])
	return res, in, nil
}

func parseString(in []rune) (string, []rune, error) {
	if len(in) == 0 {
		return "", in, fmt.Errorf("input ended unexpectedly")
	} else if in[0] == '"' {
		return readString(in[1:])
	}
	return "", in, fmt.Errorf("invalid character: '%c'", in[0])
}

func readString(in []rune) (string, []rune, error) {
	escape := false
	res := make([]rune, 0, 10)
	if len(in) == 0 {
		return "", in, fmt.Errorf("input ended unexpectedly")
	}
	for in[0] != '"' || escape {
		if escape {
			switch in[0] {
			case '\\':
				res = append(res, '\\')
			case 'n':
				res = append(res, '\n')
			case '"':
				res = append(res, '"')
			case 'r':
				res = append(res, '\r')
			case 't':
				res = append(res, '\t')
			default:
				return "", in, fmt.Errorf("invalid escape code: '\\%c'", in[0])
			}
			escape = false
		} else {
			if in[0] == '\\' {
				escape = true
			} else {
				res = append(res, in[0])
			}
		}
		in = in[1:]
		if len(in) == 0 {
			return "", in, fmt.Errorf("input ended unexpectedly")
		}
	}
	return string(res), in[1:], nil
}

func parseValue(in []rune) (any, []rune, error) {
	in = skipWhitespaces(in)
	if len(in) == 0 {
		return nil, in, fmt.Errorf("input ended unexpectedly")
	}
	if in[0] == 't' || in[0] == 'f' || in[0] == 'n' {
		return readBoolOrNull(in)
	} else if isDigit(in[0]) {
		return readNumber(in)
	} else if in[0] == '"' {
		return parseString(in)
	} else if in[0] == '{' {
		return parseStructure(in)
	} else if in[0] == '[' {
		return parseArray(in)
	} else {
		return nil, in, fmt.Errorf("invalid propriety")
	}
}

func isDigit(c rune) bool {
	return (c >= '0' && c <= '9') || c == '.'
}

func readBoolOrNull(in []rune) (any, []rune, error) {
	str, in, err := readTillCommaOrBracket(in)
	if err != nil {
		return "", in, err
	}
	if str == "true" {
		return true, in, nil
	} else if str == "false" {
		return false, in, nil
	} else if str == "null" {
		return nil, in, nil
	} else {
		return false, in, fmt.Errorf("invalid propriety: %v", str)
	}
}

func readNumber(in []rune) (any, []rune, error) {
	var res, prev []rune = make([]rune, 0, 10), nil
	for isDigit(in[0]) {
		if in[0] == '.' {
			if prev == nil {
				prev = res
				res = make([]rune, 0, 10)
			} else {
				return nil, in, fmt.Errorf("invalid number")
			}
		}
		res = append(res, in[0])
		in = in[1:]
		if len(in) == 0 {
			return nil, in, fmt.Errorf("input ended unexpectedly")
		}
	}
	if prev == nil {
		number, err := strconv.ParseInt(string(res), 10, 64)
		// this shoulden't happen
		if err != nil {
			return nil, in, err
		}
		return number, in, nil
	} else {
		number, err := strconv.ParseFloat(string(prev)+string(res), 64)
		// this shoulden't happen
		if err != nil {
			return nil, in, err
		}
		return number, in, nil
	}
}

func readTillCommaOrBracket(in []rune) (string, []rune, error) {
	res := make([]rune, 0, 10)
	for in[0] != ',' && in[0] != '}' && in[0] != ']' && !isWhitespace(in[0]) {
		res = append(res, in[0])
		in = in[1:]
		if len(in) == 0 {
			return "", in, fmt.Errorf("input ended unexpectedly")
		}
	}
	return string(res), in, nil
}

func isWhitespace(c rune) bool {
	return c == ' ' || c == '\n' || c == '\r' || c == '\t'
}
