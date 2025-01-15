package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Постоянные значения размеров поля
const (
	fieldH = 20 // Высота игрового поля
	fieldW = 60 // Ширина игрового поля
)

// Структура, представляющая параметры мяча
type ballOptions struct {
	x, y   int // Текущие координаты мяча
	dx, dy int // Скорость мяча по осям X и Y
}

// Структура, представляющая игрока
type Player struct {
	y int // Координата Y игрока (позиция на поле)
}

// Структура, представляющая счет игры
type Score struct {
	leftScore, rigthScore int // Счет левого и правого игрока
}

// clearScreen очищает консольный экран в зависимости от операционной системы.
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

// drawField рисует игровое поле, включая игроков, мяч и текущий счет.
func drawField(ballLocation ballOptions, leftPlayer Player, rigthPlayer Player, matchScore Score) {
	clearScreen() // Очищаем экран перед отрисовкой
	for i := 0; i <= fieldH+3; i++ {
		for j := 0; j <= fieldW; j++ {
			// Отрисовка границ и игроков
			if (i == 0 && j != 0) || (i == fieldH && j != 0) {
				fmt.Print("-") // Верхняя и нижняя границы
			} else if j == 3 {
				if i == leftPlayer.y || i == leftPlayer.y-1 || i == leftPlayer.y+1 {
					fmt.Print("X") // Левый игрок
				} else {
					fmt.Print(" ")
				}
			} else if j == fieldW-3 {
				if i == rigthPlayer.y || i == rigthPlayer.y-1 || i == rigthPlayer.y+1 {
					fmt.Print("X") // Правый игрок
				} else {
					fmt.Print(" ")
				}
			} else if i == ballLocation.x && j == ballLocation.y {
				fmt.Print("*") // Мяч
			} else {
				fmt.Print(" ")
			}
			// Отрисовка вертикальных границ
			if j == 0 && i <= fieldH {
				fmt.Print("|")
			}
			if j == fieldW && i <= fieldH {
				fmt.Print("|\n")
			}
			if j == fieldW && i > fieldH {
				fmt.Print("\n")
			}
			// Отображение счета
			if i == fieldH+2 && j == fieldW/2-15 {
				fmt.Print("Left:", matchScore.leftScore)
			}
			if i == fieldH+2 && j == fieldW/2+3 {
				fmt.Print("Rigth:", matchScore.rigthScore)
			}
			if i == fieldH+3 && j == fieldW/2-20 {
				fmt.Print("Управление: WS & KM для ракеток, Q - выход")
			}
		}
	}
}

// gamepad обрабатывает ввод пользователя и обновляет позицию игроков в зависимости от нажатых клавиш.
func gamepad(input string, leftPlayer *Player, rigthPlayer *Player) {
	switch input {
	case "w":
		leftPlayer.y-- // Движение вверх для левого игрока
		if leftPlayer.y == 1 {
			leftPlayer.y++ // Ограничение на верхнюю границу
		}
	case "s":
		leftPlayer.y++ // Движение вниз для левого игрока
		if leftPlayer.y == 19 {
			leftPlayer.y-- // Ограничение на нижнюю границу
		}
	case "k":
		rigthPlayer.y-- // Движение вверх для правого игрока
		if rigthPlayer.y == 1 {
			rigthPlayer.y++ // Ограничение на верхнюю границу
		}
	case "m":
		rigthPlayer.y++ // Движение вниз для правого игрока
		if rigthPlayer.y == 19 {
			rigthPlayer.y-- // Ограничение на нижнюю границу
		}
	default:
		break // Игнорировать другие вводы
	}
}

// ballMovement обновляет положение мяча на основе его скорости и проверяет столкновения с игроками и границами.
func ballMovement(ballSpeed *ballOptions, ballLocation *ballOptions, leftPlayer Player, rigthPlayer Player) {
	ballLocation.x += ballSpeed.dx // Обновление позиции по X
	ballLocation.y += ballSpeed.dy // Обновление позиции по Y

	// Проверка столкновения с верхней и нижней границами
	if ballLocation.x == 2 || ballLocation.x == fieldH-1 {
		ballSpeed.dx *= -1 // Изменение направления по X
	}
	// Проверка столкновения с левым игроком
	if ballLocation.y == 3 {
		if ballLocation.x == leftPlayer.y-1 || ballLocation.x == leftPlayer.y || ballLocation.x == leftPlayer.y+1 {
			ballSpeed.dy *= -1 // Изменение направления по Y
		}
	}
	// Проверка столкновения с правым игроком
	if ballLocation.y == fieldW-3 {
		if ballLocation.x == rigthPlayer.y-1 || ballLocation.x == rigthPlayer.y || ballLocation.x == rigthPlayer.y+1 {
			ballSpeed.dy *= -1 // Изменение направления по Y
		}
	}
}

// scorecheck проверяет, не вышел ли мяч за границы поля, и обновляет счет.
func scorecheck(ballLocation *ballOptions, ballSpeed *ballOptions, matchScore *Score) {
	// Проверка, забил ли левый игрок
	if ballLocation.y == 1 {
		matchScore.leftScore++
		ballLocation.y = fieldW / 2 // Сброс мяча в центр
		ballLocation.x = fieldH / 2
		ballSpeed.dy *= -1 // Изменение направления мяча
		ballSpeed.dx *= -1
	}
	// Проверка, забил ли правый игрок
	if ballLocation.y == fieldW-2 {
		matchScore.rigthScore++
		ballLocation.y = fieldW / 2 // Сброс мяча в центр
		ballLocation.x = fieldH / 2
		ballSpeed.dy *= -1 // Изменение направления мяча
		ballSpeed.dx *= -1
	}
}

// main является основной функцией, запускающей игру.
func main() {
	leftPlayer := Player{y: fieldH / 2}                       // Инициализация левого игрока
	rigthPlayer := Player{y: fieldH / 2}                      // Инициализация правого игрока
	ballLocation := ballOptions{x: fieldH / 2, y: fieldW / 2} // Инициализация мяча
	ballSpeed := ballOptions{dx: 1, dy: 1}                    // Инициализация скорости мяча
	matchScore := Score{leftScore: 0, rigthScore: 0}          // Инициализация счета

	winner := false // Переменная для отслеживания победителя
	for winner != true {
		drawField(ballLocation, leftPlayer, rigthPlayer, matchScore) // Отрисовка игрового поля
		var input string
		fmt.Scanln(&input) // Считывание ввода от пользователя
		if input == "q" || input == "Q" {
			clearScreen()                          // Очистка экрана перед выходом
			fmt.Println("Спасибо за игру! Удачи!") // Сообщение перед выходом
			winner = true
		}
		gamepad(input, &leftPlayer, &rigthPlayer)                        // Обработка ввода
		ballMovement(&ballSpeed, &ballLocation, leftPlayer, rigthPlayer) // Движение мяча
		scorecheck(&ballLocation, &ballSpeed, &matchScore)               // Проверка счета
		// Проверка условий победы
		if matchScore.leftScore == 11 {
			fmt.Println("Левый игрок победил!") // Объявление победителя
			winner = true
		} else if matchScore.rigthScore == 11 {
			fmt.Println("Правый игрок победил!") // Объявление победителя
			winner = true
		}
	}
}
