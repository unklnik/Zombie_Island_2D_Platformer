package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unsafe"

	z "github.com/gen2brain/raylib-go/raylib"
)

var (
	roundrec           float32 = 0.2
	palette            bool
	infotxtT, warntxtT int32
	infotxt, warntxt   string
)

// MARK: ANIM
func mANIMframes(t z.Texture2D, w, h, fps float32, framesnum int) ANIM {
	var anm ANIM
	anm.fps = fps
	anm.tex = t
	var x, y float32
	for range framesnum {
		anm.frames = append(anm.frames, RH(x, y, w, h))
		x += w
	}
	return anm
}
func (a *ANIM) DrawFrames(r z.Rectangle, once, flip bool, c z.Color) bool {
	// 1. Advance timer by delta time (frame time in seconds)
	a.timer += z.GetFrameTime()

	frameDuration := 1.0 / a.fps
	ended := false

	// 2. Advance frames when timer passes duration
	for a.timer >= frameDuration {
		a.timer -= frameDuration
		a.frameActive++

		if a.frameActive >= len(a.frames) {
			if once {
				a.frameActive = len(a.frames) - 1
				ended = true
				break
			} else {
				a.frameActive = 0
			}
		}
	}
	dcolor := c
	if a.alpha != 0 {
		dcolor = ca(z.White, a.alpha)
	}
	// 3. Draw the active frame texture
	if flip {
		dtrrflip(a.frames[a.frameActive], r, a.tex, dcolor)
	} else {
		dtrr(a.frames[a.frameActive], r, a.tex, dcolor)
	}

	if debug && len(a.rC) > 0 {
		/*
			scale := rrScale(a.frames[0], r)
			if len(a.rC) > 0 {
				x2 := r.X + (a.rC[a.frameActive].X * scale)
				y2 := r.Y + (a.rC[a.frameActive].Y * scale)
				w := a.rC[a.frameActive].Width * scale
				h := a.rC[a.frameActive].Height * scale
				drl(RH(x2, y2, w, h), 2, z.Magenta)
			}
		*/
		drl(a.rC[a.frameActive], 2, z.Magenta)
	}

	return ended
}
func (a *ANIM) DrawTexNums(r z.Rectangle, once, flip bool) bool {
	// 1. Advance timer by delta time (frame time in seconds)
	a.timer += z.GetFrameTime()

	frameDuration := 1.0 / a.fps
	ended := false

	// 2. Advance frames when timer passes duration
	for a.timer >= frameDuration {
		a.timer -= frameDuration
		a.frameActive++

		if a.frameActive >= len(a.framesTexNums) {
			if once {
				a.frameActive = len(a.framesTexNums) - 1
				ended = true
				break
			} else {
				a.frameActive = 0
			}
		}
	}
	dcolor := z.White
	if a.alpha != 0 {
		dcolor = ca(z.White, a.alpha)
	}
	// 3. Draw the active frame texture
	if flip {
		dtflip(r, etcIM[a.framesTexNums[a.frameActive]], dcolor)
	} else {
		dt(r, etcIM[a.framesTexNums[a.frameActive]], dcolor)
	}

	return ended
}
func mANIMtexnums(tNums []int, fps float32) ANIM {
	return ANIM{framesTexNums: tNums, fps: fps}
}
func dANMoneframe(a ANIM, r z.Rectangle) {
	dtrr(a.frames[0], r, a.tex, z.White)
}
func dANMallframesRC(a ANIM, x, y, siz float32) {
	r := R(x, y, siz)
	scale := rrScale(a.frames[0], r)
	for i := range a.frames {
		dtrr(a.frames[i], r, a.tex, z.White)
		r.X += siz
	}
	r = R(x, y, siz)
	if len(a.rC) > 0 {
		for i := range a.rC {
			x2 := r.X + (a.rC[i].X * scale)
			y2 := r.Y + (a.rC[i].Y * scale)
			w := a.rC[i].Width * scale
			h := a.rC[i].Height * scale
			drl(RH(x2, y2, w, h), 2, z.Magenta)
			dtxtdef(fmt.Sprint(i), x2, y2+h)
			r.X += siz
		}
	}

}
func deltaspd(spd float32) float32 {
	return (spd * float32(FPS)) * delta
}
func mRCmultiframe(path2png string, wFrame, hFrame float32) []z.Rectangle {
	var rc []z.Rectangle
	file, err := os.Open(path2png)
	if err == nil {
		defer file.Close()
		imgData, _ := png.Decode(file)
		texture := z.LoadTexture(path2png)
		var x float32 = 0
		for x+wFrame <= float32(texture.Width) {
			//	ims.r = append(ims.r, RH(x, 0, w, h))
			tightBox := getTightBounds(imgData, int(x), 0, int(wFrame), int(hFrame))
			rc = append(rc, tightBox)
			x += wFrame
		}
	}
	return rc
}

