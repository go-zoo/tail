package tail

import (
	"io/ioutil"
)

func ReadAssetFile(src string) ([]byte, error) {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}
	return data, nil
}
