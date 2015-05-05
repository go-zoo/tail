package tail

// Cache is a interface that is responsible for set or get data from de template
type Cache interface {
	Get(id string) []byte
	Set(id string, data []byte) error
}
