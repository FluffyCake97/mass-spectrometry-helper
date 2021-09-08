package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

func getInput(msg string, expectNum bool) (str string) {
	fmt.Println(msg)
	fmt.Scanln(&str)
	if expectNum {
		_, err := strconv.Atoi(str)
		if err != nil {
			getInput("You have inputted characters that cannot be converted into Int, try again", expectNum)
		}
	}
	return str
}

func main() {
	fmt.Println("Loading configuration file...")
	bytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Validating data files...")
	config, err = ValidationFileStructure(config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Loading data files...")
	res, err := ReadFile(config)

	fmt.Println("Handling data files...")
	var finalResults results
	for _, obj := range res {
		for _, o := range obj {
			finalResults = append(finalResults, o)
		}
	}
	finalResults.Sort()

	saveFileToDisk(config, finalResults, config.columnTitle)
	_ = getInput("The final results had been save to your disk... Don't forget buy me a coca cola.", false)
	return
}
