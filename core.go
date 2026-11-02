package main

import (
	"math"

	z "github.com/gen2brain/raylib-go/raylib"
)

var (

	//GAME
	invenActiv                 bool
	gameBlokNumW, gameBlokNumH = 1000, 100

	pl              PLAYER
	bgIM            []IM
	spc, blokw      = un8th, un
	cloudnums       = []int{21, 22, 23, 24, 25, 26}
	cloudsActiv     = true
	quicklen        = 10
	fx              []ANIM
	ychngOBJarrow   float32
	ychngOBJarrowSW bool

	bloksGround, bloks, blokTypes []BLOK

	clouds, obj, objList, quickOBJ, projPL []OBJ

	bloksSCR []*BLOK
	enmSCR   []*OBJ

	grav, springH                float32 = 4, -(un10 * 2)
	invtotal, levnum, quickActiv int

	gameW, gameH, groundH = un * float32(gameBlokNumW), un * float32(gameBlokNumH), float32(SCRH) - un2

	//FLAGS
	flags, flagsDropped           []OBJ
	flagsnxtnum                   int
	flagdropActiv, flagsAllPlaced bool

	//ANIM
	seagullanm   ANIM
	seagullR     z.Rectangle
	seagullT     int32
	seagullLR    bool
	seagullActiv = true
	anms         []ANIM

	//ENEMY
	enmList, enms, npcList, npc []OBJ

	//IMG
	etcIM, splatsIM []z.Texture2D

	//UNITS
	un float32 = 64

	un3q, unh, unq, un8th, un16th, un2, un3, un4, un5, un6, un7, un8, un9, un10 = (un / 4) * 3, un / 2, un / 4, un / 8, un / 16, un * 2, un * 3, un * 4, un * 5, un * 6, un * 7, un * 8, un * 9, un * 10

	//FONTS
	fnsml, fn, fnmed z.Font
	diclist          bool

	//CORE
	CAM         z.Camera2D
	ms, msnocam z.Vector2
	dbgstring   string
	fadeBlink   uint8 = 200
	dbgnum      int

	fadeSW, camChange, mslc, msrc, msld, msrd, debug bool
)

type ANIM struct {
	frameActive   int
	fps, timer    float32
	frames, rC    []z.Rectangle
	framesTexNums []int
	tex           z.Texture2D
	rD            z.Rectangle
	alpha         uint8

	once, off bool
}
type IM struct {
	t z.Texture2D
	r z.Rectangle
}
type PLAYER struct {
	r, rC, rTool      z.Rectangle
	cnt, cntPlatCheck z.Vector2
	anm               []ANIM
	inv               []OBJ

	catActiv, anmActiv, invActiv, invlen int

	jump, lr, moving, useitem, useweap, onplat, onmovingplat, spring bool

	roweap, rotool, vely, velx, velmax, siz, jumph, jumpy, rCxoff, rCyoff, rToolxoff, rToolyoff float32
}
type OBJ struct {
	nm    string
	r, rC z.Rectangle
	t     z.Texture2D
	anm   []ANIM
	txt   []string
	cnt   z.Vector2
	c     z.Color

	txNumActiv, numinven, cat, dmg, hp int

	txT, useT, useTdef int32

	velx, vely, velmax, hitTimer, spd, rCxRatio, rCyRatio, rCwRatio, rChRatio float32

	intool, lr, flip, off, collect, nograv, danim bool
}
type BLOK struct {
	r, rsplat z.Rectangle
	c, csplat z.Color
	t, tsplat z.Texture2D

	solid, splat, flag, off, moving, powerup bool

	velx, vely, moveHmax, moveVmax float32
}

