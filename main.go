package main

import (
	"fmt"
	"math"

	z "github.com/gen2brain/raylib-go/raylib"
)

var (

	//UI
	msInUI, editName bool

	//CORE
	FPS, SCRW, SCRH, DRAWW, DRAWH int32 = 60, 1920, 1080, 1920, 1080
	scrhf32, scrwf32                    = f32(SCRH), f32(SCRW)

	SCRREC, SCRRECNOCAM, DBGREC z.Rectangle
	CNT                         z.Vector2
	RNDTEX                      z.RenderTexture2D
	countDraw, countCull        int
	delta                       float32
	EXIT                        = false
)

func main() {
	z.SetConfigFlags(z.FlagVsyncHint)
	z.InitWindow(DRAWW, DRAWH, "ZOMBIE ISLAND 12/10/2026")
	z.HideCursor()
	z.SetExitKey(z.KeyNull)

	RNDTEX = z.LoadRenderTexture(SCRW, SCRH)
	defer z.UnloadRenderTexture(RNDTEX)
	z.SetTextureFilter(RNDTEX.Texture, z.FilterBilinear)
	CNT = z.NewVector2(float32(SCRW/2), float32(SCRH/2))

	INITIAL()

	z.SetTargetFPS(FPS)

	for !z.WindowShouldClose() && !EXIT {
		delta = z.GetFrameTime()
		UP()
		//DRAW TO TEXTURE
		z.BeginTextureMode(RNDTEX)
		z.ClearBackground(z.SkyBlue)
		z.DrawRectangleGradientV(0, 0, SCRW, SCRH, c.CeruleanTeal, z.SkyBlue) //GRADIENT BG
		z.BeginMode2D(CAM)                                                    //DRAW CAM

		//MARK: DRAW GAME DROPPED OBJ
		if len(obj) > 0 {
			clear := false
			for i := range obj {
				if obj[i].off {
					clear = true
				} else {
					o := &obj[i]
					if crr(o.r, SCRREC) {
						if o.hitTimer > 0 {
							o.hitTimer -= delta
							if o.hitTimer < 0 {
								o.hitTimer = 0
							}
							if o.hitTimer <= 0 {
								switch o.nm {
								case "crate":
									anm := anms[6]
									anm.rD = recfromcnt(rcnt(o.r), un7)
									anm.once = true
									anm.alpha = 100
									fx = append(fx, anm)
									spawnOBJ(1, o.r)
									o.off = true
								}
							}
							switch o.nm {
							case "crate", "tree", "wheat", "berry tree", "sunflower": //OBJ SWAY ON HIT
								angle := float32(math.Sin(float64(o.hitTimer*40.0))) * 6.0
								z.DrawTexturePro(
									o.t,
									RH(0, 0, float32(o.t.Width), float32(o.t.Height)),
									z.Rectangle{X: o.r.X + o.r.Width/2, Y: o.r.Y + o.r.Height, Width: o.r.Width, Height: o.r.Height},
									V2(o.r.Width/2, o.r.Height),
									angle,
									RANCOLOR(),
								)
							}
						} else { //DRAW OBJ
							if o.collect {
								r := o.r //DROPPED OBJ ARROW
								r.Y -= r.Height + un16th + ychngOBJarrow
								dt(r, etcIM[67], z.White)
							}
							if o.danim {
								o.anm[0].DrawFrames(o.r, false, false, z.White)
								switch o.nm {
								case "spring":
									if crr(pl.r, o.r) {
										pl.spring = true
									}
								}
							} else {
								dt(o.r, o.t, z.White)
							}

						}
						if pl.catActiv == 1 {
							if crr(pl.rTool, o.r) {
								o.intool = true
								if pl.useitem {
									if quickOBJ[quickActiv].nm == "axe" {
										switch o.nm {
										case "tree", "crate", "wheat", "berry tree", "sunflower":
											if o.hitTimer == 0 {
												o.hitTimer = 0.15
											}
										}
									}
								}
								if quickOBJ[quickActiv].useT == quickOBJ[quickActiv].useTdef {
									switch o.nm {
									case "crate":
										anm := anms[6]
										anm.rD = recfromcnt(rcnt(o.r), un7)
										anm.once = true
										anm.alpha = 100
										fx = append(fx, anm)
										spawnOBJ(1, o.r)
										o.off = true
									case "tree":
										spawnOBJ(0, o.r)
									case "wheat":
										spawnOBJ(2, o.r)
									case "berry tree":
										spawnOBJ(3, o.r)
									case "sunflower":
										spawnOBJ(4, o.r)
									}
									quickOBJ[quickActiv].useT = 0
								}
							} else {
								o.intool = false
							}
						}
						if debug {
							drl(o.r, 2, z.Orange)
							drl(o.rC, 2, z.Blue)
						}
					}

					//GRAVITY
					if !o.nograv {
						if o.r.Y+o.r.Height < groundH {
							found, yOBJ := cblokunder(o.r)
							if found {
								o.r.Y = yOBJ - o.r.Height
							} else {
								o.r.Y += deltaspd(grav)
							}
						} else if o.r.Y+o.r.Height > groundH {
							o.r.Y = groundH - o.r.Height
						}
					}

					//COLLECT
					if !o.off {
						if o.collect && invtotal < len(pl.inv) {
							if crr(o.r, pl.rC) {
								collect(o)
							}
						} else if o.collect && invtotal >= len(pl.inv) {
							mwarntxt("inventory is full")
						}
					}
				}
			}
			if clear {
				alive := obj[:0]
				for _, o := range obj {
					if !o.off {
						alive = append(alive, o)
					}
				}
				obj = alive
			}
		}
		//MARK: DRAW BLOKS PLATFORM
		bloksSCR = nil
		if len(bloks) > 0 {
			clear := false
			for i := len(bloks) - 1; i >= 0; i-- {
				if bloks[i].off {
					clear = true
				} else {
					b := &bloks[i]
					if crr(b.r, SCRREC) {
						bloksSCR = append(bloksSCR, b)
						if b.powerup {
							dt(b.r, b.t, ca(YELLOWBRIGHT(), fadeBlink))
						} else {
							dt(b.r, b.t, b.c)
							if b.splat {
								dt(b.r, b.tsplat, ca(b.csplat, 200))
							}
						}
						//FLAG DROP
						if !b.flag {
							if flagdropActiv {
								if cmsc(b.r) {
									drlround(b.r, 2, RANCOLOR())
									if mslc {
										dropflag(b.r)
										b.flag = true
										flagdropActiv = false
									}
								}
							}
						}
					}
					//UPDATE OFF SCREEN BLOKS
					if b.velx != 0 || b.vely != 0 {
						if b.vely != 0 {
							b.r.Y += deltaspd(b.vely)
							if b.r.Y+b.r.Height >= groundH || b.r.Y <= groundH-b.moveVmax {
								b.vely *= -1
							}
						}
					}
				}
			}
			if clear {
				alive := bloks[:0]
				for _, b := range bloks {
					if !b.off {
						alive = append(alive, b)
					}
				}
				bloks = alive
			}
		}
		//MARK: DRAW DROPPED FLAGS
		if len(flagsDropped) > 0 {
			for i := len(flagsDropped) - 1; i >= 0; i-- {
				if crr(flagsDropped[i].r, SCRREC) {
					flagsDropped[i].anm[0].DrawFrames(flagsDropped[i].r, false, false, z.White)
					if cmsc(flagsDropped[i].r) {
						r2 := flagsDropped[i].r
						r2.X -= r2.Width / 4
						drlround(r2, 2, RANCOLOR())
						minfotxt("right click destroy flag")
						if msrc {

						}
					}
				}
			}
		}
		//MARK: DRAW PLAYER PROJECTILES
		for i := range projPL {
			// Draw projectile (using a texture or solid rectangle)
			if projPL[i].t.ID != 0 {
				dt(projPL[i].r, projPL[i].t, projPL[i].c)
			} else {
				z.DrawRectangleRec(projPL[i].r, z.Yellow) // Default yellow tracer
			}

			if debug {
				drl(projPL[i].rC, 1, z.Red)
			}
		}
		//MARK: DRAW ENEMIES
		enmSCR = nil
		if len(enms) > 0 {
			clear := false
			for i := len(enms) - 1; i >= 0; i-- {
				e := &enms[i]
				if e.off {
					clear = true
				} else {
					if crr(e.r, SCRREC) {
						enmSCR = append(enmSCR, e)
						if e.hitTimer == 0 {
							e.anm[0].DrawFrames(e.r, false, false, z.White)
						} else {
							if ROLL(6) > 3 {
								e.anm[0].DrawFrames(e.r, false, false, RANREDBRIGHT())
							} else {
								e.anm[0].DrawFrames(e.r, false, false, z.White)
							}
						}

						if debug {
							drl(e.r, 2, z.Blue)
							drl(e.rC, 2, z.Magenta)
							if cmsc(e.anm[0].rC[e.anm[0].frameActive]) {
								dtxtdef(fmt.Sprint(e.hp), e.r.X, e.r.Y-f32(fnsml.BaseSize))
							}
						}
					}
					uOBJFrameRC(e, 0, e.r)
				}
				if e.hitTimer > 0 {
					e.hitTimer--
				}
			}

			if clear {
				alive := enms[:0]
				for _, e := range enms {
					if !e.off {
						alive = append(alive, e)
					}
				}
				enms = alive
			}
		}
		//MARK: DRAW NPC
		if len(npc) > 0 {
			for i := len(npc) - 1; i >= 0; i-- {
				if !npc[i].off {
					if crr(npc[i].r, SCRREC) {
						npc[i].anm[0].DrawTexNums(npc[i].r, false, false)
						if debug {
							drl(npc[i].r, 2, z.Blue)
							drl(npc[i].rC, 2, z.Magenta)
						}

						if crr(pl.rC, npc[i].rC) {
							switch npc[i].nm {
							case "intro dude":
								txx := npc[i].r.X + npc[i].r.Width/2
								txy := npc[i].r.Y - (float32(fnmed.BaseSize + 2))
								dtxtcntrv2BGrecround(fnmed, npc[i].txt[npc[i].txNumActiv], V2(txx, txy), 10, c.VoidBlack, c.IvoryWhite, c.VoidBlack, true)
								if npc[i].txT > 0 {
									npc[i].txT--
								} else {
									if npc[i].txNumActiv < len(npc[i].txt)-1 {
										npc[i].txNumActiv++
										npc[i].txT = FPS
									} else {
										npc[i].txNumActiv = 0
										npc[i].txT = FPS
									}
								}
							}
						} else {
							npc[i].txNumActiv = 0
						}
					}
				}
			}
		}

		//MARK: DRAW FX
		if len(fx) > 0 {
			clear := false
			for i := range fx {
				if len(fx[i].framesTexNums) > 0 {
					ended := fx[i].DrawTexNums(fx[i].rD, fx[i].once, false)
					if ended {
						fx[i].off = true
					}
				} else {
					ended := fx[i].DrawFrames(fx[i].rD, fx[i].once, false, z.White)
					if ended {
						fx[i].off = true
					}
				}
				if fx[i].off {
					clear = true
				}
			}
			if clear {
				alive := fx[:0]
				for _, f := range fx {
					if !f.off {
						alive = append(alive, f)
					}
				}
				fx = alive
			}
		}
		//MARK: DRAW PLAYER
		switch pl.catActiv {
		case 1: // TOOL
			if pl.useitem {
				rTool := pl.r
				rTool.X += pl.r.Width / 2
				rTool.Y += pl.r.Width / 2
				if ms.X < pl.cnt.X {
					pl.lr = true
				} else {
					pl.lr = false
				}
				switch quickOBJ[quickActiv].nm {
				case "axe":

				}

				if pl.lr {
					dtfliproBR(rTool, quickOBJ[quickActiv].t, z.White, pl.rotool)
				} else {
					dtroBL(rTool, quickOBJ[quickActiv].t, z.White, pl.rotool)
				}
				quickOBJ[quickActiv].useT++
			} else {
				quickOBJ[quickActiv].useT = 0
			}
		case 2: //WEAPON
			weapR := rlrgr(pl.r, 8)
			weapR.Y += 4
			v2 := pl.cnt
			if pl.useitem {
				if pl.lr {
					v2.X = weapR.X + weapR.Width
					v2.Y = weapR.Y + weapR.Height/2
					weapR.X -= (weapR.Width - (unq + 2))
					ro := anglp(v2, ms)
					dtroRCNT(weapR, quickOBJ[quickActiv].t, z.White, ro)
					if debug {
						drl(weapR, 2, z.Magenta)
					}
				} else {
					v2.X = weapR.X
					v2.Y = weapR.Y + weapR.Height/2
					weapR.X += (weapR.Width - (unq + 2))
					ro := anglp(v2, ms)
					dtroLCNT(weapR, quickOBJ[quickActiv].t, z.White, ro)
					if debug {
						drl(weapR, 2, z.Magenta)
					}
				}
				//quickOBJ[quickActiv].useT++
				//	if quickOBJ[quickActiv].useT == quickOBJ[quickActiv].useTdef {
				useweap()
				//	quickOBJ[quickActiv].useT = 0
				//}
			}
			if quickOBJ[quickActiv].useT > 0 {
				quickOBJ[quickActiv].useT--
			}
		}

		if pl.moving || pl.useitem {
			if pl.moving && pl.useitem || pl.moving && !pl.useitem {
				if pl.lr {
					pl.anm[1].DrawTexNums(pl.r, false, true)
				} else {
					pl.anm[1].DrawTexNums(pl.r, false, false)
				}
			} else if !pl.moving && pl.useitem {
				if pl.lr {
					pl.anm[2].DrawTexNums(pl.r, false, true)
				} else {
					pl.anm[2].DrawTexNums(pl.r, false, false)
				}
			}

		} else {
			pl.anm[0].DrawTexNums(pl.r, false, false)
		}

		if debug {
			drl(pl.r, 2, z.Blue)
			drl(pl.rTool, 2, z.Orange)
			drl(pl.rC, 2, z.Magenta)
			dr(recfromcnt(pl.cnt, 4), z.Magenta)
		}
		//MARK: DRAW BLOKS GROUND
		if len(bloksGround) > 0 {
			for i := len(bloksGround) - 1; i >= 0; i-- {
				if crr(bloksGround[i].r, SCRREC) {
					b := bloksGround[i]
					dt(b.r, b.t, b.c)
				}
			}
		}
		//MARK: DRAW GRASS
		if len(bgIM) > 0 {
			for i := range bgIM {
				if crr(bgIM[i].r, SCRREC) {
					dia(bgIM[i], 180)
				}
			}
		}

		z.EndMode2D()
		//DRAW NO CAM IN RENDERTEX

		switch levnum {
		case 0: //PLAINS
			//SEAGULL
			if seagullT == 0 {
				if seagullLR {
					seagullanm.DrawFrames(seagullR, false, false, z.White)
				} else {
					seagullanm.DrawFrames(seagullR, false, true, z.White)
				}

			}

			//CLOUDS
			if len(clouds) > 0 && cloudsActiv {
				cloudson := false
				for i := len(clouds) - 1; i >= 0; i-- {
					if !clouds[i].off {
						cloudson = true
						if crr(clouds[i].r, SCRRECNOCAM) {
							if clouds[i].flip {
								dtflip(clouds[i].r, clouds[i].t, z.White)
							} else {
								dt(clouds[i].r, clouds[i].t, z.White)
							}
						}
						if clouds[i].lr {
							clouds[i].r.X += deltaspd(clouds[i].spd)
							if clouds[i].r.X > scrwf32 {
								clouds[i].off = true
							}
						} else {
							clouds[i].r.X -= deltaspd(clouds[i].spd)
							if clouds[i].r.X+clouds[i].r.Width < 0 {
								clouds[i].off = true
							}
						}
					}
				}

				if !cloudson {
					cloudsActiv = false
					mCLOUDS()
					cloudsActiv = true
				}
			}
		}
		//MARK: DRAW INVENTORY
		if invenActiv {
			siz := un + unh
			w := f32(len(pl.inv)/4) * (siz + spc)
			w -= spc
			h := (siz + spc) * 4
			h -= spc
			x := CNT.X - w/2
			ox := x
			y := CNT.Y - h/2
			count := 0
			for i := range pl.inv {
				r := R(x, y, siz)
				dr(r, ca(c.VoidBlack, 120))
				if pl.inv[i].nm != "" {
					if cms(r) {
						drl(r, 2, RANCOLOR())
						minfotxt(pl.inv[i].nm)
						if mslc {
							if len(quickOBJ) < 10 {
								inven2quick(i)
							} else {
								mwarntxt("quick slots full")
							}
						}
					}
					r2 := rsmlr(r, 4)
					dt(r2, pl.inv[i].t, z.White)
					if pl.inv[i].numinven > 0 {
						dtxtc(fnsml, fmt.Sprint(pl.inv[i].numinven+1), r.X+4, r.Y, z.White)
					}
				}
				x += siz + spc
				count++
				if count == 5 {
					count = 0
					x = ox
					y += siz + spc
				}

			}
		}
		//QUICK SLOTS
		msInUI = false
		pl.catActiv = -1
		var x, y, siz float32 = un8th, un8th, unh + unq
		for i := range quicklen { //QUICK ICONS
			r := R(x, y, siz)
			clear := false
			clearnum := 0
			if cms(r) {
				drlround(r, 2, RANCOLOR())
				if mslc {
					quickActiv = i
				}
			}
			drround(r, ca(c.VoidBlack, 50))
			if quickActiv == i {
				drlround(r, 2, z.Magenta)
			}
			if len(quickOBJ) > 0 {
				for j := range quickOBJ {
					if i == j {
						if quickActiv == i {
							pl.catActiv = quickOBJ[j].cat
						}
						dt(r, quickOBJ[j].t, z.White)
						if cms(r) {
							msInUI = true
							minfotxt(quickOBJ[j].nm)
							drlround(r, 2, RANCOLOR())
							if invenActiv {
								if mslc {
									added := addquick2inven(j)
									if added {
										if quickActiv == i {
											pl.catActiv = -1
											quickActiv = 0
										}

										clear = true
										clearnum = j
									} else {
										mwarntxt("inventory full")
									}
								}

							}
						}
						if quickOBJ[j].numinven > 1 {
							dquicknuminven(r, quickOBJ[j].numinven)
						}
					}
				}
			}
			if clear {
				quickOBJ, _ = addremOBJ(quickOBJ, quickOBJ[clearnum], 10)
			}

			y += siz + spc
		}
		//INVENTORY FLAGS
		r := R(x, y, siz)
		if cms(r) {
			if flagsAllPlaced {
				minfotxt("all flags dropped")
			} else {
				minfotxt("drop marker flag")
				if mslc {
					flagdropActiv = true
				}
			}
			drlround(r, 2, RANCOLOR())
		}
		drround(r, ca(c.VoidBlack, 50))
		r2 := rsmlr(r, 4)
		if flagsAllPlaced {
			dt(r2, etcIM[40], z.Red)
		} else {
			r2.X += 8
			flags[flagsnxtnum].anm[0].DrawFrames(r2, false, false, z.White)
		}
		//MARK: INFOTXT WARNTXT
		if infotxtT > 0 {
			dinfotxt()
			infotxtT--
		}
		if warntxtT > 0 {
			dwarntxt()
			warntxtT--
		}

		z.EndTextureMode() //END DRAW RENDERTEX

		//DRAW TO SCREEN
		z.BeginDrawing()

		z.ClearBackground(z.Black)
		z.DrawTexturePro(RNDTEX.Texture, RH(0, 0, float32(RNDTEX.Texture.Width), -float32(RNDTEX.Texture.Height)), RH(0, 0, float32(DRAWW), float32(DRAWH)), V2(0, 0), 0, z.White)

		//DRAW NO CAM NO RENDERTEX

		//PALETTE
		if palette {
			txt := ""
			drhxy(0, 0, DRAWW, DRAWH, z.Black)
			x, y := unq, unq
			for i := range COLORS {
				r := R(x, y, un)
				dr(r, COLORS[i])
				if cms(r) {
					txt = COLORNAMES[i] + " " + fmt.Sprint(i)
				}
				x += un + spc
				if x+un+spc >= float32(DRAWW) {
					x = unq
					y += un + spc
				}
			}
			if txt != "" {
				dtxtdef(txt, un, float32(DRAWH)-un2)
			}
		}
		//MARK: DEBUG
		if debug {
			dbgnum = invtotal
			drl(SCRREC, 4, z.Green) //SCREEN REC
			dr(DBGREC, ca(z.Maroon, 150))
			x, y := DBGREC.X+un8th, un8th
			dtxtdef("ms.x "+fmt.Sprintf("%.0f", ms.X)+" ms.y "+fmt.Sprintf("%.0f", ms.Y), x, y)
			y += float32(fnsml.BaseSize)
			dtxtdef("pl.cnt.X "+fmt.Sprintf("%.0f", pl.cnt.X)+" pl.cnt.Y "+fmt.Sprintf("%.0f", pl.cnt.Y), x, y)
			y += float32(fnsml.BaseSize)
			dtxtdef("CAM.Zoom "+fmt.Sprint(CAM.Zoom)+" dbgstring "+dbgstring, x, y)
			y += float32(fnsml.BaseSize)
			dtxtdef("dbgnum "+fmt.Sprint(dbgnum)+" dbgstring "+dbgstring, x, y)
			y += float32(fnsml.BaseSize)
			dtxtdef("len(projPL) "+fmt.Sprint(len(projPL))+" pl.onplat "+fmt.Sprint(pl.onplat), x, y)
			y += float32(fnsml.BaseSize)
			dtxtdef("pl.catActiv "+fmt.Sprint(pl.catActiv)+" len(bloks) "+fmt.Sprint(len(bloks)), x, y)
			y += float32(fnsml.BaseSize)
			dtxtdef("len(projPL) "+fmt.Sprint(len(projPL))+" len(fx) "+fmt.Sprint(len(fx)), x, y)
			y += float32(fnsml.BaseSize)
			dtxtdef("len(enmSCR) "+fmt.Sprint(len(enmSCR))+" len(enms) "+fmt.Sprint(len(enms)), x, y)
			y += float32(fnsml.BaseSize)
			dtxtdef("invtotal "+fmt.Sprint(invtotal)+" len(enms) "+fmt.Sprint(len(enms)), x, y)
			y += float32(fnsml.BaseSize)
		}

		//MARK: CURSOR
		if flagdropActiv {
			siz := un
			flags[flagsnxtnum].anm[0].DrawFrames(R(msnocam.X, msnocam.Y-siz, siz), false, false, z.White)
		}
		dt(R(msnocam.X, msnocam.Y, unh), etcIM[0], c.IvoryWhite)

		z.EndDrawing()
	}

	//MARK: UNLOAD
	z.UnloadFont(fn)
	z.UnloadFont(fnsml)
	z.UnloadFont(fnmed)
	for i := range etcIM {
		z.UnloadTexture(etcIM[i])
	}
	for i := range splatsIM {
		z.UnloadTexture(splatsIM[i])
	}
	z.CloseWindow()
}
