/*
Author: Luuminous, Chengyang, Yichen
Date: 30, Oct, 2018
*/

package main

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/layout"
	//"github.com/fyne-io/fyne/theme"
	"github.com/fyne-io/fyne/dialog"
	"github.com/fyne-io/fyne/widget"
	"github.com/fyne-io/fyne/examples/apps"
	"strconv"
	//"desktop"
	//"reflect"
	"fmt"
	"math/rand"
	"math"
	"time"
	"os"
	"bufio"
	"strings"
)

type Player struct {
	Hands                       []Card //This is the holehands for the player
	SeatPosition                int
	GamePosition                int
	Chips, Bet                  int
	Fold, AllIn, OK, Eliminated bool
	PositionName                string
	Pattern                     string
	Name                        string
	PlayerType                  string
	OpponentEvaluation			map[int]float64
}

type Card struct {
	Num   int
	Color int
}

type Current struct {
	Pool          []Card
	CommunityCard []Card
	ChipPool      int
	Players       []Player
	StartPlayer   int
	CurrentBet    int
	GameCount     int
	Stage         string
	PreEventsList []string
	HumanPosition int
}

var current Current
var numPlayers int
var initialChips int
var numTurns int

/*
0: Human player, 1: Random bot, 2: Probability based bot, 
3: Conventional rules based bot, 4: Persona based bot, 5: Minimum regret based bot
*/
var numOfEachTypePlayers [6]int
var numOfEachTypePlayersEntry [6]*widget.Label

/*
Store the human player's name. If the number of human player is 0, playerName has no effect.
*/
var playerName string
var botIndexList []int
var numPlayersEntry *widget.Label
var initialChipsEntry *widget.Label
var numTurnsEntry *widget.Label
var playerNameEntry *widget.Label
var humanActionList []string

//This is to store all the player's information to be displayed.
var playerInformationEntry [][]*widget.Label
var conditionInformationEntry []*widget.Label

//preFlopProbability is a map[string][]string, 1D is a string made by two start cards, 2D is a index for numPlayers.
var preFlopProbability map[string][]float64

//coefficient is a [][]float64
var coefficient map[int][]float64

//personaMap is to record the action for each players
//alphaMap is a evaluation for each players
var personaMap map[int][4]int
var alphaMap []float64

/*
Author: Luuminous
Date: Nov. 22nd, 2018
*/

/*
This function is to refresh all the information on the board.
*/

func Refresh() {
	/*
	Refresh human player.
	*/

	//Chips:
	playerInformationEntry[0][0].SetText("Chips: " + strconv.Itoa(current.Players[current.HumanPosition].Chips))
	//Bets:
	playerInformationEntry[0][1].SetText("Bet: " + strconv.Itoa(current.Players[current.HumanPosition].Bet))
	//PositionName:
	playerInformationEntry[0][2].SetText("Position: " + current.Players[current.HumanPosition].PositionName)
	//Hands:
	playerInformationEntry[0][3].SetText("Hands: " + HandsToString(current.Players[current.HumanPosition].Hands))

	/*
	Refresh bot player.
	*/

	for i := 0; i < len(botIndexList); i++ {
		//Chips:
		playerInformationEntry[i + 1][0].SetText("Chips: " + strconv.Itoa(current.Players[botIndexList[i]].Chips))
		//Bets:
		playerInformationEntry[i + 1][1].SetText("Bet: " + strconv.Itoa(current.Players[botIndexList[i]].Bet))
		//PositionName:
		playerInformationEntry[i + 1][2].SetText("Position: " + current.Players[botIndexList[i]].PositionName)
	}

	/*
	Refresh condition information.
	*/

	conditionInformationEntry[0].SetText("NumberOfTurn: " + strconv.Itoa(current.GameCount) + "\nStage: " + current.Stage + "\nCurrentBet: " + strconv.Itoa(current.CurrentBet) + "\nCommunityCard: " + HandsToString(current.CommunityCard))
	conditionInformationEntry[1].SetText(SetInformationScroll(current.PreEventsList, 4))
}

/*
Author: Luuminous
Date: Nov. 4th, 2018
*/

/*
This function is to create a window for information display and enter.
*/

func InformationEnter(app fyne.App) {
	infoW := app.NewWindow("Welcome!")

	//Set information entry, default is numPlayers = 6, initialChips = 1000, numTurns = 10, playerName = Rosemary
	formNumPlayers := widget.NewEntry()
	formNumPlayers.Text = strconv.Itoa(numPlayers)

	formInitialChips := widget.NewEntry()
	formInitialChips.Text = strconv.Itoa(initialChips)

	formNumTurns := widget.NewEntry()
	formNumTurns.Text = strconv.Itoa(numTurns)

	formPlayerName := widget.NewEntry()
	formPlayerName.Text = playerName

	//Set information entry for the number of each type of players.
	formHumanPlayer := widget.NewEntry()
	formHumanPlayer.Text = strconv.Itoa(numOfEachTypePlayers[0])

	formRandomBot := widget.NewEntry()
	formRandomBot.Text = strconv.Itoa(numOfEachTypePlayers[1])

	formProbabilityBot := widget.NewEntry()
	formProbabilityBot.Text = strconv.Itoa(numOfEachTypePlayers[2])

	formConventionalBot := widget.NewEntry()
	formConventionalBot.Text = strconv.Itoa(numOfEachTypePlayers[3])

	formPersonaDrawerBot := widget.NewEntry()
	formPersonaDrawerBot.Text = strconv.Itoa(numOfEachTypePlayers[4])

	formSuperTalentedBot := widget.NewEntry()
	formSuperTalentedBot.Text = strconv.Itoa(numOfEachTypePlayers[5])

	form := &widget.Form {
		OnSubmit: func() {
			if CheckNumPlayers(formNumPlayers.Text) && CheckInitialChips(formInitialChips.Text) && CheckNumTurns(formNumTurns.Text) && CheckPlayerName(formPlayerName.Text) && CheckNumType(formHumanPlayer.Text, formRandomBot.Text, formProbabilityBot.Text, formConventionalBot.Text, formPersonaDrawerBot.Text, formSuperTalentedBot.Text) {
				numPlayersEntry.SetText("numPlayers: " + strconv.Itoa(numPlayers))
				initialChipsEntry.SetText("initialChips: " + strconv.Itoa(initialChips))
				numTurnsEntry.SetText("numTurns: " + strconv.Itoa(numTurns))
				playerNameEntry.SetText("playerName: " + playerName)
				numOfEachTypePlayersEntry[0].SetText("Human player: " + strconv.Itoa(numOfEachTypePlayers[0]))
				numOfEachTypePlayersEntry[1].SetText("Random bot: " + strconv.Itoa(numOfEachTypePlayers[1]))
				numOfEachTypePlayersEntry[2].SetText("Probability bot: " + strconv.Itoa(numOfEachTypePlayers[2]))
				numOfEachTypePlayersEntry[3].SetText("Conventional bot: " + strconv.Itoa(numOfEachTypePlayers[3]))
				numOfEachTypePlayersEntry[4].SetText("Persona drawer bot: " + strconv.Itoa(numOfEachTypePlayers[4]))
				numOfEachTypePlayersEntry[5].SetText("Super talented bot: " + strconv.Itoa(numOfEachTypePlayers[5]))
				infoW.Close()
			} else if !CheckNumPlayers(formNumPlayers.Text) {
				dialog.ShowInformation("Error", "Something wrong with the number of players!", infoW)
			} else if !CheckInitialChips(formInitialChips.Text) {
				dialog.ShowInformation("Error", "Something wrong with the number of initial chips!", infoW)
			} else if !CheckNumTurns(formNumTurns.Text) {
				dialog.ShowInformation("Error", "Something wrong with the number of turns!", infoW)
			} else if !CheckPlayerName(formPlayerName.Text) {
				dialog.ShowInformation("Error", "Behave yourself, guys!", infoW)
			} else if !CheckNumType(formHumanPlayer.Text, formRandomBot.Text, formProbabilityBot.Text, formConventionalBot.Text, formPersonaDrawerBot.Text, formSuperTalentedBot.Text) {
				dialog.ShowInformation("Error", "Something wrong with the number of each type of players! \n (We can't create more than one human!)", infoW)
			}
		},
	}

	form.Append("number of players: ", formNumPlayers)
	form.Append("number of initial chips: ", formInitialChips)
	form.Append("number of turns: ", formNumTurns)
	form.Append("your name: ", formPlayerName)
	form.Append("number of human player: ", formHumanPlayer)
	form.Append("number of random bot: ", formRandomBot)
	form.Append("number of probability bot: ", formProbabilityBot)
	form.Append("number of conventional bot: ", formConventionalBot)
	form.Append("number of persona drawer bot: ", formPersonaDrawerBot)
	form.Append("number of super talented bot: ", formSuperTalentedBot)

	infoW.SetContent(form)
	infoW.Show()
}

/*
Author: Luuminous
Date: Nov. 4th, 2018
*/

/*
This part is to check if the setting is legal.
*/

func CheckNumPlayers(s string) bool {
	result, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	if result > 6 || result < 2 {
		return false
	}
	numPlayers = result
	return true
}

func CheckInitialChips(s string) bool {
	result, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	if result > 10000 || result < 200 || result % 10 != 0 {
		return false
	}
	initialChips = result
	return true
}

func CheckNumTurns(s string) bool {
	result, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	if result > 500 || result < 1 {
		return false
	}
	numTurns = result
	return true
}

func CheckPlayerName(s string) bool {
	if s == "Fuck" || s == "Motherfucker" {
		return false
	}
	playerName = s
	return true
}

func CheckNumType(humanS, randomS, probabilityS, conventionalS, personaS, superS string) bool {
	human, err1 := strconv.Atoi(humanS)
	random, err2 := strconv.Atoi(randomS)
	probability, err3 := strconv.Atoi(probabilityS)
	conventional, err4 := strconv.Atoi(conventionalS)
	persona, err5 := strconv.Atoi(personaS)
	super, err6 := strconv.Atoi(superS)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
		return false
	}
	if human + random + probability + conventional + persona + super != numPlayers {
		return false
	}
	if human != 1 && human != 0 {
		return false
	}
	if random < 0 || probability < 0 || conventional < 0 || persona < 0 || super < 0 {
		return false
	}

	numOfEachTypePlayers[0] = human
	numOfEachTypePlayers[1] = random
	numOfEachTypePlayers[2] = probability
	numOfEachTypePlayers[3] = conventional
	numOfEachTypePlayers[4] = persona
	numOfEachTypePlayers[5] = super

	return true
}

/*
Author: Luuminous
Date: Oct. 17th, 2018
*/

/*
The function is to convert the hands to a readable string.
*/

func HandsToString(hands []Card) string {
	ans := ""
	for i := 0; i < len(hands); i++ {
		num := hands[i].Num
		if num == 14 {
			ans += "A"
		} else if num == 13 {
			ans += "K"
		} else if num == 12 {
			ans += "Q"
		} else if num == 11 {
			ans += "J"
		} else {
			ans += strconv.Itoa(num)
		}
		color := hands[i].Color
		if color == 0 {
			ans += "♠"
		} else if color == 1 {
			ans += "♥"
		} else if color == 2 {
			ans += "♣"
		} else if color == 3 {
			ans += "♦"
		}
		ans += " "
	}
	return ans
}

/*
Author: Luuminous
Date: Nov. 22nd, 2018
*/

/*
The function is to set the information of preEventList into a scroll displayment.
*/

func SetInformationScroll(event []string, size int) string {
	eSize := len(event)
	text := ""
	if eSize < size {
		for i := 0; i < size - eSize; i++ {
			text += "\n"
		}
		for i := 0; i < eSize; i++ {
			text += event[i] + "\n"
		}
	} else {
		for i := eSize - size; i < eSize; i++ {
			text += event[i] + "\n"
		}
	}
	text = text[:len(text) - 1]
	return text
}

/*
Author: Luuminous
Date: Nov. 4th, 2018
*/

/*
Play the whole game.
*/

