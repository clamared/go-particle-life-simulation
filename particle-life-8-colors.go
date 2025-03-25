package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth         = 640 * 3.3
	screenHeight        = 480 * 3
	numCircles          = 1000
	speedScale          = 0.5
	circleRadius        = 10
	drag                = float32(0.2)
	attractionThreshold = 300
	attractionForce     = 0.1
	restitution         = 0.85
	gravity             = 1
)

var colorInteractions = map[string]map[string]float64{
	"red": {
		"red":     -0.5,
		"blue":    2.0,
		"green":   -1.0,
		"yellow":  1.0,
		"purple":  -2.0,
		"cyan":    1.5,  // Moderate attraction to cyan
		"orange":  -1.5, // Moderate repulsion from orange
		"magenta": 0.5,  // Weak attraction to magenta
	},
	"blue": {
		"red":     -1.0, // Mild repulsion from red
		"blue":    -0.5,
		"green":   2.5,  // Very strong attraction to green
		"yellow":  -1.5, // Moderate repulsion from yellow
		"purple":  1.5,  // Moderate attraction to purple
		"cyan":    -2.0, // Strong repulsion from cyan
		"orange":  1.0,  // Mild attraction to orange
		"magenta": 0.0,  // Neutral interaction with magenta
	},
	"green": {
		"red":     1.0,  // Mild attraction to red
		"blue":    -1.0, // Mild repulsion from blue
		"green":   -0.5,
		"yellow":  2.0,  // Strong attraction to yellow
		"purple":  0.0,  // Neutral toward purple
		"cyan":    1.0,  // Mild attraction to cyan
		"orange":  -2.0, // Strong repulsion from orange
		"magenta": -1.5, // Moderate repulsion from magenta
	},
	"yellow": {
		"red":     2.0,  // Strong attraction to red
		"blue":    0.5,  // Weak attraction to blue
		"green":   -1.0, // Mild repulsion from green
		"yellow":  -0.5,
		"purple":  1.5,  // Moderate attraction to purple
		"cyan":    -1.0, // Mild repulsion from cyan
		"orange":  2.5,  // Very strong attraction to orange
		"magenta": -2.0, // Strong repulsion from magenta
	},
	"purple": {
		"red":     -2.0, // Strong repulsion from red
		"blue":    1.0,  // Mild attraction to blue
		"green":   0.5,  // Weak attraction to green
		"yellow":  -1.0, // Mild repulsion from yellow
		"purple":  -0.5,
		"cyan":    2.0,  // Strong attraction to cyan
		"orange":  -1.0, // Mild repulsion from orange
		"magenta": 1.5,  // Moderate attraction to magenta
	},
	"cyan": {
		"red":     1.5,  // Moderate attraction to red
		"blue":    -2.0, // Strong repulsion from blue
		"green":   1.0,  // Mild attraction to green
		"yellow":  0.0,  // Neutral toward yellow
		"purple":  -1.0, // Mild repulsion from purple
		"cyan":    -0.5,
		"orange":  2.0,  // Strong attraction to orange
		"magenta": -1.0, // Mild repulsion from magenta
	},
	"orange": {
		"red":     -1.5, // Moderate repulsion from red
		"blue":    1.0,  // Mild attraction to blue
		"green":   -2.0, // Strong repulsion from green
		"yellow":  2.0,  // Strong attraction to yellow
		"purple":  -1.0, // Mild repulsion from purple
		"cyan":    1.5,  // Moderate attraction to cyan
		"orange":  -0.5,
		"magenta": 2.5, // Very strong attraction to magenta
	},
	"magenta": {
		"red":     0.5,  // Weak attraction to red
		"blue":    0.0,  // Neutral toward blue
		"green":   -1.5, // Moderate repulsion from green
		"yellow":  -2.0, // Strong repulsion from yellow
		"purple":  2.0,  // Strong attraction to purple
		"cyan":    -1.0, // Mild repulsion from cyan
		"orange":  1.5,  // Mild attraction to orange
		"magenta": -0.5,
	},
}

type Circle struct {
	X         float32
	Y         float32
	Radius    float32
	Color     color.RGBA
	ColorType string 
	VelX      float32
	VelY      float32
	Mass      float32
}

type Game struct {
	circles []Circle
}