// MARK: UPDATE
func UP() {

	//CAM
	if camChange {
		uCAM()
		camChange = false
	}
	countDraw, countCull = 0, 0
	//TIMERS
	if ychngOBJarrowSW {
		if ychngOBJarrow > 0 {
			ychngOBJarrow -= 2
		} else {
			ychngOBJarrowSW = false
		}
	} else {
		if ychngOBJarrow < unh {
			ychngOBJarrow += 2
		} else {
			ychngOBJarrowSW = true
		}
	}

	if seagullT > 0 {
		seagullT--
	} else {
		if seagullLR {
			seagullR.X -= deltaspd(un8th)
		} else {
			seagullR.X += deltaspd(un8th)
		}
	}
	if seagullR.X > 0 && seagullR.X < scrwf32 {
		seagullActiv = true
	}
	if !crr(seagullR, SCRRECNOCAM) && seagullActiv {
		seagullActiv = false
		seagullT = int32(RINT(5, 15)) * FPS
		seagullLR = FLIP()
		if seagullLR {
			seagullR.X = scrwf32 + un
		} else {
			seagullR.X = -(seagullR.Width + un)
		}
		seagullR.Y = CNT.Y + RF32(-un5, un3)
	}

	//FADE BLINK
	if fadeSW {
		if fadeBlink > 200 {
			fadeBlink -= 5
		} else {
			fadeSW = false
		}
	} else {
		if fadeBlink < 255 {
			fadeBlink += 5
		} else {
			fadeSW = true
		}
	}
	//INPUT
	if z.IsKeyPressed(z.KeyEscape) {
		EXIT = true
	}
	if z.IsKeyPressed(z.KeyTab) {
		invenActiv = !invenActiv
	}

	if z.IsKeyPressed(z.KeyOne) {
		quickActiv = 0
	}
	if z.IsKeyPressed(z.KeyTwo) {
		quickActiv = 1
	}
	if z.IsKeyPressed(z.KeyThree) {
		quickActiv = 2
	}
	if z.IsKeyPressed(z.KeyFour) {
		quickActiv = 3
	}
	if z.IsKeyPressed(z.KeyFive) {
		quickActiv = 4
	}
	if z.IsKeyPressed(z.KeySix) {
		quickActiv = 5
	}
	if z.IsKeyPressed(z.KeySeven) {
		quickActiv = 6
	}
	if z.IsKeyPressed(z.KeyEight) {
		quickActiv = 7
	}
	if z.IsKeyPressed(z.KeyNine) {
		quickActiv = 8
	}
	if z.IsKeyPressed(z.KeyZero) {
		quickActiv = 9
	}

	//MOUSE
	msnocam = z.GetMousePosition()
	ms = z.GetScreenToWorld2D(z.GetMousePosition(), CAM)
	mslc, msrc = false, false
	mslc = z.IsMouseButtonPressed(z.MouseButtonLeft)
	msrc = z.IsMouseButtonPressed(z.MouseButtonRight)
	msld = z.IsMouseButtonDown(z.MouseButtonLeft)
	msrd = z.IsMouseButtonDown(z.MouseButtonRight)

	//FUNCTION
	if z.IsKeyPressed(z.KeyF1) {
		debug = !debug
	}
	if z.IsKeyPressed(z.KeyF2) {
		diclist = !diclist
	}
	if z.IsKeyPressed(z.KeyF3) {
		palette = !palette
	}

	//CAM ZOOM
	if z.IsKeyPressed(z.KeyKpAdd) {
		switch CAM.Zoom {
		case 1:
			CAM.Zoom = 1.5
		case 1.5:
			CAM.Zoom = 0.5
		case 0.5:
			CAM.Zoom = 1

		}
		camChange = true
	}
	if z.IsKeyPressed(z.KeyKpSubtract) {
		switch CAM.Zoom {
		case 1:
			CAM.Zoom = 0.5
		case 0.5:
			CAM.Zoom = 1.5
		case 1.5:
			CAM.Zoom = 1
		}
		camChange = true
	}
	if camChange {
		uCAM()
	}

	//PLAYER INVEN
	invtotal = 0
	for i := range pl.inv {
		if pl.inv[i].nm != "" {
			invtotal++
		}
	}

	//MARK: PLAYER KEYS
	var moveX, moveY float32
	moved := false

	// --- 1. HORIZONTAL INPUT (LEFT / RIGHT) ---
	if z.IsKeyDown(z.KeyLeft) || z.IsKeyDown(z.KeyA) {
		moved = true
		moveX -= 1.0
		pl.lr = true
	}
	if z.IsKeyDown(z.KeyRight) || z.IsKeyDown(z.KeyD) {
		moved = true
		moveX += 1.0
		pl.lr = false
	}

	// --- 2. JUMP INPUT (UP / W) ---
	// Tune these two variables to taste:
	jumpImpulse := float32(1800.0) // Upward snap (higher = higher jump)

	// Jump press:
	if (z.IsKeyPressed(z.KeyUp) || z.IsKeyPressed(z.KeyW)) && !pl.jump {
		pl.vely = -jumpImpulse
		pl.jump = true
		pl.onplat = false
	}

	// Short-hop / Jump cut: Letting go of jump early cuts upward momentum
	if (z.IsKeyReleased(z.KeyUp) || z.IsKeyReleased(z.KeyW)) && pl.vely < 0 {
		pl.vely *= 0.45 // Instantly cuts the float on button release
	}
	if z.IsKeyDown(z.KeyDown) || z.IsKeyDown(z.KeyS) {
	}

	if moveX != 0 && moveY != 0 {
		length := float32(math.Sqrt(float64(moveX*moveX + moveY*moveY)))
		moveX /= length
		moveY /= length
	}

	// --- 3. HORIZONTAL VELOCITY & FRICTION (X AXIS) ---
	if moved {
		// Accelerate in the direction pressed
		pl.velx += moveX * (pl.velmax * 8) * delta

		// Clamp between -velmax and +velmax (fixes running left infinitely accelerating!)
		if pl.velx > pl.velmax {
			pl.velx = pl.velmax
		}
		if pl.velx < -pl.velmax {
			pl.velx = -pl.velmax
		}
		//}

	} else {
		// Apply friction when not pressing Left/Right
		friction := float32(pl.velmax * 10 * delta)
		if pl.velx > 0 {
			pl.velx -= friction
			if pl.velx < 0 {
				pl.velx = 0
			}
		} else if pl.velx < 0 {
			pl.velx += friction
			if pl.velx > 0 {
				pl.velx = 0
			}
		}
	}
	// --- 2. ASYMMETRIC GRAVITY (Snappy falling) ---
	baseGravity := float32(2200.0)  // Strong base gravity (rising)
	fallMultiplier := float32(2.2)  // 1.8x heavier when falling down
	maxFallSpeed := float32(1500.0) // Terminal velocity

	if pl.jump {
		if pl.vely < 0 {
			// Rising upward
			pl.vely += baseGravity * delta
		} else {
			// Falling downward: much stronger gravity pulls you down fast!
			pl.vely += (baseGravity * fallMultiplier) * delta
		}

		if pl.vely > maxFallSpeed {
			pl.vely = maxFallSpeed
		}
	}

	// --- 3. HORIZONTAL PREDICTIVE MOVEMENT & WALL COLLISION ---
	nextX := pl.r.X + (pl.velx * delta)

	// Create a test box at the projected X position with small vertical insets
	// (The vertical insets prevent floor/ceiling blocks from acting like side walls)
	skinY := float32(3.0)
	testBox := z.NewRectangle(
		nextX,
		pl.r.Y+skinY,
		pl.r.Width,
		pl.r.Height-(skinY*2),
	)

	if pl.velx > 0 {
		// Moving RIGHT: check for walls to the right
		for i := range bloksSCR {
			if !bloksSCR[i].solid {
				continue
			}

			b := bloksSCR[i].r

			// Check collision with the projected position
			if crr(testBox, b) {
				// Only collide if player was previously to the left of the block
				if pl.r.X+pl.r.Width <= b.X+4 {
					nextX = b.X - pl.r.Width // Snap flush against the left face of the wall
					pl.velx = 0              // Stop horizontal momentum
					testBox.X = nextX
				}
			}
		}
	} else if pl.velx < 0 {
		// Moving LEFT: check for walls to the left
		for i := range bloksSCR {
			if !bloksSCR[i].solid {
				continue
			}

			b := bloksSCR[i].r

			// Check collision with the projected position
			if crr(testBox, b) {
				// Only collide if player was previously to the right of the block
				if pl.r.X >= b.X+b.Width-4 {
					nextX = b.X + b.Width // Snap flush against the right face of the wall
					pl.velx = 0           // Stop horizontal momentum
					testBox.X = nextX
				}
			}
		}
	}

	// Apply validated X position
	pl.r.X = nextX

	// Calculate vertical next positions
	currBottom := pl.r.Y + pl.r.Height
	nextY := pl.r.Y + (pl.vely * delta)
	nextBottom := nextY + pl.r.Height

	landed := false

	// --- PLATFORM COLLISION (Only when falling downwards) ---
	if pl.spring {
		if pl.r.Y > springH {
			pl.r.Y -= un3q
		} else {
			pl.spring = false
			pl.vely = (baseGravity * 20) * delta
		}
	} else {
		if pl.vely >= 0 {
			checkV2 := pl.cnt
			checkV2.Y = pl.r.Y + pl.r.Height + 8
			pl.onmovingplat = false
			for i := range bloksSCR {
				if crp(bloksSCR[i].r, checkV2) {
					if bloksSCR[i].moving {
						pl.onmovingplat = true
					}
				}
			}

			for i := range bloksSCR {
				if !bloksSCR[i].solid {
					continue
				}

				b := bloksSCR[i].r

				// Check horizontal overlap between player and platform
				horizOverlap := (pl.r.X+pl.r.Width > b.X) && (pl.r.X < b.X+b.Width)

				// Check if player's feet cross from above the platform to on/below it
				if horizOverlap && currBottom <= b.Y+8 && nextBottom >= b.Y {
					pl.r.Y = b.Y - pl.r.Height
					pl.onplat = true
					landed = true
					pl.jump = false
					if !pl.onmovingplat {
						pl.vely = 0

					}
					break
				}
			}
		}

		// --- GROUND FLOOR COLLISION ---
		if !landed {
			if nextBottom >= groundH {
				pl.r.Y = groundH - pl.r.Height
				pl.vely = 0
				pl.jump = false
				pl.onplat = false
				landed = true
			} else {
				// In mid-air
				pl.r.Y = nextY
			}
		}

		// --- LEDGE CHECK (Walking off edges) ---
		if !pl.jump && pl.onplat {
			// Verify still supported by a solid platform beneath feet
			feetY := pl.r.Y + pl.r.Height
			stillOnPlatform := false

			for i := range bloksSCR {
				if !bloksSCR[i].solid {
					continue
				}
				b := bloksSCR[i].r
				horizOverlap := (pl.r.X+pl.r.Width > b.X) && (pl.r.X < b.X+b.Width)

				// If feet are resting right on top of this block
				if horizOverlap && math.Abs(float64(feetY-b.Y)) <= 2.0 {
					stillOnPlatform = true
					break
				}
			}

			if !stillOnPlatform {
				pl.onplat = false
				pl.jump = true // Started falling off ledge
			}
		}
	}

	// Update player center and secondary hitboxes:
	pl.cnt = rcnt(pl.r)
	uPLrecs()

	camChange = true

	// MARK: UP PLAYER ANIM & TOOLS
	pl.anm[0].timer += delta
	if pl.velx == 0 && pl.vely == 0 {
		pl.moving = false
	} else if pl.velx == 0 && pl.vely != 0 && pl.onmovingplat {
		pl.moving = false
	} else {
		pl.moving = true
	}
	if pl.lr {
		pl.rotool -= 8
	} else {
		pl.rotool += 8
	}
	pl.useitem = msld
	if pl.useitem {
		if ms.X < pl.cnt.X {
			pl.lr = true
		} else {
			pl.lr = false
		}
	}

	if len(projPL) > 0 {
		uProjPL(delta)
	}
}
func uPLrecs() { //MARK: UP PLAYER RECS
	pl.rTool.X = pl.r.X - pl.rToolxoff
	pl.rTool.Y = pl.r.Y - pl.rToolyoff
	pl.rC.X = pl.r.X - pl.rCxoff
	pl.rC.Y = pl.r.Y - pl.rCyoff
}
func uProjPL(delta float32) { //MARK: UP PLAYER PROJ
	// 1. Tick down weapon cooldown timer
	if pl.catActiv == 2 {
		if quickOBJ[quickActiv].useT > 0 {
			quickOBJ[quickActiv].useT--
		}
	}

	// 3. Update & Collide all active projectiles
	for i := range projPL {
		p := &projPL[i]
		if p.off {
			continue
		}

		// Move projectile
		p.r.X += p.velx * delta
		p.r.Y += p.vely * delta
		uOBJrecs(p) // Updates rC position & dimensions
		p.cnt = rcnt(p.r)

		// Ground Collision: hit ground floor
		if p.rC.Y+p.rC.Height >= groundH {
			p.off = true
			continue
		}
		if !crr(p.rC, SCRREC) { //OFF SCREEN
			p.off = true
			continue
		}

		// Wall/Platform Collision: hit solid blocks on screen
		if len(bloksSCR) > 0 {
			for j := range bloksSCR {
				if bloksSCR[j].solid && crr(p.rC, bloksSCR[j].r) {
					p.off = true
					//FX BLOK COLLIS
					anm := anms[6]
					anm.rD = recfromcnt(rcnt(p.r), un4)
					anm.once = true
					anm.alpha = 30
					fx = append(fx, anm)
					break
				}
			}
		}
		if len(enmSCR) > 0 {
			for j := range enmSCR {
				e := enmSCR[j]
				if crr(p.rC, e.rC) {
					if e.hitTimer == 0 {
						e.hitTimer = float32(FPS)
						e.hp -= p.dmg
						if e.hp <= 0 {
							enmDestroy(e)
						}
					}
					p.off = true
					break
				}
			}
		}
	}

	// 4. Remove dead projectiles (off == true) efficiently in-place
	alive := projPL[:0]
	for _, p := range projPL {
		if !p.off {
			alive = append(alive, p)
		}
	}
	projPL = alive
}

