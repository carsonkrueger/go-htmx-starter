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
	XL3
	XL4
	XL5
	XL6
	XL7
	XL8
	XL9
	XL10
	XL11
	XL12
)

type IconType int

const (
	LeafIcon IconType = iota
	FlexIcon
	AntioxidantIcon
)

type Color int

const (
	White Color = iota
	Black
	Gray
	Blue
	Red
	Yellow
	Primary
	PrimaryHi
	PrimaryLo
	Secondary
	SecondaryHi
	SecondaryLo
)

var sizeClassMap map[Size]string = map[Size]string{
	XS:   "size-3",
	SM:   "size-4",
	MD:   "size-5",
	LG:   "size-6",
	XL:   "size-7",
	XL2:  "size-8",
	XL3:  "size-9",
	XL4:  "size-10",
	XL5:  "size-11",
	XL6:  "size-12",
	XL7:  "size-13",
	XL8:  "size-14",
	XL9:  "size-15",
	XL10: "size-16",
	XL11: "size-17",
	XL12: "size-18",
}

var fillClassMap map[Color]string = map[Color]string{
	White:       "fill-white",
	Black:       "fill-black",
	Gray:        "fill-gray-500",
	Blue:        "fill-blue-500",
	Red:         "fill-red-500",
	Yellow:      "fill-yellow-500",
	Primary:     "fill-primary",
	PrimaryHi:   "fill-primary-hi",
	PrimaryLo:   "fill-primary-lo",
	Secondary:   "fill-secondary",
	SecondaryLo: "fill-secondary-lo",
	SecondaryHi: "fill-secondary-hi",
}

var borderClassMap map[Color]string = map[Color]string{
	White:       "",
	Black:       "border-black",
	Gray:        "border-gray-500",
	Blue:        "border-blue-500",
	Red:         "border-red-500",
	Yellow:      "border-yellow-500",
	Primary:     "border-primary",
	PrimaryHi:   "border-primary-hi",
	PrimaryLo:   "border-primary-lo",
	Secondary:   "border-secondary",
	SecondaryHi: "border-secondary-hi",
	SecondaryLo: "border-secondary-lo",
}
