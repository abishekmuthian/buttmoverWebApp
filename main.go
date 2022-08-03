package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	resources "github.com/hajimehoshi/ebiten/v2/examples/resources/images/flappy"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"strings"
	"syscall/js"
	"time"
)

func floorDiv(x, y int) int {
	d := x / y
	if d*y == x || x >= 0 {
		return d
	}
	return d - 1
}

func floorMod(x, y int) int {
	return x - floorDiv(x, y)*y
}

const (
	screenWidth   = 640
	screenHeight  = 480
	tileSize      = 32
	bulldozerSize = 72
	chairSize     = 48
	buttSize      = 48
	fontSize      = 24
	titleFontSize = fontSize
	smallFontSize = fontSize / 2
	clockFontSize = fontSize * 4

	frameOX = 0
	frameOY = 0

	bulldozerFrameWidth  = 72
	bulldozerFrameHeight = 72
	bulldozerFrameNum    = 5

	chairFrameWidth  = 48
	chairFrameHeight = 48

	buttFrameWidth  = 48
	buttFrameHeight = 48
	buttFrameNum    = 5
)

var (
	bgColor                           = color.Black
	fontColor                         = color.White
	tilesImage                        *ebiten.Image
	titleArcadeFont                   font.Face
	arcadeFont                        font.Face
	smallArcadeFont                   font.Face
	clockArcadeFont                   font.Face
	munroSmallFont                    font.Face
	ticker25                          *time.Ticker
	doneTicker25                      chan bool
	ticker30                          *time.Ticker
	doneTicker30                      chan bool
	ticker5                           *time.Ticker
	doneTicker5                       chan bool
	connectInfo                       = []string{"Click &", "", "Press S to Start/Stop", "the timer manually.", "", "Or", "", "Press C to connect", "to the Butt Trigger."}
	connectedText                     = "Connected"
	disconnectedText                  = "Disconnected"
	gameError                         = false
	gameErrorText                     string
	portErrorInfo                     = []string{"Butt Trigger not selected", "Check the guide below", "", "Press R to Restart"}
	buttTriggerErrorInfo              = []string{"Butt Trigger has been lost", "Check the guide below", "", "Press R to Restart"}
	notStartedText                    = []string{"Stand up and sit down", "to begin"}
	gameOverInfoConnected             = []string{"GAME OVER!", "", "Stand up and sit down", "to begin"}
	gameOVerInfoNotConnected          = []string{"GAME OVER!", "", "Press S to Start/Stop", "", "the task timer manually."}
	colorMode                         string
	timeText                          string
	connected                         = false
	moverYourButt                     = false
	moveYourButtTextConnected         = []string{"Move Your Butt"}
	moveYourButtTextNotConnected      = []string{"Move Your Butt", "", "Press S to Start/Stop", "", "break timer"}
	bringYourButtBack                 = false
	bringYourButtBackTextConnected    = []string{"Bring Your Butt Back"}
	bringYourButtBackTextNotConnected = []string{"Bring Your Butt Back", "", "Press S to Start/Stop", "", "task timer"}
	clockStarted                      = false
	productivityScore                 int
	healthScore                       int
	lifeScore                         = 100
	appName                           = "Butt Mover"
	//go:embed "bulldozer_sprite.png"
	bulldozerImageFile []byte

	//go:embed "chair.png"
	chairImageFile []byte

	//go:embed "butt_sprite.png"
	buttImageFile []byte

	bulldozerImage   *ebiten.Image
	bulldozerStep    int
	bulldozerStarted bool
	chairImage       *ebiten.Image
	buttImage        *ebiten.Image
	buttonTrigger    bool
	newGame          bool
	timerLabel       string
)

/*func init() {

	fontBase64 := "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"

	fontBytes, err := base64.StdEncoding.DecodeString(fontBase64)

	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	munroSmallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    18,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}*/

