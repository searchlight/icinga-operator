package monitoring

type Receiver struct {
	// For which state notification will be sent
	State string

	// To whom notification will be sent
	To []string

	// How this notification will be sent
	Notifier string
}
