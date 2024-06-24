//go:generate enumer -type Tag -json -sql -text -transform=snake
package logging

type Tag int

const (
	_ Tag = iota
	Service
	App
	DB
	Env
	Error
	Panic
	Args
	Address
	Formname
	Filename
	Size
	ID
	Len
)
