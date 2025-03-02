package libdnsregistry

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/libdns/libdns"
)

type Provider interface {
	libdns.RecordGetter
	libdns.RecordSetter
	libdns.RecordAppender
	libdns.RecordDeleter
}

type (
	RegistryStore    map[string]*RegistryProvider
	RegistryProvider struct {
		Init        func(conf [][]byte) (Provider, error)
		URL         string
		Description string
	}
)

func initProvider[T interface{}](confs [][]byte) (*T, error) {
	provider := new(T)
	for _, conf := range confs {
		err := json.Unmarshal(conf, provider)
		if err != nil {
			return nil, fmt.Errorf("decoding provider config failed: %w", err)
		}
	}

	return provider, nil
}

func List() []string {
	providers := []string{}
	for k := range registry {
		providers = append(providers, k)
	}

	slices.Sort(providers)

	return providers
}

// Init a new libdns provider with the given configuration.
//
//nolint:ireturn // interacted by purpose
func New(providerName string, confs [][]byte) (Provider, error) {
	reg, exist := registry[providerName]
	if !exist {
		return nil, fmt.Errorf("libdns provider '%s' does not exist", providerName)
	}

	provider, err := reg.Init(confs)
	if err != nil {
		return nil, fmt.Errorf("failed to create new libdns provider %s: %w", providerName, err)
	}

	return provider, nil
}
