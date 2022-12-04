package main

import (
	"github.com/dnsge/puzzle-bot"
	"time"
)

type region struct {
	startRow, startCol, endRow, endCol uint16
}

func waitDelay() {
	time.Sleep(*delayFlag)
}

func combineRegion(session *puzzle.Session, state *puzzle.SessionState, reg region, x, y float32, force bool) {
	rows := state.Room.Rows()
	cols := state.Room.Columns()
	numPieces := rows * cols

	root := rc(reg.startRow, reg.startCol, cols)
	if !force {
		rootGroup := state.Room.GroupByID(root)
		if rootGroup == nil {
			panic("could not find root group")
		} else if rootGroup.Locked {
			panic("root group is locked")
		}
	}

	for row := reg.startRow; row <= reg.endRow; row += 1 {
		for col := reg.startCol; col <= reg.endCol; col += 1 {
			i := rc(row, col, cols)
			if i > numPieces {
				break
			}

			if i != root {
				session.CombinePieces(root, i, x, y)
				waitDelay()
			}
		}
	}
}

func combineRegionTopLeft(session *puzzle.Session, state *puzzle.SessionState, reg region, force bool) {
	pieceWidth := state.Room.PieceWidth()
	pieceHeight := state.Room.PieceHeight()
	regionWidth := float32(reg.endCol - reg.startCol + 1)
	regionHeight := float32(reg.endRow - reg.startRow + 1)

	centerX := (regionWidth * pieceWidth) / 2
	centerY := (regionHeight * pieceHeight) / 2
	combineRegion(session, state, reg, centerX, centerY, force)
}

func combineRegionCenter(session *puzzle.Session, state *puzzle.SessionState, reg region, force bool) {
	centerX := float32(state.Room.BoardWidth) / 2
	centerY := float32(state.Room.BoardHeight) / 2
	combineRegion(session, state, reg, centerX, centerY, force)
}

func solveEdges(session *puzzle.Session, state *puzzle.SessionState) {
	rows := state.Room.Rows()
	cols := state.Room.Columns()

	centerX := float32(state.Room.BoardWidth) / 2
	centerY := float32(state.Room.BoardHeight) / 2

	root := uint16(1)
	for col := uint16(0); col < cols; col++ {
		i := rc(0, col, cols)
		if i != root {
			session.CombinePieces(root, i, centerX, centerY)
			waitDelay()
		}
	}

	for row := uint16(1); row < rows-1; row++ {
		i := rc(row, 0, cols)
		if i != root {
			session.CombinePieces(root, i, centerX, centerY)
			waitDelay()
		}
		i = rc(row, cols-1, cols)
		if i != root {
			session.CombinePieces(root, i, centerX, centerY)
			waitDelay()
		}
	}

	for col := uint16(0); col < cols; col++ {
		i := rc(rows-1, col, cols)
		if i != root {
			session.CombinePieces(root, i, centerX, centerY)
			waitDelay()
		}
	}
}

func solveComplete(session *puzzle.Session, state *puzzle.SessionState) {
	rows := state.Room.Rows()
	cols := state.Room.Columns()
	centerX := float32(state.Room.BoardWidth) / 2
	centerY := float32(state.Room.BoardHeight) / 2

	reg := region{
		0, 0, rows, cols,
	}

	combineRegion(session, state, reg, centerX, centerY, true)
}

func rc(row, col, cols uint16) uint16 {
	return col + row*cols + 1
}
