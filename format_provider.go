package gojson

type FormatDataProvider interface {
	AddOptions(options ...FormatDataOption)
	UpdateTemplate(rawTemplate []byte) error
	FormatJSONSchema(data []byte) ([]byte, error)
}

type FormatKeyProvider interface {
	AddOptions(options ...FormatKeyOption)
	FormatJSONSchema(data []byte) ([]byte, error)
}

func NewFormatDataProvider(rawTemplate []byte, options ...FormatDataOption) (FormatDataProvider, error) {
	return newFormatDataImpl(rawTemplate, options...)
}

func NewDefaultFormatDataProvider(rawTemplate []byte) (FormatDataProvider, error) {
	return NewFormatDataProvider(rawTemplate, DefaultOptions...)
}

func (fdi *formatDataImpl) AddOptions(options ...FormatDataOption) {
	fdi.addOptions(options...)
}

func (fdi *formatDataImpl) UpdateTemplate(rawTemplate []byte) error {
	return fdi.updateTemplate(rawTemplate)
}

func (fdi *formatDataImpl) FormatJSONSchema(data []byte) ([]byte, error) {
	return fdi.formatJSONSchema(data)
}

func NewFormatKeyProvider(options ...FormatKeyOption) FormatKeyProvider {
	return newFormatKeyImpl(options...)
}

func (fki *formatKeyImpl) AddOptions(options ...FormatKeyOption) {
	fki.addOptions(options...)
}

func (fki *formatKeyImpl) FormatJSONSchema(data []byte) ([]byte, error) {
	return fki.formatJSONSchema(data)
}
