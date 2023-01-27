package puzzle

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"strconv"
)

const (
	messagePing      = 11
	messagePickUp    = 1
	messageMove      = 2
	messagePutDown   = 3
	messageCombine   = 6
	messageChallenge = 12
	messageUserID    = 15
)

type Message interface {
	Encode(state *SessionState) (int, []byte, error)
	Name() string
}

type PongMessage struct{}

func (m *PongMessage) Encode(state *SessionState) (int, []byte, error) {
	buf := make([]byte, 3)
	view := DataView(buf)

	view.PutUint8(messagePing, 0)
	view.PutUint16(state.UserID, 1)
	return websocket.BinaryMessage, buf, nil
}

func (m *PongMessage) Name() string {
	return "PongMessage"
}

type JoinMessage struct {
	UserName string
	Color    string
	Room     string
	Secret   *string
}

func (m *JoinMessage) Encode(state *SessionState) (int, []byte, error) {
	data, err := json.Marshal(struct {
		Type   string  `json:"type"`
		Name   string  `json:"name"`
		Color  string  `json:"color"`
		Room   string  `json:"room"`
		Secret *string `json:"secret"`
	}{
		Type:   "user",
		Name:   m.UserName,
		Color:  m.Color,
		Room:   m.Room,
		Secret: m.Secret,
	})

	if err != nil {
		return 0, nil, err
	}

	return websocket.TextMessage, data, nil
}

func (m *JoinMessage) Name() string {
	return "JoinMessage"
}

type PickUpPieceMessage struct {
	ID uint16
	X  float32
	Y  float32
}

func (m *PickUpPieceMessage) Encode(state *SessionState) (int, []byte, error) {
	buf := make([]byte, 13)
	view := DataView(buf)

	view.PutUint8(messagePickUp, 0)
	view.PutUint16(state.UserID, 1)
	view.PutUint16(m.ID, 3)
	view.PutFloat32(m.X, 5)
	view.PutFloat32(m.Y, 9)

	return websocket.BinaryMessage, buf, nil
}

func (m *PickUpPieceMessage) Name() string {
	return "PickUpPieceMessage"
}

type MovePieceMessage struct {
	ID uint16
	X  float32
	Y  float32
}

func (m *MovePieceMessage) Encode(state *SessionState) (int, []byte, error) {
	buf := make([]byte, 13)
	view := DataView(buf)

	view.PutUint8(messageMove, 0)
	view.PutUint16(state.UserID, 1)
	view.PutUint16(m.ID, 3)
	view.PutFloat32(m.X, 5)
	view.PutFloat32(m.Y, 9)

	return websocket.BinaryMessage, buf, nil
}

func (m *MovePieceMessage) Name() string {
	return "MovePieceMessage"
}

type PutDownPieceMessage struct {
	ID uint16
	X  float32
	Y  float32
}

func (m *PutDownPieceMessage) Encode(state *SessionState) (int, []byte, error) {
	buf := make([]byte, 13)
	view := DataView(buf)

	view.PutUint8(messagePutDown, 0)
	view.PutUint16(state.UserID, 1)
	view.PutUint16(m.ID, 3)
	view.PutFloat32(m.X, 5)
	view.PutFloat32(m.Y, 9)

	return websocket.BinaryMessage, buf, nil
}

func (m *PutDownPieceMessage) Name() string {
	return "PutDownPieceMessage"
}

type CombinePiecesMessage struct {
	FirstID  uint16
	SecondID uint16
	X        float32
	Y        float32
}

func (m *CombinePiecesMessage) Encode(state *SessionState) (int, []byte, error) {
	buf := make([]byte, 15)
	view := DataView(buf)

	view.PutUint8(messageCombine, 0)
	view.PutUint16(state.UserID, 1)
	view.PutUint16(m.FirstID, 3)
	view.PutUint16(m.SecondID, 5)
	view.PutFloat32(m.X, 7)
	view.PutFloat32(m.Y, 11)

	return websocket.BinaryMessage, buf, nil
}

func (m *CombinePiecesMessage) Name() string {
	return "CombinePiecesMessage"
}

type ChallengeResponseMessage struct {
	Value uint32
}

func (m *ChallengeResponseMessage) Encode(*SessionState) (int, []byte, error) {
	strRepresentation := strconv.FormatInt(int64(m.Value), 10)

	buf := make([]byte, 3+len(strRepresentation))
	view := DataView(buf)

	view.PutUint8(messageChallenge, 0)
	view.PutRawBytes([]byte(strRepresentation), 3)
	return websocket.BinaryMessage, buf, nil
}

func (m *ChallengeResponseMessage) Name() string {
	return "ChallengeResponseMessage"
}
