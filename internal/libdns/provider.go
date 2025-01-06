package libdns

import (
	"encoding/json"

	"github.com/libdns/libdns"
)

type LibdnsProvider interface {
	libdns.RecordGetter
	libdns.RecordSetter
	libdns.RecordAppender
	libdns.RecordDeleter
}

func initProvider[T interface{}](confs [][]byte) (*T, error) {
	provider := new(T)
	for _, conf := range confs {
		err := json.Unmarshal(conf, provider)
		if err != nil {
			return nil, err
		}
	}
	return provider, nil
}
