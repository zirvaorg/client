package helpers

var Commands []*Command

type Command struct {
	Use         string
	Description string
	Run         func()
}

func AddCommand(cmd *Command) {
	Commands = append(Commands, cmd)
	CommandMap[cmd.Use] = cmd
}