func getTightBounds(imgData image.Image, fx, fy, fw, fh int) z.Rectangle {
	minX, maxX := fw, 0
	minY, maxY := fh, 0
	hasPixels := false

	for y := 0; y < fh; y++ {
		for x := 0; x < fw; x++ {
			// Get pixel alpha (0 - 65535 in Go image)
			_, _, _, a := imgData.At(fx+x, fy+y).RGBA()
			if a > 0 { // Pixel is not completely transparent
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
				hasPixels = true
			}
		}
	}

	// Fallback if the whole frame is empty
	if !hasPixels {
		return z.Rectangle{X: 0, Y: 0, Width: float32(fw), Height: float32(fh)}
	}

	return z.Rectangle{
		X:      float32(minX),
		Y:      float32(minY),
		Width:  float32(maxX - minX + 1),
		Height: float32(maxY - minY + 1),
	}
}

// MARK: IMG TEX
func di(im IM) {
	dt(im.r, im.t, z.White)
}
func dia(im IM, a uint8) {
	dt(im.r, im.t, ca(z.White, a))
}
func tr(t z.Texture2D) z.Rectangle {
	return RH(0, 0, float32(t.Width), float32(t.Height))
}
func trscale(t z.Texture2D, scale float32) z.Rectangle {
	return RH(0, 0, float32(t.Width)*scale, float32(t.Height)*scale)
}
func dt(r z.Rectangle, t z.Texture2D, c z.Color) {
	z.DrawTexturePro(t, RH(0, 0, float32(t.Width), float32(t.Height)), r, z.Vector2Zero(), 0, c)
}
func dtrr(rsrc, rdest z.Rectangle, t z.Texture2D, c z.Color) {
	z.DrawTexturePro(t, rsrc, rdest, z.Vector2Zero(), 0, c)
}
func dtrrflip(rsrc, rdest z.Rectangle, t z.Texture2D, c z.Color) {
	src := RH(rsrc.X, rsrc.Y, -float32(rsrc.Width), float32(rsrc.Height))
	z.DrawTexturePro(t, src, rdest, z.Vector2Zero(), 0, c)
}
func dtroBL(r z.Rectangle, t z.Texture2D, c z.Color, ro float32) {
	z.DrawTexturePro(t, RH(0, 0, float32(t.Width), float32(t.Height)), r, V2(0, r.Height), ro, c)
}
func dtroLCNT(r z.Rectangle, t z.Texture2D, c z.Color, ro float32) {
	r.Y += r.Height / 2
	z.DrawTexturePro(t, RH(0, 0, float32(t.Width), float32(t.Height)), r, V2(0, r.Height/2), ro, c)
}
func dtroRCNT(r z.Rectangle, t z.Texture2D, c z.Color, ro float32) {
	src := RH(0, 0, -float32(t.Width), float32(t.Height))
	r.Y += r.Height / 2
	r.X += r.Width
	z.DrawTexturePro(t, src, r, V2(r.Width, r.Height/2), ro+180, c)
}
func dtfliproBR(r z.Rectangle, t z.Texture2D, c z.Color, ro float32) {
	src := RH(0, 0, -float32(t.Width), float32(t.Height))
	z.DrawTexturePro(t, src, r, V2(r.Width, r.Height), ro, c)
}
func dtflip(r z.Rectangle, t z.Texture2D, c z.Color) {
	src := RH(0, 0, -float32(t.Width), float32(t.Height))
	z.DrawTexturePro(t, src, r, z.Vector2Zero(), 0, c)
}
func dtflipex(r z.Rectangle, t z.Texture2D, flipX, flipY bool, c z.Color) {
	w := float32(t.Width)
	h := float32(t.Height)
	if flipX {
		w = -w
	}
	if flipY {
		h = -h
	}
	src := RH(0, 0, w, h)
	z.DrawTexturePro(t, src, r, z.Vector2Zero(), 0, c)
}
func mIMpath(path string) []z.Texture2D {
	var t []z.Texture2D
	files, _ := os.ReadDir(path)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext == ".png" {
			fullPath := filepath.Join(path, file.Name())
			texture := z.LoadTexture(fullPath)
			t = append(t, texture)
		}
	}
	return t
}
func remIMfromslice(i int, slice []IM) []IM {
	// Guard against out of bounds
	if i < 0 || i >= len(slice) {
		return slice
	}

	// 1. Move the last element into index i
	slice[i] = slice[len(slice)-1]

	// 2. Truncate the last element
	return slice[:len(slice)-1]
}

