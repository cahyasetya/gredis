package parser

import (
	"encoding/binary"
	"github.com/cahyasetya/gredis/constants" // Import the constants package
	"github.com/cahyasetya/gredis/types"     // Import the types package
)

func SplitCommand(message types.Message) []string {
	cmdCount := binary.LittleEndian.Uint32(message[0:constants.LenSize]) // Use constants.LenSize
	commands := make([]string, 0, cmdCount)
	i := constants.LenSize
	for cmdCount > 0 {
		command, length := getCommand(message[i:])
		commands = append(commands, command)
		i += int(length) + constants.LenSize
		cmdCount--
	}
	return commands
}

func getCommand(message types.Message) (string, uint32) {
	length := binary.LittleEndian.Uint32(message[0:constants.LenSize]) // Use constants.LenSize
	return string(message[constants.LenSize : constants.LenSize+length]), length
}
