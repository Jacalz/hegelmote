package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Jacalz/hegelmote/device"
	"github.com/Jacalz/hegelmote/remote"
)

func runInteractiveMode(control *remote.Control) {
	input := bufio.NewScanner(os.Stdin)
	input.Split(bufio.ScanLines)

	for input.Scan() {
		line := input.Text()

		commands := strings.Split(line, " ")
		if len(commands) < 2 {
			fmt.Fprintln(os.Stderr, "Invalid command. Exiting...")
			os.Exit(2)
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
			fmt.Fprintln(os.Stderr, "Invalid command. Exiting...")
			os.Exit(2)
		}
	}
}

func handlePowerCommand(subcommands []string, control *remote.Control) {
	switch subcommands[0] {
	case "on":
		control.SetPower(true)
	case "off":
		control.SetPower(false)
	case "toggle":
		control.TogglePower()
	case "get":
		control.GetPower()
	}
}

func handleVolumeCommand(subcommands []string, control *remote.Control) {
	switch subcommands[0] {
	case "up":
		control.VolumeUp()
	case "down":
		control.VolumeDown()
	case "set":
		if len(subcommands) > 1 {
			percentage, _ := strconv.ParseUint(subcommands[1], 10, 8)
			control.SetVolume(uint8(percentage))
		}
	case "mute":
		control.SetVolumeMute(true)
	case "unmute":
		control.SetVolumeMute(false)
	case "get":
		control.GetVolume()
	case "muted":
		control.GetVolumeMute()
	}
}

func handleInputCommand(subcommands []string, control *remote.Control) {
	switch subcommands[0] {
	case "set":
		if len(subcommands) > 1 {
			control.SetSourceInput(device.H95, subcommands[1])
		}
	case "get":
		control.GetSourceInput(device.H95)
	}
}

func handleResetCommand(subcommands []string, control *remote.Control) {
	switch subcommands[0] {
	case "stop":
		control.StopResetDelay()
	case "get":
		control.GetResetDelay()
	default:
		delay, _ := strconv.ParseUint(subcommands[0], 10, 8)
		control.SetResetDelay(uint8(delay))
	}
}
