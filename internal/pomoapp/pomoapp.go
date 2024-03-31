package pomoapp

import (
	"encoding/json"
	"log"
	"os"
)

type Settings struct {
	SettingsPath string
	LibraryPath  string
	AutoPlay     bool
	Shuffle      bool
	LinkPlayers  bool
}

func NewSettings(
	settingsPath string,
	libraryPath string,
	autoPlay bool,
	shuffle bool,
	linkPlayers bool,
) *Settings {
	return &Settings{
		SettingsPath: settingsPath,
		LibraryPath:  libraryPath,
		AutoPlay:     autoPlay,
		Shuffle:      shuffle,
		LinkPlayers:  linkPlayers,
	}
}

func (settings *Settings) Save(
	libraryPath string,
	autoPlayChecked bool,
	shuffleChecked bool,
	linkPlayersChecked bool,
) {
	settings.LibraryPath = libraryPath
	settings.AutoPlay = autoPlayChecked
	settings.Shuffle = shuffleChecked
	settings.LinkPlayers = linkPlayersChecked

	// TODO(map) Handle errors gracefully
	file, _ := json.MarshalIndent(settings, "", "    ")
	_ = os.WriteFile(settings.SettingsPath, file, 0644)
}

func (settings *Settings) Load() {
	// func loadSettings(pathToSettings string) {
	settingsFile, _ := os.ReadFile(settings.SettingsPath)
	err := json.Unmarshal(settingsFile, &settings)
	if err != nil {
		log.Print("Failure in unmarshelling the settings data")
	}
}
