package types

type Config map[string]interface{}

func (c Config) SetInt(key string, value int) Config {
	c[key] = value
	return c
}
func (c Config) SetString(key string, value string) Config {
	c[key] = value
	return c
}

func (c Config) GetString(key string, def string) string {
	if v, ok := c[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}
func (c Config) GetInt(key string, def int) int {
	if v, ok := c[key]; ok {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return def
}
func (c Config) GetStringArray(key string, def []string) []string {
	if v, ok := c[key]; ok {
		if s, ok := v.([]string); ok {
			return s
		}
	}
	return def
}
func (c Config) GetBool(key string, def bool) bool {
	if v, ok := c[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return def
}
func (c Config) GetChannel(key string) chan interface{} {
	if v, ok := c[key]; ok {
		if ch, ok := v.(chan interface{}); ok {
			return ch
		}
	}
	return nil
}