// MARK: OBJ
func mCRATES() { //MARK: MAKE CRATES
	o := objList[6]
	o.r.X = pl.r.X + un5
	obj = append(obj, o)
}
func mSPRINGS() { //MARK: MAKE SPRINGS
	x := un10
	for x < gameW {
		x += RF32(un, un10*2)
		if ROLL(12) > 10 {
			o := objList[10]
			o.r.X = x
			o.anm = make([]ANIM, len(objList[10].anm))
			copy(o.anm, objList[10].anm)
			obj = append(obj, o)
		}
	}

}
func mOBJ() { //MARK: MAKE OBJ
	//TREES
	o := OBJ{}
	o.nm = "tree" //000
	o.t = etcIM[17]
	o.r = trscale(o.t, 20)
	o.rC = o.r
	o.rC.Width = o.r.Width / 4
	o.rC.X = (o.r.X + o.r.Width/2) - o.rC.Width/2
	objList = append(objList, o)
	o = OBJ{}
	o.nm = "tree" //001
	o.t = etcIM[18]
	o.r = trscale(o.t, 20)
	o.rC = o.r
	o.rC.Width = o.r.Width / 4
	o.rC.X = (o.r.X + o.r.Width/2) - o.rC.Width/2
	objList = append(objList, o)
	o = OBJ{}
	o.nm = "tree" //002
	o.t = etcIM[19]
	o.r = trscale(o.t, 20)
	o.rC = o.r
	o.rC.Width = o.r.Width / 4
	o.rC.X = (o.r.X + o.r.Width/2) - o.rC.Width/2
	objList = append(objList, o)
	o = OBJ{}
	o.nm = "tree" //003
	o.t = etcIM[20]
	o.r = trscale(o.t, 20)
	o.rC = o.r
	o.rC.Width = o.r.Width / 4
	o.rC.X = (o.r.X + o.r.Width/2) - o.rC.Width/2
	objList = append(objList, o)

	o = OBJ{} //AXE 004
	o.nm = "axe"
	o.t = etcIM[16]
	o.r = trscale(o.t, 2)
	o.collect = true
	o.cat = 1 //TOOL
	o.useTdef = FPS / 2
	objList = append(objList, o)

	o = OBJ{} //SHOTGUN 005
	o.nm = "shotgun"
	o.t = etcIM[41]
	o.r = trscale(o.t, 2)
	o.collect = true
	o.cat = 2 //WEAPON
	o.useTdef = FPS
	objList = append(objList, o)

	o = OBJ{} //CRATE 006
	o.nm = "crate"
	o.t = etcIM[46]
	o.r = trscale(o.t, 2)
	o.r.Y = groundH - o.r.Height
	o.collect = false
	objList = append(objList, o)

	o = OBJ{} //CRATE 007
	o.nm = "sunflower"
	o.t = etcIM[54]
	o.r = trscale(o.t, 4)
	o.r.Y = groundH - o.r.Height
	o.rC = o.r
	o.rC.Width = o.r.Width / 4
	o.rC.X = (o.r.X + o.r.Width/2) - o.rC.Width/2
	o.collect = false
	o.nograv = true
	objList = append(objList, o)

	o = OBJ{} //CRATE 008
	o.nm = "wheat"
	o.t = etcIM[57]
	o.r = trscale(o.t, 3)
	o.r.Y = groundH - o.r.Height
	o.rC = o.r
	o.rC.Width = o.r.Width / 4
	o.rC.X = (o.r.X + o.r.Width/2) - o.rC.Width/2
	o.collect = false
	o.nograv = true
	objList = append(objList, o)

	o = OBJ{} //CRATE 009
	o.nm = "berry tree"
	o.t = etcIM[59]
	o.r = trscale(o.t, 10)
	o.r.Y = groundH - o.r.Height
	o.rC = o.r
	o.rC.Width = o.r.Width / 4
	o.rC.X = (o.r.X + o.r.Width/2) - o.rC.Width/2
	o.collect = false
	o.nograv = true
	objList = append(objList, o)

	o = OBJ{} //SPRING 010
	o.nm = "spring"
	o.danim = true
	o.anm = append(o.anm, mANIMframes(etcIM[45], 16, 16, 7, 3))
	o.r = R(0, 0, un)
	o.r.Y = groundH - o.r.Height
	o.rC = o.r
	o.collect = false
	o.nograv = true
	objList = append(objList, o)
}
func uOBJrecs(o *OBJ) {
	o.rC.Width = o.r.Width * o.rCwRatio
	o.rC.Height = o.r.Height * o.rChRatio
	o.rC.X = o.r.X + (o.r.Width * o.rCxRatio)
	o.rC.Y = o.r.Y + (o.r.Height * o.rCyRatio)
}
func objRCratio(o *OBJ) {
	if o.r.Width > 0 && o.r.Height > 0 {
		o.rCwRatio = o.rC.Width / o.r.Width
		o.rChRatio = o.rC.Height / o.r.Height
		o.rCxRatio = (o.rC.X - o.r.X) / o.r.Width
		o.rCyRatio = (o.rC.Y - o.r.Y) / o.r.Height
	}
}
func dropflag(r z.Rectangle) {
	o := OBJ{}
	o.anm = append(o.anm, flags[flagsnxtnum].anm[0])
	if crblokrs(r, bloks) {
		check := 100
		for {
			r.Y -= r.Height
			found := crblokrs(r, bloks)
			if !found {
				break
			}
			check--
			if check <= 0 {
				break
			}
		}
	}
	r.X += r.Width / 4
	o.r = r
	flagsDropped = append(flagsDropped, o)
	if flagsnxtnum < len(flags)-1 {
		flagsnxtnum++
	} else {
		flagsAllPlaced = true
	}
}