func (g *Game) Update() error {
	for i := range g.circles {
		g.circles[i].VelX *= (1 - float32(drag)/100)
		g.circles[i].VelY = g.circles[i].VelY*(1-float32(drag)/100) + gravity

		g.circles[i].X += g.circles[i].VelX * speedScale
		g.circles[i].Y += g.circles[i].VelY * speedScale
	}

	for i := 0; i < len(g.circles); i++ {
		for j := i + 1; j < len(g.circles); j++ {
			dx := g.circles[j].X - g.circles[i].X
			dy := g.circles[j].Y - g.circles[i].Y
			distSq := dx*dx + dy*dy
			distance := float32(math.Sqrt(float64(distSq)))

			if distance == 0 {
				distance = 0.1
				dx = 0.1
				dy = 0
			}

			nx := dx / distance
			ny := dy / distance

			minDist := g.circles[i].Radius + g.circles[j].Radius

			if distSq < minDist*minDist {
				resolveCollision(&g.circles[i], &g.circles[j])
			} else if distance < attractionThreshold {
				interactionStrength := 0.0

				if interaction, exists := colorInteractions[g.circles[i].ColorType][g.circles[j].ColorType]; exists {
					interactionStrength = interaction
				}

				forceStrength := float32(interactionStrength) * attractionForce * (1 - distance/attractionThreshold)

				g.circles[i].VelX += nx * forceStrength
				g.circles[i].VelY += ny * forceStrength
				g.circles[j].VelX -= nx * forceStrength
				g.circles[j].VelY -= ny * forceStrength
			}
		}
	}

	for i := range g.circles {
		if g.circles[i].X < g.circles[i].Radius {
			g.circles[i].X = g.circles[i].Radius
			g.circles[i].VelX = -g.circles[i].VelX * restitution
		}
		if g.circles[i].X > screenWidth-g.circles[i].Radius {
			g.circles[i].X = screenWidth - g.circles[i].Radius
			g.circles[i].VelX = -g.circles[i].VelX * restitution
		}
		if g.circles[i].Y < g.circles[i].Radius {
			g.circles[i].Y = g.circles[i].Radius
			g.circles[i].VelY = -g.circles[i].VelY * restitution
		}
		if g.circles[i].Y > screenHeight-g.circles[i].Radius {
			g.circles[i].Y = screenHeight - g.circles[i].Radius
			g.circles[i].VelY = -g.circles[i].VelY * restitution
		}
	}

	return nil
}

func resolveCollision(c1, c2 *Circle) {
	dx := c2.X - c1.X
	dy := c2.Y - c1.Y
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if distance == 0 {
		dx = 1.0
		dy = 0.0
		distance = 1.0
	}

	nx := dx / distance
	ny := dy / distance

	overlap := c1.Radius + c2.Radius - distance

	totalMass := c1.Mass + c2.Mass
	c1MassRatio := c2.Mass / totalMass
	c2MassRatio := c1.Mass / totalMass

	c1.X -= nx * overlap * c1MassRatio
	c1.Y -= ny * overlap * c1MassRatio
	c2.X += nx * overlap * c2MassRatio
	c2.Y += ny * overlap * c2MassRatio

	rvx := c2.VelX - c1.VelX
	rvy := c2.VelY - c1.VelY

	velAlongNormal := rvx*nx + rvy*ny

	if velAlongNormal > 0 {
		return
	}

	j := -(1 + restitution) * velAlongNormal
	j /= 1/c1.Mass + 1/c2.Mass

	c1.VelX -= j * nx / c1.Mass
	c1.VelY -= j * ny / c1.Mass
	c2.VelX += j * nx / c2.Mass
	c2.VelY += j * ny / c2.Mass

	randAngle := rand.Float32() * 2 * math.Pi
	perturbation := float32(0.01)
	c1.VelX += perturbation * float32(math.Cos(float64(randAngle)))
	c1.VelY += perturbation * float32(math.Sin(float64(randAngle)))
	c2.VelX -= perturbation * float32(math.Cos(float64(randAngle)))
	c2.VelY -= perturbation * float32(math.Sin(float64(randAngle)))
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, circle := range g.circles {
		vector.DrawFilledCircle(screen, circle.X, circle.Y, circle.Radius, circle.Color, true)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Number of Circles: %d", len(g.circles)))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Particle System - Ethan Annane (:")

	rand.Seed(time.Now().UnixNano())

	circleColors := []color.RGBA{
		{R: 255, G: 0, B: 0, A: 255},   // Red
		{R: 0, G: 0, B: 255, A: 255},   // Blue
		{R: 0, G: 255, B: 0, A: 255},   // Green
		{R: 255, G: 255, B: 0, A: 255}, // Yellow
		{R: 128, G: 0, B: 128, A: 255}, // Purple
		{R: 0, G: 255, B: 255, A: 255}, // Cyan
		{R: 255, G: 165, B: 0, A: 255}, // Orange
		{R: 255, G: 0, B: 255, A: 255}, // Magenta
	}

	colorNames := []string{"red", "blue", "green", "yellow", "purple", "cyan", "orange", "magenta"}

	circles := make([]Circle, numCircles)
	for i := range circles {
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float32()*5 + 2
		min := 5
		max := 15

		colorIndex := rand.Intn(len(circleColors))
		circleColor := circleColors[colorIndex]
		colorType := colorNames[colorIndex]

		radius := float32(min + rand.Intn(max-min+1))

		x := radius + rand.Float32()*(screenWidth-2*radius)
		y := radius + rand.Float32()*(screenHeight-2*radius)

		circles[i] = Circle{
			X:         x,
			Y:         y,
			Radius:    radius,
			Color:     circleColor,
			ColorType: colorType,
			VelX:      float32(math.Cos(angle) * float64(speed)),
			VelY:      float32(math.Sin(angle) * float64(speed)),
			Mass:      radius * radius,
		}
	}

	for i := 0; i < len(circles); i++ {
		for j := 0; j < i; j++ {
			dx := circles[i].X - circles[j].X
			dy := circles[i].Y - circles[j].Y
			distSq := dx*dx + dy*dy
			minDist := circles[i].Radius + circles[j].Radius

			if distSq < minDist*minDist {
				dist := float32(math.Sqrt(float64(distSq)))
				nx := dx / dist
				ny := dy / dist
				overlap := minDist - dist + 1.0

				circles[i].X += nx * overlap / 2
				circles[i].Y += ny * overlap / 2
				circles[j].X -= nx * overlap / 2
				circles[j].Y -= ny * overlap / 2
			}
		}
	}

	game := &Game{
		circles: circles,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
