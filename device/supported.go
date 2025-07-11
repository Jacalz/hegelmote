package device

var supportedDeviceNames = [...]string{"Röst", "H95", "H120", "H190", "H390", "H590", "H190V"}

var (
	InputsRöst = InputsH120

	InputsH95 = [...]string{
		"Analog 1",
		"Analog 2",
		"Coaxial",
		"Optical 1",
		"Optical 2",
		"Optical 3",
		"USB",
		"Network",
	}

	InputsH120 = [...]string{
		"Balanced",
		"Analog 1",
		"Analog 2",
		"Coaxial",
		"Optical 1",
		"Optical 2",
		"Optical 3",
		"USB",
		"Network",
	}

	InputsH190 = InputsH120

	InputsH190V = []string{
		"XLR",
		"Analog 1",
		"Analog 2",
		"Coaxial",
		"Optical 1",
		"Optical 2",
		"Optical 3",
		"USB",
		"Network",
		"Phono",
	}

	InputsH390 = [...]string{
		"XLR",
		"Analog 1",
		"Analog 2",
		"BNC",
		"Coaxial",
		"Optical 1",
		"Optical 2",
		"Optical 3",
		"USB",
		"Network",
	}

	InputsH590 = [...]string{
		"XLR 1",
		"XLR 2",
		"Analog 1",
		"Analog 2",
		"BNC",
		"Coaxial",
		"Optical 1",
		"Optical 2",
		"Optical 3",
		"USB",
		"Network",
	}
)
