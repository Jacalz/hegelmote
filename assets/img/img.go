package img

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	//go:embed power.svg
	powerIconContents []byte

	// PowerIcon is a fyne resource used to indicate a power on/off button.
	PowerIcon = theme.NewThemedResource(&fyne.StaticResource{StaticName: "power.svg", StaticContent: powerIconContents})

	//go:embed icon-512.png
	iconContents []byte

	// AppIcon holds the application icon.
	AppIcon = &fyne.StaticResource{
		StaticName:    "icon.png",
		StaticContent: iconContents,
	}
)
