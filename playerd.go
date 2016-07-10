package main

import (
	"github.com/alyu/configparser"
	"log"
	"strings"
)

type Config struct {
	runDir            string
	playerCmd         string
	stations          []Station
	stationNames      []string
	availableStations map[string]Station
}

var runDir string
var config Config

func init() {
	config.readFromFile("config")
	runDir = config.runDir

	log.Printf("Loaded config with %d stations available and %d set.", len(config.availableStations), len(config.stations))
}

func (c *Config) readFromFile(filePath string) error {
	config, err := configparser.Read(filePath)
	if err != nil {
		return err
	}

	section, err := config.Section("Global")
	if err != nil {
		return err
	}

	c.playerCmd = section.Options()["player"]
	c.runDir = section.Options()["dir"]
	//c.stationNames = []string{"dummy"}
	c.stationNames = append([]string{"dummy"}, strings.Split(section.Options()["stations"], " ")...)

	sections, _ := config.AllSections()
	c.availableStations = make(map[string]Station)
	for _, sec := range sections {
		if sec.Exists("name") && sec.Exists("url") {
			c.availableStations[sec.Name()] = Station{Name: sec.Options()["name"], URL: sec.Options()["url"]}
		}
	}

	c.stations = make([]Station, len(c.stationNames))
	for i, shortName := range c.stationNames {
		c.stations[i] = c.availableStations[shortName]
	}

	return nil
}

func main() {

	// for name, station := range config.availableStations {
	// 	println(name + " " + station.Name)
	// }

	println(-1 % 9)

	status, err := readStatusFromFiles()
	if err == nil {
		println(status.String())
	}

	startServer()
	//client()
}
