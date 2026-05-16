package ports

type StateGenerator interface {
	GenerateState(returnTo string) (string, error)
	VerifyAndExtractState(state string) (string, error)
}
