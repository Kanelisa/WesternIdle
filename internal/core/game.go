package core

import (
	_ "embed"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/opentype"

	"WesternIdle/internal/system"
	"WesternIdle/internal/ui"
)

// ---------------- Шрифт ----------------

//go:embed ..\..\assets\fonts\PrincetownSolid.ttf
var fontBytes []byte

var GameFont text.Face

func init() {
	rand.Seed(time.Now().UnixNano())

	tt, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}
	goFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 16,
		DPI:  72,
		//Hinting: opentype.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	GameFont = text.NewGoXFace(goFace)
}

// ---------------- Константы ----------------
const (
	ScreenWidth  = 900
	ScreenHeight = 600
)

// ---------------- Game ----------------
type Game struct {
	State      *system.GameState
	UI         *ui.UI
	LastUpdate time.Time
}

// NewGameInstance создаёт игру и инициализирует UI
func NewGameInstance() *Game {
	state := system.LoadGame()  // создаем пустой GameState с ресурсами
	system.InitLocations(state) // теперь map Locations создан и CurrentLocation установлен

	gameUI := ui.NewUI(state, GameFont, func(id string, dur float64) {
		state.StartAction(id)
	})

	return &Game{
		State:      state,
		UI:         gameUI,
		LastUpdate: time.Now(),
	}
}

// Update вызывается Ebiten каждый кадр
func (g *Game) Update() error {
	now := time.Now()
	delta := now.Sub(g.LastUpdate).Seconds()
	g.LastUpdate = now

	g.State.Update(delta)
	g.UI.Update() // обработка кликов, прогресс, события кнопок

	return nil
}

// Draw вызывается Ebiten для отрисовки
func (g *Game) Draw(screen *ebiten.Image) {
	g.UI.Draw(screen)
}

// Layout сообщает размеры окна
func (g *Game) Layout(w, h int) (int, int) { return ScreenWidth, ScreenHeight }
