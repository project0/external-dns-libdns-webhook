package externaldns

// Changes holds lists of actions to be executed by dns providers.
type Changes struct {
	// Records that need to be created
	Create Endpoints `json:"create,omitempty"`
	// Records that need to be updated (current data)
	UpdateOld Endpoints `json:"updateOld,omitempty"`
	// Records that need to be updated (desired data)
	UpdateNew Endpoints `json:"updateNew,omitempty"`
	// Records that need to be deleted
	Delete Endpoints `json:"delete,omitempty"`
}