// MARK: COLLIS
func crp(r z.Rectangle, p z.Vector2) bool { return z.CheckCollisionPointRec(p, r) }
func cmsc(r z.Rectangle) bool             { return z.CheckCollisionPointRec(ms, r) }
func cms(r z.Rectangle) bool              { return z.CheckCollisionPointRec(msnocam, r) }
func crr(r1, r2 z.Rectangle) bool         { return z.CheckCollisionRecs(r1, r2) }
func crblokrs(r z.Rectangle, b []BLOK) bool {
	collides := false
	for i := range b {
		if crr(r, b[i].r) {
			collides = true
			break
		}
	}
	return collides
}

// MARK: UI
func close(val bool, x, y float32) bool {
	r := R(x-(unh+spc), y+spc, unh)
	//dtxtrcnt(tablerui, icoui[0].ic, r, z.White)
	if cms(r) {
		if mslc {
			val = !val
		}
	}
	return val
}

// MARK: TEXT
func dquicknuminven(r z.Rectangle, num int) {
	x := r.X + r.Width + spc/2
	y := (r.Y + r.Height/2) - float32(fnsml.BaseSize)
	dtxtc(fnsml, fmt.Sprint(num+1), x, y, c.VoidBlack)
}
func dtxtc(f z.Font, t string, x, y float32, c z.Color) {
	z.DrawTextEx(f, t, V2(x, y), float32(f.BaseSize), 0, c)
}
func dinfotxt() {
	x := CNT.X
	y := float32(fn.BaseSize) / 2
	dtxtcntrv2BGrecroundLRG(fn, infotxt, V2(x, y), 12, c.VoidBlack, c.IvoryWhite, c.VoidBlack, true)
}
func dwarntxt() {
	x := CNT.X - txlen(fn, warntxt)/2
	y := un
	z.DrawTextEx(fn, warntxt, V2(x, y), float32(fn.BaseSize), 0, c.Crimson)
}
func minfotxt(t string) {
	infotxtT = FPS
	infotxt = t
}
func mwarntxt(t string) {
	warntxtT = FPS
	warntxt = t
}
func dtxtdef(t string, x, y float32) {
	z.DrawTextEx(fnsml, t, V2(x, y), float32(fnsml.BaseSize), 0, z.White)
}
func dtxtdefshadow(t string, x, y float32) {
	z.DrawTextEx(fnsml, t, V2(x-1, y+1), float32(fnsml.BaseSize), 0, z.Black)
	z.DrawTextEx(fnsml, t, V2(x, y), float32(fnsml.BaseSize), 0, z.White)
}
func dtxtrcnt(f z.Font, t string, r z.Rectangle, c z.Color) {
	cnt := rcnt(r)
	w := txlen(f, t)
	z.DrawTextEx(f, t, V2(cnt.X-w/2, cnt.Y-float32(f.BaseSize)/2), float32(f.BaseSize), 0, c)
}
func dtxtcntrv2(f z.Font, t string, v2 z.Vector2, c z.Color) {
	w := txlen(f, t)
	x := v2.X - w/2
	y := v2.Y - float32(f.BaseSize)/2
	z.DrawTextEx(f, t, V2(x, y), float32(f.BaseSize), 0, c)
}
func dtxtcntrv2BGrecround(f z.Font, t string, v2 z.Vector2, pad float32, ctxt, cbg, coutline z.Color, shadow bool) {
	w := txlen(f, t)
	wr := w + pad*2
	hr := float32(f.BaseSize)
	x := v2.X - w/2
	y := (v2.Y - float32(f.BaseSize)/2) - 2
	r := rechfromcnt(v2, wr, hr)
	if shadow {
		r2 := r
		r2.X -= un16th
		r2.Y += un16th
		drroundsml(r2, ca(c.VoidBlack, 100))
	}
	drroundsml(r, cbg)
	drlroundsml(r, 1, coutline)
	z.DrawTextEx(f, t, V2(x, y), float32(f.BaseSize), 0, ctxt)
}
func dtxtcntrv2BGrecroundLRG(f z.Font, t string, v2 z.Vector2, pad float32, ctxt, cbg, coutline z.Color, shadow bool) {
	w := txlen(f, t)
	wr := w + pad*2
	hr := float32(f.BaseSize) - float32(f.BaseSize)/4
	x := v2.X - w/2
	y := (v2.Y - float32(f.BaseSize)/2) - 6
	r := rechfromcnt(v2, wr, hr)
	if shadow {
		r2 := r
		r2.X -= un16th
		r2.Y += un16th
		drroundsml(r2, ca(c.VoidBlack, 100))
	}
	drroundsml(r, cbg)
	drlroundsml(r, 1, coutline)
	z.DrawTextEx(f, t, V2(x, y), float32(f.BaseSize), 0, ctxt)
}
func dtxtrcntshadow(f z.Font, t string, r z.Rectangle, c, cShadow z.Color, offset float32) {
	cnt := rcnt(r)
	z.DrawTextEx(f, t, V2((cnt.X-float32(f.BaseSize)/2)-offset, (cnt.Y-float32(f.BaseSize)/2)+offset), float32(f.BaseSize), 0, cShadow)
	z.DrawTextEx(f, t, V2(cnt.X-float32(f.BaseSize)/2, cnt.Y-float32(f.BaseSize)/2), float32(f.BaseSize), 0, c)
}
func txlen(f z.Font, t string) float32 {
	v2 := z.MeasureTextEx(f, t, float32(f.BaseSize), 0)
	return v2.X
}

