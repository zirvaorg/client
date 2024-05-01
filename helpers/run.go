package helpers

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

var CommandMap = make(map[string]*Command)

func Run() error {
	if len(CommandMap) == 0 {
		return fmt.Errorf("no commands found")
	}

	if _, exists := CommandMap["help"]; !exists {
		return fmt.Errorf("no help command defined")
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		d := color.New(color.FgHiGreen, color.Bold)
		d.Print(">  ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading command: %v", err)
		}
		cmdString = strings.TrimSpace(cmdString)

		if command, exists := CommandMap[cmdString]; exists {
			command.Run()
		} else {
			CommandMap["help"].Run()
		}
	}
}
