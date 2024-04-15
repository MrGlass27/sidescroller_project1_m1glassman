package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type scrollDemo struct {
	player          playerSprite
	background      *ebiten.Image
	backgroundYView int
	projectiles     []*projectile
	enemies         []*enemyShip
	enemySpawnTimer int
	shootCooldown   int
	score           int
}

type playerSprite struct {
	sprite *ebiten.Image
	xLoc   int
	yLoc   int
	width  int
	height int
}

type projectile struct {
	sprite *ebiten.Image
	xLoc   int
	yLoc   int
	speed  int
	active bool
}

type enemyShip struct {
	sprite *ebiten.Image
	xLoc   int
	yLoc   int
	speed  int
	active bool
}

func (demo *scrollDemo) Update() error {
	const (
		playerSpeed     = 5
		minX            = 0
		maxX            = 436
		projectileSpeed = -10
	)

	// Move player horizontally
	if demo.score < 30 {
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			demo.player.xLoc -= playerSpeed
			if demo.player.xLoc < minX {
				demo.player.xLoc = minX
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			demo.player.xLoc += playerSpeed
			if demo.player.xLoc > maxX {
				demo.player.xLoc = maxX
			}
		}
	}

	// Shoot projectile
	if demo.score < 30 {
		if ebiten.IsKeyPressed(ebiten.KeySpace) && demo.shootCooldown == 0 {
			projectilePict, _, err := ebitenutil.NewImageFromFile("28.png")
			if err != nil {
				fmt.Println("Unable to load projectile sprite image:", err)
				return err
			}

			// Adjust the projectile's initial position to be at the center of the player
			projectileX := demo.player.xLoc + (demo.player.width / 2)
			projectileY := demo.player.yLoc

			demo.projectiles = append(demo.projectiles, &projectile{
				sprite: projectilePict,
				xLoc:   projectileX,
				yLoc:   projectileY,
				speed:  projectileSpeed,
				active: true,
			})
			// Reset the cooldown timer
			demo.shootCooldown = 15 // Adjust this value to set the desired cooldown duration
		}
	}

	// Decrement the cooldown timer if it's greater than 0
	if demo.shootCooldown > 0 {
		demo.shootCooldown--
	}

	// Update projectiles
	for i := 0; i < len(demo.projectiles); i++ {
		p := demo.projectiles[i]
		p.yLoc += p.speed

		// Remove the projectile if it goes off-screen
		if p.yLoc < 0 {
			demo.projectiles = append(demo.projectiles[:i], demo.projectiles[i+1:]...)
			i--
		}
	}
	// Spawn enemy ships
	if demo.score < 30 {
		demo.enemySpawnTimer++
		if demo.enemySpawnTimer >= 200 {
			demo.enemySpawnTimer = 0
			enemySprite, _, err := ebitenutil.NewImageFromFile("Spaceship_05_Orange.png")
			if err != nil {
				fmt.Println("Unable to load enemy ship sprite image:", err)
				return err
			}
			enemyX := rand.Intn(436)
			demo.enemies = append(demo.enemies, &enemyShip{
				sprite: enemySprite,
				xLoc:   enemyX,
				yLoc:   -64,
				speed:  2,
				active: true,
			})
		}
	}

	// Update enemy ships
	for i := 0; i < len(demo.enemies); i++ {
		e := demo.enemies[i]
		e.yLoc += e.speed

		// Remove the enemy ship if it goes off-screen
		if e.yLoc > 500 {
			demo.enemies = append(demo.enemies[:i], demo.enemies[i+1:]...)
			i--
		}
	}

	// Check for collisions between projectiles and enemy ships
	for i := 0; i < len(demo.projectiles); i++ {
		p := demo.projectiles[i]
		for j := 0; j < len(demo.enemies); j++ {
			e := demo.enemies[j]
			if demo.checkCollision(p, e) {
				// Remove the enemy ship
				demo.enemies = append(demo.enemies[:j], demo.enemies[j+1:]...)
				j--

				// Remove the projectile
				demo.projectiles = append(demo.projectiles[:i], demo.projectiles[i+1:]...)
				i--

				// Increment score
				demo.score++
				break
			}
		}
	}
	// Update background scrolling
	backgroundHeight := demo.background.Bounds().Dy()
	maxY := backgroundHeight * 2
	demo.backgroundYView += 4
	demo.backgroundYView %= maxY

	return nil
}

func (demo *scrollDemo) checkCollision(p *projectile, e *enemyShip) bool {
	// Calculate the bounding boxes for the projectile and the enemy ship
	pLeft := p.xLoc
	pRight := p.xLoc + p.sprite.Bounds().Dx()
	pTop := p.yLoc
	pBottom := p.yLoc + p.sprite.Bounds().Dy()

	eLeft := e.xLoc
	eRight := e.xLoc + e.sprite.Bounds().Dx()
	eTop := e.yLoc
	eBottom := e.yLoc + e.sprite.Bounds().Dy()

	// Check if the bounding boxes overlap
	if pLeft < eRight && pRight > eLeft && pTop < eBottom && pBottom > eTop {
		return true
	}
	return false
}

func (demo *scrollDemo) Draw(screen *ebiten.Image) {
	// Clear the screen
	screen.Fill(color.White)

	// Draw background
	const repeat = 3
	backgroundHeight := demo.background.Bounds().Dy()
	for count := 0; count < repeat; count++ {
		drawOps := ebiten.DrawImageOptions{}
		drawOps.GeoM.Translate(0, float64(demo.backgroundYView%backgroundHeight)-float64(count)*float64(backgroundHeight))
		screen.DrawImage(demo.background, &drawOps)
	}

	// Draw player sprite
	playerDrawOps := ebiten.DrawImageOptions{}
	playerDrawOps.GeoM.Translate(float64(demo.player.xLoc), float64(demo.player.yLoc))
	screen.DrawImage(demo.player.sprite, &playerDrawOps)

	// Draw projectiles
	for _, p := range demo.projectiles {
		projectileDrawOps := ebiten.DrawImageOptions{}
		projectileDrawOps.GeoM.Translate(float64(p.xLoc), float64(p.yLoc))
		screen.DrawImage(p.sprite, &projectileDrawOps)
	}
	// Draw enemy ships
	for _, e := range demo.enemies {
		enemyDrawOps := ebiten.DrawImageOptions{}
		enemyDrawOps.GeoM.Translate(float64(e.xLoc), float64(e.yLoc))
		screen.DrawImage(e.sprite, &enemyDrawOps)
	}
	if demo.score == 30 {
		ebitenutil.DebugPrint(screen, "Congratulations, you've hit 30 points!\nThe game is now over.")
	}
}

func (s scrollDemo) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(500, 500)
	ebiten.SetWindowTitle("Top-Down Scroller Example")

	backgroundPict, _, err := ebitenutil.NewImageFromFile("Nebula Aqua-Pink.png")
	if err != nil {
		fmt.Println("Unable to load background image:", err)
	}

	playerPict, _, err := ebitenutil.NewImageFromFile("Spaceship_05_ORANGE.png")
	if err != nil {
		fmt.Println("Unable to load player sprite image:", err)
	}

	demo := scrollDemo{
		player: playerSprite{
			sprite: playerPict,
			xLoc:   250,
			yLoc:   436,
			width:  64,
			height: 64,
		},
		background:      backgroundPict,
		projectiles:     make([]*projectile, 0),
		enemies:         make([]*enemyShip, 0),
		enemySpawnTimer: 0,
		score:           0,
	}

	err = ebiten.RunGame(&demo)
	if err != nil {
		fmt.Println("Failed to run game", err)
	}
}
