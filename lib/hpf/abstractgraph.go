// Package abstractgraph constructs and manages the abstract node space corresponding to a TileMap object.
// This package will be used ast he underlying topology for hiearchical A* searching.
package abstractgraph

import (
	// "math"

	// rtscpb "github.com/cripplet/rts-pathing/lib/proto/constants_go_proto"
	// rtsspb "github.com/cripplet/rts-pathing/lib/proto/structs_go_proto"

	"github.com/cripplet/rts-pathing/lib/hpf/abstractedgemap"
	"github.com/cripplet/rts-pathing/lib/hpf/abstractnodemap"
	// "github.com/cripplet/rts-pathing/lib/hpf/astar"
	// "github.com/cripplet/rts-pathing/lib/hpf/cluster"
	// "github.com/cripplet/rts-pathing/lib/hpf/entrance"
	// "github.com/cripplet/rts-pathing/lib/hpf/tile"
	// "github.com/cripplet/rts-pathing/lib/hpf/utils"
	// "github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	notImplemented = status.Error(
		codes.Unimplemented, "function not implemented")
)

// AbstractGraph contains the necessary state information to make an efficient
// path planning call on very large maps via hierarchical A* search, as
// described in Botea 2004.
type AbstractGraph struct {
	// Level is the maximum hierarchy of AbstractNodes in this graph;
	// this is a positive, non-zero integer. The 0th level here loosely
	// refers to the underlying base map.
	Level int32

	// NodeMap contains a Level: AbstractNodeMap dict representing the
	// AbstractNodes per Level. As per AbstractGraph.ClusterMap, there
	// is a corresponding AbstractNodeMap object per level. Nodes within
	// a specific AbstractNodeMap may move between levels, and may be
	// deleted when the underlying terrain changes.
	NodeMap map[int32]*abstractnodemap.AbstractNodeMap

	// EdgeMap contains a Level: AbstractEdgeMap dict representing the
	// AbstractEdges per Level. Edges may move between levels and may
	// be deleted when the underlying terrain changes.
	EdgeMap map[int32]abstractedgemap.AbstractEdgeMap
}
