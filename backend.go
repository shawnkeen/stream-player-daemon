package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	StationFileName = "station"
	IDFileName      = "id"
	PIDFileName     = "pid"
	URLFileName     = "url"
	TagFileName     = "tag"
	ChangeVolCmd    = "./chvol"
)

var (
	URLFilePath     string
	PIDFilePath     string
	TagFilePath     string
	StationFilePath string
	IDFilePath      string
)

func init() {
	log.Println(config.runDir)
	_, err := os.Stat(config.runDir)
	if os.IsNotExist(err) {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	URLFilePath = filepath.FromSlash(config.runDir + "/" + URLFileName)
	PIDFilePath = filepath.FromSlash(config.runDir + "/" + PIDFileName)
	TagFilePath = filepath.FromSlash(config.runDir + "/" + TagFileName)
	StationFilePath = filepath.FromSlash(config.runDir + "/" + StationFileName)
	IDFilePath = filepath.FromSlash(config.runDir + "/" + IDFileName)
}

func writeToStatusFile(filePath string, entry string, doAppend bool) error {
	var file *os.File
	var err error
	if doAppend {
		file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0660))
		if err != nil {
			log.Println(err.Error())
			return err
		}
		defer file.Close()
	} else {
		file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0660))
		if err != nil {
			log.Println(err.Error())
			return err
		}
		defer file.Close()
	}
	if file == nil {
		log.Printf("No file opened for '%s'", filePath)
		return nil
	}
	log.Printf("Writing to file '%s': '%s'", filePath, entry)
	_, err = file.WriteString(entry + "\n")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func readFromStatusFile(filePath string) ([]string, error) {
	out := make([]string, 0)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return out, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return out, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		out = append(out, scanner.Text())
	}
	return out, scanner.Err()
}

func readStatusFromFiles() (*Status, error) {
	pids, err := readFromStatusFile(PIDFilePath)
	stationNames, err := readFromStatusFile(StationFilePath)
	urls, err := readFromStatusFile(URLFilePath)
	tags, err := readFromStatusFile(TagFilePath)
	// id, err := readFromStatusFile(IDFilePath)
	if err != nil {
		return nil, err
	}

	var index int

	log.Printf("%v, %d", pids, len(pids))

	if len(pids) == 0 {
		return &Status{Volume: -1, CurrStationID: 0, Tag: ""}, nil
	}

	if len(stationNames) == 0 {
		return nil, errors.New("Could not get station name.")
	}

	found := false
	for i, station := range config.stations {
		if station.Name == stationNames[0] && station.URL == urls[0] {
			index = i
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("Currently playing station name does not match any name in config.")
	}

	volume, err := backendGetVolume()
	if err != nil {
		log.Println(err.Error())
	}

	var tag string
	if len(tags) == 0 {
		tag = ""
	} else {
		tag = tags[0]
	}

	return &Status{Volume: volume, CurrStationID: index, Tag: tag}, nil
}

func backendGetStatus() (*Status, error) {
	return readStatusFromFiles()
}

func backendGetVolume() (int, error) {
	volString, err := exec.Command(ChangeVolCmd).Output()
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(string(strings.Replace(string(volString[:]), "\n", "", -1)))
}

func backendSetVolume(vol int) error {
	return exec.Command(ChangeVolCmd, "set", strconv.Itoa(vol)).Run()
}

func backendIncVolume(perc int) error {
	return exec.Command(ChangeVolCmd, "inc", strconv.Itoa(perc)).Run()
}

func backendDecVolume(perc int) error {
	return exec.Command(ChangeVolCmd, "dec", strconv.Itoa(perc)).Run()
}

func backendStopPlayback() error {
	pids, err := readFromStatusFile(PIDFilePath)
	if err != nil {
		return err
	}
	for _, pidString := range pids {
		pid, err := strconv.Atoi(pidString)
		if err == nil {
			syscall.Kill(pid, syscall.SIGKILL)
		}
	}
	writeToStatusFile(TagFilePath, "", false)
	return os.Remove(PIDFilePath)
}

func backendPlayStation(station *Station) error {
	if station == nil {
		return nil
	}
	//player := config.playerCmd
	cmd := exec.Command(config.playerCmd, "-t", TagFilePath, station.URL)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmderr := cmd.Start()
	if cmderr != nil {
		return cmderr
	}
	pid := strconv.Itoa(cmd.Process.Pid)
	writeToStatusFile(PIDFilePath, pid, true)
	writeToStatusFile(StationFilePath, station.Name, false)
	writeToStatusFile(URLFilePath, station.URL, false)
	return nil
}
