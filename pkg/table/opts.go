package table

// WithHeader option can be supplied to indicate the presence of
// a header when using ReadCSV or supplied otherwise to indicate
// whether the header should be rendered
func WithHeader(value bool) Option {
	return func(t *Table) {
		t.header = value
	}
}

func WithFooter(value bool) Option {
	return func(t *Table) {
		t.footer = value
	}
}

func WithOffsetLimit(offset, limit uint) Option {
	return func(t *Table) {
		t.offset, t.limit = offset, limit
	}
}

func WithMergeCells() Option {
	return func(t *Table) {
		t.merge = true
	}
}
