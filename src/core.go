/*
 * @Author: V4 Games
 * @Date: 2018-11-14 02:26:53
 * @Last Modified by:   Dominik Madarász (zaklaus@madaraszd.net)
 * @Last Modified time: 2018-11-14 02:26:53
 */

package main // madaraszd.net/zaklaus/rurik

var (
	// MainCamera is the primary camera used for the viewport
	MainCamera *Object

	// LocalPlayer is player's object
	LocalPlayer *Object

	// DebugMode switch
	DebugMode = true

	// TimeScale is game update time scale
	TimeScale = 1
)

const (
	// GameVersion describes itself
	GameVersion = "1.0.0"
)

// InitCore initializes the game engine
func InitCore() {
	initObjectTypes()
	InitDatabase()
}
