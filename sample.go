/*
 * go-kdtree/sample: stuff for quickly finding a median by sampling
 *
 * Copyright (C) 2018 Pawel Foremski <pjf@foremski.pl>
 * Licensed to you under GNU GPL v3
 */

package kdtree

import "sort"

type sample_points struct {
	points Points
	axis   int
}

func (a sample_points) Len() int {
	return len(a.points)
}

func (a sample_points) Swap(i,j int) {
	a.points[i],a.points[j] = a.points[j],a.points[i]
}

func (a sample_points) Less(i,j int) bool {
	return (*a.points[i])[a.axis] < (*a.points[j])[a.axis]
}

// sample_median() finds median on given axis using uniform sampling
func sample_median(points Points, axis int, sample_size int) *Point {
	sample := sample_points{ axis: axis }

	// how big sample?
	size := len(points)
	if size < 3 { return points[size-1] }
	if size > sample_size { size = sample_size }
	sample.points = make(Points, 0, size)

	// how big step through points?
	plen := float64(len(points))
	step := plen / float64(size)

	// take sample
	for fi := 0.0; fi < plen; fi += step {
		sample.points = append(sample.points, points[int(fi)])
	}

	// sort it
	sort.Sort(sample)

	// take median
	return sample.points[size / 2]
}
