/*
 * go-kdtree: a kd-tree implementation in Golang
 *
 * Copyright (C) 2018 Pawel Foremski <pjf@foremski.pl>
 * Licensed to you under GNU GPL v3
 */

package kdtree

import "math"

// KDNode represents a kd-tree
type KDNode struct {
	point   *Point   // point stored at node
	axis    int      // splitting axis
	left    *KDNode  // left child:  points less than this point (on axis)
	right   *KDNode  // right child: points greater than this point (on axis)
	count   int      // number of all points stored in node and both children
}

// Point represents an abstract point in n-dimensional space
type Point []float64

// Points is a collection of memory pointers to points
type Points []*Point

// Range represents a range in n-dimensional space (a rectangle)
type Range struct {
	min    []float64
	max    []float64
}

// ------------------------------------------

// NewPoint() creates a new point
func NewPoint(vals ...float64) *Point {
	ret := make(Point, 0, len(vals))
	for i := range vals { ret = append(ret, vals[i]) }
	return &ret
}

// NewInfiniteRange() creates a range that contains everything
// Parameter axes specifies dimensionality
func NewInfiniteRange(axes int) Range {
	r := Range{}
	r.min = make([]float64, axes)
	r.max = make([]float64, axes)

	for axis := 0; axis < axes; axis++ {
		r.min[axis] = math.Inf(-1)
		r.max[axis] = math.Inf(1)
	}

	return r
}

// NewKDTree() creates a new kd-tree with given points
func NewKDTree(points Points) *KDNode {
	return insert(points, 0)
}

// ------------------------------------------

// Search() performs range search for [reference - margin, reference + margin],
// starting at given kd-tree node, returning a slice of pointers to matching points
func (node *KDNode) Search(reference *Point, margin []float64) Points {
	// translate reference+margin into range
	query := NewInfiniteRange(len(*reference))
	for axis := 0; axis < len(*reference) && axis < len(margin); axis++ {
		if margin[axis] >= 0 {
			query.min[axis] = (*reference)[axis] - margin[axis]
			query.max[axis] = (*reference)[axis] + margin[axis]
		}
	}

	// prepare info on current's node worldview
	world := NewInfiniteRange(len(*reference))

	// query
	points := make(Points, 0, 32) // NB: pre-allocate for 32 results
	points = query.search(node, world, points)

	return points
}

// Dump() returns a slice of pointers to all points stored in a given kd-tree node,
// and all of it's children (left / right)
func (node *KDNode) Dump() Points {
	points := make(Points, 0, node.count)
	return node.dump(points)
}

// ------------------------------------------

func (node *KDNode) dump(points Points) Points {
	points = append(points, node.point)
	if node.left  != nil { points = node.left.dump(points) }
	if node.right != nil { points = node.right.dump(points) }
	return points
}

func insert(points Points, depth int) *KDNode {
	if len(points) == 0 { return nil }

	node := &KDNode{
		axis: (depth % len(*points[0])),
		count: len(points),
	}

	// find median by sampling on given axis
	median := sample_median(points, node.axis, 250)

	// divide
	points_below := make(Points, 0, len(points)/2)
	points_above := make(Points, 0, len(points)/2)
	for i := range points {
		if points[i] == median {
			continue
		} else if (*points[i])[node.axis] < (*median)[node.axis] {
			points_below = append(points_below, points[i])
		} else {
			points_above = append(points_above, points[i])
		}
	}

	node.point = median
	node.left = insert(points_below, depth + 1)
	node.right = insert(points_above, depth + 1)

	return node
}
