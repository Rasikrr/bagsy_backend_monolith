package service

//go:generate enumer -type=Color -json -trimprefix Color -transform=snake -output color_enumer.go

type Color uint8

const (
	ColorBlue Color = iota
	ColorGreen
	ColorRed
	ColorYellow
	ColorPurple
	ColorOrange
	ColorGray
)