func useweap() { //MARK: USE WEAPON MAKE PROJECTILES
	weap := &quickOBJ[quickActiv]

	if weap.useT > 0 {
		return
	}

	// 1. Aim direction vector
	dx := ms.X - pl.cnt.X
	dy := ms.Y - pl.cnt.Y
	dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	if dist == 0 {
		dist = 1.0
	}
	dirX := dx / dist
	dirY := dy / dist

	// 2. MUZZLE SPAWN POINT: push bullet spawn outside player hitbox
	// (e.g. 24-32 pixels forward from center in the direction of aim)
	barrelLength := float32(pl.r.Width * 0.7)
	muzzleX := pl.cnt.X + (dirX * barrelLength)
	muzzleY := pl.cnt.Y + (dirY * barrelLength)

	switch weap.nm {
	case "shotgun":
		pelletCount := 5
		weap.useT = weap.useTdef

		// Snappy platformer bullet speed (1400-1600 px/sec)
		bulletSpeed := float32(1500.0)

		for i := 0; i < pelletCount; i++ {
			p := OBJ{}
			p.nm = "bullet"
			p.dmg = 1
			p.t = etcIM[43]
			p.c = c.VoidBlack

			projSize := un8th
			// Spawn at the barrel tip, NOT at the body center
			p.r = z.NewRectangle(muzzleX-projSize/2, muzzleY-projSize/2, projSize, projSize)
			p.rC = p.r
			p.cnt = rcnt(p.r)
			objRCratio(&p)

			// Random spread (-12° to +12°)
			spreadAngle := float32(z.GetRandomValue(-12, 12)) * (math.Pi / 180.0)
			cosA := float32(math.Cos(float64(spreadAngle)))
			sinA := float32(math.Sin(float64(spreadAngle)))
			finalDirX := dirX*cosA - dirY*sinA
			finalDirY := dirX*sinA + dirY*cosA

			// INHERIT PLAYER VELOCITY: Bullet speed + player momentum
			p.velx = (finalDirX * bulletSpeed) + pl.velx
			p.vely = (finalDirY * bulletSpeed) + (pl.vely * 0.3) // partial Y momentum for cleaner aim

			projPL = append(projPL, p)
		}

	case "pistol":
		weap.useT = weap.useTdef

		bulletSpeed := float32(1800.0)
		projSize := float32(8.0)

		p := OBJ{}
		p.nm = "bullet"
		p.dmg = 25
		p.r = z.NewRectangle(muzzleX-projSize/2, muzzleY-projSize/2, projSize, projSize)
		p.rC = p.r
		p.cnt = rcnt(p.r)
		objRCratio(&p)

		// Inherit player momentum
		p.velx = (dirX * bulletSpeed) + pl.velx
		p.vely = (dirY * bulletSpeed) + (pl.vely * 0.3)

		projPL = append(projPL, p)
	}
}
func addquick2inven(num int) bool {
	found := false
	for i := range pl.inv {
		if pl.inv[i].nm == "" {
			pl.inv[i] = quickOBJ[num]
			found = true
			break
		}
	}
	return found
}
func inven2quick(num int) {
	o := pl.inv[num]
	found := false
	if len(quickOBJ) > 0 {
		for i := range quickOBJ {
			if o.nm == quickOBJ[i].nm {
				quickOBJ[i].numinven += o.numinven
				pl.inv[num] = OBJ{}
				found = true
				break
			}
		}
	}
	if !found {
		quickOBJ = append(quickOBJ, pl.inv[num])
		pl.inv[num] = OBJ{}
	}
}
func addremOBJ(list []OBJ, item OBJ, maxLen int) ([]OBJ, bool) {
	// 1. Search if the item is already in the slice
	foundIndex := -1
	for i := range list {
		if list[i].nm == item.nm { // Matches by name/ID
			foundIndex = i
			break
		}
	}

	// 2. If found, REMOVE it (preserving order)
	if foundIndex != -1 {
		list = append(list[:foundIndex], list[foundIndex+1:]...)
		return list, true // Successfully removed
	}

	// 3. If NOT found, ADD it (only if under maxLen limit)
	if len(list) < maxLen {
		list = append(list, item)
		return list, true // Successfully added
	}

	// Already at max capacity (10), did not add
	return list, false
}
func spawnOBJ(num int, r z.Rectangle) { //MARK: SPAWN OBJ
	o := OBJ{}
	switch num {
	case 4: //SUNFLOWER
		o = OBJ{
			r:       R(pl.rTool.X, pl.rTool.Y, unh+unq),
			t:       etcIM[65],
			collect: true,
			nm:      "seeds",
		}
		if !pl.lr {
			o.r.X = pl.rTool.X + pl.rTool.Width
		}
	case 3: //BERRY
		o = OBJ{
			r:       R(pl.rTool.X, pl.rTool.Y, unh+unq),
			t:       etcIM[64],
			collect: true,
			nm:      "big berry",
		}
		if !pl.lr {
			o.r.X = pl.rTool.X + pl.rTool.Width
		}
	case 2: //WHEAT
		o = OBJ{
			r:       R(pl.rTool.X, pl.rTool.Y, unh+unq),
			t:       etcIM[66],
			collect: true,
			nm:      "wheat",
		}
		if !pl.lr {
			o.r.X = pl.rTool.X + pl.rTool.Width
		}
	case 1: //CRATE
		o = objList[5]
		o.r = rlrgr(r, un8th)
		o.r = rcntrinr(o.r, r)
		o.r.Y = (r.Y + r.Height) - o.r.Height
	case 0: //TREE LOG
		o = OBJ{
			r:       R(pl.rTool.X, pl.rTool.Y, unh+unq),
			t:       etcIM[27],
			collect: true,
			nm:      "log",
		}
		if !pl.lr {
			o.r.X = pl.rTool.X + pl.rTool.Width
		}
	}
	obj = append(obj, o) //APPEND OBJ
}
func collect(o *OBJ) { //MARK: COLLECT OBJ
	found := false
	for i := range pl.inv {
		if pl.inv[i].nm == o.nm {
			pl.inv[i].numinven++
			found = true
			break
		}
	}
	if !found {
		for i := range pl.inv {
			if pl.inv[i].nm == "" {
				pl.inv[i] = *o
				break
			}
		}
	}
	minfotxt("collected " + o.nm)
	o.off = true
}

