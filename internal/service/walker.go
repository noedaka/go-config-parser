package service

func walk(data any, cb func(key string, value any)) {
	switch v := data.(type) {
	case map[string]any:
		for key, val := range v {
			cb(key, val)
			walk(val, cb)
		}
	case []any:
		for _, item := range v {
			walk(item, cb)
		}
	}
}
