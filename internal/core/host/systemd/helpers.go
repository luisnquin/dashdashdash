package systemd

func getAllUnitStatuses() []string {
	return []string{
		"active",
		"reloading",
		"inactive",
		"failed",
		"activating",
		"deactivating",
		"maintenance",
	}
}

func getAllUnitTypes() []string {
	return []string{
		"service",
		"socket",
		"target",
		"device",
		"mount",
		"automount",
		"swap",
		"timer",
		"path",
		"slice",
		"scope",
	}
}

func parseListUnitsScopeParam(scope string) UnitScope {
	switch scope {
	case "all":
		return ALL_ONLY
	case "user":
		return USER_ONLY
	case "system":
		return SYSTEM_ONLY
	default:
		return ALL_ONLY
	}
}
