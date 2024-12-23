package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

/*
очистка поля
отрисовка поля
упраление жойстиком
отрисовка мяча
*/

const (
	fieldH = 20
	fieldW = 60
)

type ballOptions struct {
	x, y   int
	dx, dy int
}

type Player struct {
	y int
}

type Score struct {
	leftScore, rigthScore int
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls") // Для Windows
	} else {
		cmd = exec.Command("clear") // Для Unix-подобных систем
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func drawField(ballLocation ballOptions, leftPlayer Player, rigthPlayer Player, matchScore Score) {

	clearScreen()
	for i := 0; i <= fieldH+3; i++ {
		for j := 0; j <= fieldW; j++ {
			if (i == 0 && j != 0) || (i == fieldH && j != 0) {
				fmt.Print("-")
			} else if j == 3 {
				if i == leftPlayer.y || i == leftPlayer.y-1 || i == leftPlayer.y+1 {
					fmt.Print("X")
				} else {
					fmt.Print(" ")
				}
			} else if j == fieldW-3 {
				if i == rigthPlayer.y || i == rigthPlayer.y-1 || i == rigthPlayer.y+1 {
					fmt.Print("X")
				} else {
					fmt.Print(" ")
				}
			} else if i == ballLocation.x && j == ballLocation.y {
				fmt.Print("*")

			} else {
				fmt.Print(" ")
			}
			if j == 0 && i <= fieldH {
				fmt.Print("|")
			}
			if j == fieldW && i <= fieldH {
				fmt.Print("|\n")
			}
			if j == fieldW && i > fieldH {
				fmt.Print("\n")
			}
			if i == fieldH+2 && j == fieldW/2-15 {
				fmt.Print("Left:", matchScore.leftScore)
			}
			if i == fieldH+2 && j == fieldW/2+3 {
				fmt.Print("Rigth:", matchScore.rigthScore)
			}
			if i == fieldH+3 && j == fieldW/2-20 {
				fmt.Print("Управление: WS & KM для ракеток, q - выход")
			}

		}
	}
}

func gamepad(input string, leftPlayer *Player, rigthPlayer *Player) {
	switch input {
	case "w":
		leftPlayer.y--
		if leftPlayer.y == 1 {
			leftPlayer.y++
		}
	case "s":
		leftPlayer.y++
		if leftPlayer.y == 19 {
			leftPlayer.y--
		}
	case "k":
		rigthPlayer.y--
		if rigthPlayer.y == 1 {
			rigthPlayer.y--
		}
	case "m":
		rigthPlayer.y++
		if rigthPlayer.y == 19 {
			rigthPlayer.y--
		}

	default:
		break
	}
}

func ballMovement(ballSpeed *ballOptions, ballLocation *ballOptions, leftPlayer Player, rigthPlayer Player) {
	ballLocation.x += ballSpeed.dx
	ballLocation.y += ballSpeed.dy

	if ballLocation.x == 2 || ballLocation.x == fieldH-1 {
		ballSpeed.dx *= -1
	}
	if ballLocation.y == 3 {
		if ballLocation.x == leftPlayer.y-1 || ballLocation.x == leftPlayer.y || ballLocation.x == leftPlayer.y+1 {
			ballSpeed.dy *= -1
		}
	}

	if ballLocation.y == fieldW-3 {
		if ballLocation.x == rigthPlayer.y-1 || ballLocation.x == rigthPlayer.y || ballLocation.x == rigthPlayer.y+1 {
			ballSpeed.dy *= -1
		}
	}
}

func scorecheck(ballLocation *ballOptions, ballSpeed *ballOptions, matchScore *Score) {
	if ballLocation.y == 1 {
		matchScore.leftScore++
		ballLocation.y = fieldW / 2
		ballLocation.x = fieldH / 2
		ballSpeed.dy *= -1
		ballSpeed.dx *= -1
	}
	if ballLocation.y == fieldW-2 {
		matchScore.rigthScore++
		ballLocation.y = fieldW / 2
		ballLocation.x = fieldH / 2
		ballSpeed.dy *= -1
		ballSpeed.dx *= -1
	}
}

func main() {
	leftPlayer := Player{y: fieldH / 2}
	rigthPlayer := Player{y: fieldH / 2}
	ballLocation := ballOptions{x: fieldH / 2, y: fieldW / 2}
	ballSpeed := ballOptions{dx: 1, dy: 1}
	matchScore := Score{leftScore: 0, rigthScore: 0}

	winner := false
	for winner != true {
		drawField(ballLocation, leftPlayer, rigthPlayer, matchScore)
		var input string
		fmt.Scanln(&input)
		if input == "q" || input == "Q" {
			clearScreen()
			fmt.Println("Lol bye loozers")
			winner = true
		}
		gamepad(input, &leftPlayer, &rigthPlayer)
		ballMovement(&ballSpeed, &ballLocation, leftPlayer, rigthPlayer)
		scorecheck(&ballLocation, &ballSpeed, &matchScore)
	}
}