func init() {
	/*	data, err := file.ReadFile("bulldozer_sprite.png")
		if err != nil {
			log.Fatal(err)
		}*/
	img, _, err := image.Decode(bytes.NewReader(bulldozerImageFile))
	if err != nil {
		log.Fatal(err)
	}
	bulldozerImage = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(chairImageFile))
	if err != nil {
		log.Fatal(err)
	}
	chairImage = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(buttImageFile))
	if err != nil {
		log.Fatal(err)
	}
	buttImage = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(resources.Tiles_png))
	if err != nil {
		log.Fatal("Error while loading tiles image, ", err)
	}
	tilesImage = ebiten.NewImageFromImage(img)
}

func init() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	titleArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    titleFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	smallArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    smallFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	clockArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    clockFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Mode int

const (
	ModeTitle Mode = iota
	ModeConnect
	ModeGame
	ModeGameAnnouncement
	ModeGameOver
	ModeGameError
)

type Game struct {
	mode Mode
	keys []ebiten.Key

	// Camera
	cameraX int
	cameraY int

	audioContext *audio.Context
	jumpPlayer   *audio.Player
	hitPlayer    *audio.Player

	// Bulldozer Animation
	count int

	// Bulldozer position
	bX float64
	bY float64

	// Chair position
	cX  float64
	cY  float64
	cVY float64

	// Butt position
	buX  float64
	buY  float64
	buVY float64
}

