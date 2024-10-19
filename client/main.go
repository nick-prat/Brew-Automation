package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type BrewConfig struct {
	APIKey string `json:"api_key"`
}

type TempLog struct {
	TargetTemp float64 `json:"target_temp"`
	SetTemp    float64 `json:"set_temp"`
	FermentRun int     `json:"ferment_run_id"`
}

func main() {
	file, err := os.Open("config.json")
	if err != nil {
		panic(fmt.Sprintf("Couldn't open conifg %s", err.Error()))
	}

	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Sprintf("Couldn't read config file %s", err.Error()))
	}

	config := BrewConfig{}
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		panic(fmt.Sprintf("Couldn't parse config file %s", err.Error()))
	}

	fmt.Printf("init,%d\n", 1)

	scanner := bufio.NewScanner(os.Stdin)
	client := http.Client{}
	for scanner.Scan() {
		txt := string(scanner.Text())
		fmt.Println(txt)

		strs := strings.Split(txt, ",")
		targetTemp, err := strconv.ParseFloat(strs[0], 64)
		if err != nil {
			panic("Couldn't parse target temp")
		}

		setTemp, err := strconv.ParseFloat(strs[1], 64)
		if err != nil {
			panic("Couldn't parse set temp")
		}

		body, err := json.Marshal(TempLog{
			TargetTemp: targetTemp,
			SetTemp:    setTemp,
			FermentRun: 1,
		})

		if err != nil {
			println("Hi")
			panic("Failed parsing temp data for API request")
		}

		req, err := http.NewRequest("POST", "localhost:3333/temp-log", bytes.NewBuffer(body))

		if err != nil {
			println("Hi")
			panic("Failed creating request")
		}

		client.Do(req)

		fmt.Println("Temps: {}\n", targetTemp)
	}
}
