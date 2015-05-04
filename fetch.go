package tail

type Fetch interface {
	GetData(id string) []byte
	SetData(id string, data interface{}) error
}
