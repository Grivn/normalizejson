package normalizejson

type FormatProvider interface {
	AddOptions(options ...FormatOption)
	UpdateTemplate(rawTemplate []byte) error
	FormatJSONSchema(data []byte) ([]byte, error)
	Reset()
}

func NewFormatProvider(rawTemplate []byte, options ...FormatOption) (FormatProvider, error) {
	return NewFormatSchemaProvider(rawTemplate, options...)
}
