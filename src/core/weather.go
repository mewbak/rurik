/*
   Copyright 2019 Dominik Madarász <zaklaus@madaraszd.net>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package core

import (
	"strconv"

	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/system"
)

var (
	// SkyColor is the tint color used for drawn sprites/tiles
	SkyColor rl.Color

	// WeatherTimeScale specifies time cycle scale
	WeatherTimeScale float64

	weatherIsCollapsed = true
)

type weatherStage struct {
	Name     string
	Color    rl.Vector3
	Duration float64
}

// Weather represents the map time and weather
type Weather struct {
	UseTimeCycle    bool
	SkyStageName    string
	SkyTime         float64
	SkyTargetTime   float64
	SkyStageIndex   int
	SkyStages       []weatherStage
	SkyLastColor    rl.Vector3
	SkyCurrentColor rl.Vector3
	SkyTargetColor  rl.Vector3
}

// WeatherInit sets up the mood by initializing Sky color tint and other properties
func (w *Weather) WeatherInit(cmap *Map) {
	var err error
	w.SkyCurrentColor, err = GetColorFromHex(cmap.tilemap.Properties.GetString("skyColor"))

	if err != nil {
		SkyColor = rl.White
	} else {
		SkyColor = Vec3ToColor(w.SkyCurrentColor)
	}

	w.SkyStages = []weatherStage{}

	w.appendSkyStage(cmap, "skyRiseColor", "riseDuration")
	w.appendSkyStage(cmap, "skyDayColor", "dayDuration")
	w.appendSkyStage(cmap, "skyDawnColor", "dawnDuration")
	w.appendSkyStage(cmap, "skyNightColor", "nightDuration")

	if len(w.SkyStages) > 0 {
		w.SkyLastColor = w.SkyCurrentColor
		w.SkyStageName = w.SkyStages[0].Name
		w.SkyTargetColor = w.SkyStages[0].Color
		w.SkyTime = w.SkyStages[0].Duration
		w.SkyTargetTime = w.SkyTime
		SkyColor = Vec3ToColor(w.SkyCurrentColor)
		w.SkyStageIndex = 0

		if err != nil {
			w.SkyCurrentColor = w.SkyTargetColor
			w.nextSkyStage()
		}
	}

	weatherIsCollapsed = true
}

// UpdateWeather updates the time cycle and weather effects
func (w *Weather) UpdateWeather() {
	if w.UseTimeCycle {
		if w.SkyTime <= 0 {
			w.nextSkyStage()
		} else {
			w.SkyTime -= float64(system.FrameTime) * WeatherTimeScale
		}

		if w.SkyTargetTime != 0 {
			w.SkyCurrentColor = LerpColor(w.SkyLastColor, w.SkyTargetColor, 1-w.SkyTime/w.SkyTargetTime)
		} else {
			w.SkyCurrentColor = w.SkyTargetColor
		}

		SkyColor = Vec3ToColor(w.SkyCurrentColor)
	}
}

// DrawWeather draws weather effects
func (w *Weather) DrawWeather() {

}

func (w *Weather) nextSkyStage() {
	w.SkyStageIndex++

	if w.SkyStageIndex >= len(w.SkyStages) {
		w.SkyStageIndex = 0
	}

	stage := w.SkyStages[w.SkyStageIndex]
	w.SkyStageName = stage.Name
	w.SkyTime = stage.Duration
	w.SkyTargetTime = w.SkyTime
	w.SkyTargetColor = stage.Color
	w.SkyLastColor = w.SkyCurrentColor
}

func (w *Weather) appendSkyStage(cmap *Map, SkyName, stageName string) {
	color, err := GetColorFromHex(cmap.tilemap.Properties.GetString(SkyName))

	if err == nil {
		w.UseTimeCycle = true
	} else {
		return
	}

	duration, _ := strconv.ParseFloat(cmap.tilemap.Properties.GetString(stageName), 10)

	w.SkyStages = append(w.SkyStages, weatherStage{
		Name:     SkyName,
		Color:    color,
		Duration: duration * 60,
	})
}
