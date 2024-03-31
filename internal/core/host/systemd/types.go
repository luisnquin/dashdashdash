package systemd

type UnitScope uint8

// Systemd unit scopes.
const (
	ALL_ONLY = iota
	USER_ONLY
	SYSTEM_ONLY
)
