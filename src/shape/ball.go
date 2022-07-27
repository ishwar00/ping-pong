package shape

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ishwar00/simple-ping-pong-game/utils"
)

type Ball struct {
	textureImg *ebiten.Image
	Radius     int
	color      color.Color
	Position   utils.Vector2d
	Velocity   utils.Vector2d
}

func (b *Ball) Init(clr color.Color) {
	img := ebiten.NewImage(3, 3)
	img.Fill(clr)

	*b = Ball{
		textureImg: img,
		Position: utils.Vector2d{
			X: float32(utils.ScreenHeight) / 2,
			Y: float32(utils.ScreenHeight) / 2,
		},
		color: clr,
		Velocity: utils.Vector2d{
			X: utils.RandSign() * 6,
			Y: utils.RandSign() * 6,
		},
	}
}

func (b *Ball) HandleCollisionWithWalls() (bool, bool) {
	leftWall := b.Position.X <= float32(b.Radius) && b.Velocity.X < 0
	rightWall := (float32(utils.ScreenWidth)-b.Position.X) <= float32(b.Radius) && b.Velocity.X > 0
	if leftWall || rightWall {
		b.Velocity.X *= -1
		return leftWall, rightWall
	}

	topWall := b.Position.Y <= float32(b.Radius) && b.Velocity.Y < 0
	bottomWall := (float32(utils.ScreenHeight)-b.Position.Y) <= float32(b.Radius) && b.Velocity.Y > 0
	if topWall || bottomWall {
		b.Velocity.Y *= -1
	}
	return false, false
}

func (b *Ball) IncreaseVelocity(change float32) {
	b.Velocity.X += change * (b.Velocity.X / float32(math.Abs(float64(b.Velocity.X))))
	b.Velocity.Y += change * (b.Velocity.Y / float32(math.Abs(float64(b.Velocity.Y))))
}

func (b *Ball) MoveTo(x, y float32) {
	b.Position.X, b.Position.Y = x, y
}

func (b *Ball) Recolor(clr color.Color) {
	b.textureImg.Fill(clr)
}

func (b *Ball) Draw(screen *ebiten.Image) {
	b.Radius = (utils.ScreenHeight * 4) / 100
	var path vector.Path
	path.Arc(b.Position.X, b.Position.Y, float32(b.Radius), 0, 2*math.Pi, vector.Clockwise)

	R, G, B, A := b.color.RGBA()
	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vertices {
		v := &vertices[i]
		v.SrcX = 1
		v.SrcY = 1
		v.ColorA = float32(A) / 0xffff
		v.ColorB = float32(B) / 0xffff
		v.ColorG = float32(G) / 0xffff
		v.ColorR = float32(R) / 0xffff
	}
	screen.DrawTriangles(vertices, indices, b.textureImg, nil)
}
