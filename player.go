package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Player struct {
	x, y         float64
	w, h         float64
	angle        float64
	lasers       []*Laser
	tempRewards  []*TempReward
	img          *ebiten.Image
	ySpeed       float64
	xSpeed       float64
	points       int
	maxPoints    int
	healthBar    HealthBar
	invincible   bool
	instaKill    bool
	vampire      bool
	score        int
	fireRate     int
	laserSpeed   float64
	speed        float64
	acceleration float64
	damage       int
	gun          string
	ammo         int
}

func (p *Player) shoot() {
	// Adapt variables to weapon
	var fireRate int
	if p.ammo == 0 {
		p.gun = playerDefaultGun
		p.ammo = -1
	}
	switch p.gun {
	case "Shotgun":
		fireRate = p.fireRate * 3
	case "Exploding Lasers":
		fireRate = p.fireRate * 3
	case "Piercing Lasers":
		fireRate = p.fireRate * 2
	default:
		fireRate = p.fireRate
	}
	if playerFireFrameCount%fireRate == 0 {
		switch p.gun {
		case "Exploding Lasers":
			p.lasers = append(p.lasers, &Laser{
				x:         p.x,
				y:         p.y,
				angle:     p.angle,
				speed:     p.laserSpeed,
				color:     playerExplodingLaserColor,
				duration:  laserDuration,
				size:      laserSize,
				damage:    p.damage,
				exploding: true,
			})
		case "Double Lasers":
			laserDist := 10.0
			p.lasers = append(p.lasers, &Laser{
				x:        p.x + math.Sin(p.angle)*laserDist,
				y:        p.y - math.Cos(p.angle)*laserDist,
				angle:    p.angle,
				speed:    p.laserSpeed,
				color:    playerLaserColor,
				duration: laserDuration,
				size:     laserSize,
				damage:   p.damage,
			})
			p.lasers = append(p.lasers, &Laser{
				x:        p.x - math.Sin(p.angle)*laserDist,
				y:        p.y + math.Cos(p.angle)*laserDist,
				angle:    p.angle,
				speed:    p.laserSpeed,
				color:    playerLaserColor,
				duration: laserDuration,
				size:     laserSize,
				damage:   p.damage,
			})
		case "Piercing Lasers":
			p.lasers = append(p.lasers, &Laser{
				x:        p.x,
				y:        p.y,
				angle:    p.angle,
				speed:    p.laserSpeed,
				color:    playerPiercingLaserColor,
				duration: laserDuration,
				size:     laserSize,
				damage:   p.damage,
				piercing: true,
			})
		case "Homing Lasers":
			p.lasers = append(p.lasers, &Laser{
				x:                 p.x,
				y:                 p.y,
				angle:             p.angle,
				speed:             p.laserSpeed,
				color:             playerHomingLaserColor,
				duration:          laserDuration,
				size:              laserSize,
				damage:            p.damage,
				homing:            true,
				homingTargetIndex: -1,
				homingRange:       float64(int(math.Min(screenWidth, screenHeight))+rand.Intn(int(math.Max(screenWidth, screenHeight))-int(math.Min(screenWidth, screenHeight)))) / 4,
			})
		case "Shotgun":
			for i := -math.Pi / 12; i < math.Pi/12; i += math.Pi / 36 {
				p.lasers = append(p.lasers, &Laser{
					x:        p.x,
					y:        p.y,
					angle:    p.angle + i,
					speed:    p.laserSpeed * (0.7 + rand.Float64()*0.3),
					color:    playerLaserColor,
					duration: laserDuration,
					size:     laserSize,
					damage:   p.damage,
				})
			}
		default:
			p.lasers = append(p.lasers, &Laser{
				x:        p.x,
				y:        p.y,
				angle:    p.angle,
				speed:    p.laserSpeed,
				color:    playerLaserColor,
				duration: laserDuration,
				size:     laserSize,
				damage:   p.damage,
			})
		}
		playerFireFrameCount = 0
		if p.ammo > 0 {
			p.ammo -= 1
		}
	}
}