// uOBJFrameRC computes the world-space collision box for the current active animation frame
func uOBJFrameRC(o *OBJ, anmNum int, r z.Rectangle) {
	if len(o.anm) == 0 || len(o.anm[0].rC) == 0 {
		o.rC = o.r
		return
	}

	anm := &o.anm[anmNum]
	frameIdx := anm.frameActive
	if frameIdx >= len(anm.rC) {
		frameIdx = 0
	}

	// 1. Get the local 0,0 box for this frame
	localBox := anm.rC[frameIdx]

	// 2. Calculate scaling (in case o.r is not 64x64)
	frameW := anm.frames[0].Width  // base frame width
	frameH := anm.frames[0].Height // base frame height
	scaleX := r.Width / frameW
	scaleY := r.Height / frameH

	// 3. Scale the collision dimensions
	boxW := localBox.Width * scaleX
	boxH := localBox.Height * scaleY

	// 4. Translate local X, Y to World X, Y (handling flip!)
	var boxX float32
	if o.flip {
		// If flipped horizontally: mirror the local X offset
		boxX = r.X + (frameW-localBox.X-localBox.Width)*scaleX
	} else {
		boxX = r.X + (localBox.X * scaleX)
	}
	boxY := r.Y + (localBox.Y * scaleY)

	// 5. Store directly into the enemy's world collision box
	o.rC = z.NewRectangle(boxX, boxY, boxW, boxH)
}
func mANM() { //MARK: MAKE ANIM
	//ANMS
	anms = append(anms, mANIMframes(etcIM[29], 64, 64, 10, 10)) //000 SLIME
	anms[0].rC = mRCmultiframe("res/im/etc/029.png", 64, 64)
	for i := 1; i < len(anms[0].rC); i++ { //ADJUST COLLISION REC SIZES
		anms[0].rC[i].Height = anms[0].rC[0].Height
	}
	anms[0].rC[7].Height = (anms[0].rC[0].Height / 2) + 8
	anms[0].rC[8].Height = (anms[0].rC[0].Height / 2) + 8
	anms[0].rC[9].Height = (anms[0].rC[0].Height / 2) + 8

	anms = append(anms, mANIMtexnums([]int{33, 34, 35}, 5))     //001 INTRO DUDE
	anms = append(anms, mANIMframes(etcIM[36], 16, 16, 10, 4))  //002 FLAG MAGENTA
	anms = append(anms, mANIMframes(etcIM[37], 16, 16, 10, 4))  //003 FLAG RED
	anms = append(anms, mANIMframes(etcIM[38], 16, 16, 10, 4))  //004 FLAG YELLOW
	anms = append(anms, mANIMframes(etcIM[39], 16, 16, 10, 4))  //005 FLAG CYAN
	anms = append(anms, mANIMframes(etcIM[42], 64, 64, 30, 11)) //006 PROJ SMOKE FX
	anms = append(anms, mANIMframes(etcIM[44], 64, 64, 30, 9))  //007 GREEN SLIME DEATH FX
	anms = append(anms, mANIMframes(etcIM[47], 64, 64, 30, 12)) //008 CRATE SMOKE FX
	anms = append(anms, mANIMframes(etcIM[48], 64, 64, 45, 12)) //009 COLLECT OBJ FX

	//SEAGULL
	seagullanm = mANIMframes(etcIM[28], 32, 32, 10, 8)
	seagullR = R(-(seagullR.Width + un), CNT.Y-un2, un2)
	seagullR.Y = CNT.Y + RF32(-un5, un3)
	seagullT = int32(RINT(5, 15)) * FPS
}

