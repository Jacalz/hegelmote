package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Jacalz/hegelmote/device"
	"github.com/Jacalz/hegelmote/remote"
)

var errInvalidCommand = errors.New("invalid command format")

func runInteractiveMode(control *remote.Control) {
	input := bufio.NewScanner(os.Stdin)
	input.Split(bufio.ScanLines)

	for input.Scan() {
		line := input.Text()

		commands := strings.Split(line, " ")
		if len(commands) < 2 {
			exitWithError(errInvalidCommand)
		}

		switch commands[0] {
		case "power":
			handlePowerCommand(commands[1:], control)
		case "volume":
			handleVolumeCommand(commands[1:], control)
		case "input":
			handleInputCommand(commands[1:], control)
		case "reset":
			handleResetCommand(commands[1:], control)
		case "exit", "quit":
			return
		default:
			exitWithError(errInvalidCommand)
		}
	}
}

func handlePowerCommand(subcommands []string, control *remote.Control) {
	switch subcommands[0] {
	case "on":
		err := control.SetPower(true)
		if err != nil {
			exitWithError(err)
		}
	case "off":
		err := control.SetPower(false)
		if err != nil {
			exitWithError(err)
		}
	case "toggle":
		err := control.TogglePower()
		if err != nil {
			exitWithError(err)
		}
	case "get":
		on, err := control.GetPower()
		if err != nil {
			exitWithError(err)
		}

		fmt.Println("Powered on:", on)
	}
}

func handleVolumeCommand(subcommands []string, control *remote.Control) {
	switch subcommands[0] {
	case "up":
		err := control.VolumeUp()
		if err != nil {
			exitWithError(err)
		}
	case "down":
		err := control.VolumeDown()
		if err != nil {
			exitWithError(err)
		}
	case "set":
		if len(subcommands) > 1 {
			percentage, _ := strconv.ParseUint(subcommands[1], 10, 8)
			err := control.SetVolume(uint8(percentage))
			if err != nil {
				exitWithError(err)
			}
		}
	case "mute":
		err := control.SetVolumeMute(true)
		if err != nil {
			exitWithError(err)
		}
	case "unmute":
		err := control.SetVolumeMute(false)
		if err != nil {
			exitWithError(err)
		}
	case "get":
		volume, err := control.GetVolume()
		if err != nil {
			exitWithError(err)
		}

		fmt.Println("Volume percentage:", volume)
	case "muted":
		muted, err := control.GetVolumeMute()
		if err != nil {
			exitWithError(err)
		}

		fmt.Println("Volume muted:", muted)
	}
}

func handleInputCommand(subcommands []string, control *remote.Control) {
	switch subcommands[0] {
	case "set":
		if len(subcommands) != 2 {
			exitWithError(errInvalidCommand)
		}

		err := control.SetSourceInput(device.H95, subcommands[1])
		if err != nil {
			exitWithError(err)
		}
	case "get":
		input, err := control.GetSourceInput(device.H95)
		if err != nil {
			exitWithError(err)
		}

		fmt.Println("Selected input:", input)
	}
}

func handleResetCommand(subcommands []string, control *remote.Control) {
	switch subcommands[0] {
	case "stop":
		err := control.StopResetDelay()
		if err != nil {
			exitWithError(err)
		}
	case "get":
		delay, stopped, err := control.GetResetDelay()
		if err != nil {
			exitWithError(err)
		}

		if stopped {
			fmt.Println("Reset timeout: stopped")
			return
		}

		fmt.Println("Time until reset:", delay)
	default:
		delay, _ := strconv.ParseUint(subcommands[0], 10, 8)
		err := control.SetResetDelay(uint8(delay))
		if err != nil {
			exitWithError(err)
		}
	}
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}
