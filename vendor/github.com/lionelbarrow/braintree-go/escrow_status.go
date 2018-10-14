package braintree

type EscrowStatus string

const (
	EscrowStatusHoldPending    EscrowStatus = "hold_pending"
	EscrowStatusHeld           EscrowStatus = "held"
	EscrowStatusReleasePending EscrowStatus = "release_pending"
	EscrowStatusReleased       EscrowStatus = "released"
	EscrowStatusRefunded       EscrowStatus = "refunded"
)
