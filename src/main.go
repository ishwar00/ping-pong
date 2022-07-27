package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/ishwar00/simple-ping-pong-game/fonts"
	"github.com/ishwar00/simple-ping-pong-game/shape"
	"github.com/ishwar00/simple-ping-pong-game/utils"
)

func init() {
	game.began = false
	game.durationOfGame = 7200 // 7200 frames is quivalent to 2 minutes
	game.ball.Init(color.NRGBA{243, 120, 120, 255})
	game.leftBlock = shape.Block{
		Position:    utils.Vector2d{X: float32(utils.ScreenWidth) * 0.03, Y: 50},
		Color:       color.NRGBA{246, 70, 104, 255},
		Sensitivity: 3,
	}

	game.rightBlock = shape.Block{
		Position:    utils.Vector2d{X: float32(utils.ScreenWidth) * 0.97, Y: 50},
		Color:       color.NRGBA{254, 150, 119, 255},
		Sensitivity: 3,
	}

	game.obstacleBlock = shape.Block{
		Color:    color.NRGBA{135, 128, 94, 255},
		Position: utils.Vector2d{X: 50, Y: 50},
		Width:    80,
		Height:   80,
	}
}

type Game struct {
	ball               shape.Ball
	leftBlock          shape.Block
	rightBlock         shape.Block
	obstacleBlock      shape.Block
	began              bool
	showRightCollision uint8
	showLeftCollision  uint8 // it's kind a duration of red splash on right or left wall
	leftPlayerScore    int
	rightPlayerScore   int
	durationOfGame     int
}

var game Game

func (g *Game) Update() error {
	utils.ScreenWidth, utils.ScreenHeight = ebiten.WindowSize()
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.began {
		g.began = true
	}
	if g.began && g.durationOfGame > 0 {
		if g.durationOfGame%1200 == 0 { // every 20 seconds
			g.ball.IncreaseVelocity(2)
			g.leftBlock.IncreaseSensitivity(3)
			g.rightBlock.IncreaseSensitivity(3)
		}
		g.obstacleBlock.Width = utils.ScreenWidth / 3
		g.obstacleBlock.Height = utils.ScreenHeight / 90
		g.obstacleBlock.Position.X = float32(utils.ScreenWidth-g.obstacleBlock.Width) / 2
		g.obstacleBlock.Position.Y = float32(utils.ScreenHeight-g.obstacleBlock.Height) / 2

		g.leftBlock.Width = (utils.ScreenWidth * 1) / 130
		g.leftBlock.Height = (utils.ScreenHeight * 25) / 100

		g.rightBlock.Width = (utils.ScreenWidth * 1) / 130
		g.rightBlock.Height = (utils.ScreenHeight * 25) / 100

		g.leftBlock.Position.X = float32(utils.ScreenWidth) * 0.01
		g.rightBlock.Position.X = float32(utils.ScreenWidth+g.rightBlock.Width) * 0.9778

		g.rightBlock.HandleKeyPress(ebiten.KeyArrowUp, ebiten.KeyArrowDown)
		g.leftBlock.HandleKeyPress(ebiten.KeyW, ebiten.KeyS)
		leftCollision, rightCollision := g.HandleCollision()

		if leftCollision {
			g.showLeftCollision = 40
			g.rightPlayerScore++
		}

		if rightCollision {
			g.showRightCollision = 40
			g.leftPlayerScore++
		}

		g.ball.Position.X += g.ball.Velocity.X
		g.ball.Position.Y += g.ball.Velocity.Y
		g.durationOfGame--
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{217, 248, 196, 255})
	if g.began && g.durationOfGame > 0 {
		leftScore := fmt.Sprintf("%v", g.leftPlayerScore)
		text.Draw(
			screen,
			leftScore,
			fonts.GoMonoFace[2],
			(utils.ScreenWidth*40)/100,
			45,
			color.NRGBA{246, 70, 104, 255},
		)

		rightScore := fmt.Sprintf("%v", g.rightPlayerScore)
		text.Draw(
			screen,
			rightScore,
			fonts.GoMonoFace[2],
			(utils.ScreenWidth*55)/100,
			45,
			color.NRGBA{254, 150, 119, 255},
		)

		timeLeft := fmt.Sprintf("%v seconds left", g.durationOfGame/60)
		text.Draw(
			screen,
			timeLeft,
			fonts.GoMonoFace[1],
			(utils.ScreenWidth*40)/100,
			utils.ScreenHeight-15,
			color.NRGBA{135, 100, 69, 255},
		)

		g.ball.Draw(screen)
		g.leftBlock.Draw(screen)
		g.rightBlock.Draw(screen)
		g.obstacleBlock.Draw(screen)

		if g.showLeftCollision > 0 {
			g.ShowCollision(screen, 0, 0, g.showLeftCollision)
			g.showLeftCollision--
		}

		if g.showRightCollision > 0 {
			g.ShowCollision(screen, float64(utils.ScreenWidth)*0.991, 0, g.showRightCollision)
			g.showRightCollision--
		}
	} else if !g.began {
		text.Draw(
			screen,
			"PRESS SPACE TO BEGIN",
			fonts.GoMonoFace[2],
			(utils.ScreenWidth*35)/100,
			utils.ScreenHeight/2,
			color.NRGBA{135, 100, 69, 255},
		)
	}

	if g.durationOfGame <= 0 {
		g.PrintTheResult(screen)
	}
}

