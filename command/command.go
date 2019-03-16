package command

// Command is what the client will request server
type Command struct {
	cmdtype int
	para1   string
	para2   string
}

const (
	// CMDSTOP : QUIT PROGRAM
	CMDSTOP = iota
	// CMDWeather : request the weather of a certain city
	CMDWeather
	// CMDEND : end of cmd type
	CMDEND
)

// CmdtypeStr : each command type string
var CmdtypeStr = [...]string{
	"quit",
	"weather",
}
var commands = map[string]int{
	CmdtypeStr[0]: 0,
	CmdtypeStr[1]: 1,
}

// GetCmdType turns string into int
func GetCmdType(name string) int {
	if cmdtype, ok := commands[name]; ok {
		return cmdtype
	}
	return -1
}