// MARK: GEOM
func rcntrinr(rInner, rOuter z.Rectangle) z.Rectangle {
	cnt := rcnt(rOuter)
	return RH(cnt.X-rInner.Width/2, cnt.Y-rInner.Height/2, rInner.Width, rInner.Height)
}
func drhxy(x, y, w, h int32, c z.Color) {
	z.DrawRectangle(x, y, w, h, c)
}
func drl(r z.Rectangle, lineW float32, c z.Color) { z.DrawRectangleLinesEx(r, lineW, c) }
func drroundsml(r z.Rectangle, c z.Color)         { z.DrawRectangleRounded(r, 0.4, 12, c) }
func drlroundsml(r z.Rectangle, lineW float32, c z.Color) {
	z.DrawRectangleRoundedLinesEx(r, 0.4, 12, lineW, c)
}
func drround(r z.Rectangle, c z.Color) { z.DrawRectangleRounded(r, roundrec, 10, c) }
func drlround(r z.Rectangle, lineW float32, c z.Color) {
	z.DrawRectangleRoundedLinesEx(r, roundrec, 10, lineW, c)
}
func R(x, y, w float32) z.Rectangle     { return z.NewRectangle(x, y, w, w) }
func RH(x, y, w, h float32) z.Rectangle { return z.NewRectangle(x, y, w, h) }
func recfromcnt(cnt z.Vector2, w float32) z.Rectangle {
	return z.NewRectangle(cnt.X-w/2, cnt.Y-w/2, w, w)
}
func rechfromcnt(cnt z.Vector2, w, h float32) z.Rectangle {
	return z.NewRectangle(cnt.X-w/2, cnt.Y-h/2, w, h)
}
func rlrgr(r z.Rectangle, offset float32) z.Rectangle {
	return RH(r.X-offset, r.Y-offset, r.Width+offset*2, r.Height+offset*2)
}
func rsmlr(r z.Rectangle, inset float32) z.Rectangle {
	return RH(r.X+inset, r.Y+inset, r.Width-inset*2, r.Height-inset*2)
}
func drrcnt(cnt z.Vector2, r z.Rectangle, c z.Color) {
	r.X = cnt.X - r.Width/2
	r.Y = cnt.Y - r.Height/2
	dr(r, c)
}
func dr(r z.Rectangle, c z.Color) { z.DrawRectangleRec(r, c) }
func anglp(p1, p2 z.Vector2) float32 {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y

	// math.Atan2 returns radians in range [-Pi, Pi]
	radians := math.Atan2(float64(dy), float64(dx))

	// Convert to degrees for Raylib
	degrees := float32(radians * (180.0 / math.Pi))

	// Optional: normalize to 0° - 360° range
	if degrees < 0 {
		degrees += 360.0
	}

	return degrees
}

// MARK: COLOR
func ca(c z.Color, a uint8) z.Color {
	c.A = a
	return c
}
func RANREDBRICK() z.Color {
	// Hue: 4° to 16° (warm earthy red with a subtle rusty-orange clay undertone)
	hue := float32(z.GetRandomValue(4, 16))

	// Saturation: 58% to 76% (earthy and rich; avoids looking like neon plastic or bright cherry)
	sat := float32(z.GetRandomValue(58, 76)) / 100.0

	// Value: 46% to 66% (medium-dark fired clay tone without being pitch black)
	val := float32(z.GetRandomValue(46, 66)) / 100.0

	return z.ColorFromHSV(hue, sat, val)
}

