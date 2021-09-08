package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func Contains(slice []string, s string) int {
	for index, value := range slice {
		if value == s {
			return index
		}
	}
	return -1
}

func ValidationFileStructure(config Config) (newConfig Config, err error) {
	newConfig = config
	var columnTitle []string
	var x_index []int
	var y_index []int
	var z_index []int
	var flag = true
	for _, obj := range config.Files {
		fi, err := os.Open("./" + obj.Filepath)
		if err != nil {
			return config, err
		}
		defer fi.Close()
		br := bufio.NewReader(fi)
		a, _, _ := br.ReadLine()
		lineSlice := strings.Split(string(a), "\t")
		if string(lineSlice[0]) == "Start_push" {
			columnTitle = append(columnTitle, string(a))
			x_index = append(x_index, Contains(lineSlice, "X"))
			y_index = append(y_index, Contains(lineSlice, "Y"))
			z_index = append(z_index, Contains(lineSlice, "Z"))
		}
	}
	flag = validateIsNumArrayElementSame(x_index)
	flag = validateIsNumArrayElementSame(y_index)
	flag = validateIsNumArrayElementSame(z_index)
	flag = validateIsStrArrayElementSame(columnTitle)
	if flag {
		newConfig.columnTitle = columnTitle[0]
		newConfig.x_index = x_index[0]
		newConfig.y_index = y_index[0]
		newConfig.z_index = z_index[0]
		columnTitleHeader := strings.Split(columnTitle[0], "\t")
		for i, obj := range config.Filters {
			newConfig.Filters[i].marker_index = Contains(columnTitleHeader, obj.ColumnHeader)
		}
		return newConfig, nil
	} else {
		return config, errors.New("an error occurred during the verification of the file column header, Please check all the files")
	}
}

func validateIsNumArrayElementSame(array []int) bool {
	var temp int
	flag := true
	for i, obj := range array {
		if i == 0 {
			temp = obj
		} else {
			if obj != temp {
				flag = false
			}
		}
	}
	return flag
}
func validateIsStrArrayElementSame(array []string) bool {
	var temp string
	flag := true
	for i, obj := range array {
		if i == 0 {
			temp = obj
		} else {
			if obj != temp {
				flag = false
			}
		}
	}
	return flag
}

func ReadFile(config Config) (res []results, err error) {
	ColumnNum := len(strings.Split(config.columnTitle, "\t"))
	for _, o := range config.Files {
		var data results = nil
		fi, err := os.Open("./" + o.Filepath)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return res, err
		}
		defer fi.Close()
		br := bufio.NewReader(fi)
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			lineSlice := strings.Split(string(a), "\t")
			if len(lineSlice) != ColumnNum {
				return res, errors.New("there is some wrong in file " + o.Filepath)
			}
			if string(lineSlice[0]) == "Start_push" {
				continue
			}
			var row Result
			for index, slice := range lineSlice {
				switch {
				case index < config.x_index:
					number, err := strconv.Atoi(slice)
					if err != nil {
						return res, err
					}
					row.info = append(row.info, number)
				case index == config.x_index:
					number, err := strconv.Atoi(slice)
					if err != nil {
						return res, err
					}
					row.X = number + o.X_offset
				case index == config.y_index:
					number, err := strconv.Atoi(slice)
					if err != nil {
						return res, err
					}
					row.Y = number + o.Y_offset
				case index == config.z_index:
					number, err := strconv.Atoi(slice)
					if err != nil {
						return res, err
					}
					row.Z = number + o.Z_offset
				case index > config.z_index:
					number, err := strconv.ParseFloat(slice, 64)
					if err != nil {
						return res, err
					}
					for _, filter := range config.Filters {
						if index == filter.marker_index {
							if number > filter.UpperThreshold || number < filter.LowerThreshold {
								number = filter.Default
							}
						}
					}
					row.marker = append(row.marker, number)
				}
			}
			data = append(data, row)
		}
		res = append(res, data)
	}
	return res, nil
}

//
func checkFileIsExist(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return
	}
	os.Remove(filename)
}

func saveFileToDisk(config Config, data results, columnTitle string) {
	filename := "./final_results.txt"
	checkFileIsExist(filename)
	f, _ := os.Create(filename)
	defer f.Close()
	_, _ = io.WriteString(f, columnTitle+"\n")
	for _, o := range data {
		content := ""
		for i, cell := range o.info {
			if i != 0 {
				content = content + "\t"
			}
			content = content + strconv.Itoa(cell)
		}
		content = content + "\t" + strconv.Itoa(o.X) + "\t" + strconv.Itoa(o.Y) + "\t" + strconv.Itoa(o.Z)
		for _, cell := range o.marker {
			str := strconv.FormatFloat(cell, 'f', 3, 64)
			content = content + "\t" + str
		}
		content = content + "\n"
		_, _ = io.WriteString(f, content)
	}
}
