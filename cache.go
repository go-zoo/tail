package tail

// Cache is a interface that is responsible of get, set and delete data.
type Cache interface {
	Get(id string) []byte
	Set(id string, data []byte) error
	Update(id string, data []byte) error
	Del(id string) error
}