func PlayGame(app fyne.App) {
	current = InitialGame(numPlayers, initialChips)
	
	//To store all the bot index.
	botIndexList = make([]int, 0)

	//The main game window.
	gameW := app.NewWindow("Texas Holdem")

	//To set all the condition informations.
	conditionInformationEntry = make([]*widget.Label, 3)
	conditionInformationEntry[0] = &widget.Label{Text: "NumberOfTurn: " + strconv.Itoa(current.GameCount) + "\nStage: " + current.Stage + "\nCurrentBet: " + strconv.Itoa(current.CurrentBet) + "\nCommunityCard:                             "}
	conditionInformationEntry[1] = &widget.Label{Text: SetInformationScroll(current.PreEventsList, 4)}
	conditionInformationEntry[2] = &widget.Label{Text: "Welcome to Texas Holdem!"}

	//This channel is for the button communication.
	c := make(chan string) 

	//This entry is for the raise money.
	raiseEntry := widget.NewEntry()
	raiseEntry.SetText("0")

	helpMeButton := &widget.Button{Text: "Help me!", OnTapped: func() {
		if numOfEachTypePlayers[0] == 1 && len(humanActionList) != 0 {
			index := current.HumanPosition
			action := ProbabilityBot(index)
			c <- action
		}
	}}

	if numPlayers == 2 {

		/*
		This is the game for 2 players.
		*/
		for i := range current.Players {
			if i != current.HumanPosition {
				botIndexList = append(botIndexList, i)
			}
		}

		playerInformationEntry = make([][]*widget.Label, 2)

		//Set human player information entry:
		playerInformationEntry[0] = make([]*widget.Label, 4)
		//Chips:
		playerInformationEntry[0][0] = &widget.Label{Text: "Chips: " + strconv.Itoa(current.Players[current.HumanPosition].Chips)}
		//Bets:
		playerInformationEntry[0][1] = &widget.Label{Text: "Bet: " + strconv.Itoa(current.Players[current.HumanPosition].Bet)}
		//PositionName:
		playerInformationEntry[0][2] = &widget.Label{Text: "Position: " + current.Players[current.HumanPosition].PositionName}
		//Hands:
		playerInformationEntry[0][3] = &widget.Label{Text: "Hands: " + HandsToString(current.Players[current.HumanPosition].Hands)}

		//Set bot information entry:
		playerInformationEntry[1] = make([]*widget.Label, 4)
		//Chips:
		playerInformationEntry[1][0] = &widget.Label{Text: "Chips: " +strconv.Itoa(current.Players[botIndexList[0]].Chips)}
		//Bets:
		playerInformationEntry[1][1] = &widget.Label{Text: "Bet: " + strconv.Itoa(current.Players[botIndexList[0]].Bet)}
		//PositionName:
		playerInformationEntry[1][2] = &widget.Label{Text: "Position: " + current.Players[botIndexList[0]].PositionName}
		//Hands:
		playerInformationEntry[1][3] = &widget.Label{Text: "Hands: " + HandsToString(current.Players[botIndexList[0]].Hands)}

		gameW.SetContent(&widget.Box{Children: []fyne.CanvasObject{
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				&layout.Spacer{},
				widget.NewGroup(current.Players[botIndexList[0]].Name, []fyne.CanvasObject{
					playerInformationEntry[1][0],
					playerInformationEntry[1][1],
					playerInformationEntry[1][2],
					playerInformationEntry[1][3],
				}...),
				&layout.Spacer{},
			),
			//&widget.Label{Text: "\n"},
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				&layout.Spacer{},
				conditionInformationEntry[2],
				&layout.Spacer{},			
			),
			fyne.NewContainerWithLayout(layout.NewGridLayout(4),
				&layout.Spacer{},
				conditionInformationEntry[0],

				conditionInformationEntry[1],
				&layout.Spacer{},			
			),
			
			//&widget.Label{Text: "\n"},
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				&layout.Spacer{},
				widget.NewGroup(current.Players[current.HumanPosition].Name, []fyne.CanvasObject{
					playerInformationEntry[0][0],
					playerInformationEntry[0][1],
					playerInformationEntry[0][2],
					playerInformationEntry[0][3],
				}...),
				&layout.Spacer{},

			),
			fyne.NewContainerWithLayout(layout.NewGridLayout(5),
				&layout.Spacer{},
				&widget.Button{Text: "AllIn", OnTapped: func() {
					if InStr("AllIn", humanActionList) {
						c <- "AllIn"
					} else {
						dialog.ShowInformation("Error", "You can't all in now!", gameW)
					}
				}},
				&widget.Button{Text: "Call", OnTapped: func() {
					if InStr("Call", humanActionList) {
						c <- "Call"
					} else {
						dialog.ShowInformation("Error", "You can't call now!", gameW)
					}
				}},
				&widget.Button{Text: "Fold", OnTapped: func() {
					if InStr("Fold", humanActionList) {
						c <- "Fold"
					} else {
						dialog.ShowInformation("Error", "You can't fold now!", gameW)
					}
				}},
				&layout.Spacer{},

				&layout.Spacer{},
				&widget.Button{Text: "Check", OnTapped: func() {
					if InStr("Check", humanActionList) {
						c <- "Check"
					} else {
						dialog.ShowInformation("Error", "You can't check now!", gameW)
					}
				}},
				fyne.NewContainerWithLayout(layout.NewGridLayout(2),
					&widget.Button{Text: "Raise", OnTapped: func() {
					if InStr("Raise", humanActionList) {
						money, err := strconv.Atoi(raiseEntry.Text)
						lowerBound := current.CurrentBet - current.Players[current.HumanPosition].Bet + 10
						upperBound := current.Players[current.HumanPosition].Chips
						if err != nil || money > upperBound || money < lowerBound || money % 10 != 0 {
							dialog.ShowInformation("Error", "Wrong raise number!", gameW)
						} else {
							c <- "Raise" + raiseEntry.Text
							//fmt.Println(raiseEntry.Text)
						}
					} else {
						dialog.ShowInformation("Error", "You can't raise now!", gameW)
					}
				}},
				raiseEntry,
				),
				helpMeButton,
				&layout.Spacer{},

			),
			
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				&layout.Spacer{},
				&widget.Button{Text: "Close", OnTapped: func () {
					gameW.Close()
				}},
				&layout.Spacer{},
			),
		}})

	} else if numPlayers == 3 {

		/*
		This is the game for 3 players.
		*/
		for i := range current.Players {
			if i != current.HumanPosition {
				botIndexList = append(botIndexList, i)
			}
		}

		playerInformationEntry = make([][]*widget.Label, 3)

		//Set human player information entry:
		playerInformationEntry[0] = make([]*widget.Label, 4)
		//Chips:
		playerInformationEntry[0][0] = &widget.Label{Text: "Chips: " + strconv.Itoa(current.Players[current.HumanPosition].Chips)}
		//Bets:
		playerInformationEntry[0][1] = &widget.Label{Text: "Bet: " + strconv.Itoa(current.Players[current.HumanPosition].Bet)}
		//PositionName:
		playerInformationEntry[0][2] = &widget.Label{Text: "Position: " + current.Players[current.HumanPosition].PositionName}
		//Hands:
		playerInformationEntry[0][3] = &widget.Label{Text: "Hands: " + HandsToString(current.Players[current.HumanPosition].Hands)}

		//Set bot information entry:

		for i := range botIndexList {
			playerInformationEntry[i + 1] = make([]*widget.Label, 4)
			//Chips:
			playerInformationEntry[i + 1][0] = &widget.Label{Text: "Chips: " +strconv.Itoa(current.Players[botIndexList[i]].Chips)}
			//Bets:
			playerInformationEntry[i + 1][1] = &widget.Label{Text: "Bet: " + strconv.Itoa(current.Players[botIndexList[i]].Bet)}
			//PositionName:
			playerInformationEntry[i + 1][2] = &widget.Label{Text: "Position: " + current.Players[botIndexList[i]].PositionName}
			//Hands:
			playerInformationEntry[i + 1][3] = &widget.Label{Text: "Hands: " + HandsToString(current.Players[botIndexList[i]].Hands)}
		}

		gameW.SetContent(&widget.Box{Children: []fyne.CanvasObject{
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				widget.NewGroup(current.Players[botIndexList[0]].Name, []fyne.CanvasObject{
					playerInformationEntry[1][0],
					playerInformationEntry[1][1],
					playerInformationEntry[1][2],
					playerInformationEntry[1][3],
				}...),
				&layout.Spacer{},
				widget.NewGroup(current.Players[botIndexList[1]].Name, []fyne.CanvasObject{
					playerInformationEntry[2][0],
					playerInformationEntry[2][1],
					playerInformationEntry[2][2],
					playerInformationEntry[2][3],
				}...),
			),
			//&widget.Label{Text: "\n"},
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				&layout.Spacer{},
				conditionInformationEntry[2],
				&layout.Spacer{},			
			),
			fyne.NewContainerWithLayout(layout.NewGridLayout(4),
				&layout.Spacer{},
				conditionInformationEntry[0],

				conditionInformationEntry[1],
				&layout.Spacer{},			
			),
			
			//&widget.Label{Text: "\n"},
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				&layout.Spacer{},
				widget.NewGroup(current.Players[current.HumanPosition].Name, []fyne.CanvasObject{
					playerInformationEntry[0][0],
					playerInformationEntry[0][1],
					playerInformationEntry[0][2],
					playerInformationEntry[0][3],
				}...),
				&layout.Spacer{},

			),
			fyne.NewContainerWithLayout(layout.NewGridLayout(5),
				&layout.Spacer{},
				&widget.Button{Text: "AllIn", OnTapped: func() {
					if InStr("AllIn", humanActionList) {
						c <- "AllIn"
					} else {
						dialog.ShowInformation("Error", "You can't all in now!", gameW)
					}
				}},
				&widget.Button{Text: "Call", OnTapped: func() {
					if InStr("Call", humanActionList) {
						c <- "Call"
					} else {
						dialog.ShowInformation("Error", "You can't call now!", gameW)
					}
				}},
				&widget.Button{Text: "Fold", OnTapped: func() {
					if InStr("Fold", humanActionList) {
						c <- "Fold"
					} else {
						dialog.ShowInformation("Error", "You can't fold now!", gameW)
					}
				}},
				&layout.Spacer{},

				&layout.Spacer{},
				&widget.Button{Text: "Check", OnTapped: func() {
					if InStr("Check", humanActionList) {
						c <- "Check"
					} else {
						dialog.ShowInformation("Error", "You can't check now!", gameW)
					}
				}},
				fyne.NewContainerWithLayout(layout.NewGridLayout(2),
					&widget.Button{Text: "Raise", OnTapped: func() {
					if InStr("Raise", humanActionList) {
						money, err := strconv.Atoi(raiseEntry.Text)
						lowerBound := current.CurrentBet - current.Players[current.HumanPosition].Bet + 10
						upperBound := current.Players[current.HumanPosition].Chips
						if err != nil || money > upperBound || money < lowerBound || money % 10 != 0 {
							dialog.ShowInformation("Error", "Wrong raise number!", gameW)
						} else {
							c <- "Raise" + raiseEntry.Text
							//fmt.Println(raiseEntry.Text)
						}
					} else {
						dialog.ShowInformation("Error", "You can't raise now!", gameW)
					}
				}},
				raiseEntry,
				),
				helpMeButton,
				&layout.Spacer{},

			),
			
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				&layout.Spacer{},
				&widget.Button{Text: "Close", OnTapped: func () {
					gameW.Close()
				}},
				&layout.Spacer{},
			),
		}})

	} else if numPlayers == 4 {

		/*
		This is the game for 4 players.
		*/
		for i := range current.Players {
			if i != current.HumanPosition {
				botIndexList = append(botIndexList, i)
			}
		}

		playerInformationEntry = make([][]*widget.Label, 4)

		//Set human player information entry:
		playerInformationEntry[0] = make([]*widget.Label, 4)
		//Chips:
		playerInformationEntry[0][0] = &widget.Label{Text: "Chips: " + strconv.Itoa(current.Players[current.HumanPosition].Chips)}
		//Bets:
		playerInformationEntry[0][1] = &widget.Label{Text: "Bet: " + strconv.Itoa(current.Players[current.HumanPosition].Bet)}
		//PositionName:
		playerInformationEntry[0][2] = &widget.Label{Text: "Position: " + current.Players[current.HumanPosition].PositionName}
		//Hands:
		playerInformationEntry[0][3] = &widget.Label{Text: "Hands: " + HandsToString(current.Players[current.HumanPosition].Hands)}

		//Set bot information entry:

		for i := range botIndexList {
			playerInformationEntry[i + 1] = make([]*widget.Label, 4)
			//Chips:
			playerInformationEntry[i + 1][0] = &widget.Label{Text: "Chips: " +strconv.Itoa(current.Players[botIndexList[i]].Chips)}
			//Bets:
			playerInformationEntry[i + 1][1] = &widget.Label{Text: "Bet: " + strconv.Itoa(current.Players[botIndexList[i]].Bet)}
			//PositionName:
			playerInformationEntry[i + 1][2] = &widget.Label{Text: "Position: " + current.Players[botIndexList[i]].PositionName}
			//Hands:
			playerInformationEntry[i + 1][3] = &widget.Label{Text: "Hands: " + HandsToString(current.Players[botIndexList[i]].Hands)}
		}

		gameW.SetContent(&widget.Box{Children: []fyne.CanvasObject{
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				widget.NewGroup(current.Players[botIndexList[0]].Name, []fyne.CanvasObject{
					playerInformationEntry[1][0],
					playerInformationEntry[1][1],
					playerInformationEntry[1][2],
					playerInformationEntry[1][3],
				}...),

				widget.NewGroup(current.Players[botIndexList[1]].Name, []fyne.CanvasObject{
					playerInformationEntry[2][0],
					playerInformationEntry[2][1],
					playerInformationEntry[2][2],
					playerInformationEntry[2][3],
				}...),

				widget.NewGroup(current.Players[botIndexList[2]].Name, []fyne.CanvasObject{
					playerInformationEntry[3][0],
					playerInformationEntry[3][1],
					playerInformationEntry[3][2],
					playerInformationEntry[3][3],
				}...),
			),
			//&widget.Label{Text: "\n"},
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				&layout.Spacer{},
				conditionInformationEntry[2],
				&layout.Spacer{},			
			),
			fyne.NewContainerWithLayout(layout.NewGridLayout(4),
				&layout.Spacer{},
				conditionInformationEntry[0],

				conditionInformationEntry[1],
				&layout.Spacer{},			
			),
			
			//&widget.Label{Text: "\n"},
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				&layout.Spacer{},
				widget.NewGroup(current.Players[current.HumanPosition].Name, []fyne.CanvasObject{
					playerInformationEntry[0][0],
					playerInformationEntry[0][1],
					playerInformationEntry[0][2],
					playerInformationEntry[0][3],
				}...),
				&layout.Spacer{},

			),
			fyne.NewContainerWithLayout(layout.NewGridLayout(5),
				&layout.Spacer{},
				&widget.Button{Text: "AllIn", OnTapped: func() {
					if InStr("AllIn", humanActionList) {
						c <- "AllIn"
					} else {
						dialog.ShowInformation("Error", "You can't all in now!", gameW)
					}
				}},
				&widget.Button{Text: "Call", OnTapped: func() {
					if InStr("Call", humanActionList) {
						c <- "Call"
					} else {
						dialog.ShowInformation("Error", "You can't call now!", gameW)
					}
				}},
				&widget.Button{Text: "Fold", OnTapped: func() {
					if InStr("Fold", humanActionList) {
						c <- "Fold"
					} else {
						dialog.ShowInformation("Error", "You can't fold now!", gameW)
					}
				}},
				&layout.Spacer{},

				&layout.Spacer{},
				&widget.Button{Text: "Check", OnTapped: func() {
					if InStr("Check", humanActionList) {
						c <- "Check"
					} else {
						dialog.ShowInformation("Error", "You can't check now!", gameW)
					}
				}},
				fyne.NewContainerWithLayout(layout.NewGridLayout(2),
					&widget.Button{Text: "Raise", OnTapped: func() {
					if InStr("Raise", humanActionList) {
						money, err := strconv.Atoi(raiseEntry.Text)
						lowerBound := current.CurrentBet - current.Players[current.HumanPosition].Bet + 10
						upperBound := current.Players[current.HumanPosition].Chips
						if err != nil || money > upperBound || money < lowerBound || money % 10 != 0 {
							dialog.ShowInformation("Error", "Wrong raise number!", gameW)
						} else {
							c <- "Raise" + raiseEntry.Text
							//fmt.Println(raiseEntry.Text)
						}
					} else {
						dialog.ShowInformation("Error", "You can't raise now!", gameW)
					}
				}},
				raiseEntry,
				),
				helpMeButton,
				&layout.Spacer{},

			),
			
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				&layout.Spacer{},
				&widget.Button{Text: "Close", OnTapped: func () {
					gameW.Close()
				}},
				&layout.Spacer{},
			),
		}})

	} else if numPlayers == 5 {

		/*
		This is the game for 5 players.
		*/
		for i := range current.Players {
			if i != current.HumanPosition {
				botIndexList = append(botIndexList, i)
			}
		}

		playerInformationEntry = make([][]*widget.Label, 5)

		//Set human player information entry:
		playerInformationEntry[0] = make([]*widget.Label, 4)
		//Chips:
		playerInformationEntry[0][0] = &widget.Label{Text: "Chips: " + strconv.Itoa(current.Players[current.HumanPosition].Chips)}
		//Bets:
		playerInformationEntry[0][1] = &widget.Label{Text: "Bet: " + strconv.Itoa(current.Players[current.HumanPosition].Bet)}
		//PositionName:
		playerInformationEntry[0][2] = &widget.Label{Text: "Position: " + current.Players[current.HumanPosition].PositionName}
		//Hands:
		playerInformationEntry[0][3] = &widget.Label{Text: "Hands: " + HandsToString(current.Players[current.HumanPosition].Hands)}

		//Set bot information entry:

		for i := range botIndexList {
			playerInformationEntry[i + 1] = make([]*widget.Label, 4)
			//Chips:
			playerInformationEntry[i + 1][0] = &widget.Label{Text: "Chips: " +strconv.Itoa(current.Players[botIndexList[i]].Chips)}
			//Bets:
			playerInformationEntry[i + 1][1] = &widget.Label{Text: "Bet: " + strconv.Itoa(current.Players[botIndexList[i]].Bet)}
			//PositionName:
			playerInformationEntry[i + 1][2] = &widget.Label{Text: "Position: " + current.Players[botIndexList[i]].PositionName}
			//Hands:
			playerInformationEntry[i + 1][3] = &widget.Label{Text: "Hands: " + HandsToString(current.Players[botIndexList[i]].Hands)}
		}

		gameW.SetContent(&widget.Box{Children: []fyne.CanvasObject{
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				widget.NewGroup(current.Players[botIndexList[1]].Name, []fyne.CanvasObject{
					playerInformationEntry[2][0],
					playerInformationEntry[2][1],
					playerInformationEntry[2][2],
					playerInformationEntry[2][3],
				}...),
				&layout.Spacer{},
				widget.NewGroup(current.Players[botIndexList[2]].Name, []fyne.CanvasObject{
					playerInformationEntry[3][0],
					playerInformationEntry[3][1],
					playerInformationEntry[3][2],
					playerInformationEntry[3][3],
				}...),
			),
			//&widget.Label{Text: "\n"},
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				&layout.Spacer{},
				conditionInformationEntry[2],
				&layout.Spacer{},			
			),
			fyne.NewContainerWithLayout(layout.NewGridLayout(4),
				&layout.Spacer{},
				conditionInformationEntry[0],

				conditionInformationEntry[1],
				&layout.Spacer{},			
			),
			
			//&widget.Label{Text: "\n"},
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				widget.NewGroup(current.Players[botIndexList[0]].Name, []fyne.CanvasObject{
					playerInformationEntry[1][0],
					playerInformationEntry[1][1],
					playerInformationEntry[1][2],
					playerInformationEntry[1][3],
				}...),
				widget.NewGroup(current.Players[current.HumanPosition].Name, []fyne.CanvasObject{
					playerInformationEntry[0][0],
					playerInformationEntry[0][1],
					playerInformationEntry[0][2],
					playerInformationEntry[0][3],
				}...),
				widget.NewGroup(current.Players[botIndexList[3]].Name, []fyne.CanvasObject{
					playerInformationEntry[4][0],
					playerInformationEntry[4][1],
					playerInformationEntry[4][2],
					playerInformationEntry[4][3],
				}...),

			),
			fyne.NewContainerWithLayout(layout.NewGridLayout(5),
				&layout.Spacer{},
				&widget.Button{Text: "AllIn", OnTapped: func() {
					if InStr("AllIn", humanActionList) {
						c <- "AllIn"
					} else {
						dialog.ShowInformation("Error", "You can't all in now!", gameW)
					}
				}},
				&widget.Button{Text: "Call", OnTapped: func() {
					if InStr("Call", humanActionList) {
						c <- "Call"
					} else {
						dialog.ShowInformation("Error", "You can't call now!", gameW)
					}
				}},
				&widget.Button{Text: "Fold", OnTapped: func() {
					if InStr("Fold", humanActionList) {
						c <- "Fold"
					} else {
						dialog.ShowInformation("Error", "You can't fold now!", gameW)
					}
				}},
				&layout.Spacer{},

				&layout.Spacer{},
				&widget.Button{Text: "Check", OnTapped: func() {
					if InStr("Check", humanActionList) {
						c <- "Check"
					} else {
						dialog.ShowInformation("Error", "You can't check now!", gameW)
					}
				}},
				fyne.NewContainerWithLayout(layout.NewGridLayout(2),
					&widget.Button{Text: "Raise", OnTapped: func() {
					if InStr("Raise", humanActionList) {
						money, err := strconv.Atoi(raiseEntry.Text)
						lowerBound := current.CurrentBet - current.Players[current.HumanPosition].Bet + 10
						upperBound := current.Players[current.HumanPosition].Chips
						if err != nil || money > upperBound || money < lowerBound || money % 10 != 0 {
							dialog.ShowInformation("Error", "Wrong raise number!", gameW)
						} else {
							c <- "Raise" + raiseEntry.Text
							//fmt.Println(raiseEntry.Text)
						}
					} else {
						dialog.ShowInformation("Error", "You can't raise now!", gameW)
					}
				}},
				raiseEntry,
				),
				helpMeButton,
				&layout.Spacer{},

			),
			
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				&layout.Spacer{},
				&widget.Button{Text: "Close", OnTapped: func () {
					gameW.Close()
				}},
				&layout.Spacer{},
			),
		}})

	} else if numPlayers == 6 {

		/*
		This is the game for 6 players.
		*/
		for i := range current.Players {
			if i != current.HumanPosition {
				botIndexList = append(botIndexList, i)
			}
		}

		playerInformationEntry = make([][]*widget.Label, 6)

		//Set human player information entry:
		playerInformationEntry[0] = make([]*widget.Label, 4)
		//Chips:
		playerInformationEntry[0][0] = &widget.Label{Text: "Chips: " + strconv.Itoa(current.Players[current.HumanPosition].Chips)}
		//Bets:
		playerInformationEntry[0][1] = &widget.Label{Text: "Bet: " + strconv.Itoa(current.Players[current.HumanPosition].Bet)}
		//PositionName:
		playerInformationEntry[0][2] = &widget.Label{Text: "Position: " + current.Players[current.HumanPosition].PositionName}
		//Hands:
		playerInformationEntry[0][3] = &widget.Label{Text: "Hands: " + HandsToString(current.Players[current.HumanPosition].Hands)}

		//Set bot information entry:

		for i := range botIndexList {
			playerInformationEntry[i + 1] = make([]*widget.Label, 4)
			//Chips:
			playerInformationEntry[i + 1][0] = &widget.Label{Text: "Chips: " +strconv.Itoa(current.Players[botIndexList[i]].Chips)}
			//Bets:
			playerInformationEntry[i + 1][1] = &widget.Label{Text: "Bet: " + strconv.Itoa(current.Players[botIndexList[i]].Bet)}
			//PositionName:
			playerInformationEntry[i + 1][2] = &widget.Label{Text: "Position: " + current.Players[botIndexList[i]].PositionName}
			//Hands:
			playerInformationEntry[i + 1][3] = &widget.Label{Text: "Hands: " + HandsToString(current.Players[botIndexList[i]].Hands)}
		}

		gameW.SetContent(&widget.Box{Children: []fyne.CanvasObject{
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				widget.NewGroup(current.Players[botIndexList[1]].Name, []fyne.CanvasObject{
					playerInformationEntry[2][0],
					playerInformationEntry[2][1],
					playerInformationEntry[2][2],
					playerInformationEntry[2][3],
				}...),
				widget.NewGroup(current.Players[botIndexList[2]].Name, []fyne.CanvasObject{
					playerInformationEntry[3][0],
					playerInformationEntry[3][1],
					playerInformationEntry[3][2],
					playerInformationEntry[3][3],
				}...),
				widget.NewGroup(current.Players[botIndexList[3]].Name, []fyne.CanvasObject{
					playerInformationEntry[4][0],
					playerInformationEntry[4][1],
					playerInformationEntry[4][2],
					playerInformationEntry[4][3],
				}...),
			),
			//&widget.Label{Text: "\n"},
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				&layout.Spacer{},
				conditionInformationEntry[2],
				&layout.Spacer{},			
			),
			fyne.NewContainerWithLayout(layout.NewGridLayout(4),
				&layout.Spacer{},
				conditionInformationEntry[0],

				conditionInformationEntry[1],
				&layout.Spacer{},			
			),
			
			//&widget.Label{Text: "\n"},
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				widget.NewGroup(current.Players[botIndexList[0]].Name, []fyne.CanvasObject{
					playerInformationEntry[1][0],
					playerInformationEntry[1][1],
					playerInformationEntry[1][2],
					playerInformationEntry[1][3],
				}...),
				widget.NewGroup(current.Players[current.HumanPosition].Name, []fyne.CanvasObject{
					playerInformationEntry[0][0],
					playerInformationEntry[0][1],
					playerInformationEntry[0][2],
					playerInformationEntry[0][3],
				}...),
				widget.NewGroup(current.Players[botIndexList[4]].Name, []fyne.CanvasObject{
					playerInformationEntry[5][0],
					playerInformationEntry[5][1],
					playerInformationEntry[5][2],
					playerInformationEntry[5][3],
				}...),

			),
			fyne.NewContainerWithLayout(layout.NewGridLayout(5),
				&layout.Spacer{},
				&widget.Button{Text: "AllIn", OnTapped: func() {
					if InStr("AllIn", humanActionList) {
						c <- "AllIn"
					} else {
						dialog.ShowInformation("Error", "You can't all in now!", gameW)
					}
				}},
				&widget.Button{Text: "Call", OnTapped: func() {
					if InStr("Call", humanActionList) {
						c <- "Call"
					} else {
						dialog.ShowInformation("Error", "You can't call now!", gameW)
					}
				}},
				&widget.Button{Text: "Fold", OnTapped: func() {
					if InStr("Fold", humanActionList) {
						c <- "Fold"
					} else {
						dialog.ShowInformation("Error", "You can't fold now!", gameW)
					}
				}},
				&layout.Spacer{},

				&layout.Spacer{},
				&widget.Button{Text: "Check", OnTapped: func() {
					if InStr("Check", humanActionList) {
						c <- "Check"
					} else {
						dialog.ShowInformation("Error", "You can't check now!", gameW)
					}
				}},
				fyne.NewContainerWithLayout(layout.NewGridLayout(2),
					&widget.Button{Text: "Raise", OnTapped: func() {
					if InStr("Raise", humanActionList) {
						money, err := strconv.Atoi(raiseEntry.Text)
						lowerBound := current.CurrentBet - current.Players[current.HumanPosition].Bet + 10
						upperBound := current.Players[current.HumanPosition].Chips
						if err != nil || money > upperBound || money < lowerBound || money % 10 != 0 {
							dialog.ShowInformation("Error", "Wrong raise number!", gameW)
						} else {
							c <- "Raise" + raiseEntry.Text
							//fmt.Println(raiseEntry.Text)
						}
					} else {
						dialog.ShowInformation("Error", "You can't raise now!", gameW)
					}
				}},
				raiseEntry,
				),
				helpMeButton,
				&layout.Spacer{},

			),
			
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				&layout.Spacer{},
				&widget.Button{Text: "Close", OnTapped: func () {
					gameW.Close()
				}},
				&layout.Spacer{},
			),
		}})

	} else {
		gameW.SetContent(&widget.Box{Children: []fyne.CanvasObject{
			fyne.NewContainerWithLayout(layout.NewGridLayout(3), 
				&layout.Spacer{},
				numPlayersEntry,
				&layout.Spacer{},

				&layout.Spacer{},
				initialChipsEntry,
				&layout.Spacer{},

				&layout.Spacer{},
				numTurnsEntry,
				&layout.Spacer{},

				&layout.Spacer{},
				playerNameEntry,
				&layout.Spacer{},
			),
			&widget.Button{Text: "Close", OnTapped: func () {
				gameW.Close()
			}},
		}})

	}
	gameW.Show()

	go BackGrounder(c, app)

}

