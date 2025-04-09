package vars_cache

type Cache interface {
	Set(key string, value string)
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
}

func NewCache() Cache {
	return new(memoryCache)
}
