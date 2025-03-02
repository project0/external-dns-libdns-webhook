package externaldns

import (
	"context"
)

type Provider interface {
	Records(ctx context.Context) (records Endpoints, err error)
	ApplyChanges(ctx context.Context, changes *Changes) (err error)
	// AdjustEndpoints canonicalizes a set of candidate endpoints.
	// It is called with a set of candidate endpoints obtained from the various sources.
	// It returns a set modified as required by the provider. The provider is responsible for
	// adding, removing, and modifying the ProviderSpecific properties to match
	// the endpoints that the provider returns in `Records` so that the change plan will not have
	// unnecessary (potentially failing) changes. It may also modify other fields, add, or remove
	// Endpoints. It is permitted to modify the supplied endpoints.
	AdjustEndpoints(ctx context.Context, endpoints Endpoints) (updated Endpoints, err error)
	GetInitialization(ctx context.Context) Initialization
}
