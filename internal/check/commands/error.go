package commands

// CheckerFailure indicates a failed check i.e. the system does not meet requirements.
type CheckerFailure struct {
	Message string
}

func (e *CheckerFailure) Error() string {
	return e.Message
}
