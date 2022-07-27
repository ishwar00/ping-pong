package shape

import (
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/ishwar00/rigidBodyDynamics/utils"
)

type Block struct {
    Position    utils.Vector2d
    Color       color.Color
    Width       int
    Height      int
    Sensitivity float32 
}

func (b *Block) Draw(screen *ebiten.Image) {
    ebitenutil.DrawRect(
        screen, 
        float64(b.Position.X), 
        float64(b.Position.Y), 
        float64(b.Width), 
        float64(b.Height), 
        b.Color,
    )
}

func (b* Block) HandleKeyPress(upKey ebiten.Key, downKey ebiten.Key) {
    if ebiten.IsKeyPressed(downKey) {
        if (b.Position.Y + 2 * b.Sensitivity + float32(b.Height)) < float32(utils.ScreenHeight) {
            b.Position.Y += b.Sensitivity
        }
    }

    if ebiten.IsKeyPressed(upKey) {
        if b.Position.Y + 2 * b.Sensitivity > 0 {
            b.Position.Y -= b.Sensitivity
        }
    }
}

func (b *Block) IncreaseSensitivity(change float32) {
    b.Sensitivity += change
}
