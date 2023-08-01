package input

type Action int

const (
	NoAction Action = iota
	Up
	UpRight
	Right
	DownRight
	Down
	DownLeft
	Left
	UpLeft
	Start
	Select
	Action1
	Action2
	Action3
	Action4
	Action5
	Action6
)

var gActionNames = map[Action]string{
	Up:        "↑",
	UpRight:   "↗",
	Right:     "→",
	DownRight: "↘",
	Down:      "↓",
	DownLeft:  "↙",
	Left:      "←",
	UpLeft:    "↖",
	Start:     "[Start]",
	Select:    "[Select]",
	Action1:   "[1]",
	Action2:   "[2]",
	Action3:   "[3]",
	Action4:   "[4]",
	Action5:   "[5]",
	Action6:   "[6]",
}

func ActionName(action Action) string {
	return gActionNames[action]
}

type State int

const (
	Idle State = iota
	Pressed
	Held
	Released
)

type Button struct {
	Action Action
	Name   string
	State  State
}