/*
Author: Luuminous
Date: Nov. 25th, 2018
*/

/*
This function is for run the program in the back ground.
*/

func BackGrounder(c chan string, app fyne.App) {

	//GameCount start from 1, so use <= rather than <.

	//Initiate for personaBot
	InitialAlphaMap()
	for current.GameCount <= numTurns && CheckGame() {
		time.Sleep(3 * 1000000000)
		StartOneGame(c)
	}

	result := ""
	for i := range current.Players {
		result += current.Players[i].Name + ": " + strconv.Itoa(current.Players[i].Chips) + "\n"
	}

	resultW := app.NewWindow("Result")
	resultW.SetContent(&widget.Box{Children: []fyne.CanvasObject{
		fyne.NewContainerWithLayout(layout.NewGridLayout(3),
			&layout.Spacer{},
			&widget.Label{Text: result},
			&layout.Spacer{},
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(3),
			&layout.Spacer{},
			&widget.Button{Text: "Close", OnTapped: func() {
				resultW.Close()
			}},
			&layout.Spacer{},
		),
	}})

	resultW.Show()


	//Hold the process until the user taps the close button.
	for true {
		time.Sleep(10 * 1000000000)
	}
}

/*
Author: Luuminous
Date: Oct. 17th, 2018
*/

/*
The function is to check if the game can continue. False means game end.
*/

func CheckGame() bool {
	numPlayers := 0 // To record how many players still alive.
	for _, player := range current.Players {
		if !player.Eliminated {
			numPlayers++
		}
	}
	if numPlayers <= 1 {
		return false
	}
	return true
}

/*
Author: Luuminous
Date: Oct. 17th, 2018
*/

/*
This function is to initiate everything after one game.
Also, set everthing to the game window.
*/

func InitialOneGame(){
	numPlayers := 0

	//Clear the information.
	conditionInformationEntry[2].SetText("")

	for i := range current.Players {
		if !current.Players[i].Eliminated {
			numPlayers++
		}
	}

	/*
	The loop is to reset several values from each uneliminated player and determine the start player.
	*/

	positionFinder := 2
	check := false
	countCheck := 0
	for (check == false) {
		countCheck++
		for index, player := range current.Players{
			if positionFinder == player.GamePosition {
				if player.Eliminated {
					current.Players[index].GamePosition = 0 
					positionFinder++
				} else {
					current.StartPlayer = player.SeatPosition
					check = true
				}
			}
			if player.Eliminated == false {
				current.Players[index].OK = false
				current.Players[index].Fold = false
				current.Players[index].AllIn = false
				current.Players[index].Bet = 0
			}
		}
		if countCheck > 6 {
			fmt.Println("-----------------------------------------------------------------")
			for i := range current.Players {
				fmt.Println(current.Players[i].Eliminated, current.Players[i].GamePosition)
			}
			panic("Death Loop!")
		}
	}

	for i := range current.Players {
		if current.Players[i].Eliminated {
			current.Players[i].GamePosition = 0
		}
	}

	if numPlayers < 2 {
		panic("There aren't enough players! What's wrong with CheckGame()?")
	}
	/*
	Generate random sequence for the card.
	 */
	cardPool := Initiation()
	current.Pool = CreateRandomCards(2 * numPlayers + 5, cardPool)
	current.ChipPool = 0
	current.Stage = "Pre-flop"
	var temp []Card
	current.CommunityCard = temp
	var tempList []string
	current.PreEventsList = tempList
	/*
	Give hands to each player.
	*/
	count := 0
	for index := range current.Players {
		if current.Players[index].Eliminated == false {
			var temp []Card
			temp = append(temp, current.Pool[count])
			count++
			temp = append(temp, current.Pool[count])
			count++
			current.Players[index].Hands = temp
		}
	}
	/*
	The loop is to change the position of each uneliminated player.
	*/
	gamePosition := 1
	for (gamePosition <= numPlayers) {
		for index, player := range current.Players {
			if (gamePosition == 1) && (player.SeatPosition == current.StartPlayer) {
				current.Players[index].GamePosition = gamePosition
				if current.Players[index].Chips <= 10 {
					AllIn(index)
				} else {
					Raise(index, 10)
				}
				gamePosition++
				if gamePosition > numPlayers {
					break
				}
				continue
			}
			if (gamePosition > 1) && (player.Eliminated == false) {
				current.Players[index].GamePosition = gamePosition
				if gamePosition == 2 {
					if current.Players[index].Chips <= 20 {
						AllIn(index)
					} else {
						Raise(index, 20)
					}
					//Though Big Blind should raise 20, it still shouldn't be OK.
					current.Players[index].OK = false
				}
				gamePosition++
				if gamePosition > numPlayers {
					break
				}
				continue
			}
		}
	}
	InitialPlayerName()
	Refresh()

	//Empty all the bots hands.
	for i := range botIndexList {
		playerInformationEntry[botIndexList[i]][3].SetText("Hands: ")
	}
	//Initiate for personaBot
	personaMap = InitialMap()
}

/* this part is the code for five actions for Texas Holdem GameBoard
including AllIn()
@Yichen Mo 0ct.16.2018
*/

/*AllIn takes the player's all chips all his current chips to the pool
 */

func AllIn(index int) {
	current.Players[index].AllIn = true
	current.ChipPool = current.ChipPool + current.Players[index].Chips
	current.Players[index].Bet = current.Players[index].Bet + current.Players[index].Chips
	current.Players[index].Chips = 0
	if current.Players[index].Bet > current.CurrentBet {
		current.CurrentBet = current.Players[index].Bet
		for i := range current.Players {
			if i != index && !current.Players[i].Eliminated && !current.Players[i].Fold {
				current.Players[i].OK = false
			}
		}

	}
	current.PreEventsList = append(current.PreEventsList, (current.Players[index].Name + " chose to all in"))
	if current.CurrentBet < current.Players[index].Bet {
		current.CurrentBet = current.Players[index].Bet
		for i := range current.Players {
			current.Players[i].OK = false
		}
	}
	current.Players[index].OK = true

}

/*The player would raise the bet
 */

func Raise(index, numBet int) {
	for i := range current.Players {
		current.Players[i].OK = false
	}
	current.Players[index].Bet = numBet + current.Players[index].Bet
	current.CurrentBet = current.Players[index].Bet
	current.ChipPool = current.ChipPool + numBet
	current.Players[index].Chips = current.Players[index].Chips - numBet
	current.Players[index].OK = true

	stringInput := current.Players[index].Name + " chose to raise "
	stringInput = stringInput + strconv.Itoa(numBet)
	current.PreEventsList = append(current.PreEventsList, stringInput)

}

/*Check is the action that take the Player and update the attributes in the Player also update the Current
 */

func Check(index int) {
	current.Players[index].OK = true
	current.PreEventsList = append(current.PreEventsList, (current.Players[index].Name + " chose to check"))
}

/*Call take the input of the player and current, for this action the player would match the bet of current board
 */

