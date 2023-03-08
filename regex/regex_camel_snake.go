package regex

import "regexp"

var regexJSONKey = regexp.MustCompile(`\"([\w_:/$]+)\":`) // JSON key

var regexCamelCaseJSONKey = regexp.MustCompile(`\"(\w+)\":`) // camel-case JSON key

var regexCamelCase = regexp.MustCompile(`(\w)([A-Z])`) // camel-case style
