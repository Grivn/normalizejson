package regex

import "regexp"

var regexJSONKey = regexp.MustCompile(`\"(\w+)\":`) // JSON key

var regexCamelCase = regexp.MustCompile(`(\w)([A-Z])`) // camel-case style

var regexSnakeCase = regexp.MustCompile(`(\w)_(\w)`) // snake-case style