func (p *Player) update(x, y float64, dots []*Dot) {
	// Move the player based on the mouse position
	mx, my := ebiten.CursorPosition()
	p.y += p.ySpeed
	p.x += p.xSpeed
	p.healthBar.update(p.x-camX-p.w/2, p.y-(p.h-p.h/3)-camY, p.points, p.maxPoints)
	p.angle = angleBetweenPoints(x, y, float64(mx), float64(my))
	for dotIndex := range dots {
		if !dots[dotIndex].eaten && dots[dotIndex].duration > 0 && (distanceBetweenPoints(p.x+p.xSpeed, p.y+p.ySpeed, float64(dots[dotIndex].x), float64(dots[dotIndex].y)) < p.w*0.8 || distanceBetweenPoints(p.x+p.xSpeed, p.y+p.ySpeed, float64(dots[dotIndex].x+len(dots[dotIndex].msg)), float64(dots[dotIndex].y)) < p.w) {
			dot := Dot{
				x:        dots[dotIndex].x,
				y:        dots[dotIndex].y,
				color:    dotHitColor,
				msg:      "+" + strconv.Itoa(pointsPerDot),
				textFont: text.NewGoXFace(hitTextFont),
			}
			setDotDrawOptions(&dot)
			p.points += pointsPerDot
			dots[dotIndex].hits = append(dots[dotIndex].hits, Hit{
				Dot:      dot,
				duration: 2 * framesPerSecond / 3,
			})
			dots[dotIndex].eaten = true
			if p.points > p.maxPoints {
				p.maxPoints = p.points
			}
		}
	}
}

func (p *Player) updateTempRewards() {
	for index := len(p.tempRewards) - 1; index >= 0; index-- {
		if p.tempRewards[index].duration < 0 {
			p.tempRewards[index] = p.tempRewards[len(p.tempRewards)-1]
			p.tempRewards = p.tempRewards[:len(p.tempRewards)-1]
			continue
		}
		p.tempRewards[index].update()
	}
}

func (p *Player) drawTempRewards(screen *ebiten.Image) {
	for index := len(p.tempRewards) - 1; index >= 0; index-- {
		if p.tempRewards[index].duration < 0 {
			p.tempRewards[index] = p.tempRewards[len(p.tempRewards)-1]
			p.tempRewards = p.tempRewards[:len(p.tempRewards)-1]
			continue
		}
		p.tempRewards[index].draw(screen)
	}
}

func (p *Player) drawStats(screen *ebiten.Image) {
	text.Draw(screen, fmt.Sprintf("Score: %d", p.score), text.NewGoXFace(scoreTextFont), scoreDrawOptions)
	scoreDrawOptions.GeoM.Translate(0, float64(newlinePadding))
	text.Draw(screen, fmt.Sprintf("Gun: %s", p.gun), text.NewGoXFace(scoreTextFont), scoreDrawOptions)
	scoreDrawOptions.GeoM.Translate(0, float64(newlinePadding))
	if p.ammo >= 0 {
		text.Draw(screen, fmt.Sprintf("Ammo: %d", p.ammo), text.NewGoXFace(scoreTextFont), scoreDrawOptions)
	} else {
		text.Draw(screen, "Ammo: Infinite", text.NewGoXFace(scoreTextFont), scoreDrawOptions)
	}
	scoreDrawOptions.GeoM = scoreGeoMatrix
}

func (p *Player) draw(screen *ebiten.Image, x float64, y float64) {
	// Draw the player
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(p.w/2), -float64(p.h/2))
	op.GeoM.Rotate(p.angle)
	op.GeoM.Translate(x, y)
	screen.DrawImage(p.img, op)
	p.healthBar.draw(screen)
}