// MARK: MAKE GAME ENEMIES PLAYER
func mGAME() {

	//BLOKS
	b := blokTypes[0] //FLOOR
	var x, y float32 = 0, float32(SCRH) - un2
	r := R(x, y, blokw)
	r2 := R(x, y+blokw, blokw)
	for range gameBlokNumW {
		b.r = r
		blokadd(&b, 2)
		b.r = r2
		blokadd(&b, 2)
		r.X += blokw
		r2.X += blokw
	}

	x, y = 0, groundH-un4 //LEVEL 1 PLATFORMS
	r = R(x, y, blokw)
	b = blokTypes[0]
	for {
		var l float32
		numl := RINT(3, 9)
		l = f32(numl+4) * blokw
		if ROLL(12) > 4 {
			r3 := r
			for range numl {
				b.r = r3
				blokadd(&b, 0)
				bloks = append(bloks, b)
				r3.X += blokw
				if r3.X >= gameW {
					break
				}
			}
		}
		if ROLL(6) == 6 { //MOVING VERT PLATFORMS
			r.X += f32(RINT(2, 6) * int(blokw))
			b2 := blokTypes[1]
			b2.r = r
			b2.moveVmax = blokw * f32(RINT(18, 35))
			b2.r.Y -= RF32(un, b2.moveVmax/2)
			for {
				b2.vely = RF32(-un8th, un8th)
				if b2.vely > un16th/2 || b2.vely < -(un16th/2) {
					break
				}
			}
			bloks = append(bloks, b2)
			b2.r.X += blokw
			b2.t = etcIM[50]
			r.X += blokw
			bloks = append(bloks, b2)
		}
		r.X += l
		if r.X >= gameW {
			break
		}
	}
	x, y = 0, groundH-un4 //LEVEL 2 BLOK SHAPES
	for {
		if ROLL(12) > 5 {
			x = mBLOKSHAPE(RINT(0, 9), x, y-un4)
		}
		x += RF32(un, un3)
		if x >= gameW {
			break
		}
	}

	x, y = 0, groundH-(un10+un7) //LEVEL 3 ISLANDS
	x += RF32(un, un3)
	for {
		x = mISLAND(x, y)
		x += RF32(un, un3)
		if x >= gameW {
			break
		}
	}

	//CLEAR VERT PLATFORM SHAFTS
	for i := range bloks {
		if bloks[i].vely != 0 {
			xl := bloks[i].r.X - un2
			xr := bloks[i].r.X + bloks[i].r.Width + un2
			for j := range bloks {
				if !bloks[j].moving {
					if bloks[j].r.X >= xl && bloks[j].r.X <= xr {
						bloks[j].off = true
					}
					if !bloks[j].off {
						if bloks[j].r.X+bloks[j].r.Width >= xl && bloks[j].r.X+bloks[j].r.Width <= xr {
							bloks[j].off = true
						}
					}
				}
			}

			keep := bgIM[:0]
			for _, img := range bgIM {
				// Keep images that are OUTSIDE the range [xl, xr]

				if img.r.X+img.r.Width < xl || img.r.X > xr {
					keep = append(keep, img)
				}

			}
			bgIM = keep

			keep2 := obj[:0]
			for _, ob := range obj {
				// Keep images that are OUTSIDE the range [xl, xr]
				if ob.r.X+ob.r.Width < xl || ob.r.X > xr {
					keep2 = append(keep2, ob)
				}
			}
			bgIM = keep

		}
	}
	//GRASS
	var imgrass []IM
	im := IM{}
	im.t = etcIM[7]
	im.r = trscale(im.t, 4)
	imgrass = append(imgrass, im)
	im.t = etcIM[8]
	im.r = trscale(im.t, 4)
	imgrass = append(imgrass, im)
	im.t = etcIM[9]
	im.r = trscale(im.t, 4)
	imgrass = append(imgrass, im)
	x, y = 0, 0
	for x < gameW {
		y = groundH + 3
		choose := RINT(0, 3)
		y -= imgrass[choose].r.Height
		imgrass[choose].r.Y = y
		imgrass[choose].r.X = x
		bgIM = append(bgIM, imgrass[choose])
		x += imgrass[choose].r.Width
	}
	//MARK: MAKE PLAYER
	pl.siz = un
	pl.r = R(un5, groundH-pl.r.Height, pl.siz)
	pl.cnt = rcnt(pl.r)
	pl.rC = rlrgr(pl.r, unq)
	pl.rTool = rlrgr(pl.r, un)
	pl.rCxoff = pl.r.X - pl.rC.X
	pl.rCyoff = pl.r.Y - pl.rC.Y
	pl.rToolxoff = pl.r.X - pl.rTool.X
	pl.rToolyoff = pl.r.Y - pl.rTool.Y
	pl.velmax = un10 * 2
	pl.jumph = un5
	pl.catActiv = -1
	pl.anm = append(pl.anm, mANIMtexnums([]int{10, 11}, 2))
	pl.anm = append(pl.anm, mANIMtexnums([]int{14, 15}, 10))
	pl.anm = append(pl.anm, mANIMtexnums([]int{14, 15}, 4))
	pl.invlen = 20
	for range pl.invlen {
		pl.inv = append(pl.inv, OBJ{})
	}
	//pl.inv[0] = objList[4]
	//invtotal++
	quickOBJ = append(quickOBJ, objList[4])

	switch levnum {
	case 0:
		//ADD TREES
		x, y = 0, 0
		for x < gameW {
			if ROLL(24) > 16 {
				y = groundH + 3
				choose := RINT(0, 4)
				y -= objList[choose].r.Height
				objList[choose].r.Y = y
				objList[choose].r.X = x
				objList[choose].rC.Y = y
				objList[choose].rC.X = (objList[choose].r.X + objList[choose].r.Width/2) - objList[choose].rC.Width/2
				obj = append(obj, objList[choose])
				x += objList[choose].r.Width
			}
			x += RF32(un, un4)
			if x >= gameW {
				break
			}
		}
	}
}
func mISLAND(x, y float32) float32 { //MARK: MAKE ISLAND
	wnum, hnum := RINT(12, 23), RINT(7, 13)
	w, h := f32(wnum)*blokw, f32(hnum)*blokw
	y -= h
	ox := x
	ex := ox + w + un
	b := blokTypes[0]

	var lev3plants []IM
	im := IM{}
	im.t = etcIM[53]
	im.r = trscale(im.t, 4)
	lev3plants = append(lev3plants, im)
	im.t = etcIM[55]
	im.r = trscale(im.t, 4)
	lev3plants = append(lev3plants, im)
	im.t = etcIM[56]
	im.r = trscale(im.t, 4)
	lev3plants = append(lev3plants, im)
	im.t = etcIM[58]
	im.r = trscale(im.t, 4)
	lev3plants = append(lev3plants, im)

	count := 0
	for range hnum {
		for i := range wnum {
			r := R(x, y, blokw)
			b.r = r
			if count == 0 {
				blokadd(&b, 3)
				if ROLL(6) > 4 && i+2 < wnum {
					y2 := y //PLANTS BG
					choose := RINT(0, 4)
					y2 -= lev3plants[choose].r.Height
					lev3plants[choose].r.Y = y2
					lev3plants[choose].r.X = x
					bgIM = append(bgIM, lev3plants[choose])

				}

				if ROLL(36) > 34 && i+2 < wnum {
					y2 := y //WHEAT SUNFLOWER BERRY TREE
					choose := RINT(7, 10)
					y2 -= objList[choose].r.Height
					objList[choose].r.Y = y2
					objList[choose].r.X = x
					objList[choose].rC.Y = y2
					objList[choose].rC.X = (objList[choose].r.X + objList[choose].r.Width/2) - objList[choose].rC.Width/2
					obj = append(obj, objList[choose])
				}
			} else {
				if ROLL(120) == 120 {
					blokadd(&b, 4)
				} else {
					blokadd(&b, 1)
				}
			}
			x += blokw
		}
		count = 1
		y += blokw
		change := RINT(2, 5)
		wnum -= change
		x = ox
		x += f32(change/2) * blokw
		ox = x
	}
	x = ex
	return x

}