func Call(index int) {

	raisedAmount := current.CurrentBet - current.Players[index].Bet
	current.Players[index].Bet = current.CurrentBet
	current.Players[index].Chips = current.Players[index].Chips - raisedAmount
	current.ChipPool = current.ChipPool + raisedAmount
	current.Players[index].OK = true

	current.PreEventsList = append(current.PreEventsList, (current.Players[index].Name + " chose to call"))
}

/*Fold is the action that player will hand out the holecards and will not participate the current board
 */

func Fold(index int) {
	current.Players[index].Fold = true
	current.PreEventsList = append(current.PreEventsList, (current.Players[index].Name + " chose to fold"))
	numActivePlayers := 0
	for _, player := range current.Players {
		if !player.Eliminated && !player.Fold {
			numActivePlayers++
		}
	}
	if numActivePlayers == 1 {
		for i, player := range current.Players {
			if !player.Eliminated && !player.Fold {
				current.Players[i].OK = true
			}
		}
	}
}

/*
Author: Luuminous
Date: Oct. 17th, 2018
*/

/*
This function is to initiate the name according to each player's position.
*/

/*
position 1: Small Blind
position 2: Big Blind
position 3: Under the Gun
position 4: Middle Position
.
.
.
position -2: Cut Off
position -1: Button
*/
		
func InitialPlayerName() {
	numPlayers := 0
	for _, player := range current.Players {
		if player.Eliminated == false {
			numPlayers++
		}
	}
	nameList := make([]string, numPlayers)
	nameList[0] = "Small Blind"
	nameList[1] = "Big Blind"
	if numPlayers == 3 {
		nameList[2] = "Button"
	} else if numPlayers == 4 {
		nameList[2] = "Under the Gun"
		nameList[3] = "Button"
	} else if numPlayers == 5 {
		nameList[2] = "Under the Gun"
		nameList[3] = "Cut Off"
		nameList[4] = "Button"
	} else if numPlayers >= 6{
		nameList[2] = "Under the Gun"
		for i := 3; i < numPlayers - 2; i++ {
			nameList[i] = "Middle Position"
		}
		nameList[numPlayers - 2] = "Cut Off"
		nameList[numPlayers - 1] = "Button"
	}
	for index, player := range current.Players {
		if current.Players[index].Eliminated == false {
			current.Players[index].PositionName = nameList[player.GamePosition - 1]
		}
	}
}

/*
Author: Luuminous
Date: Oct. 11th, 2018
*/

/*
This function is to create a new 52 cards in a slice.
*/

func Initiation() []Card {
	var ans []Card
	for i := 2; i <= 14; i++ {
		for j := 0; j <= 3; j++ {
			var temp Card
			temp.Num = i
			temp.Color = j
			ans = append(ans, temp)
		}
	}
	return ans
}

/*
  Author: Chengyang Nie
  Date: 10/11/2018
*/

/*
This function is used to create random card from the pool
The input is the number of cards and a slice of Card pool
The output is a slice of cards.

use rand package shuffle function, then choose the first n elements.

*/

func CreateRandomCards(num int, cardPool []Card) []Card {
	if num < 0 {
		panic("Generate negative number of random cards")
	}
	if num > 52 {
		panic("Generate more than 52 number of random cards")
	}
	if cardPool == nil {
		panic("The input cardPool is empty")
	}
	rand.Shuffle(len(cardPool), func(i, j int) {
		cardPool[i], cardPool[j] = cardPool[j], cardPool[i]
	})
	var result []Card
	for i := 0; i < num; i++ {
		result = append(result, cardPool[i])
	}
	return result
}

/* this is the code for the command for current round of Texas Holdem GameBoard

@Yichen 0ct.16.2018
Luuminous Nov. 25th, 2018

*/

