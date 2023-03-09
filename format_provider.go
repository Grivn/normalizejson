package njson

func NewFormatSchemaProvider(rawTemplate []byte, options ...FormatOption) (FormatProvider, error) {
	return newFormatSchemaImpl(rawTemplate, options...)
}

func NewDefaultFormatSchemaProvider(rawTemplate []byte) (FormatProvider, error) {
	return NewFormatSchemaProvider(rawTemplate, DefaultFormatDataOptions...)
}

func (fsi *formatSchemaImpl) AddOptions(options ...FormatOption) {
	fsi.addOptions(options...)
}

func (fsi *formatSchemaImpl) UpdateTemplate(rawTemplate []byte) error {
	return fsi.updateTemplate(rawTemplate)
}

func (fsi *formatSchemaImpl) FormatJSONSchema(data []byte) ([]byte, error) {
	return fsi.formatJSONSchema(data)
}

func (fsi *formatSchemaImpl) Reset() {
	fsi.reset()
}

func NewFormatDataProvider(rawTemplate []byte, options ...FormatOption) (FormatProvider, error) {
	return newFormatDataImpl(rawTemplate, options...)
}

func NewDefaultFormatDataProvider(rawTemplate []byte) (FormatProvider, error) {
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

func (fdi *formatDataImpl) Reset() {
	fdi.reset()
}

func NewFormatKeyProvider(options ...FormatOption) FormatProvider {
	return newFormatKeyImpl(options...)
}

func (fki *formatKeyImpl) AddOptions(options ...FormatOption) {
	fki.addOptions(options...)
}

func (fki *formatKeyImpl) UpdateTemplate(rawTemplate []byte) error {
	// do nothing.
	return nil
}

func (fki *formatKeyImpl) FormatJSONSchema(data []byte) ([]byte, error) {
	return fki.formatJSONSchema(data)
}

func (fki *formatKeyImpl) Reset() {
	fki.reset()
}
