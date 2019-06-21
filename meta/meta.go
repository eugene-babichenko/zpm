package meta

// The format to write the time of the last check for updates.
const LastUpdateCheckLayout string = "Mon Jan 2 15:04:05 2006"

// Internal data required for `zpm` functioning.
type Meta struct {
	// Time passed since the last check for updates.
	LastUpdateCheck string
}
