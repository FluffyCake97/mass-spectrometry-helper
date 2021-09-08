package main

import "sort"

type Result struct {
	info   []int
	marker []float64
	X      int
	Y      int
	Z      int
}

type results []Result

type File struct {
	Filepath string `json:"filepath"`
	Width    int    `json:"width,string"`
	Height   int    `json:"height,string"`
	X_offset int    `json:"x_offset,string"`
	Y_offset int    `json:"y_offset,string"`
	Z_offset int    `json:"z_offset,string"`
}

type Config struct {
	Files       Files   `json:"files"`
	Filters     Filters `json:"filters"`
	columnTitle string
	x_index     int
	y_index     int
	z_index     int
}

type Filter struct {
	ColumnHeader   string  `json:"column_header"`
	UpperThreshold float64 `json:"upper_threshold,string"`
	LowerThreshold float64 `json:"lower_threshold,string"`
	Default        float64 `json:"default,string"`
	marker_index   int
}

type Filters []Filter

type Files []File

func (p results) Len() int {
	return len(p)
}

func (p results) Less(i, j int) bool {
	if p[i].Y < p[j].Y {
		return true
	}
	if p[i].Y > p[j].Y {
		return false
	}

	return p[i].X < p[j].X
}

func (p results) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *results) Sort() {
	sort.Sort(p)
}
