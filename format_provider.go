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
	return newFormatItemImpl(rawTemplate, options...)
}

func NewDefaultFormatDataProvider(rawTemplate []byte) (FormatDataProvider, error) {
	return NewFormatDataProvider(rawTemplate, DefaultOptions...)
}

func (fii *formatItemsImpl) AddOptions(options ...FormatDataOption) {
	fii.addOptions(options...)
}

func (fii *formatItemsImpl) UpdateTemplate(rawTemplate []byte) error {
	return fii.updateTemplate(rawTemplate)
}

func (fii *formatItemsImpl) FormatJSONSchema(data []byte) ([]byte, error) {
	return fii.formatJSONSchema(data)
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
