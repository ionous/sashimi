package source

// statement options
type Options map[string]string

func (opts Options) GetOption(name, defaultValue string) (ret string) {
	if v, ok := opts[name]; ok {
		ret = v
	} else {
		ret = defaultValue
	}
	return ret
}
