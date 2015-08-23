package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// valid commands and their arity
var commands = map[string]int{
	"version":  0,
	"current":  0,
	"play":     1,
	"next":     0,
	"prev":     0,
	"stop":     0,
	"volume":   2,
	"help":     0,
	"stations": 0,
	"status":   0,
}

const (
	PROT_OK      = 100
	PROT_GENERR  = 200
	PROT_UKNCMD  = 201
	PROT_INVARG  = 202
	PROT_SYNTAX  = 203
	PROT_INTERR  = 300
	PROT_NAVAIL  = 301
	PROT_OORANGE = 302
	PROT_VERSION = "1.1.0"
)

var commandCodeMessages = map[int]string{
	200: "error",
	201: "unknown command",
	202: "invalid argument",
	203: "incorrect syntax",
	300: "internal server error",
	301: "requested property not available",
	302: "number out of range",
}

type Station struct {
	Name string
	URL  string
}

type Status struct {
	Volume        int
	CurrStationID int
	Tag           string
}

type ProtReturn struct {
	Code    int
	Message string
}

func (r ProtReturn) String() string {
	if r.Code > PROT_OK {
		var msg string
		if r.Message != "" {
			msg = r.Message
		} else {
			msg = commandCodeMessages[r.Code]
		}
		return fmt.Sprintf("ERROR %d %v", r.Code, msg)
	}
	return fmt.Sprintf("OK %d", r.Code)
}

func init() {
}

func (status *Status) JSON() ([]byte, error) {
	return json.Marshal(status)
}

func (s *Status) String() string {
	out := ""
	station := protStatusFromID(s.CurrStationID)
	out += "station: " + station.Name + "\n"
	out += "url: " + station.URL + "\n"
	out += "id: " + strconv.Itoa(s.CurrStationID) + "\n"
	out += "tag: " + s.Tag + "\n"
	out += "volume: " + strconv.Itoa(s.Volume)
	return out
}

func protHello() string {
	return "OK MSPD " + PROT_VERSION
}

func protStatusFromID(id int) *Station {
	if id < 0 || id >= len(config.stations) {
		return nil
	}
	return &config.stations[id]
}

func protCurrent() (string, ProtReturn) {
	status, err := readStatusFromFiles()
	if err != nil {
		return "", ProtReturn{PROT_GENERR, err.Error()}
	}
	tag := ""
	if status.CurrStationID != 0 {
		tag = status.Tag
	}
	return tag, ProtReturn{status.CurrStationID, ""}
}

func protPlay(stationID int) (string, ProtReturn) {
	if stationID < 0 || stationID >= len(config.stations) {
		return "", ProtReturn{PROT_OORANGE, ""}
	}
	station := config.stations[stationID]
	backendStopPlayback()
	err := backendPlayStation(&station)
	if err != nil {
		return "", ProtReturn{PROT_GENERR, fmt.Sprintf("Could not start player '%s': %s", config.playerCmd, err.Error())}
	}
	return "", ProtReturn{PROT_OK, ""}
}

func protStop() (string, ProtReturn) {
	backendStopPlayback()
	return "", ProtReturn{PROT_OK, ""}
}

func protDelta(delta int) (string, ProtReturn) {
	status, err := readStatusFromFiles()
	if err != nil {
		return "", ProtReturn{PROT_GENERR, fmt.Sprintf("Could not get status: %s", err.Error())}
	}
	// there is an offset of 1, because of the dummy station
	currentID := status.CurrStationID - 1
	length := len(config.stations) - 1
	// we do the modulo operation without the offset, and then add it again
	newID := (currentID + delta) % length
	if newID < 0 {
		newID += length
	}
	fmt.Printf("station change by %d from %d to %d\n", delta, currentID, newID)
	return protPlay(newID + 1)
}

func protHelp() (string, ProtReturn) {
	out := ""
	first := true
	for s, _ := range commands {
		if !first {
			out += "\n"
		} else {
			first = false
		}
		out += s
	}
	return out, ProtReturn{PROT_OK, ""}
}

func protStations() (string, ProtReturn) {
	out := ""
	first := true
	for i, station := range config.stations {
		if i == 0 {
			continue
		}
		if first {
			first = false
		} else {
			out += "\n"
		}
		out += strconv.Itoa(i) + " " + station.Name
	}
	return out, ProtReturn{PROT_OK, ""}
}

func protVolume(arguments []string) (string, ProtReturn) {
	vol, err := strconv.Atoi(arguments[1])
	if err != nil {
		return "", ProtReturn{PROT_INVARG, ""}
	}
	switch arguments[0] {
	case "set":
		{
			err = backendSetVolume(vol)
			if err != nil {
				return "", ProtReturn{PROT_GENERR, err.Error()}
			}
		}
	case "inc":
		{
			err = backendIncVolume(vol)
			if err != nil {
				return "", ProtReturn{PROT_GENERR, err.Error()}
			}
		}
	case "dec":
		{
			err = backendDecVolume(vol)
			if err != nil {
				return "", ProtReturn{PROT_GENERR, err.Error()}
			}
		}
		// default:
		// 	return "", ProtReturn{PROT_INVARG, ""}
	}
	return "", ProtReturn{PROT_OK, ""}
}

func decodeStationsJSON(rawString string) ([]Station, error) {
	rawBytes := []byte(rawString)
	var stations []Station
	err := json.Unmarshal(rawBytes, &stations)
	return stations, err
}

func encodeStationsJSON(stations *[]Station) ([]byte, error) {
	return json.Marshal(stations)
}

func protResponseEnd(line string) bool {
	if strings.HasPrefix(line, "OK ") || strings.HasPrefix(line, "ERROR ") {
		return true
	}
	return false
}

func handleResponse(line string, answer []string) error {
	tokens := strings.Split(line, " ")
	if len(tokens) < 2 {
		return errors.New("Invalid Response")
	}
	returnCode, err := strconv.Atoi(tokens[1])
	if err != nil {
		return err
	}
	if returnCode == PROT_OK {

	}
	return nil
}

func handleRequest(line string) (string, ProtReturn) {
	tokens := strings.Split(line, " ")

	for i := range tokens {
		fmt.Errorf("tokens: %s", tokens[i])
	}

	command := tokens[0]

	if arity, ok := commands[command]; ok {
		if len(tokens) != arity+1 {
			return "", ProtReturn{PROT_SYNTAX, fmt.Sprintf("Wrong number of arguments. '%s' takes %d arguments.", command, arity)}
		}
		switch command {
		case "version":
			return PROT_VERSION, ProtReturn{PROT_OK, ""}
		case "current":
			return protCurrent()
		case "play":
			{
				stationID, err := strconv.Atoi(tokens[1])
				if err != nil {
					return "", ProtReturn{PROT_INVARG, ""}
				}
				return protPlay(stationID)
			}
		case "next":
			return protDelta(1)
		case "prev":
			return protDelta(-1)
		case "stop":
			return protStop()
		case "volume":
			{
				protVolume(tokens[1:])
			}
		case "help":
			return protHelp()
		case "stations":
			return protStations()
		case "status":
			{
				status, err := backendGetStatus()
				if err != nil {
					return "", ProtReturn{PROT_GENERR, err.Error()}
				}
				return status.String(), ProtReturn{PROT_OK, ""}
			}
		default:
			return "", ProtReturn{PROT_INTERR, "Internal Server Error"}
		}
	} else {
		return "", ProtReturn{PROT_UKNCMD, fmt.Sprintf("Unknown command '%s'.", command)}
	}
	return "", ProtReturn{PROT_OK, ""}
}
