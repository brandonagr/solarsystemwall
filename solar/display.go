package solar

// Display that can render the field
type Display interface {
	Render(*System)
	Dispose()
}
