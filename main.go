package main

import (
	"fgtrainer/input"
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func uPanic(err error) {
	if err != nil {
		panic(err)
	}
}

var gGLFWJoysticks = []glfw.Joystick{
	glfw.Joystick1,
	glfw.Joystick2,
	glfw.Joystick3,
	glfw.Joystick4,
	glfw.Joystick5,
	glfw.Joystick6,
	glfw.Joystick7,
	glfw.Joystick8,
	glfw.Joystick9,
	glfw.Joystick10,
	glfw.Joystick11,
	glfw.Joystick12,
	glfw.Joystick13,
	glfw.Joystick14,
	glfw.Joystick15,
	glfw.Joystick16,
}

const rGLFWJoystickCount = 16

type JoystickState struct {
	Axes    []float32
	Buttons []glfw.Action
	Hats    []glfw.JoystickHatState

	GamepadAxes    [6]float32
	GamepadButtons [15]glfw.Action
}

func glfwFillJoystickState(glfwJoystick glfw.Joystick, joystickState *JoystickState) {
	joystickState.Axes = glfwJoystick.GetAxes()
	joystickState.Buttons = glfwJoystick.GetButtons()
	joystickState.Hats = glfwJoystick.GetHats()

	glfwGamepadState := glfwJoystick.GetGamepadState()
	if glfwGamepadState != nil {
		joystickState.GamepadAxes = glfwGamepadState.Axes
		joystickState.GamepadButtons = glfwGamepadState.Buttons
	}
}

func copyJoystickState(sourceJoystick, destinationJoystick *JoystickState) {
	destinationJoystick.Axes = sourceJoystick.Axes
	destinationJoystick.Buttons = sourceJoystick.Buttons
	destinationJoystick.Hats = sourceJoystick.Hats
	destinationJoystick.GamepadAxes = sourceJoystick.GamepadAxes
	destinationJoystick.GamepadButtons = sourceJoystick.GamepadButtons
}

func (joystickState *JoystickState) Update(glfwJoystick glfw.Joystick, joystickStateChanges *JoystickState) {
	var newJoystickState JoystickState
	glfwFillJoystickState(glfwJoystick, &newJoystickState)
	//fmt.Println("Current", joystickState)
	//fmt.Println("New", newJoystickState)
	//fmt.Println("Changes (Pre)", joystickStateChanges)

	for index := range newJoystickState.Axes {
		difference := joystickState.Axes[index] - newJoystickState.Axes[index]
		//fmt.Println(joystickState.Axes[index], newJoystickState.Axes[index], difference)
		if difference < 0 {
			difference = difference * -1.0
		}
		joystickStateChanges.Axes[index] = difference
	}

	for index := range newJoystickState.Buttons {
		joystickStateChanges.Buttons[index] = joystickState.Buttons[index]
	}

	for index := range newJoystickState.Hats {
		difference := joystickState.Hats[index] - newJoystickState.Hats[index]
		if difference < 0 {
			difference = difference * -1
		}
		joystickStateChanges.Hats[index] = difference
	}

	for index := range newJoystickState.GamepadAxes {
		difference := joystickState.GamepadAxes[index] - newJoystickState.GamepadAxes[index]
		if difference < 0 {
			difference = difference * -1.0
		}
		joystickStateChanges.GamepadAxes[index] = difference
	}

	for index := range newJoystickState.GamepadButtons {
		joystickStateChanges.GamepadButtons[index] = joystickState.GamepadButtons[index]
	}

	copyJoystickState(&newJoystickState, joystickState)
}

var gJosytickButtonState = map[glfw.Action]string{
	glfw.Release: "Release",
	glfw.Press:   "Press",
}

func (joystickState JoystickState) Changed(compareState *JoystickState) bool {
	for index, axe := range compareState.Axes {
		if axe > 0 {
			fmt.Println("Axe", index, axe)
			return false
		}
	}

	for index := range compareState.Buttons {
		if compareState.Buttons[index] != joystickState.Buttons[index] {
			fmt.Println("Button", index, gJosytickButtonState[compareState.Buttons[index]], gJosytickButtonState[joystickState.Buttons[index]])
			return false
		}
	}

	for index, hat := range compareState.Hats {
		if hat > 0 {
			fmt.Println("Hat", index, hat)
			return false
		}
	}

	for index, axe := range compareState.GamepadAxes {
		if axe > 0 {
			fmt.Println("Gamepad Axe", index, axe)
			return false
		}
	}

	for index := range compareState.GamepadButtons {
		if compareState.GamepadButtons[index] != joystickState.GamepadButtons[index] {
			fmt.Println("Gamepad Buttons", index, gJosytickButtonState[compareState.GamepadButtons[index]], gJosytickButtonState[joystickState.GamepadButtons[index]])
			return false
		}
	}

	return true
}

const rAxeDeadzone = .01

type Joystick struct {
	glfwJoystick glfw.Joystick

	Axes    []float32
	Buttons []glfw.Action
	Hats    []glfw.JoystickHatState

	//GamepadAxes    [6]float32
	GamepadButtons [21]glfw.Action
}

func (joystick *Joystick) UpdateOld(joystickState *JoystickState) {
	//joystick.State.Update(joystick.glfwJoystick, joystickState)
}

type ButtonMap struct {
}

func (joystick *Joystick) Update() {
	glfwJoystick := joystick.glfwJoystick
	joystick.Axes = glfwJoystick.GetAxes()
	joystick.Buttons = glfwJoystick.GetButtons()
	joystick.Hats = glfwJoystick.GetHats()

	glfwGamepadState := glfwJoystick.GetGamepadState()
	if glfwGamepadState != nil {
		gamepadButtons := glfwGamepadState.Buttons
		for index, gamepadButton := range gamepadButtons {
			joystick.GamepadButtons[index] = gamepadButton
		}

		gamepadAxes := glfwGamepadState.Axes
		for index := 0; index < len(gamepadAxes); index++ {
			if gamepadAxes[index] == 1 {
				joystick.GamepadButtons[len(gamepadButtons)+index] = glfw.Press
			} else {
				joystick.GamepadButtons[len(gamepadButtons)+index] = glfw.Release
			}
		}
	}
}

const rFPS = 60
const rFPSSleep = time.Second / rFPS

type IntArray []int

func (intArray IntArray) Contains(value int) bool {
	for _, element := range intArray {
		if element == value {
			return true
		}
	}

	return false
}

type BoolArray []bool

func (boolArray BoolArray) Copy(other BoolArray) {
	if len(boolArray) != len(other) {
		panic("Inequal bool array copy lengths")
	}

	for index := range boolArray {
		boolArray[index] = other[index]
	}
}

func (boolArray BoolArray) Compare(other BoolArray) bool {
	if len(boolArray) != len(other) {
		return false
	}

	for index := range boolArray {
		if boolArray[index] != other[index] {
			return false
		}
	}

	return true
}

type IntDLLNode struct {
	prev *IntDLLNode
	data int
	next *IntDLLNode
}

type IntDLLArray []IntDLLNode

func NewIntDLLArray(length int) IntDLLArray {
	if length < 2 {
		panic(fmt.Errorf("Invalid IntDLLNode length %d", length))
	}

	intDLLArray := make(IntDLLArray, length)
	var index = 1
	for ; index < length-1; index++ {
		intDLLArray[index].prev = &intDLLArray[index-1]
		intDLLArray[index].next = &intDLLArray[index+1]
	}

	intDLLArray[index].prev = &intDLLArray[index-1]
	intDLLArray[index].next = &intDLLArray[0]
	intDLLArray[0].next = &intDLLArray[1]
	intDLLArray[0].prev = &intDLLArray[index]

	return intDLLArray
}

/*
func (pressed Pressed) Contains(value int) bool {
	for _, element := range pressed {
		if element == value {
			return true
		}
	}

	return false
}

type PressedAction struct {
	pressed Pressed
	action  input.Action
}

type PressedActionArray []*PressedAction

func (pressedActionArray PressedActionArray) ContainsPressed(pressed Pressed) bool {
	for _, pressedAction := range pressedActionArray {
		var contains []int
		for _, press := range pressedAction.pressed {
			if pressed.Contains(press) {
				contains = append(contains, press)
			}
		}

		if len(contains) == len(pressed) {
			return true
		}
	}

	return false
}
*/

/*
func (pressedActionArray PressedActionArray) ContainsPressed(pressed Pressed) bool {
	for _, pressedAction := range pressedActionArray {
		if pressedAction.pressed == pressed {
			return true
		}
	}
	return false
}
*/

// ButtonN+
type ButtonNP struct {
	buttons []int
	action  input.Action
}

type ButtonNPArray []ButtonNP

func (buttonNPArray ButtonNPArray) Get(buttons IntArray) *ButtonNP {

	for _, buttonNP := range buttonNPArray {
		if len(buttonNP.buttons) != len(buttons) {
			continue
		}

		found := 0
		for _, button := range buttonNP.buttons {
			if buttons.Contains(button) {
				found++
			}
		}

		if found == len(buttons) {
			return &buttonNP
		}
	}

	return nil
}

func main() {
	soundFile, err := os.Open("files/GOAT.mp3")
	uPanic(err)

	nash, format, err := mp3.Decode(soundFile)
	uPanic(err)
	defer nash.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	nashloop := beep.Loop(-1, nash)

	nashvol := &effects.Volume{
		Streamer: nashloop,
		Base:     10.0,
		Volume:   -1.5,
		Silent:   false,
	}

	uPanic(glfw.Init())
	defer glfw.Terminate()

	var joysticks [rGLFWJoystickCount]Joystick
	var joystickStates [rGLFWJoystickCount]JoystickState

	for index, glfwJoystick := range gGLFWJoysticks {
		joystickState := &joystickStates[index]
		joystick := &joysticks[index]
		joystick.glfwJoystick = glfwJoystick

		glfwFillJoystickState(glfwJoystick, joystickState)
		//copyJoystickState(joystickState, &joystick.State)
		//fmt.Println(index, joystick.State)
	}

	fmt.Println("Hold any four buttons on any connected device")
	complete := false
	var joystickIndex int
	startTime := time.Now()
	for !complete {
		for index := range joysticks {
			joystick := &joysticks[index]
			glfwJoystick := joystick.glfwJoystick
			if !glfwJoystick.Present() {
				continue
			}

			joystick.Update()

			var pressed int
			for _, button := range joystick.Buttons {
				if button == glfw.Press {
					pressed++
				}
			}

			if pressed >= 4 {
				complete = true
				joystickIndex = index
				break
			}
		}

		futureTime := startTime.Add(rFPSSleep)
		sleepDuration := rFPSSleep - time.Now().Sub(startTime)
		time.Sleep(sleepDuration)
		glfw.PollEvents()
		startTime = futureTime
	}

	joystick := &joysticks[joystickIndex]

	fmt.Println(fmt.Sprintf("[%d]", joystickIndex+1), gGLFWJoysticks[joystickIndex].GetGamepadName(), "selected")
	fmt.Println("Release to cotinue")

	startTime = time.Now()
	for true {
		joystick.Update()

		var held int
		for _, button := range joystick.Buttons {
			if button == glfw.Press {
				held++
			}
		}

		if held == 0 {
			break
		}

		futureTime := startTime.Add(rFPSSleep)
		sleepDuration := rFPSSleep - time.Now().Sub(startTime)
		time.Sleep(sleepDuration)
		glfw.PollEvents()
		startTime = futureTime
	}

	actions := []input.Action{
		input.Start,
		input.Select,
		input.Up,
		input.Down,
		input.Left,
		input.Right,
		input.Action1,
		input.Action2,
		input.Action3,
		input.Action4,
		input.Action5,
		input.Action6,
	}

	actionsMap := make(map[input.Action]int)

	//buttonHistory := NewIntDLLArray(len(joystick.Buttons))
	//buttonHistoryCurrent := &buttonHistory[0]
	pressed := make(BoolArray, len(joystick.Buttons))
	actionArray := make([]input.Action, len(joystick.Buttons))

	fmt.Println("")
	fmt.Println("Map the following actions")
	for _, action := range actions {
		startTime = time.Now()

		actionName := input.ActionName(action)
		fmt.Print(actionName + " ")
		button := -1

		var complete bool
		for true {
			joystick.Update()

			for index, glfwAction := range joystick.Buttons {
				if glfwAction != glfw.Press {
					pressed[index] = false
					continue
				}

				if !pressed[index] {
					button = index
					//fmt.Printf("%d ", button)
				}

				pressed[index] = true
			}

			var held int
			for index := range pressed {
				if pressed[index] {
					held++
				}
			}

			if button > -1 && held == 0 {
				if actionArray[button] == input.NoAction {
					actionsMap[action] = button
					actionArray[button] = action
					complete = true
				}

				button = -1
			}

			futureTime := startTime.Add(rFPSSleep)
			sleepDuration := rFPSSleep - time.Now().Sub(startTime)
			time.Sleep(sleepDuration)
			glfw.PollEvents()
			startTime = futureTime

			if complete {
				break
			}
		}
	}

	fmt.Println()
	fmt.Println(actionArray)

	buttonNPs := ButtonNPArray{
		ButtonNP{
			buttons: []int{actionsMap[input.Up], actionsMap[input.Right]},
			action:  input.UpRight,
		},
		ButtonNP{
			buttons: []int{actionsMap[input.Right], actionsMap[input.Down]},
			action:  input.DownRight,
		},
		ButtonNP{
			buttons: []int{actionsMap[input.Down], actionsMap[input.Left]},
			action:  input.DownLeft,
		},
		ButtonNP{
			buttons: []int{actionsMap[input.Left], actionsMap[input.Up]},
			action:  input.UpLeft,
		},
	}

	fmt.Println()
	fmt.Printf("Enter your command.  Press %s to confirm or %s to restart\n", input.ActionName(input.Start), input.ActionName(input.Select))

	pressedShadow := make(BoolArray, len(joystick.Buttons))

	var inputArray []input.Action
	startTime = time.Now()
	var reset bool
	complete = false
	for true {
		joystick.Update()

		for index, glfwAction := range joystick.Buttons {
			press := (glfwAction == glfw.Press)
			pressed[index] = press
			if press {
				if index == actionsMap[input.Start] {
					complete = true
					break
				} else if index == actionsMap[input.Select] {
					reset = true
					break
				}
			}
		}

		active := make([]int, 0, len(pressed))
		for index := range pressed {
			if pressed[index] {
				active = append(active, index)
			}
		}

		if len(active) > 0 && !pressed.Compare(pressedShadow) {
			//fmt.Printf("%v ", active)
			action := input.NoAction
			if len(active) == 1 {
				action = actionArray[active[0]]
			} else if buttonNP := buttonNPs.Get(active); buttonNP != nil {
				action = buttonNP.action
			}

			if action != input.NoAction {
				fmt.Printf("%s ", input.ActionName(action))

				if action != input.Start && action != input.Select {
					inputArray = append(inputArray, action)
				}
			}
		} else if len(active) == 0 {
			if complete {
				break
			} else if reset {
				inputArray = []input.Action{}
				fmt.Println()
				fmt.Println()
				reset = false
			}
		}

		pressedShadow.Copy(pressed)

		futureTime := startTime.Add(rFPSSleep)
		sleepDuration := rFPSSleep - time.Now().Sub(startTime)
		time.Sleep(sleepDuration)
		glfw.PollEvents()
		startTime = futureTime
	}

	fmt.Println()
	fmt.Println()
	speaker.Play(nashvol)

	var inputArrayPosition int
	var failure bool
	pressOrder := make(IntArray, 0, len(pressed))
	startTime = time.Now()
	for true {
		joystick.Update()

		for index, glfwAction := range joystick.Buttons {
			press := (glfwAction == glfw.Press)
			pressed[index] = press
			if press {
				if index == actionsMap[input.Start] {
					complete = true
					break
				} else if index == actionsMap[input.Select] {
					reset = true
					break
				}
			}
		}

		active := make([]int, 0, len(pressed))
		for index := range pressed {
			if pressed[index] {
				active = append(active, index)
				if !pressOrder.Contains(index) {
					pressOrder = append(pressOrder, index)
				}
			}
		}

		if len(active) > 0 && !pressed.Compare(pressedShadow) {
			//fmt.Printf("%v ", active)
			action := input.NoAction
			if buttonNP := buttonNPs.Get(active); buttonNP != nil {
				action = buttonNP.action
			} else {
				action = actionArray[pressOrder[len(pressOrder)-1]]
			}

			fmt.Printf("%s ", input.ActionName(action))

			if action != inputArray[inputArrayPosition] {
				if !failure {
					fmt.Print("[FAILURE] ")
					inputArrayPosition = 0
					failure = true
				}
			} else {
				inputArrayPosition++
			}
		} else if len(active) == 0 {
			if failure {
				fmt.Println()
				failure = false
			}
			pressOrder = pressOrder[:0]
		}

		if inputArrayPosition == len(inputArray) {
			fmt.Println("[SUCCESS]")
			inputArrayPosition = 0
		}

		pressedShadow.Copy(pressed)

		futureTime := startTime.Add(rFPSSleep)
		sleepDuration := rFPSSleep - time.Now().Sub(startTime)
		time.Sleep(sleepDuration)
		glfw.PollEvents()
		startTime = futureTime
	}

	return

	for true {
		futureTime := startTime.Add(rFPSSleep)
		sleepDuration := rFPSSleep - time.Now().Sub(startTime)
		time.Sleep(sleepDuration)
		glfw.PollEvents()
		startTime = futureTime
	}

	return

	for true {
		for index := range joysticks {
			joystick := &joysticks[index]
			if !joystick.glfwJoystick.Present() {
				continue
			}

			joystick.Update()
			/*if !joystick.State.Changed(joystickState) {
				fmt.Println(index, joystick.State)
			}*/
		}

		fmt.Scanln()
		//time.Sleep(1 * time.Millisecond)
		glfw.PollEvents()
	}

}
