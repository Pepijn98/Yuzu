package yuzu

// Command defines what is considered a "command"
type Command interface {
	IsOwnerOnly() bool
	Process(Context)
	// First description, second usage
	Help() [2]string
}