func RANGREYSTONELIGHT() z.Color {
	// Hue: 195° to 225° (subtle cool slate stone undertone)
	hue := float32(z.GetRandomValue(195, 225))

	// Saturation: 3% to 8% (very low saturation gives an authentic natural mineral/rock look)
	sat := float32(z.GetRandomValue(3, 8)) / 100.0

	// Value: 72% to 88% (light grey without being pure white)
	val := float32(z.GetRandomValue(72, 88)) / 100.0

	return z.ColorFromHSV(hue, sat, val)
}
func YELLOWSUNBRIGHT() z.Color {
	return z.NewColor(255, 235, 0, 255)
}
func YELLOWBRIGHT() z.Color {
	return z.NewColor(255, 255, 0, 255)
}
func RANBROWNBRIGHT() z.Color {
	// Hue: 24° to 34° (warm amber/orange-brown spectrum)
	hue := float32(z.GetRandomValue(24, 34))

	// Saturation: 65% to 88% (rich and vivid, avoids gray/muddy tones)
	sat := float32(z.GetRandomValue(65, 88)) / 100.0

	// Value: 56% to 75% (bright enough to clearly read as brown without becoming dark chocolate/black)
	val := float32(z.GetRandomValue(56, 75)) / 100.0

	return z.ColorFromHSV(hue, sat, val)
}
func RANBEACHSAND() z.Color {
	// Hue: 42° to 52° (warm sandy cream to golden pale yellow)
	hue := float32(z.GetRandomValue(42, 52))

	// Saturation: 25% to 48% (natural sandy tone, avoids neon yellow)
	sat := float32(z.GetRandomValue(25, 48)) / 100.0

	// Value: 86% to 98% (bright, sunny beach sand)
	val := float32(z.GetRandomValue(86, 98)) / 100.0

	return z.ColorFromHSV(hue, sat, val)
}
func RANGREEN() z.Color {
	// Hue: 85° (warm yellowish-lime) to 145° (cool deep forest green)
	hue := float32(z.GetRandomValue(85, 145))

	// Saturation: 60% to 90% (rich and colorful, not washed out)
	sat := float32(z.GetRandomValue(60, 90)) / 100.0

	// Brightness (Value): 45% to 85% (prevents pitch-black or neon white)
	val := float32(z.GetRandomValue(45, 85)) / 100.0

	return z.ColorFromHSV(hue, sat, val)
}
func RANORANGEBRIGHT() z.Color {
	// Hue: 24° (red-orange) to 36° (vivid golden orange)
	hue := float32(z.GetRandomValue(24, 36))

	// Saturation: 85% to 100% (punchy, high saturation)
	sat := float32(z.GetRandomValue(85, 100)) / 100.0

	// Value: 90% to 100% (high brightness)
	val := float32(z.GetRandomValue(90, 100)) / 100.0

	return z.ColorFromHSV(hue, sat, val)
}
func RANREDBRIGHT() z.Color {
	// In HSV, Pure Red is 0° (and 360°).
	// We pick between 354° (vibrant ruby/crimson red) and 10° (punchy scarlet red).
	hue := float32(z.GetRandomValue(354, 370))
	if hue >= 360.0 {
		hue -= 360.0 // Wraps around 360° to 0°-10°
	}

	// Saturation: 85% to 100% (intense, vivid red, not washed out)
	sat := float32(z.GetRandomValue(85, 100)) / 100.0

	// Brightness (Value): 85% to 100% (punchy and bright, prevents muddy dark maroons)
	val := float32(z.GetRandomValue(85, 100)) / 100.0

	return z.ColorFromHSV(hue, sat, val)
}
func RANSEPIA() z.Color {
	// Hue: 25° (warm reddish-brown) to 42° (golden/amber sepia tone)
	hue := float32(z.GetRandomValue(25, 42))

	// Saturation: 20% to 45% (subtle and muted; keeps it vintage rather than vibrant orange)
	sat := float32(z.GetRandomValue(20, 45)) / 100.0

	// Brightness (Value): 50% to 85% (ranges from aged photo shadows to warm highlights)
	val := float32(z.GetRandomValue(50, 85)) / 100.0

	return z.ColorFromHSV(hue, sat, val)
}

