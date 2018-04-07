package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Coordinates struct {
	Lat  string `json:"lat"`
	Long string `json:"lng"`
}

func main() {
	key := os.Getenv("KEY")
	if key == "" {
		fmt.Println("key is missing")
		return
	}
	if len(strings.TrimSpace(key)) == 0 {
		fmt.Println("key is empty")
		return
	}

	j := Coordinates{}
	v, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("input read all error %s", err)
	}

	err = json.Unmarshal(v, &j)
	if err != nil {
		fmt.Printf("josn unmarshall error %s", err)
	}

	if len(j.Long) == 0 && len(j.Lat) == 0 {
		fmt.Printf("Lat and/or Lng positions error - not set")
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://maps.googleapis.com/maps/api/geocode/json?latlng=" + j.Lat + "," + j.Long + "&key=" + key)
	if err != nil {
		fmt.Printf("http get error %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("body error %s", err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}
}
