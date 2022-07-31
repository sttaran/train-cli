package main

import (
	"encoding/json"
	"errors"
	"fmt"
	colorPrint "github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

var trainTimeFormat = "15:04:05"
var criterias = map[string]struct{}{
	"price":          {},
	"arrival-time":   {},
	"departure-time": {},
}

type Trains []Train

type Train struct {
	TrainID            int       `json:"trainId"`
	DepartureStationID int       `json:"departureStationId"`
	ArrivalStationID   int       `json:"arrivalStationId"`
	Price              float32   `json:"price"`
	ArrivalTime        time.Time `json:"arrivalTime"`
	DepartureTime      time.Time `json:"departureTime"`
}

func (t *Train) UnmarshalJSON(data []byte) error {
	var err error
	type Alias Train
	aux := &struct {
		ArrivalTime   string `json:"arrivalTime"`
		DepartureTime string `json:"departureTime"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.ArrivalTime, err = time.Parse(trainTimeFormat, aux.ArrivalTime)
	if err != nil {
		return err
	}
	t.DepartureTime, err = time.Parse(trainTimeFormat, aux.DepartureTime)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var err error
	var depStation, arrStation, criteria string

	colorPrint.Green("Please enter the departure station.\n")
	_, err = fmt.Scan(&depStation)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	colorPrint.Blue("Please enter the station of arrival.\n")
	_, err = fmt.Scan(&arrStation)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	colorPrint.Green("Please enter the criteria to sort search result.\n")
	_, err = fmt.Scan(&criteria)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if _, ok := criterias[criteria]; !ok {
		fmt.Println("unsupported criteria")
		os.Exit(1)
	}

	trains, err := FindTrains(depStation, arrStation, criteria)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	top := trains[0:3]
	indent, err := json.MarshalIndent(top, "", " ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	colorPrint.Green(string(indent))

}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	var trains Trains
	a, err := strconv.Atoi(arrivalStation)
	if err != nil {
		return nil, errors.New("bad arrival station input")
	}
	d, err := strconv.Atoi(departureStation)
	if err != nil {
		return nil, errors.New("bad departure station input")
	}

	bytes, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(bytes, &trains)
	if err != nil {
		log.Fatalln(err)
	}

	var filteredTrains Trains
	for _, train := range trains {
		if train.ArrivalStationID == a && train.DepartureStationID == d {
			filteredTrains = append(filteredTrains, train)
		}
	}

	switch criteria {
	case "price":
		sort.SliceStable(filteredTrains, func(i, j int) bool {
			return filteredTrains[i].Price < filteredTrains[j].Price
		})
	case "arrival-time":
		sort.SliceStable(filteredTrains, func(i, j int) bool {
			return filteredTrains[i].ArrivalTime.Before(filteredTrains[j].ArrivalTime)
		})
	case "departure-time":
		sort.SliceStable(filteredTrains, func(i, j int) bool {
			return filteredTrains[i].DepartureTime.Before(filteredTrains[j].DepartureTime)
		})
	}

	return filteredTrains, nil
}
