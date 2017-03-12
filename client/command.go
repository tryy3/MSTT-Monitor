package client

import (
	"strings"
)

// Command är en struct för kommandon
type Command struct {
	Name   string
	Params []Argument
}

// Argument är en struct för kommand arguments/flaggor
type Argument struct {
	Name  string
	Value string
}

// ParseCommand analyserar en rå kommand sträng
// och strukterar den till en Command Struct
// Example: Input: Check_cpu -cpu=1
// Output: {
//      Name: "Check_cpu"
//      Params: [
//            Name: "-cpu"
//         Value: "1"
//      ]
// }
func ParseCommand(command string) Command {
	commandSplit := strings.Split(command, " ")
	cmd := Command{Name: commandSplit[0], Params: []Argument{}}

	// Kommandot har flaggor
	if len(commandSplit) > 1 {
		for i := 1; i < len(commandSplit); i++ {
			argumentSplit := strings.Split(commandSplit[i], "=")
			val := ""
			if len(argumentSplit) > 1 {
				val = argumentSplit[1]
			}
			cmd.Params = append(cmd.Params, Argument{Name: argumentSplit[0], Value: val})
		}
	}

	return cmd
}
