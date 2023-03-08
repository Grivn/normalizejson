package gojson

type FormatDataProvider interface {
	AddOptions(options ...FormatOption)
	UpdateTemplate(rawTemplate []byte) error
	FormatJSONSchema(data []byte) ([]byte, error)
}

type FormatKeyProvider interface {
	AddOptions(options ...FormatOption)
	FormatJSONSchema(data []byte) ([]byte, error)
}

func NewFormatDataProvider(rawTemplate []byte, options ...FormatOption) (FormatDataProvider, error) {
	return newFormatDataImpl(rawTemplate, options...)
}

func NewDefaultFormatDataProvider(rawTemplate []byte) (FormatDataProvider, error) {
	return NewFormatDataProvider(rawTemplate, DefaultFormatDataOptions...)
}

func (fdi *formatDataImpl) AddOptions(options ...FormatOption) {
	fdi.addOptions(options...)
}

func (fdi *formatDataImpl) UpdateTemplate(rawTemplate []byte) error {
	return fdi.updateTemplate(rawTemplate)
}

func (fdi *formatDataImpl) FormatJSONSchema(data []byte) ([]byte, error) {
	return fdi.formatJSONSchema(data)
}

func NewFormatKeyProvider(options ...FormatOption) FormatKeyProvider {
	return newFormatKeyImpl(options...)
}

func (fki *formatKeyImpl) AddOptions(options ...FormatOption) {
	fki.addOptions(options...)
}

func (fki *formatKeyImpl) FormatJSONSchema(data []byte) ([]byte, error) {
	return fki.formatJSONSchema(data)
}