func (p *Player) updateLasers() {
	for index := len(p.lasers) - 1; index >= 0; index-- {
		hit := false
		for _, enemy := range enemies {
			if !enemy.dead && math.Abs(float64(p.lasers[index].y+p.lasers[index].speed*math.Sin(p.lasers[index].angle))-float64(enemy.y)) < enemy.h/2 && math.Abs(float64(p.lasers[index].x+p.lasers[index].speed*math.Cos(p.lasers[index].angle))-float64(enemy.x)) < enemy.w/2 {
				damage := p.lasers[index].damage
				if p.instaKill {
					damage = enemy.points
				}
				enemy.points -= damage
				if p.vampire && p.points <= p.maxPoints {
					healing := damage
					if p.points+damage > p.maxPoints {
						healing = p.maxPoints - p.points
					}
					p.points += healing
				}
				hit = true
				dot := Dot{
					x:        int(enemy.x),
					y:        int(enemy.y - enemy.h/2),
					color:    damageColor,
					msg:      strconv.Itoa(-damage),
					textFont: text.NewGoXFace(hitTextFont),
				}
				setDotDrawOptions(&dot)
				enemy.hits = append(enemy.hits, Hit{
					Dot:      dot,
					duration: 2 * framesPerSecond / 3,
				})
			}
		}
		for _, lootBox := range lootBoxes {
			if !lootBox.broken && lootBox.duration > 0 && math.Abs(float64(p.lasers[index].y+p.lasers[index].speed*math.Sin(p.lasers[index].angle))-float64(lootBox.y)) < lootBox.h/2 && math.Abs(float64(p.lasers[index].x+p.lasers[index].speed*math.Cos(p.lasers[index].angle))-float64(lootBox.x)) < lootBox.w/2 {
				lootBox.hitpoints -= p.lasers[index].damage
				hit = true
				dot := Dot{
					x: int(lootBox.x),
					y: int(lootBox.y - lootBox.h/2),
					color: color.RGBA{
						R: 0xff,
						G: 0xff,
						B: 0xff,
						A: 0xf0,
					},
					msg:      strconv.Itoa(-p.lasers[index].damage),
					textFont: text.NewGoXFace(hitTextFont),
				}
				setDotDrawOptions(&dot)
				lootBox.hits = append(lootBox.hits, Hit{
					Dot:      dot,
					duration: 2 * framesPerSecond / 3,
				})
			}
		}
		for _, rubberDuck := range rubberDucks {
			if !rubberDuck.dead && math.Abs(float64(p.lasers[index].y+p.lasers[index].speed*math.Sin(p.lasers[index].angle))-float64(rubberDuck.y)) < rubberDuck.h/2 && math.Abs(float64(p.lasers[index].x+p.lasers[index].speed*math.Cos(p.lasers[index].angle))-float64(rubberDuck.x)) < rubberDuck.w/2 {
				damage := p.lasers[index].damage
				if p.instaKill {
					damage = rubberDuck.points
				}
				rubberDuck.points -= damage
				hit = true
				dot := Dot{
					x: int(rubberDuck.x),
					y: int(rubberDuck.y - rubberDuck.h/2),
					color: color.RGBA{
						R: 0xff,
						G: 0xff,
						B: 0xff,
						A: 0xf0,
					},
					msg:      strconv.Itoa(-p.lasers[index].damage),
					textFont: text.NewGoXFace(hitTextFont),
				}
				setDotDrawOptions(&dot)
				rubberDuck.hits = append(rubberDuck.hits, Hit{
					Dot:      dot,
					duration: 2 * framesPerSecond / 3,
				})
			}
		}
		if hit && !p.lasers[index].piercing {
			if p.lasers[index].exploding {
				for i := 0.0; i < 2*math.Pi; i += math.Pi / 24 {
					p.lasers = append(p.lasers, &Laser{
						x:        p.lasers[index].x,
						y:        p.lasers[index].y,
						angle:    p.angle + i,
						speed:    p.laserSpeed * (0.7 + rand.Float64()*0.3),
						color:    playerExplodingLaserColor,
						duration: laserDuration,
						size:     laserSize,
						damage:   p.damage / 2,
					})
				}
			}
			p.lasers[index] = p.lasers[len(p.lasers)-1]
			p.lasers = p.lasers[:len(p.lasers)-1]
			continue
		}
		p.lasers[index].update()
	}
}

func (p *Player) drawLasers(screen *ebiten.Image, camX float64, camY float64) {
	for index := len(p.lasers) - 1; index >= 0; index-- {
		if p.lasers[index].duration < 0 {
			p.lasers[index] = p.lasers[len(p.lasers)-1]
			p.lasers = p.lasers[:len(p.lasers)-1]
			continue
		}
		p.lasers[index].draw(screen, camX, camY)
	}
}