// MARK:  ENEMIES
func mENMlist() {
	o := OBJ{anm: []ANIM{anms[0]}, nm: "green slime", r: R(CNT.X, groundH-un3, un3), hp: 2}
	enmList = append(enmList, o)
}
func mENM() {
	switch levnum {
	case 0: //PLAINS
		var x, y float32 = 0, groundH - (enmList[0].r.Height - unh) //GREEN SLIME
		for x < gameW {
			if ROLL(24) > 20 {
				e := enmList[0]
				e.anm = make([]ANIM, len(enmList[0].anm))
				copy(e.anm, enmList[0].anm)
				e.anm[0].frameActive = RINT(0, len(e.anm[0].frames))
				e.r.X = x
				e.r.Y = y
				enms = append(enms, e)
			}
			x += enmList[0].r.Width
			x += RF32(un2, un5)
		}
	}

}
func enmDestroy(e *OBJ) {
	switch e.nm {
	case "green slime":
		anm := anms[7]
		anm.rD = recfromcnt(rcnt(e.r), un8)
		anm.rD.Y -= anm.rD.Height / 5
		anm.once = true
		//anm.alpha = 220
		fx = append(fx, anm)
	}

	e.off = true
}

// MARK: MAKE NPC
func mNPClist() {
	o := OBJ{anm: []ANIM{anms[1]}, nm: "intro dude", r: R(CNT.X, groundH-un2, un2), txt: []string{"welcome to zombie island", "i am intro dude", "you will die many times", "play anyway you like", "learn from mistakes"}, txT: FPS}
	o.rC = rlrgr(o.r, un)
	objRCratio(&o)
	npcList = append(npcList, o) //000 INTRO DUDE
}
func mNPC() {
	o := npcList[0] //000 INTRO DUDE
	o.r = rlrgr(pl.r, un8th)
	o.r.X -= un2
	o.r.Y = groundH - o.r.Height
	uOBJrecs(&o)
	npc = append(npc, o)
}

// MARK: UP CAM
func uCAM() {
	CAM.Target = pl.cnt
	CAM.Offset = V2(scrwf32/2, scrhf32/2)
	SCRREC = rechfromcnt(pl.cnt, float32(SCRW)+un2, float32(SCRH)+un2)
	if pl.cnt.Y > scrhf32/2 {
		CAM.Target.Y = scrhf32 / 2
		SCRREC.Y = 0
	}

}

// MARK: INITIAL
func INITIAL() {
	CAM.Target = CNT
	CAM.Offset = z.NewVector2(float32(SCRW/2), float32(SCRH/2))
	CAM.Rotation = 0.0
	CAM.Zoom = 1

	fnmed = z.LoadFontEx("res/fonts/PatrickHandSC-Regular.ttf", 36, nil)
	fnsml = z.LoadFontEx("res/fonts/PatrickHandSC-Regular.ttf", 28, nil)
	fn = z.LoadFontEx("res/fonts/PatrickHandSC-Regular.ttf", int32(un), nil)

	z.SetTextureFilter(fnsml.Texture, z.FilterBilinear)
	z.SetTextureFilter(fn.Texture, z.FilterBilinear)

	DBGREC = RH(float32(DRAWW)-300, 0, 300, float32(DRAWH))
	SCRREC = RH(0, 0, float32(SCRW), float32(SCRH))
	SCRRECNOCAM = SCRREC

	//IMG
	etcIM = mIMpath("res/im/etc")
	splatsIM = mIMpath("res/im/splats")

	//ANIM
	mANM()

	//ETC
	o := OBJ{}
	o.anm = append(o.anm, anms[2])
	o.r = R(0, 0, un)
	flags = append(flags, o) //FLAG MAGENTA
	o = OBJ{}
	o.anm = append(o.anm, anms[3])
	o.r = R(0, 0, un)
	flags = append(flags, o) //FLAG RED
	o = OBJ{}
	o.anm = append(o.anm, anms[4])
	o.r = R(0, 0, un)
	flags = append(flags, o) //FLAG YELLOW
	o = OBJ{}
	o.anm = append(o.anm, anms[5])
	o.r = R(0, 0, un)
	flags = append(flags, o) //FLAG CYAN

	//MAKE
	mBLOKTYPES()
	mOBJ()
	mGAME()
	mCRATES()
	mSPRINGS()
	uCAM()
	mCLOUDS()
	mENMlist()
	mENM()
	mNPClist()
	mNPC()

}

