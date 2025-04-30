package datadisplay

type NotificationVariant int

const (
	Success NotificationVariant = iota
	Error
	Warning
	Info
)

type Size int

const (
	XS Size = iota
	SM
	MD
	LG
	XL
	XL2
)

type IconType int

const (
	XIcon IconType = iota
)

type Color int

const (
	White Color = iota
	Black
	Gray
	Blue
	Red
	Yellow
)