func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.cameraX = -240
	g.cameraY = 0

	if g.audioContext == nil {
		g.audioContext = audio.NewContext(48000)
	}

	jumpD, err := vorbis.Decode(g.audioContext, bytes.NewReader(raudio.Jump_ogg))
	if err != nil {
		log.Fatal(err)
	}
	g.jumpPlayer, err = g.audioContext.NewPlayer(jumpD)
	if err != nil {
		log.Fatal(err)
	}

	jabD, err := wav.Decode(g.audioContext, bytes.NewReader(raudio.Jab_wav))
	if err != nil {
		log.Fatal(err)
	}
	g.hitPlayer, err = g.audioContext.NewPlayer(jabD)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	g.count++
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	switch g.mode {
	case ModeTitle:
		//fmt.Println("Inside ModeTitle")
		if inpututil.IsKeyJustPressed(ebiten.KeyC) || inpututil.IsKeyJustReleased(ebiten.KeyC) {
			g.mode = ModeConnect
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			buttonTrigger = true
			js.Global().Call("buttonPressed")
			g.mode = ModeGame
		}

		//Chair position
		g.cX = 430

		// Butt position
		g.buX = 425
		g.buY = 1.23
	case ModeConnect:
		//fmt.Println("Inside ModeConnect")
		js.Global().Call("callConnectButtTrigger")
		if gameError {
			//fmt.Println("Inside connection Error in ModeConnect")
			g.mode = ModeGameError
		} else if connected {
			//fmt.Println("Inside Game loop in ModeConnect")
			g.mode = ModeGame
		}
	case ModeGame:
		//fmt.Println("Inside ModeGame")
		if connected && gameError {
			fmt.Println("Inside connection Error in ModeGame")
			g.mode = ModeGameError
		} else {
			if inpututil.IsKeyJustPressed(ebiten.KeyS) && !buttonTrigger && !connected {
				js.Global().Call("buttonPressed")
				buttonTrigger = true
			} else if inpututil.IsKeyJustPressed(ebiten.KeyS) && buttonTrigger && !connected {
				buttonTrigger = false
				js.Global().Call("buttonReleased")
			}

			if moverYourButt || bringYourButtBack {
				g.hitPlayer.Rewind()
				g.hitPlayer.Play()
				g.mode = ModeGameAnnouncement
			}

			if clockStarted {
				err := g.jumpPlayer.Rewind()
				if err != nil {
					return err
				}
				g.jumpPlayer.Play()
				clockStarted = false
			}

			if buttonTrigger {
				bulldozerStarted = true
				g.bX = float64(bulldozerStep) / 4.16

				// Chair position
				g.cX = 430
				g.cVY = 0

				// Butt position
				g.buX = 425
				g.buY = 1.23
			}

			if !newGame && lifeScore < 1 {
				js.Global().Call("clearBreak")
				js.Global().Call("clearTask")
				g.hitPlayer.Rewind()
				g.hitPlayer.Play()
				g.mode = ModeGameOver

				if buttonTrigger && connected {
					buttonTrigger = false
				}
			}
		}
	case ModeGameAnnouncement:
		//fmt.Println("Inside ModeGameAnnouncement")
		if connected && gameError {
			//fmt.Println("Inside connection Error in ModeAnnouncement")
			g.mode = ModeGameError
		} else {
			if !moverYourButt && !bringYourButtBack {
				if timerLabel == "Task" {
					g.bX = 0
					bulldozerStep = 0
					g.cX = 420 // Handle chair position
					g.cVY = 1.5
					g.cX = 460

					// Handle butt position
					//g.buVY = 3.15
					g.buX = 515
					g.buY = 1.18
				}
				g.mode = ModeGame
			} else {
				// Handle chair position
				g.cVY = 1.5
				g.cX = 460

				// Handle butt position
				//g.buVY = 3.15
				g.buX = 515
				g.buY = 1.18
			}
			bulldozerStarted = false
			if inpututil.IsKeyJustPressed(ebiten.KeyS) && !buttonTrigger && !connected {
				buttonTrigger = true
				js.Global().Call("buttonPressed")
			} else if inpututil.IsKeyJustPressed(ebiten.KeyS) && buttonTrigger && !connected {
				buttonTrigger = false
				js.Global().Call("buttonReleased")
			}
		}
	case ModeGameOver:
		//fmt.Println("Inside ModeGameOver")
		if connected && gameError {
			//fmt.Println("Inside connection Error in ModeGameOver")
			g.mode = ModeGameError
		} else {
			timerLabel = ""
			bulldozerStarted = false

			if connected {
				js.Global().Call("setGameOverFlag", "true")
			}

			// Handle chair position
			g.cVY = 1.5
			g.cX = 460

			// Handle butt position
			//g.buVY = 3.15
			g.buX = 515
			g.buY = 1.18

			if inpututil.IsKeyJustPressed(ebiten.KeyS) && !connected {
				js.Global().Call("reset")
				buttonTrigger = true
				js.Global().Call("buttonPressed")
				g.mode = ModeGame
			}

			if buttonTrigger && connected {
				js.Global().Call("setGameOverFlag", "false")
				bulldozerStep = 0
				g.mode = ModeGame
			}
		}
	case ModeGameError:
		//fmt.Println("Inside ModeGameError")
		if inpututil.IsKeyJustPressed(ebiten.KeyR) && gameError {
			g.mode = ModeTitle
		} else if !gameError {
			g.mode = ModeGame
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(bgColor)
	g.drawTiles(screen)
	g.drawBulldozer(screen)
	g.drawButt(screen)
	g.drawChair(screen)

	keyStrs := []string{}
	for _, p := range g.keys {
		keyStrs = append(keyStrs, p.String())
	}
	ebitenutil.DebugPrint(screen, strings.Join(keyStrs, ", "))

	if connected {
		text.Draw(screen, connectedText, smallArcadeFont, (screenWidth-len(connectedText)*smallFontSize)-24, smallFontSize+6, fontColor)
	} else {
		text.Draw(screen, disconnectedText, smallArcadeFont, (screenWidth-len(disconnectedText)*smallFontSize)-24, smallFontSize+6, fontColor)
	}

	switch g.mode {
	case ModeTitle:
		// Draw the connect text
		for i, l := range connectInfo {
			x := (screenWidth - len(l)*titleFontSize) / 2
			text.Draw(screen, l, titleArcadeFont, x, (i+5)*titleFontSize, fontColor)
		}
	case ModeGame:
		if timeText != "" {

			text.Draw(screen, timerLabel, smallArcadeFont, (screenWidth-len(timerLabel)*smallFontSize)-24, smallFontSize+32, fontColor)

			x := (screenWidth - len(timeText)*clockFontSize) / 2
			text.Draw(screen, timeText, clockArcadeFont, x, 3*clockFontSize, fontColor)
			g.drawBulldozer(screen)
		} else if connected {
			for i, l := range notStartedText {
				x := (screenWidth - len(l)*titleFontSize) / 2
				text.Draw(screen, l, titleArcadeFont, x, (i+8)*titleFontSize, fontColor)
			}
		}
	case ModeGameAnnouncement:
		x := (screenWidth - len(timeText)*clockFontSize) / 2
		text.Draw(screen, timeText, clockArcadeFont, x, 3*clockFontSize, fontColor)

		if moverYourButt {
			if connected {
				for i, l := range moveYourButtTextConnected {
					x := (screenWidth - len(l)*titleFontSize) / 2
					text.Draw(screen, l, titleArcadeFont, x, (i+5)*titleFontSize, fontColor)
				}
			} else {
				for i, l := range moveYourButtTextNotConnected {
					x := (screenWidth - len(l)*titleFontSize) / 2
					if i > 0 {
						text.Draw(screen, l, titleArcadeFont, x, (i+12)*titleFontSize, fontColor)
					} else {
						text.Draw(screen, l, titleArcadeFont, x, (i+5)*titleFontSize, fontColor)
					}
				}
			}

		} else if bringYourButtBack {
			if connected {
				for i, l := range bringYourButtBackTextConnected {
					x := (screenWidth - len(l)*titleFontSize) / 2
					text.Draw(screen, l, titleArcadeFont, x, (i+5)*titleFontSize, fontColor)
				}
			} else {
				for i, l := range bringYourButtBackTextNotConnected {
					x := (screenWidth - len(l)*titleFontSize) / 2
					if i > 0 {
						text.Draw(screen, l, titleArcadeFont, x, (i+12)*titleFontSize, fontColor)
					} else {
						text.Draw(screen, l, titleArcadeFont, x, (i+5)*titleFontSize, fontColor)
					}
				}
			}
		}
		g.drawBulldozer(screen)
	case ModeGameOver:
		if connected {
			for i, l := range gameOverInfoConnected {
				x := (screenWidth - len(l)*titleFontSize) / 2
				if i > 0 {
					text.Draw(screen, l, titleArcadeFont, x, (i+8)*titleFontSize, fontColor)
				} else {
					text.Draw(screen, l, titleArcadeFont, x, (i+5)*titleFontSize, fontColor)
				}
			}
		} else {
			for i, l := range gameOVerInfoNotConnected {
				x := (screenWidth - len(l)*titleFontSize) / 2
				if i > 0 {
					text.Draw(screen, l, titleArcadeFont, x, (i+8)*titleFontSize, fontColor)
				} else {
					text.Draw(screen, l, titleArcadeFont, x, (i+5)*titleFontSize, fontColor)
				}
			}
		}

	case ModeGameError:
		if strings.Contains(gameErrorText, "Error in selecting the port") ||
			strings.Contains(gameErrorText, "Browser not supported") {
			for i, l := range portErrorInfo {
				x := (screenWidth - len(l)*titleFontSize) / 2
				text.Draw(screen, l, titleArcadeFont, x, (i+8)*titleFontSize, fontColor)
			}
		}

		if strings.Contains(gameErrorText, "Butt Trigger has been lost") {
			for i, l := range buttTriggerErrorInfo {
				x := (screenWidth - len(l)*titleFontSize) / 2
				text.Draw(screen, l, titleArcadeFont, x, (i+8)*titleFontSize, fontColor)
			}
		}

	}

	productivityScoreStr := fmt.Sprintf("Productivity:%d", productivityScore)
	text.Draw(screen, productivityScoreStr, smallArcadeFont, (screenWidth-len(productivityScoreStr)*smallFontSize)/24, smallFontSize+6, fontColor)

	healthScoreStr := fmt.Sprintf("Health:%d", healthScore)
	text.Draw(screen, healthScoreStr, smallArcadeFont, (screenWidth-len(healthScoreStr)*smallFontSize)/(24+2), smallFontSize+32, fontColor)

	lifeScoreStr := fmt.Sprintf("Life:%d", lifeScore) + "%"
	text.Draw(screen, lifeScoreStr, smallArcadeFont, (screenWidth-len(lifeScoreStr)*smallFontSize)/(24+4), smallFontSize+58, fontColor)

	text.Draw(screen, appName, smallArcadeFont, (screenWidth-len(appName)*smallFontSize)/2, smallFontSize+6, fontColor)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func setErrorStatus(this js.Value, args []js.Value) interface{} {
	gameError = true
	gameErrorText = args[0].String()
	fmt.Println("Error status set", gameErrorText)
	return nil
}

func setConnectionStatus(this js.Value, args []js.Value) interface{} {
	connected = args[0].Bool()
	gameError = false
	return nil
}

func setColorMode(this js.Value, args []js.Value) interface{} {
	colorMode = args[0].String()
	if colorMode == "dark" {
		bgColor = color.Black
		fontColor = color.White
	} else if colorMode == "light" {
		bgColor = color.White
		fontColor = color.Black
	}
	return nil
}

func setClock(this js.Value, args []js.Value) interface{} {
	clockStarted = true
	return nil
}

func setScore(this js.Value, args []js.Value) interface{} {
	healthScore = args[0].Int()
	productivityScore = args[1].Int()
	lifeScore = args[2].Int()
	newGame = args[3].Bool()

	//fmt.Println("Health score is: ", healthScore)
	//fmt.Println("Productivity score is: ", productivityScore)
	//fmt.Println("Life  score is: ", lifeScore)
	//fmt.Println("NewGame is : ", newGame)

	return nil
}

func setButtonTrigger(this js.Value, args []js.Value) interface{} {
	buttonTrigger = args[0].Bool()

	if buttonTrigger {
		bulldozerStarted = false
	}

	return nil
}

func setTick(this js.Value, args []js.Value) interface{} {
	/*	moverYourButt = false
		bringYourButtBack = false
		ticker25 = time.NewTicker(time.Second)
		doneTicker25 = make(chan bool)*/
	var minute, second int

	if !args[0].IsUndefined() {
		minute = args[0].Int()
	}

	if !args[1].IsUndefined() {
		second = args[1].Int()
	}
	moverYourButt = args[2].Bool()
	bringYourButtBack = args[3].Bool()
	timerLabel = args[4].String()

	if !args[0].IsUndefined() && !args[1].IsUndefined() {
		timeText = fmt.Sprintf("%02d:%02d", minute, second)
	}

	if !moverYourButt && !bringYourButtBack {
		bulldozerStep++
	}

	//fmt.Println("Move Your Butt", moverYourButt)
	//fmt.Println("Bring Your Butt Back", bringYourButtBack)

	return timeText
}

/*func tick25() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		moverYourButt = false
		bringYourButtBack = false
		ticker25 = time.NewTicker(time.Second)
		doneTicker25 = make(chan bool)
		var timer25Min, timer25Sec int

		go func() {
			for {
				select {
				case <-doneTicker25:
					return
				case t := <-ticker25.C:
					timer25Sec++
					//fmt.Println("\n Tick at", t)
					// Defining duration
					// of Seconds method
					if timer25Sec > 60 {
						timer25Sec = 0
						timer25Min++
					}

					timeText = //fmt.Sprintf("%02d:%02d", timer25Min, timer25Sec)
				}
			}
		}()

		// Handler for the Promise: this is a JS function
		// It receives two arguments, which are JS functions themselves: resolve and reject
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			resolve.Invoke("Tick25 started")
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}*/

/*func tick30() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		moverYourButt = false
		bringYourButtBack = false
		ticker30 = time.NewTicker(time.Second)
		doneTicker30 = make(chan bool)
		var timer30Min, timer30Sec int

		go func() {
			for {
				select {
				case <-doneTicker25:
					return
				case t := <-ticker25.C:
					timer30Sec++
					//fmt.Println("\n Tick at", t)
					// Defining duration
					// of Seconds method
					if timer30Sec > 60 {
						timer30Sec = 0
						timer30Min++
					}

					timeText = //fmt.Sprintf("%02d:%02d", timer30Min, timer30Sec)
				}
			}
		}()

		// Handler for the Promise: this is a JS function
		// It receives two arguments, which are JS functions themselves: resolve and reject
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			resolve.Invoke("Tick30 started")
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}*/

/*func tick5() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		moverYourButt = false
		bringYourButtBack = false
		ticker5 = time.NewTicker(time.Second)
		doneTicker5 = make(chan bool)
		var timer5Min, timer5Sec int

		go func() {
			for {
				select {
				case <-doneTicker25:
					return
				case t := <-ticker25.C:
					timer5Sec++
					//fmt.Println("\n Tick at", t)
					// Defining duration
					// of Seconds method
					if timer5Sec > 60 {
						timer5Sec = 0
						timer5Min++
					}

					timeText = //fmt.Sprintf("%02d:%02d", timer5Min, timer5Sec)
				}
			}
		}()

		// Handler for the Promise: this is a JS function
		// It receives two arguments, which are JS functions themselves: resolve and reject
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			resolve.Invoke("Tick5 started")
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}*/

/*func stopTick25() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		if ticker25 != nil {
			ticker25.Stop()
			doneTicker25 <- true

			// Handler for the Promise: this is a JS function
			// It receives two arguments, which are JS functions themselves: resolve and reject
			handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				resolve := args[0]
				resolve.Invoke("Tick25 stopped")
				return nil
			})

			// Create and return the Promise object
			promiseConstructor := js.Global().Get("Promise")
			return promiseConstructor.New(handler)
		}
		return nil
	})
}*/

/*func stopTick30() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		if ticker30 != nil {
			ticker30.Stop()
			doneTicker30 <- true

			// Handler for the Promise: this is a JS function
			// It receives two arguments, which are JS functions themselves: resolve and reject
			handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				resolve := args[0]
				resolve.Invoke("Tick30 stopped")
				return nil
			})

			// Create and return the Promise object
			promiseConstructor := js.Global().Get("Promise")
			return promiseConstructor.New(handler)
		}
		return nil
	})
}*/

/*func stopTick5() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		if ticker5 != nil {
			ticker5.Stop()
			doneTicker5 <- true

			// Handler for the Promise: this is a JS function
			// It receives two arguments, which are JS functions themselves: resolve and reject
			handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				resolve := args[0]
				resolve.Invoke("Tick30 stopped")
				return nil
			})

			// Create and return the Promise object
			promiseConstructor := js.Global().Get("Promise")
			return promiseConstructor.New(handler)
		}
		return nil
	})
}*/

/*func mover25() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		moverYourButt = true
		bringYourButtBack = false

		// Handler for the Promise: this is a JS function
		// It receives two arguments, which are JS functions themselves: resolve and reject
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			resolve.Invoke(moveYourAssText)
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
*/
/*func break5() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		moverYourButt = false
		bringYourButtBack = true

		stopTick5()

		// Handler for the Promise: this is a JS function
		// It receives two arguments, which are JS functions themselves: resolve and reject
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			resolve.Invoke("Bring your Ass back")
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}*/

/*func break30() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		moverYourButt = false
		bringYourButtBack = true

		stopTick30()

		// Handler for the Promise: this is a JS function
		// It receives two arguments, which are JS functions themselves: resolve and reject
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			resolve.Invoke(bringYourAssBackText)
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}*/

/*func break0() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		moverYourButt = false
		bringYourButtBack = true

		stopTick25()

		// Handler for the Promise: this is a JS function
		// It receives two arguments, which are JS functions themselves: resolve and reject
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			resolve.Invoke(bringYourAssBackText)
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
*/
func (g *Game) drawTiles(screen *ebiten.Image) {
	const (
		nx = screenWidth / tileSize
		ny = screenHeight / tileSize
	)

	op := &ebiten.DrawImageOptions{}
	for i := -2; i < nx+1; i++ {
		// ground
		op.GeoM.Reset()
		op.GeoM.Translate(float64(i*tileSize-floorMod(g.cameraX, tileSize)),
			float64((ny-1)*tileSize-floorMod(g.cameraY, tileSize)))
		screen.DrawImage(tilesImage.SubImage(image.Rect(0, 0, tileSize, tileSize)).(*ebiten.Image), op)

	}
}

func (g *Game) drawBulldozer(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	//op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate((screenWidth/bulldozerSize)+g.bX, screenHeight/1.25)
	var i int
	if bulldozerStarted {
		i = (g.count / 10) % bulldozerFrameNum
	} else {
		i = 0
	}

	//fmt.Println("G Count is: ", g.count)
	//fmt.Println("The bulldozer sx is: ", i)
	//fmt.Println("The butt sx is: ", i)

	sx, sy := frameOX+i*bulldozerFrameWidth, frameOY
	screen.DrawImage(bulldozerImage.SubImage(image.Rect(sx, sy, sx+bulldozerFrameWidth, sy+bulldozerFrameHeight)).(*ebiten.Image), op)
	//screen.DrawImage(bulldozerImage.SubImage(image.Rect(0, 0, 72, 72)).(*ebiten.Image), op)
	//screen.DrawImage(bulldozerImage, nil)
}

func (g *Game) drawChair(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Rotate(g.cVY)
	op.GeoM.Translate((screenWidth/chairSize)+g.cX, screenHeight/1.15)
	sx, sy := frameOX, frameOY
	screen.DrawImage(chairImage.SubImage(image.Rect(sx, sy, sx+chairFrameWidth, sy+chairFrameHeight)).(*ebiten.Image), op)
}

func (g *Game) drawButt(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	//op.GeoM.Rotate(g.buVY)
	op.GeoM.Translate((screenWidth/buttSize)+g.buX, screenHeight/g.buY)
	var i int
	if bulldozerStarted {
		i = (g.count / 15) % 2
	} else if moverYourButt || bringYourButtBack {
		t := (g.count / 15) % 2
		if t == 1 {
			i = 2
		} else if t == 0 {
			i = 3
		}
	} else {
		i = 0
	}

	////fmt.Println("G Count is: ", g.count)
	////fmt.Println("The butt sx is: ", i)

	sx, sy := frameOX+i*buttFrameWidth, frameOY
	screen.DrawImage(buttImage.SubImage(image.Rect(sx, sy, sx+buttFrameWidth, sy+buttFrameHeight)).(*ebiten.Image), op)
}

func main() {

	//fmt.Println("Go Web Assembly")
	/*	js.Global().Set("tick25", tick25())
		js.Global().Set("stopTick25", stopTick25())
		js.Global().Set("tick30", tick30())
		js.Global().Set("stopTick30", stopTick30())
		js.Global().Set("tick5", tick5())
		js.Global().Set("stopTick5", stopTick5())
		js.Global().Set("mover25", mover25())
		js.Global().Set("break5", break5())
		js.Global().Set("break30", break30())
		js.Global().Set("break0", break0())*/
	js.Global().Set("setColorMode", js.FuncOf(setColorMode))
	js.Global().Set("setErrorStatus", js.FuncOf(setErrorStatus))
	js.Global().Set("setConnectionStatus", js.FuncOf(setConnectionStatus))
	js.Global().Set("setClock", js.FuncOf(setClock))
	js.Global().Set("setTick", js.FuncOf(setTick))
	js.Global().Set("setScore", js.FuncOf(setScore))
	js.Global().Set("setButtonTrigger", js.FuncOf(setButtonTrigger))

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Butt Mover")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}

	<-make(chan bool)
}
