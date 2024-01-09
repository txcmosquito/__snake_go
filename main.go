package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

// Direction represents the movement direction of the snake.
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

// Point represents a 2D coordinate.
type Point struct {
	X, Y int
}

// Snake represents the snake in the game.
type Snake struct {
	Body      []Point
	Direction Direction
}

// Game represents the game state.
type Game struct {
	Screen   tcell.Screen
	Snake    Snake
	Food     Point
	Score    int
	Speed    time.Duration
	GameOver bool
}

func (g *Game) initialize() {
	g.Snake = Snake{
		Body:      []Point{{5, 5}, {5, 6}, {5, 7}},
		Direction: Right,
	}

	g.spawnFood()
	g.Score = 0
	g.Speed = 100 * time.Millisecond
	g.GameOver = false
}

func (g *Game) spawnFood() {
	width, height := g.Screen.Size()
	g.Food = Point{
		X: rand.Intn(width),
		Y: rand.Intn(height),
	}
}

func (g *Game) draw() {
	g.Screen.Clear()

	for _, p := range g.Snake.Body {
		g.Screen.SetContent(p.X, p.Y, 'O', nil, tcell.StyleDefault)
	}

	g.Screen.SetContent(g.Food.X, g.Food.Y, 'F', nil, tcell.StyleDefault)

	scoreText := fmt.Sprintf("Score: %d", g.Score)
	g.drawText(1, 0, scoreText)

	g.Screen.Show()
}

func (g *Game) drawText(x, y int, text string) {
	for i, char := range text {
		style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
		g.Screen.SetContent(x+i, y, char, nil, style)
	}
}


func (g *Game) handleEvents() {
	ev := g.Screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyUp:
			g.Snake.Direction = Up
		case tcell.KeyDown:
			g.Snake.Direction = Down
		case tcell.KeyLeft:
			g.Snake.Direction = Left
		case tcell.KeyRight:
			g.Snake.Direction = Right
		case tcell.KeyCtrlC, tcell.KeyEsc:
			g.Screen.Fini()
			fmt.Println("Game Over!")
			os.Exit(0)
		case tcell.KeyRune:
			if ev.Rune() == 'R' || ev.Rune() == 'r' {
				// Allow restarting the game at any time
				g.initialize()
			}
		}
	}
}


func (g *Game) move() {
    if g.GameOver {
        return
    }

    head := g.Snake.Body[0]

    var newHead Point
    switch g.Snake.Direction {
    case Up:
        newHead = Point{head.X, head.Y - 1}
    case Down:
        newHead = Point{head.X, head.Y + 1}
    case Left:
        newHead = Point{head.X - 1, head.Y}
    case Right:
        newHead = Point{head.X + 1, head.Y}
    }

    // Check for collisions
    width, height := g.Screen.Size()
    if newHead.X < 0 || newHead.X >= width || newHead.Y < 0 || newHead.Y >= height {
        g.GameOver = true
        return
    }

    for _, segment := range g.Snake.Body[1:] {
        if newHead.X == segment.X && newHead.Y == segment.Y {
            g.GameOver = true
            return
        }
    }

    // Check for collision with food
    if newHead.X == g.Food.X && newHead.Y == g.Food.Y {
        g.spawnFood()
        g.Score++
        g.Speed -= 5 * time.Millisecond // Increase speed
    } else {
        // Remove the last tail segment
        g.Snake.Body = g.Snake.Body[:len(g.Snake.Body)-1]
    }

    g.Snake.Body = append([]Point{newHead}, g.Snake.Body...)
}


func main() {
	// Initialize tcell screen
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Println("Error initializing screen:", err)
		return
	}
	if err := screen.Init(); err != nil {
		fmt.Println("Error initializing screen:", err)
		return
	}
	defer screen.Fini()

	// Set up the game
	game := &Game{Screen: screen}
	game.initialize()

	// Start a goroutine for handling keyboard events
	go func() {
		for !game.GameOver {
			game.handleEvents()
			time.Sleep(10 * time.Millisecond) // Small delay to avoid busy-waiting
		}
	}()

	// Main game loop
	for !game.GameOver {
		game.move()
		game.draw()

		time.Sleep(game.Speed)
	}
}
