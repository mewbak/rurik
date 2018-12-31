package main

import (
	"fmt"
	"math"

	jsoniter "github.com/json-iterator/go"
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

type demoGameMode struct {
	playState int
	textWave  int32
}

func (g *demoGameMode) Init() {
	core.LoadPlaylist("tracklist.txt")
	core.LoadNextTrack()

	// test class
	err := core.RegisterClass("demo_testclass", "TestClass", NewTestClass)

	// player class
	core.RegisterClass("player", "Player", NewPlayer)

	if err != nil {
		fmt.Printf("Custom type registration has failed: %s", err.Error())
	}

	g.playState = stateMenu

	if fmapload {
		g.playState = statePlay
		core.LoadMap(playMapName)
		core.InitMap()
	}

	initShaders()
}

func (g *demoGameMode) Shutdown() {}

func (g *demoGameMode) Update() {
	switch g.playState {
	case stateMenu:
		g.textWave = int32(math.Round(math.Sin(float64(rl.GetTime()) * 10)))

		if system.IsKeyPressed("use") {
			core.LoadMap("start")
			core.InitMap()

			g.playState = statePlay
		}

	case statePlay:
		core.UpdateMaps()
	}

	updateInternals(g)
}

func (g *demoGameMode) Serialize() string {
	data := demoGameSaveData{
		ObjectCounter: dynobjCounter,
	}

	ret, _ := jsoniter.MarshalToString(data)
	return ret
}

func (g *demoGameMode) Deserialize(data string) {
	var saveData demoGameSaveData
	jsoniter.UnmarshalFromString(data, &saveData)

	dynobjCounter = saveData.ObjectCounter
}

type demoGameSaveData struct {
	ObjectCounter int
}

func (g *demoGameMode) Draw() {
	drawBackground()

	rl.BeginMode2D(core.RenderCamera)
	{
		core.DrawMap(true)
	}
	rl.EndMode2D()
}

func (g *demoGameMode) DrawUI() {
	switch g.playState {
	case stateMenu:
		core.DrawTextCentered("Rurik Framework", system.ScreenWidth/2, system.ScreenHeight/2-20+g.textWave, 24, rl.RayWhite)
		core.DrawTextCentered("Press E/ENTER to continue", system.ScreenWidth/2, system.ScreenHeight/2+5+g.textWave, 14, rl.White)

	case statePlay:
		core.DrawMapUI()

		if core.CurrentMap.Name != "start" {

			rl.DrawText("Press F5 at any time to go back to the menu.", 5, 40, 12, rl.RayWhite)
			rl.DrawText("Press F2 to save and F3 to load a game state.", 5, 52, 12, rl.RayWhite)
			rl.DrawText("Press F9 to spawn a light object.", 5, 64, 12, rl.RayWhite)
			rl.DrawText("Use your mousewheel to zoom in/out.", 5, 76, 12, rl.RayWhite)
		} else {
			core.DrawTextCentered("Rurik Framework", system.ScreenWidth/2, system.ScreenHeight/2-20, 24, rl.RayWhite)
		}

		if core.CurrentMap.Name == "village" {
			// draw a minimap
			{
				rl.DrawRectangle(system.ScreenWidth-105, 5, 100, 100, rl.Blue)
				rl.DrawTexturePro(
					minimap.RenderTexture.Texture,
					rl.NewRectangle(0, 0,
						float32(minimap.RenderTexture.Texture.Width),
						float32(-minimap.RenderTexture.Texture.Height)),
					rl.NewRectangle(float32(system.ScreenWidth)-102, 8, 94, 94),
					rl.Vector2{},
					0,
					rl.White,
				)
			}

			// draw shadertoy example
			{
				rl.DrawRectangle(system.ScreenWidth-105, 110, 100, 100, rl.Fade(rl.Red, 0.6))
				rl.DrawTexturePro(
					shadertoy.RenderTexture.Texture,
					rl.NewRectangle(0, 0,
						float32(shadertoy.RenderTexture.Texture.Width),
						float32(shadertoy.RenderTexture.Texture.Height)),
					rl.NewRectangle(float32(system.ScreenWidth)-102, 113, 94, 94),
					rl.Vector2{},
					0,
					rl.White,
				)
			}
		}
	}
}

func (g *demoGameMode) PostDraw() {

	switch g.playState {
	case stateMenu:

	case statePlay:
		// Generates and applies the lightmaps
		core.UpdateLightingSolution()

		if core.CurrentMap.Name == "village" {
			bloom.Apply()
			minimap.Apply()
			shadertoy.Apply()
		}
	}

}

func (g *demoGameMode) LoadMap(mapName string) {
	core.FlushMaps()
	core.LoadMap(mapName)
	core.InitMap()
}