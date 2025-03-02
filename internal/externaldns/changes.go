package externaldns

// Changes holds lists of actions to be executed by dns providers.
type Changes struct {
	// Records that need to be created
	Create Endpoints `json:"create"`
	// Records that need to be updated (current data)
	UpdateOld Endpoints `json:"updateOld"`
	// Records that need to be updated (desired data)
	UpdateNew Endpoints `json:"updateNew"`
	// Records that need to be deleted
	Delete Endpoints `json:"delete"`
}
