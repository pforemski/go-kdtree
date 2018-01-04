/*
 * go-kdtree/range: code implementing range search
 *
 * Copyright (C) 2018 Pawel Foremski <pjf@foremski.pl>
 * Licensed to you under GNU GPL v3
 */

package kdtree

// range_search() returns all node's points that match given query range
// world is an information on what the node and all children can contain
func (query *Range) search(node *KDNode, world Range, out Points) Points {
	// check if query intersects with world
	switch query.has_range(world) {
	case 0: // no intersection
		return out

	case 1: // some intersection
		// check if node's point is contained
		if query.has_point(node.point) {
			out = append(out, node.point)
		}

		// check in the left child
		if node.left != nil {
			out = query.search(node.left, world.limit_left(node), out)
		}

		// check in the right child
		if node.right != nil {
			out = query.search(node.right, world.limit_right(node), out)
		}

	case 2: // fully contained
		out = node.dump(out)
	}

	return out
}

// has_point() returns true if given point is within the query
func (query *Range) has_point(point *Point) bool {
	for axis := 0; axis < len(point.V); axis++ {
		if point.V[axis] < query.min[axis] { return false }
		if point.V[axis] > query.max[axis] { return false }
	}
	return true
}

// has_range() checks intersection between query and given world range
// 0 means no intersection, 1 means partial, 2 means query fully contains r
func (query *Range) has_range(world Range) int {
	ret := 2

	for axis := 0; axis < len(query.min); axis++ {
		if query.min[axis] >= world.max[axis] { return 0 } // no intersection possible
		if query.max[axis] <  world.min[axis] { return 0 } // no intersection possible

		if ret == 2 {
			if query.max[axis] < world.max[axis] { ret = 1 } // partial intersection
			if query.min[axis] > world.min[axis] { ret = 1 } // partial intersection
		}
	}

	return ret
}

// limit_left() updates worldview to reflect the left child of given node
func (world *Range) limit_left(parent *KDNode) Range {
	r := Range{}
	r.min = world.min // bottom limits don't change

	// upper limits: leave all but parent.axis
	r.max = make([]float64, len(world.max))
	copy(r.max, world.max)
	r.max[parent.axis] = parent.point.V[parent.axis] // all elements < median

	return r
}

// limit_right() updates worldview to reflect the right child of given node
func (world *Range) limit_right(parent *KDNode) Range {
	r := Range{}
	r.max = world.max // upper limits don't change

	// botton limits: leave all but parent.axis
	r.min = make([]float64, len(world.min))
	copy(r.min, world.min)
	r.min[parent.axis] = parent.point.V[parent.axis] // all elements >= median

	return r
}
