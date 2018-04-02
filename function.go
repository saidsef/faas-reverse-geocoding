package main

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Coordinates struct {
	Lat  string `json:"lat"`
	Long string `json:"long"`
}

type Resp struct {
	Results []string
}

func main() {
	key := os.Getenv("KEY")
	j := Coordinates{}
	v, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("io error %s", err.Error())
	}

	err = json.Unmarshal(v, &j)
	if err != nil {
		log.Fatalf("Unmarshal Coordinates Error! %s", err.Error())
	}

	resp, err := http.Get("https://maps.googleapis.com/maps/api/geocode/json?latlng=" + j.Lat + "," + j.Long + "&key=" + key)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("request Get error! %s", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("request ReadAll error! %s", err.Error())
	}

	binary.Write(os.Stdout, binary.LittleEndian, body)
}
