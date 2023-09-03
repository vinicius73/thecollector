package errors

var (
	ErrConfigNotFound      = Business("config not found", "TC:000")
	ErrFailToLoadConfig    = System(nil, "fail to load config", "TC:001")
	ErrFailEnsureConfig    = System(nil, "fail to ensure config", "TC:002")
	ErrFailToMarshalConfig = System(nil, "fail to marshal config", "TC:004")
	ErrInvalidAction       = Business("invalid action: %s", "TC:005")
)