func Command(startPosition int, c chan string) {
	numActivePlayers := 0
	for _, player := range current.Players {
		if !player.Eliminated {
			numActivePlayers++
		}
	}
	currentPos := startPosition

	for currentPos <= numActivePlayers {

		for index, _ := range current.Players {

			if current.Players[index].GamePosition == currentPos {

				if current.Players[index].Eliminated {
					continue
				} else if current.Players[index].Fold {
					currentPos++
					continue
				} else if !current.Players[index].AllIn && !current.Players[index].OK {

					/*
					Now is the action player.
					*/

					if current.Players[index].PlayerType == "Random Bot" {

						action := RandomBot(index)
						Record(index, action)
						if action[0] == 'R' {
							mon := action[5:]
							intMon, _ := strconv.Atoi(mon)
							Raise(index, intMon)
							
							time.Sleep(2 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Check" {
							Check(index)

							time.Sleep(2 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Call" {
							Call(index)

							time.Sleep(2 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "AllIn" {
							AllIn(index)

							time.Sleep(2 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Fold" {
							Fold(index)

							time.Sleep(2 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						}
					}

					if current.Players[index].PlayerType == "Human Player" {

						conditionInformationEntry[2].SetText("Now it's your turn!")
						//print the previous players' actions
						humanActionList = ChooseActions(index)

						/*
						Here need to response human's action.
						*/
						action := <- c
						Record(index, action)
						
						humanActionList = make([]string, 0)

						conditionInformationEntry[2].SetText("")

						if action[0] == 'R' {
							mon := action[5:]
							intMon, _ := strconv.Atoi(mon)
							//fmt.Println(mon, intMon)
							Raise(index, intMon)

							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Check" {
							Check(index)

							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Call" {
							Call(index)

							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "AllIn" {
							AllIn(index)

							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Fold" {
							Fold(index)

							Refresh()
							time.Sleep(2 * 1000000000)

						}

					}

					if current.Players[index].PlayerType == "Probability Bot" {

						action := ProbabilityBot(index)
						Record(index, action)
						if action[0] == 'R' {
							mon := action[5:]
							intMon, _ := strconv.Atoi(mon)
							Raise(index, intMon)
							
							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Check" {
							Check(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Call" {
							Call(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "AllIn" {
							AllIn(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Fold" {
							Fold(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						}
					}

					if current.Players[index].PlayerType == "Conventional Bot" {

						action := ConventionalBot(index)
						Record(index, action)
						
						if action[0] == 'R' {
							mon := action[5:]
							intMon, _ := strconv.Atoi(mon)
							Raise(index, intMon)
							
							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Check" {
							Check(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Call" {
							Call(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "AllIn" {
							AllIn(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Fold" {
							Fold(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						}
					}

					if current.Players[index].PlayerType == "Persona Drawer Bot" {

						action := PersonaBot(index)
						Record(index, action)
						
						if action[0] == 'R' {
							mon := action[5:]
							intMon, _ := strconv.Atoi(mon)
							Raise(index, intMon)
							
							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Check" {
							Check(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Call" {
							Call(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "AllIn" {
							AllIn(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Fold" {
							Fold(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						}
					}

					if current.Players[index].PlayerType == "Super Talented Bot" {

						action := MinimumRegretBot(index)
						Record(index, action)
						
						if action[0] == 'R' {
							mon := action[5:]
							intMon, _ := strconv.Atoi(mon)
							Raise(index, intMon)
							
							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Check" {
							Check(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Call" {
							Call(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "AllIn" {
							AllIn(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						} else if action == "Fold" {
							Fold(index)

							time.Sleep(1.5 * 1000000000)
							Refresh()
							time.Sleep(2 * 1000000000)

						}
					}

					currentPos++

				} else {
					currentPos++
				}
			}
		}
	}
}

/*
Author: Luuminous
Date: Nov. 25th, 2018
*/

//This function is to check the string is in the slice or not.

func InStr(s string, list []string) bool {
	for _, value := range list {
		if value == s {
			return true
		}
	}
	return false
}

/*
Author: Luuminous
Date: Nov. 25th, 2018
*/

//This function is to check the integer is in the slice or not.

func InInt(i int, list []int) bool {
	for _, value := range list {
		if value == i {
			return true
		}
	}
	return false
}

/*
Author: Luuminous
Date: Nov. 28th, 2018
*/

//This function is set a integer to integer tens digit.

func SetToTen(i int) int {
	remainder := i % 10
	if remainder < 5 {
		return i / 10 * 10
	} else {
		return (i / 10 + 1) * 10
	}
}

/*
Author: Yichen Mo
Date: Oct. 31st, 2018
 */

func ChooseActions(index int) []string {
	actionList := []string{}
	raiseMoney := current.CurrentBet - current.Players[index].Bet
	if raiseMoney == 0 {
		//check raise
		actionList = append(actionList, "AllIn", "Raise", "Check")
	} else {
		if raiseMoney >= current.Players[index].Chips {
			actionList = append(actionList, "AllIn", "Fold")
		} else {
			actionList = append(actionList, "AllIn", "Raise", "Call", "Fold")
		}

	}
	return actionList

}


/*
Author: Luuminous
Date: Oct. 17th, 2018
*/

/*
This function is to check out the state of the players and the game and split the chip pool.
*/

func CheckOut() int {
	activePlayers := 0 // Record the active players.
	for _, player := range current.Players {
		if !player.Eliminated{
			activePlayers++
		}
	}
	numPlayers := 0 // numPlayers check the number of players which isn't fold.
	for _, player := range current.Players {
		if !player.Eliminated && !player.Fold {
			numPlayers++
		}
	}
	if numPlayers == 1 {
		return 2
	}
	// More than 1 player not fold.
	for _, player := range current.Players {
		if !player.Eliminated && !player.Fold && !(player.OK || player.AllIn) {
			return 1
		}
	}
	return 0
}

/*
Author: Luuminous
Date: Oct. 30th, 2018
 */

func PrintFoldResult() {
	winnerIndex := 0
	for index, player := range current.Players {
		if !player.Eliminated && !player.Fold {
			winnerIndex = index
		}
	}
	conditionInformationEntry[2].SetText("The Chip Pool now is: " + strconv.Itoa(current.ChipPool) + "\n" + current.Players[winnerIndex].Name + " wins all the chips!!")

	Refresh()
}

/*
  Author: Chengyang Nie
  Date: 10/15/2018 12:45
*/

/*
Flop function is used to generate the three cards for the community []Cards.
The input of the function is the reference of Current.
*/

func Flop() {
	//OutputFlop()
	current.Stage = "Flop"
	var tempList []string
	current.PreEventsList = tempList
	handCards := len(current.Pool) - 5
	endOfFlop := len(current.Pool) - 2
	temp := current.Pool[handCards:endOfFlop]
	for _, card := range temp {
		current.CommunityCard = append(current.CommunityCard, card)
	}
	for i := 0; i < len(current.Players); i++ {
		current.Players[i].OK = false
	}

}

/*
  Author: Chengyang Nie
  Date: 10/15/2018 12:58
*/

/*
Flop function is used to generate the turn card for the community []Cards.
The input of the function is the reference of Current.
*/

func Turn() {
	//OutputTurn()
	current.Stage = "Turn"
	var tempList []string
	current.PreEventsList = tempList
	turnCardIndex := len(current.Pool) - 2
	current.CommunityCard = append(current.CommunityCard, current.Pool[turnCardIndex])

	for i := 0; i < len(current.Players); i++ {
		current.Players[i].OK = false
	}
}

/*
  Author: Chengyang Nie
  Date: 10/15/2018 1:04
*/

/*
Flop function is used to generate the river card for the community []Cards.
The input of the function is the reference of Current.
*/

func River() {
	//OutputRiver()
	current.Stage = "River"
	var tempList []string
	current.PreEventsList = tempList
	riverCardIndex := len(current.Pool) - 1
	current.CommunityCard = append(current.CommunityCard, current.Pool[riverCardIndex])

	for i := 0; i < len(current.Players); i++ {
		current.Players[i].OK = false
	}
}

/*
Author: Luuminous
Date: Oct. 31, 2018
 */

func CheckOutFold() {
	for index, player := range current.Players {
		if !player.Eliminated && !player.Fold{
			current.Players[index].Chips += current.ChipPool
			current.ChipPool = 0
		}
	}
}

/*
  Author: Chengyang Nie
  Date: 10/11/2018
*/

/*
This function is used to genenrate 5 cards from 7 cards
The input is a slice of card.
The output is a two D slice of card.
I will use brute force method to list all the combinations.
*/

func Generate5CardsFrom7Cards(cardSlice []Card) [][]Card {
    if cardSlice == nil {
        panic("The input 7 cards slice is empty")
    }
    var slice [][]int
    if len(cardSlice) > 7 {
    	cardSlice = cardSlice[:7]
    }
    if len(cardSlice) == 7 {
    	
		for i := 0; i < len(cardSlice) - 1; i++ {
			for j := i + 1; j < len(cardSlice); j++ {
				temp := [] int {i, j} 
				slice = append(slice, temp)
			}
		}
    }
    if len(cardSlice) == 6 {
    	for i := 0; i < len(cardSlice); i++ {
    		temp := [] int {i}
    		slice = append(slice, temp)
    	}
    }
    if len(cardSlice) == 5 {
    	temp := make([][]Card, 0)
    	temp = append(temp, cardSlice)
    	return temp
    }
    if len(cardSlice) < 5 {
    	panic("Lalala")
    }
    
    var newSlice [][]Card
    for i := 0; i < len(slice); i++ {
    	var tempSlice []Card
    	for j := 0; j < len(cardSlice); j++ {
    		if (j != slice[i][0]) && (j != slice[i][1]) {
    			tempSlice = append(tempSlice, cardSlice[j])
    		}
    	}
    	newSlice = append(newSlice, tempSlice)
    }
    return newSlice
}


/*
This function is to convert a given integer to a string.
Add zero to one digital innteger.
*/

func ConvertIntToStr(x int) string {
	ans := strconv.Itoa(x)
	if len(ans) == 1 {
		ans = "0" + ans
	}
	return ans
}

/*
  Author: Chengyang Nie
  Date: 10/11/2018
*/

/*
This function is used to check if the 5 nums is a straight.
The input of the function is the sorted five nums in a slice.
The output of the funciton is a bool.

Check, if the distance of two nums is 1, if not return false.
*/

func IsStraight(sortedSlice []int) bool {
	//check if this is bicycle
	if sortedSlice == nil {
		panic("The input sortedSlice is empty")
	}
	if sortedSlice[0] == 14 && sortedSlice[1] == 5 && sortedSlice[2] == 4 && sortedSlice[3] == 3 && sortedSlice[4] == 2{
		return true
	}
	for i := 0; i < len(sortedSlice) - 1; i++ {
		//check the distance between two is 1
		if sortedSlice[i] - sortedSlice[i + 1] == 1 {
			continue
		} else {
			return false
		}
	}
	return true
}

/*
  Author: Chengyang Nie
  Date: 10/11/2018 1:02
 */


/*
This function is used to check if these five cards are the same color
The input is a slice with 5 numbers which represents the color
The output is a bool

*/

func IsFlush(colorSlice []int) bool {
	if colorSlice == nil {
		panic("colorSlice is empty")
	}
	for i := 0; i < len(colorSlice); i++ {
		if colorSlice[0] != colorSlice[i] {
			return false
		}
	}
	return true
}

/*Duplicate check whether there is the duplicate in the card
input is the slice of int, return two slices the value and its occurences
@Yichen Oct.11.2018
*/

func Duplicate(cards []int) ([]int, []int) {
	valueSlice := []int{}
	occurSlice := []int{}
	occur := make(map[int]int) //this is the map that key is the card number and value is the occuurences
	for i := 0; i < len(cards); i++ {
		if _, ok := occur[cards[i]]; ok {
			occur[cards[i]]++
		} else {
			occur[cards[i]] = 1
		}
	}
	maxNum := 0           // the max card value
	maxOccur := 0         // the occurences of Value
	for len(occur) != 0 { //the map is not empty
		maxNum, maxOccur = 0, 0
		for key, times := range occur {
			if times > maxOccur {
				maxOccur, maxNum = times, key
			} else if times == maxOccur && key > maxNum {
				maxNum, maxOccur = key, times
			}
		}
		valueSlice = append(valueSlice, maxNum)
		occurSlice = append(occurSlice, maxOccur)
		delete(occur, maxNum)
	}
	return occurSlice, valueSlice
}

/*
  Auther: Chengyang Nie
  Date: 10/11/2018 12:23
*/

/*
This function is used to sort five cards' num.
The input is a slice of 5 nums of 5 cards.
The output is a slice of sorted 5 nums.

*/

func Sort(numSlice []int) []int {
	if numSlice == nil {
		panic("The input slice for sorting is empty")
	}
	for i := 0; i < len(numSlice) - 1; i++ {
		for j := i + 1; j < len(numSlice); j++ {
			if numSlice[i] < numSlice[j] {
				temp := numSlice[j]
				numSlice[j] = numSlice[i]
				numSlice[i] = temp
			}
		}
	}
	return numSlice
}

/*
Author: Luuminous
Date: Oct. 11th, 2018
*/

/*
This function is to judge the pattern of the given five cards.
Input is a slice of card, output is a string, which can be compared by alphabet order.
*/

func Convert5CardsToPatterns(hands []Card) string {
	var numList []int   // Record the num of five card.
	var colorList []int // Record the color of five card.
	for _, handCard := range hands {
		numList = append(numList, handCard.Num)
		colorList = append(colorList, handCard.Color)
	}
	sortedNumList := Sort(numList)
	if IsStraight(sortedNumList) {
		if IsFlush(colorList) {
			if (sortedNumList[0] == 14) && (sortedNumList[1] == 13) {
				// Royal flush
				ans := "9"
				return ans
			} else if sortedNumList[0] == 14 {
				// Straight flush bicycle
				ans := "805"
				return ans
			} else {
				// Straight flush
				ans := "8"
				ans += ConvertIntToStr(sortedNumList[0])
				return ans
			}
		} else {
			if (sortedNumList[0] == 14) && (sortedNumList[1] == 13) {
				// Straight
				ans := "414"
				return ans
			} else if sortedNumList[0] == 14 {
				// Straight bicycle
				ans := "405"
				return ans
			} else {
				// Straight
				ans := "4"
				ans += ConvertIntToStr(sortedNumList[0])
				return ans
			}
		}
	} else {
		duplicate, numDuplicate := Duplicate(sortedNumList)
		if duplicate[0] == 1 {
			if IsFlush(colorList) {
				// Flush
				ans := "5"
				for _, value := range sortedNumList {
					ans += ConvertIntToStr(value)
				}
				return ans
			} else {
				// High card
				ans := "0"
				for _, value := range sortedNumList {
					ans += ConvertIntToStr(value)
				}
				return ans
			}
		} else if duplicate[0] == 4 {
			// Four of a kind
			ans := "7"
			for _, value := range numDuplicate {
				ans += ConvertIntToStr(value)
			}
			return ans
		} else if duplicate[0] == 3 {
			if duplicate[1] == 2 {
				// Full house
				ans := "6"
				for _, value := range numDuplicate {
					ans += ConvertIntToStr(value)
				}
				return ans
			} else {
				// Three of a kind
				ans := "3"
				for _, value := range numDuplicate {
					ans += ConvertIntToStr(value)
				}
				return ans
			}
		} else if duplicate[0] == 2 {
			if duplicate[1] == 2 {
				// Two pair
				ans := "2"
				for _, value := range numDuplicate {
					ans += ConvertIntToStr(value)
				}
				return ans
			} else {
				// One pair
				ans := "1"
				for _, value := range numDuplicate {
					ans += ConvertIntToStr(value)
				}
				return ans
			}
		}
	}
	panic("no pattern matched")
}

/*
Author: Luuminous
Date: Oct. 11th, 2018
*/

/*
This function is calculate the max pattern from 7 cards.
*/

func MaxPattern(totalCards []Card) (string, []Card) {
	possibleCardsList := Generate5CardsFrom7Cards(totalCards)
	max := "00"
	var tempCards []Card
	for _, cards := range possibleCardsList {
		temp := Convert5CardsToPatterns(cards)
		if temp > max {
			tempCards = cards
			max = temp
		}
	}
	
	return max, tempCards
}

/*
Author: Chengyang Nie
Date: Nov. 25th, 2018
*/

func PrintWinner(maxSeatNumber []int) {
	winnerList := ""

	for i := 0; i < len(maxSeatNumber); i++ {
		for j := 0; j < len(current.Players); j++ {

			if current.Players[j].SeatPosition == maxSeatNumber[i] {
				var sevenCardSlice []Card
				sevenCardSlice = append(sevenCardSlice, current.Players[j].Hands[0], current.Players[j].Hands[1])
				for m := 0; m < len(current.CommunityCard); m++ {
					sevenCardSlice = append(sevenCardSlice, current.CommunityCard[m])
				}
				//_, tempCards := MaxPattern(sevenCardSlice)

				winnerList += ", " + current.Players[j].Name
				//fmt.Print(" and the max pattern is ")
				//fmt.Println(HandsToString(tempCards))
			}
		}
	}
	conditionInformationEntry[2].SetText("The winner is" + winnerList[1:])

	time.Sleep(8 * 1000000000)
}


/*
  Author: Chengyang Nie
  Date: 10/15/2018 1:20
*/

/*
ShowDown function is used to show the two cards of each player
The input is the reference of Current
The ouput is also the reference of Current
*/

func ShowDown() []int {

	var playerSlice []Player

	for i := 0; i < len(current.Players); i++ {
		if !current.Players[i].Fold && !current.Players[i].Eliminated {
			playerSlice = append(playerSlice, current.Players[i])
		}
	}

	PrintShowDown()

	CalculateAlpha()

	curMax := "00"
	var maxSeatNumber []int
	for i := 0; i < len(playerSlice); i++ {
		var sevenCardSlice []Card
		sevenCardSlice = append(sevenCardSlice, playerSlice[i].Hands[0], playerSlice[i].Hands[1])
		for j := 0; j < len(current.CommunityCard); j++ {
			sevenCardSlice = append(sevenCardSlice, current.CommunityCard[j])
		}

		compareString, _ := MaxPattern(sevenCardSlice)

		if compareString > curMax {
			curMax = compareString
		}
	}

	for i := 0; i < len(playerSlice); i++ {
		var sevenCardSlice []Card
		sevenCardSlice = append(sevenCardSlice, playerSlice[i].Hands[0], playerSlice[i].Hands[1])
		for j := 0; j < len(current.CommunityCard); j++ {
			sevenCardSlice = append(sevenCardSlice, current.CommunityCard[j])
		}
		compareString, _ := MaxPattern(sevenCardSlice)
		if curMax == compareString {
			maxSeatNumber = append(maxSeatNumber, playerSlice[i].SeatPosition)

		}
	}

	PrintWinner(maxSeatNumber)
	/*
		This loop is to set fold to all loser.
	*/
	for index, player := range current.Players {
		if InInt(player.SeatPosition, maxSeatNumber) {
			current.Players[index].Fold = false
		} else {
			current.Players[index].Fold = true
		}
	}

	return maxSeatNumber
}

/*
PrintShowDown will print the ShowDown results for player @Chengyang
*/

func PrintShowDown() {
	var playerSlice []int

	for i := 0; i < len(current.Players); i++ {
		if !current.Players[i].Fold && !current.Players[i].Eliminated {
			playerSlice = append(playerSlice, i)
		}
	}

	for i := range botIndexList {
		if InInt(botIndexList[i], playerSlice) {
			playerInformationEntry[i + 1][3].SetText("Hands: " + HandsToString(current.Players[botIndexList[i]].Hands))
		} 
	}

	time.Sleep(3 * 1000000000)
}

/*
Author: Luuminous
Date: Oct. 17th, 2018
*/

func FinalCheckOut(winner []int) {
	//Make the all-in loser eliminated.
	for index, player := range current.Players {
		if !player.Eliminated && player.AllIn && !InInt(player.SeatPosition, winner){
			current.Players[index].Eliminated = true
			current.Players[index].PositionName = "Loser"
			current.Players[index].Bet = 0
			current.Players[index].Hands = []Card{}
		}
	}
	//Calculate how many chips each players can win.
	totalBetAmount := 0
	for _, player := range current.Players {
		if InInt(player.SeatPosition, winner) {
			totalBetAmount += player.Bet
		}
	}
	for index, player := range current.Players {
		if InInt(player.SeatPosition, winner) {
			amount := SetToTen(current.ChipPool * player.Bet / totalBetAmount)
			current.Players[index].Chips += amount
			conditionInformationEntry[2].SetText(current.Players[index].Name + " wins " + strconv.Itoa(amount) + " chips!!")

			time.Sleep(2 * 1000000000)
		}
	}
}

/*
  Author: Chengyang Nie, Luuminous
  Date: 10/14/2018 11:49
*/

/*
This function is used to start one Texas Holdem game.
The input of the function is the Current Consturct.
This function will contain subroutines of one game.
*/

func StartOneGame(c chan string) {

	numPlayers := 0
	for _, player := range current.Players {
		if !player.Eliminated {
			numPlayers++
		}
	}

	InitialOneGame()

	//Wait for 0.7 second for 
	time.Sleep(0.7 * 1000000000)

	isFirst := true
	for CheckOut() == 1 {
		if isFirst{
			if numPlayers >= 3 {
				Command(3, c)
				isFirst = false
			} else {
				Command(1, c)
				isFirst = false
			}
		} else {
			Command(1, c)
		}
	}
	if CheckOut() == 2 {
		PrintFoldResult()

		time.Sleep(5 * 1000000000)
		CheckOutFold()
		current.GameCount++
		return
	}
	
	/*
	Flop.
	*/
	Flop()

	Refresh()

	time.Sleep(1 * 1000000000)

	for CheckOut() == 1 {
		Command(1, c)
	}
	if CheckOut() == 2 {
		PrintFoldResult()

		time.Sleep(5 * 1000000000)
		CheckOutFold()
		current.GameCount++
		return
	}

	/*
	Turn.
	*/

	Turn()

	Refresh()

	time.Sleep(1 * 1000000000)
	
	for CheckOut() == 1 {
		Command(1, c)
	}
	if CheckOut() == 2 {
		PrintFoldResult()

		time.Sleep(5 * 1000000000)
		CheckOutFold()
		current.GameCount++
		return
	}
	
	/*
	River.
	*/

	River()

	Refresh()

	time.Sleep(1 * 1000000000)
	
	for CheckOut() == 1 {
		Command(1, c)
	}

	if CheckOut() == 2 {
		PrintFoldResult()

		time.Sleep(5 * 1000000000)
		CheckOutFold()
		current.GameCount++
		return
	}
	winner := ShowDown()

	FinalCheckOut(winner)

	Refresh()
	current.GameCount++

}

/*
Author: Luuminous
Date: Nov. 22nd, 2018
*/

/*
This function is to the main platform of the whole game.
*/

func Information(app fyne.App) {
	w := app.NewWindow("Texas Holdem")
	//Set default
	numPlayers = 6
	initialChips = 1000
	numTurns = 10
	playerName = "Rosemary"
	numOfEachTypePlayers[0] = 1 // Human players
	numOfEachTypePlayers[1] = 2 // Random bot
	numOfEachTypePlayers[2] = 3 // Probability based bot
	numOfEachTypePlayers[3] = 0 // Conventional rules based bot
	numOfEachTypePlayers[4] = 0 // Persona based bot
	numOfEachTypePlayers[5] = 0 // Minimum regret based bot

	numPlayersEntry = &widget.Label{Text: "numPlayers: " + strconv.Itoa(numPlayers),}
	initialChipsEntry = &widget.Label{Text: "initialChips: " + strconv.Itoa(initialChips),}
	numTurnsEntry = &widget.Label{Text: "numTurns: " + strconv.Itoa(numTurns),}
	playerNameEntry = &widget.Label{Text: "playerName: " + playerName,}
	numOfEachTypePlayersEntry[0] = &widget.Label{Text: "Human player: " + strconv.Itoa(numOfEachTypePlayers[0])}
	numOfEachTypePlayersEntry[1] = &widget.Label{Text: "Random bot: " + strconv.Itoa(numOfEachTypePlayers[1])}
	numOfEachTypePlayersEntry[2] = &widget.Label{Text: "Probability bot: " + strconv.Itoa(numOfEachTypePlayers[2])}
	numOfEachTypePlayersEntry[3] = &widget.Label{Text: "Conventional bot: " + strconv.Itoa(numOfEachTypePlayers[3])}
	numOfEachTypePlayersEntry[4] = &widget.Label{Text: "Persona drawer bot: " + strconv.Itoa(numOfEachTypePlayers[4])}
	numOfEachTypePlayersEntry[5] = &widget.Label{Text: "Super talented bot: " + strconv.Itoa(numOfEachTypePlayers[5])}

	//fmt.Println(reflect.TypeOf(w).String())
	w.SetContent(&widget.Box{Children: []fyne.CanvasObject{
		fyne.NewContainerWithLayout(layout.NewGridLayout(3), 
			&layout.Spacer{},
			numPlayersEntry,
			&layout.Spacer{},

			&layout.Spacer{},
			initialChipsEntry,
			&layout.Spacer{},

			&layout.Spacer{},
			numTurnsEntry,
			&layout.Spacer{},

			&layout.Spacer{},
			playerNameEntry,
			&layout.Spacer{},
		),
		fyne.NewContainerWithLayout(layout.NewGridLayout(5), 
			&layout.Spacer{},
			numOfEachTypePlayersEntry[0],
			numOfEachTypePlayersEntry[1],
			numOfEachTypePlayersEntry[2],
			&layout.Spacer{},

			&layout.Spacer{},
			numOfEachTypePlayersEntry[3],
			numOfEachTypePlayersEntry[4],
			numOfEachTypePlayersEntry[5],
			&layout.Spacer{},

			&layout.Spacer{},
			&widget.Button{Text: "Settings", OnTapped: func() {
				InformationEnter(app)
			}},
			&layout.Spacer{},
			&widget.Button{Text: "Play Game!", OnTapped: func() {
				PlayGame(app)
			}},
			&layout.Spacer{},

		),
	}})
	w.ShowAndRun()
}

/*
Author: Luuminous
Date: Oct. 17th, 2018
*/

/*
This function is to initiate the current.
*/
func InitialGame(numPlayers, initialChips int) Current{
	var newCurrent Current
	newCurrent.ChipPool = 0
	var newPlayerList []Player

	//Decide which one is to start.
	decide := rand.Intn(numPlayers)

	//Start from 1.
	for i := 1; i <= numPlayers; i++ {
		var newPlayer Player
		newPlayer.SeatPosition = i
		newPlayer.GamePosition = i + decide
		if newPlayer.GamePosition > numPlayers {
			newPlayer.GamePosition -= numPlayers
		}
		newPlayer.Chips = initialChips
		newPlayer.Bet = 0
		newPlayer.Fold = false
		newPlayer.AllIn = false
		newPlayer.OK = false
		newPlayer.Eliminated = false
		newPlayerList = append(newPlayerList, newPlayer)
	}
	newCurrent.Players = newPlayerList
	newCurrent.StartPlayer = 1
	newCurrent.CurrentBet = 0
	var nilCard []Card
	newCurrent.CommunityCard = nilCard
	newCurrent.Pool = nilCard
	newCurrent.GameCount = 1
	
	//Here name each player and decide their type.
	botNameList := []string {"Anchovy", "Buffalo", "Coyote", "Dromedary", "Eel", "Ferret", "Gibbon", "Heron", 
								"Iguana", "Jackal", "Kingfisher", "Lobster", "Magpie", "Numbat", "Oyster", "Python", 
								"Quail", "Rhinoceros", "Salmon", "Tarpon", "Urchin", "Vole", "Walrus", "Xeme", "Yak", "Zebra"}

	rand.Shuffle(len(botNameList), func(i, j int) {
		botNameList[i], botNameList[j] = botNameList[j], botNameList[i]
	})

	typeList := make([]string, 0)
	typeNameList := []string {"Human Player", "Random Bot", "Probability Bot", 
							"Conventional Bot", "Persona Drawer Bot", "Super Talented Bot"}
	
	if numOfEachTypePlayers[0] == 0 {
		/*
		There is no human player.
		*/

		for i := 1; i < 6; i++ {
			for j := 0; j < numOfEachTypePlayers[i]; j++ {
				typeList = append(typeList, typeNameList[i])
			}
		}

		rand.Shuffle(len(typeList), func(i, j int) {
			typeList[i], typeList[j] = typeList[j], typeList[i]
		})

		for i := 0; i < numPlayers; i++ {
			newCurrent.Players[i].Name = botNameList[i]
			newCurrent.Players[i].PlayerType = typeList[i]
		}

	} else {
		/*
		There is one human player.
		*/
		//The human player will always be the Players[0].

		newCurrent.Players[0].Name = playerName
		newCurrent.Players[0].PlayerType = "Human Player"

		for i := 1; i < 6; i++ {
			for j := 0; j < numOfEachTypePlayers[i]; j++ {
				typeList = append(typeList, typeNameList[i])
			}
		}

		rand.Shuffle(len(typeList), func(i, j int) {
			typeList[i], typeList[j] = typeList[j], typeList[i]
		})

		for i := 1; i < numPlayers; i++ {
			newCurrent.Players[i].Name = botNameList[i - 1]
			newCurrent.Players[i].PlayerType = typeList[i - 1]
		}
		
	}

	for i := range newCurrent.Players {
		fmt.Println(newCurrent.Players[i].Name + ": " + newCurrent.Players[i].PlayerType)
	}

	//This is not only for recording the human position, also, it may refer to the host bot position.
	newCurrent.HumanPosition = 0

	return newCurrent
}

func main() {
	rand.Seed(time.Now().UnixNano())

	//Initiate for conventionalBot
	InitialConventionalBot()
	//Initiate for probabilityBot
	coefficient = ReadCoefficient("Rules.txt")

	fmt.Println("Let's play Texas Holdem!")

	/*
	This part is for the make settings and start game.
	*/
	//app is the platform of the whole game!.
	app := apps.NewApp()
	Information(app)
	
}

/*
Author: Chengyang Nie
Date: 11/12/2018
This function is used for creating a bot using random strategy that is randomly choosing an action from action list and random betting
*/

func RandomBot(index int) string {
	actionList := ChooseActions(index)
	playerActionIndex := rand.Intn(len(actionList))
	playerAction := actionList[playerActionIndex]
	if playerAction == "Raise" {
		lowerBound := current.CurrentBet - current.Players[index].Bet + 10
		upperBound := current.Players[index].Chips - 10
		if upperBound <= lowerBound {
			return "AllIn"
		}
		money := SetToTen(rand.Intn(upperBound - lowerBound) + lowerBound)
		playerAction = playerAction + strconv.Itoa(money)
	}
	return playerAction
}

/*this is the function for the probability bot
Nov.12.2018 @Yichen, Luuminous*/

/*ProbabilityBot analyze the current bot's cards and do the simulations return the probability given the current known community card and bot's cards*/

func ProbabilityBot(index int) string {

	numActivePlayers := 0
	//this for loop returns the num of active players
	for _, player := range current.Players {
		if !player.Eliminated && !player.Fold {
			numActivePlayers++
		}
	}

	//Set default winning probability.
	prob := float64(1) / float64(numActivePlayers)

	length := len(current.CommunityCard)
	var CC []Card
	for i := 0; i < length; i++ {
		CC = append(CC, current.CommunityCard[i])
	}
	prob = MCsimulation(length, CC, 1000, current.Players[index].Hands, numActivePlayers)
	lowerBound := current.CurrentBet - current.Players[index].Bet + 10
	upperBound := current.Players[index].Chips - 10

	dreamChips := GetChipsByProb(coefficient[numActivePlayers], prob, current.Players[index].Chips + current.Players[index].Bet)

	/*
	We decide our raise amount by bot current chips.
	*/

	var actionList []string
	actionList = ChooseActions(index)
	check := false

	/*
	We only need to at most maxmoney chips to make sure other players will all in.
	*/
	maxMoney := 0
	for i := range current.Players {
		if !current.Players[i].Fold && maxMoney < current.Players[i].Chips && i != index {
			maxMoney = current.Players[i].Chips
		}
	}

	if maxMoney > upperBound {
		maxMoney = upperBound
	}
	if maxMoney < lowerBound {
		maxMoney = lowerBound
		check = true
	}

	if len(actionList) == 2 {

		//AllIn or Fold.
		if current.Players[index].Bet + current.Players[index].Chips <= dreamChips {
			return "AllIn"
		} else {
			return "Fold"
		}
	}

	if len(actionList) == 3 {

		//AllIn, Raise or Check.
		if dreamChips <= current.CurrentBet || check{
			return "Check"
		} else if current.Players[index].Bet + current.Players[index].Chips <= dreamChips {
			return "AllIn"
		} else {
			raiseMoney := dreamChips - current.Players[index].Bet
			decide := rand.Intn(4)
			if decide == 0 {
				return "Check"
			} else {
				decideMoney := SetToTen((raiseMoney - lowerBound) / decide) + lowerBound
				if maxMoney < decideMoney {
					return "Raise" + strconv.Itoa(maxMoney)
				} else {
					return "Raise" + strconv.Itoa(decideMoney)
				}		
			}
		}
	}

	if len(actionList) == 4 {

		//AllIn, Raise, Call or Fold.
		if dreamChips < current.CurrentBet {
			return "Fold"
		} else if dreamChips == current.CurrentBet || check{
			return "Call"
		} else if current.Players[index].Bet + current.Players[index].Chips <= dreamChips {
			return "AllIn"
		} else {
			raiseMoney := dreamChips - current.Players[index].Bet
			decide := rand.Intn(6)
			if decide == 0 {
				return "Call"
			} else {
				decideMoney := SetToTen((raiseMoney - lowerBound) / decide) + lowerBound
				if maxMoney < decideMoney {
					return "Raise" + strconv.Itoa(maxMoney)
				} else {
					return "Raise" + strconv.Itoa(decideMoney)
				}
			}
		}
	}
	
	//Default
	return "AllIn"
}

/*This function is to calculate the dream chips determined by a polynomial curve (ordered 2)*/

func GetChipsByProb(coefficient []float64, prob float64, max int) int {
	if prob * prob * coefficient[0] + prob * coefficient[1] + coefficient[2] >= 1.0{
		return initialChips * numPlayers //Set to max!
	}
	return SetToTen(int(float64(max) * (prob * prob * coefficient[0] + prob * coefficient[1] + coefficient[2])))
}

/*MCsimulation simulate the numTrials of simualtions and calculate the probability*/

func MCsimulation(numCC int, knownCC []Card, numTrials int, bothands []Card, numActivePlayers int) float64 {
	win := 0

	for i := 0; i < numTrials; i++ {
		if WhetherWin(numCC, knownCC, bothands, numActivePlayers) {
			win++
		}
	}
	//fmt.Println(win)
	return float64(win) / float64(numTrials)
}

/*WhetherWin() return whether the player could win by random generate the cards for other players and unkown generate community cards. and compare the max pattern. returns true if for this one simulation, bot wins.
umCC num of CommunityCard, knownCC known CommunityCard
*/

func WhetherWin(numCC int, knownCC []Card, bothands []Card, numActivePlayers int) bool {

	initiationPool := Initiation()
	//fmt.Println(HandsToString(bothands))
	deletedOwnCard := DeleteFromPool(initiationPool, bothands)
	deletedOwnCard = DeleteFromPool(deletedOwnCard, knownCC)
	//fmt.Println("This is the bot's hands : ", HandsToString(bothands))
	//delete the bot's holehands from the pool
	RandomgenecommunityCard := CreateRandomCards(5 - numCC, deletedOwnCard)
	//random generate the unknown communityCard
	communityCard := append(knownCC, RandomgenecommunityCard...)
	//fmt.Println("This is the community card : ", HandsToString(communityCard))

	botCards := append(bothands, communityCard...)
	maxbotCards, _ := MaxPattern(botCards)

	leftcards := DeleteFromPool(deletedOwnCard, RandomgenecommunityCard)
	//fmt.Println(HandsToString(leftcard))
	//fmt.Println(len(leftcard))
	//fmt.Println("Left cards pool : ", HandsToString(leftcards), len(leftcards))
	for i := 1; i < numActivePlayers; i++ {
		playercards := CreateRandomCards(2, leftcards)
		//fmt.Println("Guess the player's hands : ", HandsToString(playercards))
		cards := append(playercards, communityCard...)
		maxplayer, _ := MaxPattern(cards)
		if maxplayer > maxbotCards {
			return false
		}
		leftcards = DeleteFromPool(leftcards, playercards)
		//fmt.Println("Left cards pool : ", HandsToString(leftcards), len(leftcards))
	}
	return true
}

/*ReadCoefficient will read the coefficient text according to the probability, the bot will choose corresponding strategies.
The file column names should be as following*/

func ReadCoefficient(filename string) map[int][]float64 {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: something went wrong opening the strategy file.")
		os.Exit(1)
	}
	var lines []string
	lines = make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if scanner.Err() != nil {
		fmt.Println("Sorry: there was some kind of error during the file reading")
		os.Exit(1)
	}

	file.Close()

	thresholds := make(map[int][]float64)
	//thresholds is the map key is the int of num of active players. value is slice of float64
	for i := range lines {
		var items []string
		items = strings.Split(lines[i], " ")
		numPlayers, err1 := strconv.Atoi(items[0])
		if err1 != nil {
			fmt.Println("There is something wrong with converting numActivePlayers to int")
		}
		thresholds[numPlayers] = make([]float64, 0)
		for j := 1; j < len(items); j++ {
			bar, err2 := strconv.ParseFloat(items[j], 64)
			if err2 == nil {
				thresholds[numPlayers] = append(thresholds[numPlayers], bar)
			}
		}
	}

	return thresholds

}


/*
Author: Luuminous
Date: Oct. 11th, 2018
*/

/*
This function is to delete the out cards from the pool cards.
*/
func DeleteFromPool(pool []Card, out []Card) []Card {
	if len(out) >= len(pool) {
		panic("Invalid input: out cards' length is greater than pool")
	}
	var ans []Card
	for i := 0; i < len(pool); i++ {
		check := true
		for j := 0; j < len(out); j++ {
			if (pool[i].Num == out[j].Num) && (pool[i].Color == out[j].Color) {
				check = false
			}
		}
		if check {
			ans = append(ans, pool[i])
		}
	}
	return ans
}

/*
Author: Chengyang Nie
Date: 11/22/2018

This function is used to create a bot that implements conventional rule strategy.

The conventional rules are:
1. range for palying cards, such as 40%, 30%, and I will create a random factor in the boundary using double curve algorithm to blur the boudary.
2. Based on the expected earnings. The rule will be: winpercentage x winmoney + losepercentage x (-losemoney)
3. Four-Two rule: At Flop, the probability to get a card you want at river is 4N%; At Turn, the probability to get a card you want at river phase is 2N%.
4. The range for palying card can be affected by the playing position in each game.
5. The range for palying card can be affected by the chips you have.

*/

func InitialConventionalBot() {
	file, err := os.Open("conventionalProbability.txt")
	if err != nil {
		fmt.Println("Wrong in opening the file!")
		os.Exit(1)
	}
	line := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = append(line, scanner.Text())
	}

	file.Close()
	preFlopProbability = PutInformationToMap(line)
}

func ConventionalBot(index int) string {
	file, err1 := os.Open("conventionalRule.txt")
	if err1 != nil {
		fmt.Println("Wrong in opening the file!")
		os.Exit(1)
	}
	//The following is used to read the file
	lines := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if scanner.Err() != nil {
		fmt.Println("Sorry: there was some kind of error during the file reading")
		os.Exit(1)
	}

	file.Close()

	key := strconv.Itoa(current.Players[index].Hands[0].Color) + strconv.Itoa(current.Players[index].Hands[0].Num) + strconv.Itoa(current.Players[index].Hands[1].Color) + strconv.Itoa(current.Players[index].Hands[1].Num)
	var liveCount int
	for i := 0; i < len(current.Players); i++ {
		if !current.Players[i].Eliminated {
			liveCount++
		}
	}
	handWinPercentage := preFlopProbability[key][liveCount - 2] //float64, the percentage to win.

	//Following is to read the strategy file.

	var items []string
	items = strings.Split(lines[0], " ")
	initialRange1, _ := strconv.ParseFloat(items[1], 64) 

	items = strings.Split(lines[1], " ")
	initialRange2, _ := strconv.ParseFloat(items[1], 64) 
	raiseMoney1, _ := strconv.Atoi(items[4])

	items = strings.Split(lines[2], " ")
	gaussianFactor, _ := strconv.ParseFloat(items[1], 64) 

	items = strings.Split(lines[3], " ")
	positionFactor, _ := strconv.ParseFloat(items[1], 64) 

	items = strings.Split(lines[4], " ")
	earnings1, _ := strconv.Atoi(items[1])

	items = strings.Split(lines[5], " ")
	earnings2, _ := strconv.Atoi(items[1])
	raiseMoney2, _ := strconv.Atoi(items[4])

	var result string
	if current.Stage == "Pre-flop" {
		result = CheckRange(index, initialRange1, initialRange2, raiseMoney1, gaussianFactor, positionFactor, handWinPercentage)
		return result
	} else {
		//Flop, Turn, River
		result = CheckEarnings(index, earnings1, earnings2, raiseMoney2)
		return result
	}

}

func CheckRange(index int, initialRange1 float64, initialRange2 float64, raiseMoney1 int, gaussianFactor float64, positionFactor float64, handWinPercentage float64) string {
	//read the file according to the livePlayerNumber
	//initialRange1's boundary becomes a range according to curve
	//initialRange2's boundary becomes a range according to curve
	
	//range through the players to get the unelimiated player numbers
	var liveCount int
	for i := 0; i < len(current.Players); i++ {
		if !current.Players[i].Eliminated {
			liveCount++
		}
	}


	initialRange1 -= (positionFactor / float64(liveCount)) * float64(current.Players[index].GamePosition + 1)
	initialRange2 -= (positionFactor / float64(liveCount)) * float64(current.Players[index].GamePosition + 1)
	//fmt.Println(initialRange1)
	//fmt.Println(initialRange2)
	//range through the txt to find the winning percentage of the hands
	
	winningPercentage := handWinPercentage
	
	probaCall := GetPro(winningPercentage, initialRange1, gaussianFactor)
	probaRaise := GetPro(winningPercentage, initialRange2, gaussianFactor)
	//fmt.Println(probaCall)
	//fmt.Println(probaRaise)

	actionList := ChooseActions(index)
	randNumber := rand.Float64()
	if len(actionList) == 2 {

		//AllIn, Fold
		if randNumber <= probaCall {
			return "AllIn"
		} else {
			return "Fold"
		}
	}

	if len(actionList) == 3 {

		//AllIn, Check, Raise
		if randNumber <= probaRaise {
			if current.Players[index].Chips > current.CurrentBet - current.Players[index].Bet + raiseMoney1 {
				moneyString := strconv.Itoa(current.CurrentBet - current.Players[index].Bet + raiseMoney1)
				return "Raise" + moneyString 
			} else {
				return "AllIn"
			}
		}
		return "Check"
	}

	if len(actionList) == 4 {

		//AllIn, Raise, Fold, Call
		if randNumber > probaRaise && randNumber <= probaCall {
			return "Call"
		}
		if randNumber <= probaRaise {
			if current.Players[index].Chips > current.CurrentBet - current.Players[index].Bet + raiseMoney1 {
				moneyString := strconv.Itoa(current.CurrentBet - current.Players[index].Bet + raiseMoney1)
				return "Raise" + moneyString 
			} else {
				return "AllIn"
			}
		}
		return "Fold"
	}
	
	return "AllIn" //default
}


func GetPro(x float64, mean float64, factor float64) float64 {

	//e^((x - u) / f) / (1 + e^((x - u)) / f)
	upper := math.Exp((x - mean) / factor)
	base := 1 + math.Exp((x - mean) / factor)
	result := upper / base
	return result
}


func CheckEarnings(index int, earnings1 int, earnings2 int, raiseMoney2 int) string {
	if current.Stage == "Flop" {
		result := InitialCheck(index) //Result = true when you have super power hands
		actionList := ChooseActions(index)
		
		if result {
			if len(actionList) == 2 {

				//AllIn, Fold.
				return "AllIn"
			}

			if len(actionList) == 3 {

				//AllIn, Raise, Check
				/*
				Slow play!!! Lure the opponents into bets
				*/
				return "Check"
			}

			if len(actionList) == 4 {

				//AllIn, Call, Raise, Fold
				/*
				Slow play!!! Lure the opponents into bets
				*/
				return "Call"
			}
		}

		//totalOut is all the outs for straight, flush and pair
		
		totalOut := FindOut(index)
		//fmt.Println(totalOut)

		winningPercentage := float64(4 * totalOut) * 0.01
		losingPercentage := 1.0 - winningPercentage
		winningMoney := float64(current.ChipPool - current.Players[index].Bet)
		losingMoney := float64(current.Players[index].Bet)
		
		finalDecision := winningPercentage * winningMoney - losingPercentage * losingMoney

		//fmt.Println(finalDecision)

		if len(actionList) == 2 {

			//AllIn, Fold
			if finalDecision >= float64(earnings1) {
				return "AllIn"
			} else {
				return "Fold"
			}
		}

		if len(actionList) == 3 {

			//AllIn, Check, Raise
			if finalDecision >= float64(earnings2)  {
				if current.Players[index].Chips <= current.CurrentBet - current.Players[index].Bet + raiseMoney2 {
					return "AllIn"
				}
				raiseNumber := strconv.Itoa(current.CurrentBet - current.Players[index].Bet + raiseMoney2)
				return "Raise" + raiseNumber
			} else {
				return "Check"
			}
		}

		if len(actionList) == 4 {

			//AllIn, Fold, Call, Raise
			if finalDecision >= float64(earnings2) {
				if current.Players[index].Chips <= current.CurrentBet - current.Players[index].Bet + raiseMoney2 {
					return "AllIn"
				}
				raiseNumber := strconv.Itoa(current.CurrentBet - current.Players[index].Bet + raiseMoney2)
				return "Raise" + raiseNumber
			} else if finalDecision > float64(earnings1) {
				return "Call"
			} else {
				return "Fold"
			}
		}	
	}

	if current.Stage == "Turn" {
		result := InitialCheck2(index) //true when has staight, flush, 3 card or quad
		actionList := ChooseActions(index)

		if result {
			return "AllIn"
		}

		totalOut2 := FindOut(index)
		winningPercentage2 := float64(2 * totalOut2) * 0.01
		losingPercentage2 := 1 - winningPercentage2
		winningMoney2 := float64(current.ChipPool - current.Players[index].Bet)
		losingMoney2 := float64(current.Players[index].Bet)
		
		finalDecision := winningPercentage2 * winningMoney2 - losingPercentage2 * losingMoney2

		if len(actionList) == 2 {

			//AllIn, Fold
			if finalDecision >= float64(earnings2) {
				return "AllIn"
			} else {
				return "Fold"
			}
		}

		if len(actionList) == 3 {

			//AllIn, Check, Raise
			if finalDecision >= float64(earnings2)  {
				if current.Players[index].Chips <= current.CurrentBet - current.Players[index].Bet + raiseMoney2 {
					return "AllIn"
				}
				raiseNumber := strconv.Itoa(current.CurrentBet - current.Players[index].Bet + raiseMoney2)
				return "Raise" + raiseNumber
			} else {
				return "Check"
			}
		}

		if len(actionList) == 4 {

			//AllIn, Call, Raise, Fold
			if finalDecision >= float64(earnings2) {
				if current.Players[index].Chips <= current.CurrentBet - current.Players[index].Bet + raiseMoney2 {
					return "AllIn"
				}
				raiseNumber := strconv.Itoa(current.CurrentBet - current.Players[index].Bet + raiseMoney2)
				return "Raise" + raiseNumber
			} else if finalDecision > float64(earnings1) {
				return "Call"
			} else {
				return "Fold"
			}
		}
	}

	if current.Stage == "River" {
		
		actionList := ChooseActions(index)
		var totalCards []Card
		totalCards = append(totalCards, current.CommunityCard...)
		totalCards = append(totalCards, current.Players[index].Hands...) // len(totalCard) = 7
		max, finalCards := MaxPattern(totalCards)
		maxB := max[0] // The pattern type

		check := false
		for i := 0; i < len(finalCards); i++ {
			for j := 0; j < len(current.CommunityCard); j++ {
				if finalCards[i] == current.CommunityCard[j] {
					break
				}
			}
			check = true
			break
		}

		if len(actionList) == 2 {

			//AllIn, Fold
			if maxB >= '2' && check == true { // Two pairs or larger pattern
				return "AllIn"
			} else {
				return "Fold"
			}
		}

		if len(actionList) == 3 {

			//AllIn, Raise, Check
			if maxB >= '2' && check == true { // Two pairs or larger pattern
				if current.Players[index].Chips <= current.CurrentBet - current.Players[index].Bet + 200 {
					return "AllIn"
				}
				raiseNumber := strconv.Itoa(current.CurrentBet - current.Players[index].Bet + 200)
				return "Raise" + raiseNumber
			} else {
				return "Check"
			}
		}

		if len(actionList) == 4 {

			//AllIn, Raise, Call, Fold
			if maxB >= '2' && check == true { // Two pairs or larger pattern
				if current.Players[index].Chips <= current.CurrentBet - current.Players[index].Bet + 200 {
					return "AllIn"
				}
				raiseNumber := strconv.Itoa(current.CurrentBet - current.Players[index].Bet + 200)
				return "Raise" + raiseNumber
			}
			if maxB == '1' && check == true { // One pair
				return "Call"
			}
			return "Fold"
		}
	}
	return "AllIn" //default
}

//This initialCheck is used for the check during the flop phase
//If bot has straight, flush, nut or set, return true.

func InitialCheck(index int) bool {
	var totalCards []Card
	var numberList []int
	var colorList []int
	totalCards = append(totalCards, current.CommunityCard...) // len(totalCards) = 3
	totalCards = append(totalCards, current.Players[index].Hands...) // len(totalCards) = 5
	for i := 0; i < len(totalCards); i++ {
		numberList = append(numberList, totalCards[i].Num)
		colorList = append(colorList, totalCards[i].Color)
	}
	numberList = Sort(numberList)
	//check if is straight
	result1 := IsStraight(numberList)

	if result1 {
		return true
	}
	//check if is flush
	result2 := IsFlush(colorList)
	if result2 {
		return true
	}

	//check if it is nuts or set
	checkCommunity := CheckCommunityCards(index)
	if checkCommunity {
		return true
	}

	max := 1
	for i := 0; i < len(numberList) - 1; i++ {
		count := 1
		for j := i + 1; j < len(numberList); j++ {
			if numberList[i] == numberList[j] {
				count++
			} else {
				break
			}
		}
		if count > max {
			max = count
		}
	}
	if max >= 3 {
		return true
	}

	return false
}

//This check is used for the check during the Turn phase
//If bot has straight, flush, 3 card or quad, return true

func InitialCheck2(index int) bool {
	var totalCards []Card
	totalCards = append(totalCards, current.CommunityCard...)
	totalCards = append(totalCards, current.Players[index].Hands...) // len(totalCard) = 6
	cardList := GenerateCardsFromCards(totalCards, 1) //5 from 6
	//check if it is sqad or 3 cards
	checkCommunity := CheckCommunityCards(index)
	if checkCommunity {
		return true
	}
	for j := 0; j < len(cardList); j++ {
		var numberList []int
		var colorList []int
		for i := 0; i < len(cardList[j]); i++ {
			numberList = append(numberList, cardList[j][i].Num)
			colorList = append(colorList, cardList[j][i].Color)
		}
		numberList = Sort(numberList)
		//check if is straight
		result1 := IsStraight(numberList)

		if result1 {
			return true
		}
		//check if is flush
		result2 := IsFlush(colorList)
		if result2 {
			return true
		}
	
		max := 1
		for i := 0; i < len(numberList) - 1; i++ {
			count := 1
			for j := i + 1; j < len(numberList); j++ {
				if numberList[i] == numberList[j] {
					count++
				} else {
					break
				}
			}
			if count > max {
				max = count
			}
		}
		if max >= 3 {
			return true
		}
	}
	return false
}

// If the community card makes you have a nut or set, return true
func CheckCommunityCards(index int) bool {
	maxNumber := 0
	max := 1
	var resultSlice []int
	for i := 0; i < len(current.CommunityCard); i++ {
		resultSlice = append(resultSlice, current.CommunityCard[i].Num)
	}
	resultSlice = Sort(resultSlice)
	
	for i := 0; i < len(resultSlice) - 1; i++ {
		count := 1
		for j := i + 1; j < len(resultSlice); j++ {
			if resultSlice[i] == resultSlice[j] {
				count++
			} else {
				break
			}
		}
		if count > max {
			max = count
			maxNumber = resultSlice[i]
		}
	}

	if max == 4 {
		return false
	}
	if max == 3 {
		if maxNumber == current.Players[index].Hands[0].Num || maxNumber == current.Players[index].Hands[1].Num {
			return true
		}
	}
	if max == 2 {
		if maxNumber == current.Players[index].Hands[0].Num || maxNumber == current.Players[index].Hands[1].Num {
			return true
		} 
	}
	return false

}

func FindOut(index int) int {
	
	result := 0 // To count the number of outs to form straight, flush, 3 card.
	var totalCards []Card
	totalCards = append(totalCards, current.CommunityCard...)
	totalCards = append(totalCards, current.Players[index].Hands...)
	pool := Initiation()
	var resultChoices [][]Card
	if current.Stage == "Flop" { // len(totalCard) = 5
		resultChoices = GenerateCardsFromCards(totalCards, 1) // 4 from 5
	} else { // len(totalCard) = 6
		resultChoices = GenerateCardsFromCards(totalCards, 2) // 4 from 6
	}

	//Count the number of outs for straight and flush.
	for i := 0; i < len(resultChoices); i++ {
		poolCards := DeleteFromPool(pool, resultChoices[i])
		for m := 0; m < len(poolCards); m++ {
			var numberList []int
			var colorList []int
			for j := range resultChoices[i] {
				numberList = append(numberList, resultChoices[i][j].Num)
				colorList = append(colorList, resultChoices[i][j].Color)
			}
			numberList = append(numberList, poolCards[m].Num)
			colorList = append(colorList, poolCards[m].Color)
			numberList = Sort(numberList)
			if IsStraight(numberList) {
				result++
			}
			if IsFlush(colorList) {
				result++
			}
		}
	}

	//The following to loop to find out for forming pairs.
	//Not considering 3 card and quad, because if there are 3 card and quad, FindOut() won't be called.

	for i := 0; i < len(totalCards) - 1; i++ {
		for j := i + 1; j < len(totalCards); j++ {
			if totalCards[i].Num == totalCards[j].Num {
				result += 2
			}
		}
	}

	for i := 0; i < len(current.CommunityCard) - 1; i++ {
		for j := i + 1; j < len(current.CommunityCard); j++ {
			if current.CommunityCard[i] == current.CommunityCard[j] {
				result -= 2
			}
		}
	}

	return result
}

//This function is to generate all the combination of 4 cards from slice (when number = 1)
//This function is to generate all the combination of len(cards) - 2 from slice (when number = 2)

func GenerateCardsFromCards(cardSlice []Card, number int) [][]Card {

    if number == 2 {
    	var slice [][]int // Store all the combination of two cards index
		for i := 0; i < len(cardSlice) - 1; i++ {
			for j := i + 1; j < len(cardSlice); j++ {
				temp := [] int {i, j} 
				slice = append(slice, temp)
			}
		}
    
   		var newSlice [][]Card
    	for i := 0; i < len(slice); i++ {
    		var tempSlice []Card
    		for j := 0; j < len(cardSlice); j++ {
    			if (j != slice[i][0]) && (j != slice[i][1]) {
    				tempSlice = append(tempSlice, cardSlice[j])
    			}
    		}
    		newSlice = append(newSlice, tempSlice)
    	}
    	return newSlice
    }

    if number == 1 {
    	var newSlice [][]Card // Store all the arrange of two cards
    	for i := 0; i < len(cardSlice); i++ {
    		var tempSlice []Card
    		for j := 0; j < len(cardSlice); j++ {
    			if i != j {
    				tempSlice = append(tempSlice, cardSlice[j])
    			}
    		}
    		newSlice = append(newSlice, tempSlice)
    	}
    	return newSlice
    }

    return nil
}

//This function is to load the information and store it to a map
func PutInformationToMap(lines []string) map[string][]float64 {
	result := make(map[string][]float64)

	for i := 0; i < len(lines); i++ {
		var items []string
		items = strings.Split(lines[i], " ")
		var value []float64
		for i := 4; i < len(items); i++ {
			temp, _ := strconv.ParseFloat(items[i], 64)
			value = append(value, temp)
		}
		
		var key string = items[0] + items[1] + items[2] + items[3]
		result[key] = value
	}
	return result
}

//@Yichen Mo
//personaMap is a global variable
/*PersonaBot is the bot will analyze the players strategies, and return the string as the strategy*/
func PersonaBot(index int) string {
	var actionList []string
	actionList = ChooseActions(index)

	numActivePlayers := 0
	//this for loop returns the num of active players
	for _, player := range current.Players {
		if !player.Eliminated && !player.Fold {
			numActivePlayers++
		}
	}

	//Set default winning probability.
	prob := float64(1) / float64(numActivePlayers)

	length := len(current.CommunityCard)
	var CC []Card
	for i := 0; i < length; i++ {
		CC = append(CC, current.CommunityCard[i])
	}
	prob = MCsimulation(length, CC, 1000, current.Players[index].Hands, numActivePlayers)
	lowerBound := current.CurrentBet - current.Players[index].Bet + 10
	upperBound := current.Players[index].Chips - 10

	dreamChips := GetChipsByProb(coefficient[numActivePlayers], prob, current.Players[index].Chips + current.Players[index].Bet)

	maxMoney := 0
	maxModifiedCurrentBet := 0
	check := false

	for i := range current.Players {
		if !current.Players[i].Fold && maxMoney < current.Players[i].Chips && i != index {
			maxMoney = current.Players[i].Chips
		}
		if i != index && maxModifiedCurrentBet < int(float64(current.Players[i].Bet) * alphaMap[i]) {
			maxModifiedCurrentBet = SetToTen(int(float64(current.Players[i].Bet) * alphaMap[i]))
		}
	}

	if maxMoney > upperBound {
		maxMoney = upperBound
	}
	if maxMoney < lowerBound {
		maxMoney = lowerBound
		check = true
	}

	if len(actionList) == 2 {

		//AllIn or Fold.
		if current.Players[index].Bet + current.Players[index].Chips <= dreamChips || dreamChips >= maxModifiedCurrentBet{
			return "AllIn"
		} else {
			return "Fold"
		}
	}

	if len(actionList) == 3 {

		//AllIn, Raise or Check.
		if dreamChips <= maxModifiedCurrentBet || check {
			return "Check"
		} else if current.Players[index].Bet + current.Players[index].Chips <= dreamChips {
			return "AllIn"
		} else {
			raiseMoney := dreamChips - current.Players[index].Bet
			decide := rand.Intn(4)
			if decide == 0 {
				return "Check"
			} else {
				decideMoney := SetToTen((raiseMoney - lowerBound) / decide) + lowerBound
				if maxMoney < decideMoney {
					return "Raise" + strconv.Itoa(maxMoney)
				} else {
					return "Raise" + strconv.Itoa(decideMoney)
				}
			}
		}
	}

	if len(actionList) == 4 {

		//AllIn, Raise, Call or Fold.
		if dreamChips < maxModifiedCurrentBet {
			return "Fold"
		} else if dreamChips <= current.CurrentBet || check {
			return "Call"
		} else if current.Players[index].Bet + current.Players[index].Chips <= dreamChips {
			return "AllIn"
		} else {
			raiseMoney := dreamChips - current.Players[index].Bet
			decide := rand.Intn(6)
			if decide == 0 {
				return "Call"
			} else {
				decideMoney := SetToTen((raiseMoney - lowerBound) / decide) + lowerBound
				if maxMoney < decideMoney {
					return "Raise" + strconv.Itoa(maxMoney)
				} else {
					return "Raise" + strconv.Itoa(decideMoney)
				}
			}
		}
	}

	//Default
	return "AllIn"
}

func InitialMap() map[int][4]int {
	//remember to change the OpponentEvaluation to map[int][]int
	impressionMap := make(map[int][4]int)
	for i := 0; i < len(current.Players); i++ {
		impressionMap[i] = [4]int {0, 0, 0, 0}
	}
	return impressionMap
}

/*Records is the function that will record other players action only consider Raise and Allin preflop[0], flop[1], turn[2], river[3],*/
func Record(i int, action string) {

	var amount int
	if action[0] == 'R' {
		amount, _ = strconv.Atoi(action[5:])
	} else if action == "AllIn" {
		amount = current.Players[i].Chips - current.Players[i].Bet
	}
	if current.Stage == "Pre-flop" {
		UpdateArray(i, 0, amount)
	} else if current.Stage == "Flop" {
		UpdateArray(i, 1, amount)
	} else if current.Stage == "Turn" {
		UpdateArray(i, 2, amount)
	} else if current.Stage == "River" {
		UpdateArray(i, 3, amount)
	}
}

/*UpdateArray is the subroutine to update the array in the personaMap, input is the i- players' index, pos corresponding to update index of the array personaMap[i][pos]= personaMap[i][pos]+amount*/
func UpdateArray(i, pos, amount int) {
	tempt := personaMap[i]
	tempt[pos] = tempt[pos] + amount
	personaMap[i] = tempt
}

/*TraceEvaluation return the dreamchips given the communityCard and hands during the showdown phase, index is the players' index and len is length of cc that is known*/
func TraceEvaluation(index int, len int) int {
	numPlayers := 0
	//this for loop returns the num of active players
	for _, player := range current.Players {
		if !player.Eliminated {
			numPlayers++
		}
	}
	//Set default winning probability.
	prob := float64(1) / float64(numPlayers)
	var CC []Card
	for i := 0; i < len; i++ {
		CC = append(CC, current.CommunityCard[i])
	}
	prob = MCsimulation(len, CC, 1000, current.Players[index].Hands, numPlayers)
	
	dreamChips := GetChipsByProb(coefficient[numPlayers], prob, current.Players[index].Chips + current.Players[index].Bet)
	return dreamChips
}

/*CalculateAlpha is the function that will calculate Alpha value in the fomula player.bet * alpha to estimate the actual value for the player's hand This function is used when showdown*/
func CalculateAlpha() {
	
	for i := 0; i < len(current.Players); i++ {
		if !current.Players[i].Eliminated && !current.Players[i].Fold {
			preflop := 0.0
			flop := 0.0
			turn := 0.0
			river := 0.0
			evP := TraceEvaluation(i, 0)
			evF := TraceEvaluation(i, 3)
			evT := TraceEvaluation(i, 4)
			evR := TraceEvaluation(i, 5)
			if personaMap[i][0] > evP {
				preflop = float64(personaMap[i][0] - evP) * 0.01
			}
			if personaMap[i][1] > evF {
				flop = float64(personaMap[i][1] - evF) * 0.02
			}
			if personaMap[i][2] > evT {
				turn = float64(personaMap[i][2] - evT) * 0.03
			}
			if personaMap[i][3] > evR {
				river = float64(personaMap[i][3] - evR) * 0.04
			}
			score := alphaMap[i] - (preflop + flop + turn + river)
			if score < 0.0 {
				score = 0.0
			}

			alphaMap[i] = 0.5 * (score + alphaMap[i])
		}
	}
}

//This function is to initialize the alphaMap
func InitialAlphaMap() {
	alphaMap = make([]float64, numPlayers)
	for i := 0; i < numPlayers; i++ {
		alphaMap[i] = 1.0
	}
}


/*
Author: Luuminous
Date: Dec. 2nd, 2018
*/

/*
This is a bot based on maximum minimum regret.
The bot will estimate opponents' action, then calculate the regret
*/

func MinimumRegretBot(index int) string {
	
	numActivePlayers := 0

	//this for loop returns the num of active players
	for _, player := range current.Players {
		if !player.Eliminated && !player.Fold {
			numActivePlayers++
		}
	}

	lowerBound := current.CurrentBet - current.Players[index].Bet + 10
	upperBound := current.Players[index].Chips - 10

	/*
	We decide our raise amount by bot current chips.
	*/

	var actionList []string
	actionList = ChooseActions(index)

	/*
	We only need to at most maxmoney chips to make sure other players will all in.
	*/
	maxMoney := 0
	for i := range current.Players {
		if !current.Players[i].Fold && maxMoney < current.Players[i].Chips && i != index {
			maxMoney = current.Players[i].Chips
		}
	}

	if maxMoney > upperBound {
		maxMoney = upperBound
	}
	if maxMoney < lowerBound {
		for i := range actionList {
			if i < len(actionList) {
				if actionList[i] == "AllIn" {
					actionList = append(actionList[:i], actionList[i + 1:]...)
				}
			}	
		}
		for i := range actionList {
			if i < len(actionList) {
				if actionList[i][0] == 'R' {
					actionList = append(actionList[:i], actionList[i + 1:]...)
				}
			}
		}
	}

	actionList = AmplifyActionList(index, actionList, maxMoney, lowerBound)

	minimumRegret := GetRegret(index, actionList)

	maxMinimumRegret := -32000.0
	maxIndex := 0
	for i := range minimumRegret {
		if maxMinimumRegret < minimumRegret[i] {
			maxIndex = i 
			maxMinimumRegret = minimumRegret[i]
		}
	}

	return actionList[maxIndex]
}

//This function is to amplify the actionList.
func AmplifyActionList(index int, actionList []string, maxMoney, lowerBound int) []string {
	var newActionList []string

	for i := range actionList {
		if actionList[i][0] == 'R' {
			for money := lowerBound; money <= maxMoney; money += 20 {
				if money <= maxMoney {
					newActionList = append(newActionList, actionList[i] + strconv.Itoa(money))
				}
			}
		} else {
			if actionList[i] == "AllIn" {
				if current.Players[index].Chips < maxMoney {
					newActionList = append(newActionList, actionList[i])
				}
			} else {
				newActionList = append(newActionList, actionList[i])
			}
		}
	}

	return newActionList
}

func GetRegret(index int, actionList []string) []float64 {
	var newRegret []float64

	//newCurrent := CopyCurrent()
	for i := range actionList {
		if actionList[i][0] == 'R' {
			mon := actionList[i][5:]
			intMon, _ := strconv.Atoi(mon)
			
			//When the bot choose to raise forwardly, estimate other opponents reaction.
			countNumber := 0
			dreamCurrentBet := current.CurrentBet + intMon - current.Players[index].Bet
			dreamLose := current.Players[index].Bet + intMon
			dreamWin := current.ChipPool + intMon

			for j := range current.Players {
				//If the bot choose raise, everyOne else must make a decision again.
				if j != index && !current.Players[j].Fold && !current.Players[j].Eliminated {

					if dreamCurrentBet >= current.Players[j].Bet + current.Players[j].Chips {
						//AllIn, Fold
						decide := rand.Float64()
						if decide < 0.2 {
							//Guess this player will all in (20%)
							dreamWin += current.Players[j].Chips
							countNumber++
						}
					} else {
						//AllIn, Fold, Check, Raise
						decide := rand.Float64()
						if decide < 0.2 {
							//Guess this player will all in (20%)
							dreamWin += current.Players[j].Chips
							dreamCurrentBet = current.Players[j].Bet + current.Players[j].Chips
							countNumber++
							
						} else if decide < 0.5 {
							//Guess this player will raise (30%) 
							countNumber++
							randRaise := SetToTen(rand.Intn(current.Players[j].Chips - (dreamCurrentBet - current.Players[j].Bet)))
							dreamCurrentBet = current.Players[j].Bet + randRaise
							dreamWin += randRaise

						} else if decide < 0.7 {
							//Guess this player will call (20%)
							countNumber++
							dreamWin += dreamCurrentBet - current.Players[j].Bet
						}
					}
				}
			}
			var CC []Card
			for j := 0; j < len(current.CommunityCard); j++ {
				CC = append(CC, current.CommunityCard[j])
			}
			prob := MCsimulation(len(current.CommunityCard), CC, 100, current.Players[index].Hands, countNumber + 1)
			newRegret = append(newRegret, prob * float64(dreamWin) - (1 - prob) * float64(dreamLose))

		} else if actionList[i] == "Check" {

			//When bot choose to check, estimate other bots reaction.
			countNumber := 0
			dreamCurrentBet := current.CurrentBet
			dreamLose := current.Players[index].Bet
			dreamWin := current.ChipPool
			for j := range current.Players {
				if j != index && !current.Players[j].Fold && !current.Players[j].Eliminated && !current.Players[j].OK && !current.Players[j].AllIn {
					
					//If bot chooses check, there is no way the next one need to call or fold.
					//AllIn, Check, Raise
					decide := rand.Float64()
					if decide < 0.1 {
						//Guess this player will all in (10%)
						dreamCurrentBet += current.Players[j].Chips
						//Guess there will be another one fold, so countNumber--, muting the countNumber++
						//Guess there will be another one to call this crazy guy.
						dreamWin += current.Players[j].Chips * 2
						if current.Players[index].Chips < current.Players[j].Chips {
							dreamLose += current.Players[index].Chips
						} else {
							dreamLose += current.Players[j].Chips
						}
					} else if decide < 0.7 {
						//Guess this player will raise (60%) 
						countNumber++
						if current.Players[j].Chips - 10 <= 0 {
							countNumber++ 
						} else {
							randRaise := SetToTen(rand.Intn(current.Players[j].Chips - 10))
							if randRaise > current.Players[index].Chips {
								dreamLose += current.Players[index].Chips
							} else {
								dreamLose += randRaise
							}
							dreamCurrentBet += randRaise
							if countNumber >= len(current.Players) - 1 {
								dreamWin += randRaise
							} else {
								//Guess others will all call the raise
								dreamWin += randRaise * (len(current.Players) - countNumber)
							}
						}

					} else {
						//Guess this player will check (30%)
						countNumber++
					}

				} else if j != index && !current.Players[j].Fold && !current.Players[j].Eliminated {
					countNumber++
				}
			}
			var CC []Card
			for j := 0; j < len(current.CommunityCard); j++ {
				CC = append(CC, current.CommunityCard[j])
			}
			prob := MCsimulation(len(current.CommunityCard), CC, 100, current.Players[index].Hands, countNumber + 1)
			newRegret = append(newRegret, prob * float64(dreamWin) - (1 - prob) * float64(dreamLose))

		} else if actionList[i] == "Call" || (actionList[i] == "AllIn" && len(actionList) == 2){
			
			//When the bot choose to call or passively all in, estimate other opponents reaction.
			countNumber := 0

			//AllIn
			dreamCurrentBet := current.CurrentBet + current.Players[index].Chips
			dreamLose := current.Players[index].Bet + current.Players[index].Chips
			dreamWin := current.ChipPool + current.Players[index].Chips

			if actionList[i] == "Call" {
				//Call
				dreamCurrentBet = current.CurrentBet
				dreamLose = current.CurrentBet
				dreamWin = current.ChipPool + current.CurrentBet - current.Players[index].Bet
			} 
			
			for j := range current.Players {
				if j != index && !current.Players[j].Fold && !current.Players[j].Eliminated && !current.Players[j].OK && !current.Players[j].AllIn {

					if dreamCurrentBet >= current.Players[j].Chips + current.Players[j].Bet {

						//AllIn, Fold
						decide := rand.Float64()
						if decide < 0.3 {
							//Guess this player will all in (30%)
							dreamWin += current.Players[j].Chips
							countNumber++
						}
					} else {

						//AllIn, Call, Raise, Fold
						decide := rand.Float64()
						if decide < 0.1 {
							//Guess this player will all in (10%)
							dreamCurrentBet = current.Players[j].Chips + current.Players[j].Bet
							//Guess there will be another one fold, so countNumber--, muting the countNumber++
							//Guess there will be another one to call this crazy guy.
							dreamWin += 2 * current.Players[j].Chips
							remainChips := current.Players[index].Chips - (current.CurrentBet - current.Players[index].Bet)
							if remainChips < current.Players[j].Chips {
								dreamLose += remainChips
							} else {
								dreamLose += current.Players[j].Chips
							}
							
						} else if decide < 0.6 {
							//Guess this player will raise (50%) 
							countNumber++
							
							randRaise := SetToTen(rand.Intn(current.Players[j].Chips - (current.CurrentBet - current.Players[j].Bet)))
							remainChips := current.Players[index].Chips - (current.CurrentBet - current.Players[index].Bet)
							
							if remainChips < randRaise {
								dreamLose += remainChips
							} else {
								dreamLose += randRaise
							}
							dreamCurrentBet += randRaise
							if countNumber >= len(current.Players) - 1 {
								dreamWin += randRaise
							} else {
								//Guess others will all call the raise
								dreamWin += randRaise * (len(current.Players) - countNumber)
							}
						} else if decide < 0.7 {
							//Guess this player will call (10%)
							countNumber++
							dreamWin += current.CurrentBet - current.Players[j].Bet
						}
					}
					//If bot chooses call, there is no way the next one have the chance to check or do other thing.
				} else if i != index && !current.Players[j].Fold && !current.Players[j].Eliminated {
					countNumber++
				}
			}
			var CC []Card
			for j := 0; j < len(current.CommunityCard); j++ {
				CC = append(CC, current.CommunityCard[j])
			}
			prob := MCsimulation(len(current.CommunityCard), CC, 100, current.Players[index].Hands, countNumber + 1)
			newRegret = append(newRegret, prob * float64(dreamWin) - (1 - prob) * float64(dreamLose))

		} else if actionList[i] == "AllIn" {

			//When the bot choose to all in forwardly, estimate other opponents reaction.
			countNumber := 0
			dreamCurrentBet := current.CurrentBet + current.Players[index].Chips
			dreamLose := current.Players[index].Bet + current.Players[index].Chips
			dreamWin := current.ChipPool + current.Players[index].Chips

			for j := range current.Players {
				//If the bot choose all in, everyOne else must make a decision again.
				if j != index && !current.Players[j].Fold && !current.Players[j].Eliminated {

					if dreamCurrentBet >= current.Players[j].Bet + current.Players[j].Chips {
						//AllIn, Fold
						decide := rand.Float64()
						if decide < 0.2 {
							//Guess this player will all in (20%)
							dreamWin += current.Players[j].Chips
							countNumber++
						}
					} else {
						//AllIn, Fold, Check, Raise
						decide := rand.Float64()
						if decide < 0.1 {
							//Guess this player will all in (10%)
							dreamWin += current.Players[j].Chips
							dreamCurrentBet = current.Players[j].Bet + current.Players[j].Chips
							countNumber++
							
						} else if decide < 0.15 {
							//Guess this player will raise (5%) 
							countNumber++
							randRaise := SetToTen(rand.Intn(current.Players[j].Chips - (dreamCurrentBet - current.Players[j].Bet)))
							dreamCurrentBet = current.Players[j].Bet + randRaise
							dreamWin += randRaise

						} else if decide < 0.2 {
							//Guess this player will call (5%)
							countNumber++
							dreamWin += dreamCurrentBet - current.Players[j].Bet
						}
					}
				}
			}
			var CC []Card
			for j := 0; j < len(current.CommunityCard); j++ {
				CC = append(CC, current.CommunityCard[j])
			}
			prob := MCsimulation(len(current.CommunityCard), CC, 100, current.Players[index].Hands, countNumber + 1)
			newRegret = append(newRegret, prob * float64(dreamWin) - (1 - prob) * float64(dreamLose))

		} else if actionList[i] == "Fold" {
			
			dreamLose := current.Players[index].Bet
			newRegret = append(newRegret, 0.0 - float64(dreamLose))
		}
	}
	return newRegret
}