// MARK:SLICES
func addremSlice[T comparable](slice []T, val T) []T {
	if idx := slices.Index(slice, val); idx != -1 {
		return slices.Delete(slice, idx, idx+1) // Remove
	}
	return append(slice, val) // Add
}

// MARK:MATH
func rrScale(src, dest z.Rectangle) float32 {
	srcW := float32(math.Abs(float64(src.Width)))
	if srcW == 0 {
		return 1.0 // Prevent division by zero
	}
	return dest.Width / srcW
}
func rcnt(r z.Rectangle) z.Vector2 { return V2(r.X+r.Width/2, r.Y+r.Height/2) }
func V2(x, y float32) z.Vector2    { return z.NewVector2(x, y) }
func f32[T numConvert](val T) float32 {
	return float32(val)
}

type numConvert interface { //CONVERT ANY NUM TO F32
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// MARK:SCALE
func fitScale(width, height, maxSize float32) float32 {
	// Guard against zero, negative values, or invalid inputs
	if width <= 0 || height <= 0 || maxSize <= 0 {
		return 0
	}

	// Determine the dominant dimension
	maxDim := width
	if height > maxDim {
		maxDim = height
	}

	// Calculate uniform scale
	return maxSize / maxDim
}
func scale2wh(origW, origH, maxW, maxH float32) (float32, float32) {
	if origW == 0 || origH == 0 {
		return 0, 0
	}
	scaleW := maxW / origW
	scaleH := maxH / origH

	// Use smaller scale factor so it fits entirely inside
	scale := scaleW
	if scaleH < scaleW {
		scale = scaleH
	}

	return origW * scale, origH * scale
}
func scaleH(origW, origH, newW float32) float32 {
	if origW == 0 {
		return 0
	}
	return (newW / origW) * origH
}
func scaleW(origW, origH, newH float32) float32 {
	if origH == 0 {
		return 0
	}
	return (newH / origH) * origW
}

// MARK: RANDOM
func RINT(min, max int) int {
	if min >= max {
		return min
	}
	return min + rand.IntN(max-min)
}

func RI32(min, max int) int32 {
	if min >= max {
		return int32(min)
	}
	return int32(min + rand.IntN(max-min))
}

func RINT8(min, max uint8) uint8 {
	if min >= max {
		return min
	}
	return min + uint8(rand.IntN(int(max)-int(min)))
}

func RF32(min, max float32) float32 {
	if min >= max {
		return min
	}
	return min + rand.Float32()*(max-min)
}

func ROLL(sides int) int {
	if sides < 1 {
		return 1
	}
	return 1 + rand.IntN(sides) // 1 to sides inclusive
}

func FLIP() bool {
	return rand.IntN(2) == 0
}

//MARK: KEYS INPUT

var (
	activeInputPtr uintptr
	inputBuffer    string
)

func isAllowedChar(r rune, lettersonly bool) bool {
	if lettersonly {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == ' '
	}
	return r >= 32 && r <= 126
}

func KeysInpTxt(
	font z.Font,
	val *string,
	x, y, w, h float32,
	fontSize float32,
	lettersonly bool,
	maxLen int,
) string {
	if val == nil {
		return ""
	}
	if maxLen <= 0 {
		maxLen = 24
	}

	r := z.NewRectangle(x, y, w, h)
	id := uintptr(unsafe.Pointer(val))
	isFocused := activeInputPtr == id

	mousePos := z.GetMousePosition()
	isHovered := z.CheckCollisionPointRec(mousePos, r)

	if z.IsMouseButtonPressed(z.MouseButtonLeft) {
		if isHovered {
			activeInputPtr = id
			inputBuffer = *val
			isFocused = true
		} else if isFocused {
			// Clicked outside: Save & unfocus
			*val = inputBuffer
			activeInputPtr = 0
			isFocused = false
		}
	}

	if isFocused {

		for {
			char := z.GetCharPressed()
			if char == 0 {
				break
			}
			runes := []rune(inputBuffer)
			if len(runes) < maxLen && isAllowedChar(char, lettersonly) {
				inputBuffer += string(char)
			}
		}

		if z.IsKeyPressed(z.KeyBackspace) || z.IsKeyPressedRepeat(z.KeyBackspace) {
			runes := []rune(inputBuffer)
			if len(runes) > 0 {
				inputBuffer = string(runes[:len(runes)-1])
			}
		}

		if z.IsKeyPressed(z.KeyEnter) || z.IsKeyPressed(z.KeyKpEnter) {
			*val = inputBuffer
			activeInputPtr = 0
			isFocused = false
			editName = false
		}

		if z.IsKeyPressed(z.KeyEscape) {
			inputBuffer = *val // Revert back
			activeInputPtr = 0
			isFocused = false
			editName = false
		}
	}

	// --- 3. DRAWING ---
	// Background
	bgColor := z.NewColor(32, 32, 48, 255) // Pal.Void
	z.DrawRectangleRec(r, bgColor)

	// Border
	if isFocused {
		z.DrawRectangleLinesEx(r, 2.0, z.Gold) // Active highlight
	} else if isHovered {
		z.DrawRectangleLinesEx(r, 1.5, z.SkyBlue) // Hover highlight
	} else {
		z.DrawRectangleLinesEx(r, 1.0, z.DarkGray) // Normal border
	}

	// Display Text + Blinking Cursor
	displayText := *val
	if isFocused {
		displayText = inputBuffer
		// Blinks twice per second (500ms cycle)
		if int(z.GetTime()*2)%2 == 0 {
			displayText += "|"
		}
	}

	// Vertically center text in the box with 8px left padding
	//textSize := z.MeasureTextEx(font, displayText, fontSize, 0)
	textPos := z.NewVector2(
		r.X+8,
		r.Y+8,
	)

	textColor := z.White
	if !isFocused && *val == "" {
		// Placeholder styling
		displayText = "type here..."
		textColor = z.Gray
	}

	z.DrawTextEx(font, displayText, textPos, fontSize, 0, textColor)

	return *val
}

// MARK: CAMERA
func GetCameraWorldRect(camera z.Camera2D, screenWidth, screenHeight float32) z.Rectangle {
	// Screen top-left (0, 0) in world coordinates
	topLeft := z.GetScreenToWorld2D(z.NewVector2(0, 0), camera)

	// Screen bottom-right in world coordinates
	bottomRight := z.GetScreenToWorld2D(z.NewVector2(screenWidth, screenHeight), camera)

	return z.NewRectangle(
		topLeft.X,
		topLeft.Y,
		bottomRight.X-topLeft.X,
		bottomRight.Y-topLeft.Y,
	)
}

//MARK: PALETTE
//CREDIT: https://lospec.com/palette-list/pineapple-32

var c = struct {
	// --- 1. WARM REDS & OCHRE (0-7) ---
	DeepWine    z.Color // #43002a - Deep wine burgundy
	Crimson     z.Color // #890027 - Rich velvet crimson
	ScarletRed  z.Color // #d9243c - Vibrant scarlet red
	CoralBlush  z.Color // #ff6157 - Warm coral pink-red
	AmberGold   z.Color // #ffb762 - Warm amber gold
	RustCopper  z.Color // #c76e46 - Ochre rust copper
	RussetBrown z.Color // #73392e - Russet chestnut brown
	NightPlum   z.Color // #34111f - Dark plum shadow

	// --- 2. FOLIAGE, FIRE & SOLAR (8-15) ---
	VoidBlack   z.Color // #030710 - Obsidian void / near black
	DarkPine    z.Color // #273b2d - Dark pine green
	ForestGreen z.Color // #458239 - Foliage forest green
	MossOlive   z.Color // #9cb93b - Vibrant moss / olive lime
	SunYellow   z.Color // #ffd832 - Radiant sunny gold
	BlazeOrange z.Color // #ff823b - Vivid blaze orange
	BurntOrange z.Color // #d1401f - Deep burnt orange
	BloodRed    z.Color // #7c191a - Deep blood maroon

	// --- 3. EARTH, PEACH & NEUTRALS (16-23) ---
	DarkEspresso  z.Color // #310c1b - Deep espresso mahogany
	ClaySienna    z.Color // #833f34 - Earthy clay sienna
	WarmPeach     z.Color // #eb9c6e - Soft terracotta peach
	PeachCream    z.Color // #ffdaac - Warm cream highlight
	IvoryWhite    z.Color // #ffffe4 - Pale ivory buttercream
	SilverMist    z.Color // #bfc3c6 - Cool silver slate
	SlateSage     z.Color // #6d8a8d - Muted slate sage
	MidnightSlate z.Color // #293b49 - Deep midnight slate

	// --- 4. OCEAN, MINT & PURPLES (24-31) ---
	AbyssNavy     z.Color // #041528 - Deep abyss navy
	OceanTeal     z.Color // #033e5e - Rich deep ocean teal
	CeruleanTeal  z.Color // #1c92a7 - Bright cerulean aqua
	SeafoamMint   z.Color // #77d6c1 - Fresh mint turquoise
	PaleBlush     z.Color // #ffe0dc - Pale blossom blush
	RosePink      z.Color // #ff88a9 - Vibrant candy rose
	OrchidMagenta z.Color // #c03b94 - Rich orchid magenta
	RoyalViolet   z.Color // #601761 - Dark royal violet
}{
	// 1. Warm Reds & Ochre
	DeepWine:    z.NewColor(67, 0, 42, 255),
	Crimson:     z.NewColor(137, 0, 39, 255),
	ScarletRed:  z.NewColor(217, 36, 60, 255),
	CoralBlush:  z.NewColor(255, 97, 87, 255),
	AmberGold:   z.NewColor(255, 183, 98, 255),
	RustCopper:  z.NewColor(199, 110, 70, 255),
	RussetBrown: z.NewColor(115, 57, 46, 255),
	NightPlum:   z.NewColor(52, 17, 31, 255),

	// 2. Foliage, Fire & Solar
	VoidBlack:   z.NewColor(3, 7, 16, 255),
	DarkPine:    z.NewColor(39, 59, 45, 255),
	ForestGreen: z.NewColor(69, 130, 57, 255),
	MossOlive:   z.NewColor(156, 185, 59, 255),
	SunYellow:   z.NewColor(255, 216, 50, 255),
	BlazeOrange: z.NewColor(255, 130, 59, 255),
	BurntOrange: z.NewColor(209, 64, 31, 255),
	BloodRed:    z.NewColor(124, 25, 26, 255),

	// 3. Earth, Peach & Neutrals
	DarkEspresso:  z.NewColor(49, 12, 27, 255),
	ClaySienna:    z.NewColor(131, 63, 52, 255),
	WarmPeach:     z.NewColor(235, 156, 110, 255),
	PeachCream:    z.NewColor(255, 218, 172, 255),
	IvoryWhite:    z.NewColor(255, 255, 228, 255),
	SilverMist:    z.NewColor(191, 195, 198, 255),
	SlateSage:     z.NewColor(109, 138, 141, 255),
	MidnightSlate: z.NewColor(41, 59, 73, 255),

	// 4. Ocean, Mint & Purples
	AbyssNavy:     z.NewColor(4, 21, 40, 255),
	OceanTeal:     z.NewColor(3, 62, 94, 255),
	CeruleanTeal:  z.NewColor(28, 146, 167, 255),
	SeafoamMint:   z.NewColor(119, 214, 193, 255),
	PaleBlush:     z.NewColor(255, 224, 220, 255),
	RosePink:      z.NewColor(255, 136, 169, 255),
	OrchidMagenta: z.NewColor(192, 59, 148, 255),
	RoyalViolet:   z.NewColor(96, 23, 97, 255),
}

// PalList contains all 32 colors in index order (0 to 31)
var COLORS = []z.Color{
	c.DeepWine, c.Crimson, c.ScarletRed, c.CoralBlush,
	c.AmberGold, c.RustCopper, c.RussetBrown, c.NightPlum,
	c.VoidBlack, c.DarkPine, c.ForestGreen, c.MossOlive,
	c.SunYellow, c.BlazeOrange, c.BurntOrange, c.BloodRed,
	c.DarkEspresso, c.ClaySienna, c.WarmPeach, c.PeachCream,
	c.IvoryWhite, c.SilverMist, c.SlateSage, c.MidnightSlate,
	c.AbyssNavy, c.OceanTeal, c.CeruleanTeal, c.SeafoamMint,
	c.PaleBlush, c.RosePink, c.OrchidMagenta, c.RoyalViolet,
}

var COLORNAMES = []string{
	"DeepWine", "Crimson", "ScarletRed", "CoralBlush",
	"AmberGold", "RustCopper", "RussetBrown", "NightPlum",
	"VoidBlack", "DarkPine", "ForestGreen", "MossOlive",
	"SunYellow", "BlazeOrange", "BurntOrange", "BloodRed",
	"DarkEspresso", "ClaySienna", "WarmPeach", "PeachCream",
	"IvoryWhite", "SilverMist", "SlateSage", "MidnightSlate",
	"AbyssNavy", "OceanTeal", "CeruleanTeal", "SeafoamMint",
	"PaleBlush", "RosePink", "OrchidMagenta", "RoyalViolet",
}

// RandomPaletteColor returns any random color from the 64-color palette
func RANCOLOR() z.Color {
	return COLORS[RINT(0, len(COLORS))]
}