func (g *Game) PrintTheResult(screen *ebiten.Image) {
	if g.leftPlayerScore < g.rightPlayerScore {
		text.Draw(
			screen,
			"right player has won!!",
			fonts.GoMonoFace[2],
			(utils.ScreenWidth*35)/100,
			utils.ScreenHeight/2,
			color.NRGBA{243, 120, 120, 255},
		)
	} else if g.leftPlayerScore > g.rightPlayerScore {
		text.Draw(
			screen,
			"left player has won!!",
			fonts.GoMonoFace[2],
			(utils.ScreenWidth*35)/100,
			utils.ScreenHeight/2,
			color.NRGBA{243, 120, 120, 255},
		)
	} else {
		text.Draw(
			screen,
			"Huh, neither of you won :/",
			fonts.GoMonoFace[2],
			(utils.ScreenWidth*30)/100,
			utils.ScreenHeight/2,
			color.NRGBA{255, 112, 119, 255},
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) ShowCollision(screen *ebiten.Image, x, y float64, showCollision uint8) {
	ebitenutil.DrawRect(
		screen,
		x,
		y,
		float64(utils.ScreenWidth)*0.01,
		float64(utils.ScreenHeight),
		color.NRGBA{245, 72, 81, 255 - (40-showCollision)*6},
	)
}

func (g *Game) HandleCollision() (bool, bool) {
	BallBlockCollision(&g.ball, g.leftBlock)
	BallBlockCollision(&g.ball, g.rightBlock)
	BallBlockCollision(&g.ball, g.obstacleBlock)
	return g.ball.HandleCollisionWithWalls()
}

func BallBlockCollision(ball *shape.Ball, block shape.Block) {
	collisionPoint := ball.Position

	if ball.Position.X <= block.Position.X {
		collisionPoint.X = block.Position.X // left edge
	} else if ball.Position.X >= block.Position.X+float32(block.Width) {
		collisionPoint.X = block.Position.X + float32(block.Width) // right edge
	}

	if ball.Position.Y <= block.Position.Y {
		collisionPoint.Y = block.Position.Y // top edge
	} else if ball.Position.Y >= block.Position.Y+float32(block.Height) {
		collisionPoint.Y = block.Position.Y + float32(block.Height) // bottom edge
	}

	y := math.Abs(float64(ball.Position.Y - collisionPoint.Y))
	x := math.Abs(float64(ball.Position.X - collisionPoint.X))
	distance := math.Hypot(x, y)

	horizontalCollision := func() {
		if (ball.Position.X <= block.Position.X && ball.Velocity.X > 0) ||
			ball.Position.X >= (block.Position.X+float32(block.Width)) && ball.Velocity.X < 0 {
			// O->||<-O
			ball.Velocity.X *= -1
		}
	}

	verticalCollision := func() {
		if ball.Position.Y <= block.Position.Y && ball.Velocity.Y > 0 ||
			ball.Position.Y >= (block.Position.Y+float32(block.Height)) && ball.Velocity.Y < 0 {
			ball.Velocity.Y *= -1
		}
	}

	if distance <= float64(ball.Radius) {
		if x > 0 && y > 0 {
			verticalCollision()
			horizontalCollision()
		} else if x <= float64(ball.Radius) && y == 0 {
			horizontalCollision()
		} else if y <= float64(ball.Radius) && x == 0 {
			verticalCollision()
		}
	}
}

func main() {
	ebiten.SetWindowSize(utils.ScreenWidth, utils.ScreenHeight)
	ebiten.SetWindowTitle("Hello there!!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
