# Go Game
![Go Game](./resources/go_game.png)
A top down shooting game.<br>
Written in `go`.<br>
Inspired by [Risk of Rain](https://en.wikipedia.org/wiki/Risk_of_Rain).<br>
Using the [ebiten](https://github.com/hajimehoshi/ebiten) engine for the game objects and the game rendering. <br>

# Install and Run
`go install github.com/bramca/Go-Game/cmd/gogame@latest`
<br>
`gogame`

# Controls
`mouse right` hold it to shoot your gun.<br>
`w/z` hold to thrust up.<br>
`a/q` hold to thrust left.<br>
`s` hold to thrust down.<br>
`d` hold to thrust right.<br>
`p` pause the game.

# Special enemies
- ![rust](./resources/rust.png)![cpp](./resources/cpp.png)![csharp](./resources/csharp.png)![haskell](./resources/haskell.png)![java](./resources/java.png)![javascript](./resources/javascript.png)![python](./resources/python.png) <br> other programming languages are the main enemies.<br>
- ![github](./resources/github.png) lootbox with random powerup.<br>
- ![rubber duck](./resources/rubber_duck.png) a rubber duck carrying a random gun.<br>

# Permanent Powerups
- *Health* heal or increase health.<br>
- *Firerate* increase firerate.<br>
- *Movement* increase movement speed.<br>
- *Laser Speed* increase bullet speed.<br>

# Temporary Powerups
- *Detect Boxes* detect lootboxes locations form a limited time.<br>
- *Invincible* take no damage for a small period of time.<br>
- ![Insta Kill](./resources/skull.png) *Insta Kill* kill enemies instantly.<br>
- ![Vampire Mode](./resources/gopher_vampire.png) *Vampire Mode* gain back damage done to enemies.

# Guns
- *Default* standard single shot laser.<br>
- *Shotgun* shoots lasers in a scattered pattern.<br>
- *Homing Lasers* will follow the nearest enemy.<br>
- *Piercing Lasers* will pierce the enemy.<br>
- *Double Lasers* shoots 2 parallel lasers.<br>
- *Exploding Lasers* will explode in multiple lasers with damage halved at impact.
