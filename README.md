# Texas-Holdem-Game

### Author: Luuminous Chen, Chengyang Nie, Yichen Mo
### Introduction
Texas Hold'em (also known as Texas Hold'em) is a variation of the card game of poker.

Two cards are face down to each player, called holecards, and then five community cards are face up in three stages. All players are assigned the even amount of chips. Unlike other poker games, all players are trying to win the most chips. If a player wants to join the game and take the chance to win the chips, they must put at least the same amount of bet with the current bet, called “call”. If the player doesn’t have enough chips to take apart, he can still “all in” his chips to take the last bet. “All in” and “call” are actions that players can take in the game. In one game, players can take actions in turn in each stage except for show down stage. There are five stage in one game. 

The first stage in a game is “pre-flop”. Each player makes decision with the only information of their own holecards. The next stages consist of a series of three cards, called “flop”, later an additional single card, called “turn” and a final card called “river”. The game will stop at any time if there is only one player don’t quit yet.

After the four stages, if there is still up to 2 players remain battling, each player reveals their hand and seeks the best ranking-five card combination from seven cards consisting of the five community cards and their two holecards. This stage is called “show down”. The player that owns the best combination wins the whole chip pool.

Players have five possible actions: “check”, “call”, “raise”, “all in” or “fold”. “Check” is legal only when you have put the same bet with the current bet, “check” means doing nothing, just wait for the next card, or wait for others’ actions. “Call” means to raise the bet to make the total bet is equal to the current bet. “Raise” means to raise the bet to make the total bet greater than the current bet. “Raise” can give other players press to make them “fold”. “Fold” means to quit the current game, it is sometimes a wise choice when the player doesn’t have good hands. “All in” means to put all the chips you own to the pool, it can give the player the last chance to come back.

Texas Hold’em doesn’t always need show down. In fact, many games end early because many players abandon their hands to seek for the next games. 
 
### Installation
Linux (Arch, Manjaro, Debian 10+, Ubuntu 18.10+ or Fedora 27+)
For dependencies installing:

```bash -c "$(curl -fsSL https://raw.githubusercontent.com/fyne-io/bootstrap/master/bootstrap.sh)"```

You could also check https://github.com/fyne-io/fyne for more detailed downloading information for fyne package which is a UI toolkit and app API written in Go.

Mac OS X
For dependencies installing:

```/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
bash -c "$(curl -fsSL https://raw.githubusercontent.com/fyne-io/bootstrap/master/bootstrap.sh)"
```

You could also check https://github.com/fyne-io/fyne for more detailed downloading information for fyne package which is a UI toolkit and app API written in Go.

After installing the dependencies, you can install the Fyne toolkit and run our game using standard go commands.

```go get -u github.com/fyne-io/fyne
cd $GOPATH\src\CoolCoolProject\
go run main.go
```

You can play REAL FUN Texas Holdem now!

### User Interface Introduction

We use a go UI package, fine to build the whole game. After you opened our program, this would be the first interface you would see (Fig. 1). It contains the basic information of the game including number of players and bots, initial chips and number of turns which you all could adjust in the settings.

![alt text](https://raw.githubusercontent.com/Chengyang95/Texas-Holdem-Game/master/Picture1.png)

Fig. 1: The start window for our game

After you click the settings button, a new interface will appear (Fig. 2), and you can set the factors you want for the game here (Now, our game only supports maximum 6 players and 1 human player).

![alt text](https://raw.githubusercontent.com/Chengyang95/Texas-Holdem-Game/master/Picture2.png)\

Fig. 2: The setting window for our game

After setting all the parameters you want, you now could click “Play Game!” button. The game interface will appear like the above graph (Fig. 3). If you are the human player, the ‘Rosemary’ in the above graph would be you. When your turn is coming, you could click the underlying buttons to take your action. If you did not know what you should, you could click ‘Help me!’ button and the probability bot will help you make a decision.

![alt text](https://raw.githubusercontent.com/Chengyang95/Texas-Holdem-Game/master/Picture3.png)

Fig. 3: The main window for our game

### Five Bots
### Random Bot
When the bot takes action, it will randomly take an action from is legal action list.
Pseudocode:
```
RandomBot() {
	CheckLegalActions()
	RandomlyChooseOneAction()
	return action
}
```
### Probability bot
The probability bot follows the equation: 
C(p)=〖ap〗^2+bp+c
C:dream price
p: probability
```
Pseudocode:
If current bet >= max bet{
 	Choose fold/check
If current bet < max bet{
	Raise the random money/ Call (max bet-current bet)
}
```
### Persona Bot 
We would analyze the players’ hands during the showdown phase by TraceEvaluation(), which will generate alpha map by using how much bet the player initiated including raise and allin actions and the players hands’ winning probability, to analyze whether the player is featured with bluffing strategy or not. 
Using the generated alpha for each player to generate modified current bet, the max bet is the same as generated from the probability bot. For example, if the player bluffed, the alpha will be lower, the max modified current bet will smaller than the current bet. 
```
Pseudocode:
If max modified current bet >= max bet{
 	Choose fold/check
If max modified current bet < max bet{
	Raise the random money/ Call (max bet-current bet)
}
```
### Conventional Bot
The conventional bot follows three strategies to make the decision:
1. Firstly, the bot will read the winning possibility of every possible hands (two cards). It also has two thresholds one is 70% and another one is 90% which means if the winning possibility of his hands is above 70%, he will call and if above 90%, he will raise certain amount of money.  Here, we also blur boundaries using sigmoid function for assigning a random number.

2. The second rule is a practical formula which is call ‘four-two’ rule. It means if you want a specific card after flop phase, it has 4% possibility to appear in the river phase. However, if you want a specific card after turn phase, it will only has 2% possibility to appear in the river phase.

3. The third rule is based on winning money expectation to decide the action. The formula would be: win percentage x win money + lose percentage x (-lose money).
```
Pseudocode:
InitialConventionalBot() // This is for reading the hands data
ConventionalBot() {
	ReadInputStrategyFile()
	If the bot is in pre-flop phase {
		CheckRange() // Based on the first rule
		Return action
	} else {
		CheckEarnings() //Based on the second and the third rule
		Return action
	}
}
```
### Max Expectation Bot
The max expectation bot first finds all the legal actions. Then, with each action, bot will estimate other opponents’ actions after taking this action. The bot will consider the bet, the current bet, the chip pool, others’ chip and his seat position carefully. After the estimation other’s action, calculate the possible winning bet, possible losing bet and the expectation, following the equation: 
E(action)=p*estiWin+(1-p)*estiLose
```
Pseudocode:
MaxExpectationBot() {
GetIllegalAction()
AmplifyActionList() // This function is to split “raise” to “raise 20”, “raise 40” …
GetExpectation()
return max(action)
}
```

### Arena

We also upload our arena code which could be used to test the performance of each bot. The results would show the final winning money of the bots you chose.

You can run arena using the following standard go commands.

```
cd $GOPATH\src\arena\
go build
./arena 100 1000 5 0 1 1 1 1 1 (for example)
```

The second argument is the number of turns (100 here)
The third argument is the initial chips of each bot (1000 here)
The fourth argument is the number of bots (5 here)
The fifth argument is the number of human player (must be 0 because the arena does not support human player)
The sixth argument is the number of random bot(1 here)
The seventh argument is the number of probability bot(1 here)
The eighth argument is the number of conventional bot(1 here)
The ninth argument is the number of persona bot(1 here)
The tenth argument is the number of max expectation bot(1 here)

### Class Diagram
![alt text](https://raw.githubusercontent.com/Chengyang95/Texas-Holdem-Game/master/Picture4.png)
 