// MARK: BLOKS
func cblokunder(r z.Rectangle) (bool, float32) {
	found := false
	var y float32
	p := rcnt(r)
	p.Y = r.Y + r.Height + unq
	for i := range bloksSCR {
		if bloksSCR[i].solid {
			if crp(bloksSCR[i].r, p) {
				found = true
				y = bloksSCR[i].r.Y
			}
		}
	}
	return found, y
}
func bloksplat(b *BLOK) {
	b.rsplat = rsmlr(b.r, 8)
	b.splat = true
	b.csplat = RANGREEN()
	if ROLL(6) == 6 {
		b.csplat = RANCOLOR()
	}
	b.tsplat = splatsIM[RINT(0, len(splatsIM)-1)]
}
func blokadd(b *BLOK, num int) { //MARK: BLOKADD
	switch num {
	case 4: //ISLANDS LEVEL 3 POWERUP
		b.c = RANCOLOR()
		b.t = etcIM[63]
		b.powerup = true
		bloks = append(bloks, *b)
		b.powerup = false
	case 3: //ISLANDS LEVEL 3 TOP LAYER
		b.c = RANBROWNBRIGHT()
		b.t = etcIM[60]
		bloks = append(bloks, *b)
	case 2: //GROUND
		b.c = RANSEPIA()
		if ROLL(12) > 9 {
			switch RINT(1, 4) {
			case 1:
				b.t = etcIM[2]
			case 2:
				b.t = etcIM[51]
			case 3:
				b.t = etcIM[52]
			}
		}
		bloksGround = append(bloksGround, *b)
		b.t = etcIM[1]
	case 1: //ISLANDS LEVEL 3
		b.c = RANBEACHSAND()
		b.t = etcIM[62]
		if ROLL(6) == 6 {
			b.t = etcIM[61]
		}
		bloks = append(bloks, *b)
	case 0: //PLATFORMS LEVEL 1
		b.t = bloktex(0)
		b.c = RANORANGEBRIGHT()
		if ROLL(12) > 10 {
			bloksplat(b)
		}
		bloks = append(bloks, *b)
		b.splat = false
	}

}
func bloktex(num int) z.Texture2D {
	switch num {
	case 0: //PLATFORMS
		return etcIM[30+RINT(0, 3)]
	default:
		return etcIM[30]
	}
}
func mBLOKSHAPE(num int, x, y float32) float32 {
	b := blokTypes[0]
	b.r.X = x
	b.r.Y = y
	b.t = bloktex(0)
	b.c = RANORANGEBRIGHT()
	bloks = append(bloks, b)
	switch num {
	case 8: //ZIGZAG
		numw := RINT(3, 8)
		for range numw {
			b.r.Y -= blokw
			b.r.X += blokw
			blokadd(&b, 0)
			if FLIP() {
				b.r.X += blokw
				blokadd(&b, 0)
				x += blokw
			}
			b.r.Y += blokw
			b.r.X += blokw
			blokadd(&b, 0)
			x += blokw * 3
		}

	case 7: //U RAND HEIGHTS
		b.r.Y -= blokw
		blokadd(&b, 0)
		b.r.Y += blokw
		b.r.X += blokw
		blokadd(&b, 0)
		b.r.X += blokw
		blokadd(&b, 0)
		if FLIP() {
			b.r.X += blokw
			blokadd(&b, 0)
			x += blokw
		}
		if FLIP() {
			b.r.X += blokw
			blokadd(&b, 0)
			x += blokw
		}
		b.r.Y -= blokw
		blokadd(&b, 0)
		b.splat = false

		x += blokw * 5

	case 6: //U RAND HEIGHTS
		b.r.Y -= blokw
		blokadd(&b, 0)
		b.splat = false
		if FLIP() {
			b.r.Y -= blokw
			blokadd(&b, 0)
			b.r.Y += blokw
		}
		b.r.Y += blokw
		b.r.X += blokw
		blokadd(&b, 0)
		b.splat = false
		b.r.X += blokw
		blokadd(&b, 0)
		b.splat = false
		if FLIP() {
			b.r.X += blokw
			blokadd(&b, 0)
			x += blokw
		}
		if FLIP() {
			b.r.X += blokw
			blokadd(&b, 0)
			x += blokw
		}
		b.r.Y -= blokw
		blokadd(&b, 0)
		if FLIP() {
			b.r.Y -= blokw
			blokadd(&b, 0)
			b.r.Y += blokw
		}

		x += blokw * 5

	case 5: //2 BLOK LADDER
		b.r.X += blokw
		blokadd(&b, 0)
		b.r.Y -= blokw * 3
		blokadd(&b, 0)
		b.r.X -= blokw
		blokadd(&b, 0)
		b.r.Y -= blokw * 3
		blokadd(&b, 0)
		b.r.X += blokw
		blokadd(&b, 0)
		if FLIP() {
			b.r.Y -= blokw * 3
			blokadd(&b, 0)
			b.r.X -= blokw
			blokadd(&b, 0)
			if FLIP() {
				b.r.Y -= blokw * 3
				blokadd(&b, 0)
				b.r.X += blokw
				blokadd(&b, 0)
			}
		}

		x += blokw * 4

	case 4: //4 BLOK
		b.r.X += blokw
		blokadd(&b, 0)
		b.r.Y -= blokw
		blokadd(&b, 0)
		b.r.X -= blokw
		blokadd(&b, 0)

		x += blokw * 6

	case 3: //GAP PLATFORM
		b.r.X += blokw * 2
		blokadd(&b, 0)
		b.r.X += blokw * 2
		blokadd(&b, 0)
		b.r.X += blokw * 2
		blokadd(&b, 0)
		x += blokw * 8
		if FLIP() {
			b.r.X += blokw * 2
			blokadd(&b, 0)
			x += blokw * 2
		}
		if FLIP() {
			b.r.X += blokw * 2
			blokadd(&b, 0)
			x += blokw * 2
		}

	case 2: //STEPS UP RIGHT
		b.r.Y -= blokw
		b.r.X += blokw
		blokadd(&b, 0)
		b.r.Y -= blokw
		b.r.X += blokw
		blokadd(&b, 0)
		b.r.Y -= blokw
		b.r.X += blokw
		blokadd(&b, 0)
		x += blokw * 6
		if FLIP() {
			b.r.Y -= blokw
			b.r.X += blokw
			blokadd(&b, 0)
			x += blokw
		}
		if FLIP() {
			b.r.Y -= blokw
			b.r.X += blokw
			blokadd(&b, 0)
			x += blokw
		}

	case 1: //CROSS LRG
		b.r.Y -= blokw
		blokadd(&b, 0)
		b.r.Y -= blokw
		blokadd(&b, 0)
		b.r.Y -= blokw
		blokadd(&b, 0)
		b.r.Y -= blokw
		blokadd(&b, 0)
		b.r.Y += blokw * 2
		b.r.X -= blokw
		blokadd(&b, 0)
		b.r.X -= blokw
		blokadd(&b, 0)
		b.r.X += blokw * 3
		blokadd(&b, 0)
		b.r.X += blokw
		blokadd(&b, 0)

		x += blokw * 7

	case 0: //CROSS SML

		b.r.Y -= blokw
		blokadd(&b, 0)
		b.r.Y -= blokw
		blokadd(&b, 0)
		b.r.Y += blokw
		b.r.X -= blokw
		blokadd(&b, 0)
		b.r.X += blokw * 2
		blokadd(&b, 0)

		x += blokw * 5
	}
	return x
}

// MARK: MAKE BLOK TYPES
func mBLOKTYPES() {
	b := BLOK{ //GROUND PLATFORMS 00
		c:     RANSEPIA(),
		t:     etcIM[1],
		r:     R(0, 0, blokw),
		solid: true,
	}
	blokTypes = append(blokTypes, b)

	b = BLOK{ //MOVING PLATFORMS 01
		c:      z.White,
		t:      etcIM[49],
		r:      R(0, 0, blokw),
		moving: true,
		solid:  true,
	}
	blokTypes = append(blokTypes, b)
}

func mCLOUDS() {
	clouds = nil
	num := RINT(3, 8)

	for range num {
		o := OBJ{}
		o.t = etcIM[cloudnums[RINT(0, len(cloudnums))]]
		o.r = trscale(o.t, f32(RINT(5, 12)))
		o.lr = FLIP()
		o.flip = FLIP()
		o.r.Y = RF32(unh, un4)
		//o.r.Y = CNT.Y
		if o.lr {
			o.r.X = -(o.r.Width + un)
		} else {
			o.r.X = scrwf32 + un
		}
		o.spd = RF32(4, un8th)
		clouds = append(clouds, o)
	}
}
