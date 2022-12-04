package puzzle

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID    uint16 `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Room struct {
	BoardHeight   int     `json:"boardHeight"`
	BoardWidth    int     `json:"boardWidth"`
	Groups        []Group `json:"groups"`
	HidePreview   bool    `json:"hidePreview"`
	Jitter        float32 `json:"jitter"`
	Name          string  `json:"name"`
	NoLockUnlock  bool    `json:"noLockUnlock"`
	NoMultiSelect bool    `json:"noMultiSelect"`
	Pieces        uint16  `json:"pieces"`
	Rotation      bool    `json:"rotation"`
	Seed          int     `json:"seed"`
	StartTime     int     `json:"startTime"`
	TabSize       float32 `json:"tabSize"`
	Sets          []Set   `json:"sets"`
}

func (r *Room) GroupByID(id uint16) *Group {
	for i := range r.Groups {
		g := &r.Groups[i]
		if g.ID == id {
			return g
		}
	}
	return nil
}

func (r *Room) PieceWidth() float32 {
	if len(r.Sets) != 1 {
		panic(fmt.Sprintf("unexpected set length %d", len(r.Sets)))
	}
	return r.Sets[0].PieceWidth()
}

func (r *Room) PieceHeight() float32 {
	if len(r.Sets) != 1 {
		panic(fmt.Sprintf("unexpected set length %d", len(r.Sets)))
	}
	return r.Sets[0].PieceHeight()
}

func (r *Room) Rows() uint16 {
	if len(r.Sets) != 1 {
		panic(fmt.Sprintf("unexpected set length %d", len(r.Sets)))
	}
	return r.Sets[0].Rows
}

func (r *Room) Columns() uint16 {
	if len(r.Sets) != 1 {
		panic(fmt.Sprintf("unexpected set length %d", len(r.Sets)))
	}
	return r.Sets[0].Columns
}

type Set struct {
	Rows    uint16  `json:"rows"`
	Columns uint16  `json:"cols"`
	Width   float32 `json:"width"`
	Height  float32 `json:"height"`
}

func (s *Set) PieceWidth() float32 {
	return s.Width / float32(s.Columns)
}

func (s *Set) PieceHeight() float32 {
	return s.Height / float32(s.Rows)
}

type Group struct {
	ID      uint16   `json:"id"`
	IDs     []uint16 `json:"ids"`
	Indices []uint16 `json:"indices"`
	Locked  bool     `json:"locked"`
	X       float32  `json:"x"`
	Y       float32  `json:"y"`
}

type BaseJSONMessage struct {
	Type    string          `json:"type"`
	Users   json.RawMessage `json:"Users"`
	Room    json.RawMessage `json:"Room"`
	MeID    uint16          `json:"id"`
	Version string          `json:"version"`
}
