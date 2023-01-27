package puzzle

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

const (
	challengeVersion    = "x"
	expectedVersion     = "1.8.5"
	websocketGatewayURL = "wss://puzzle.aggie.io/ws"
	userAgent           = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"

	fixedA = uint32(123)
	fixedB = uint32(765)
)

var (
	pongMessage    = &PongMessage{}
	challengeRegex = regexp.MustCompile("^return a\\*(?P<A>\\d+)-!0\\+b\\*(?P<B>\\d+)-(?P<C>\\d+)$")
)

type Options struct {
	Room            string
	Secret          string
	UserName        string
	UserColor       string
	Debug           bool
	OverrideVersion bool
}

type Session struct {
	conn      *websocket.Conn
	state     SessionState
	writeChan chan Message

	ctx  context.Context
	done context.CancelFunc

	options Options

	OnJoined func(state *SessionState)
}

type SessionState struct {
	Joined bool
	UserID uint16
	Room   *Room
	Users  []User
}

// NewSession creates a new puzzle session with the given options.
func NewSession(ctx context.Context, opts Options) (*Session, error) {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, websocketGatewayURL, http.Header{
		"User-Agent": []string{userAgent},
	})

	if err != nil {
		return nil, err
	}

	return &Session{
		conn: conn,
		state: SessionState{
			Joined: false,
			UserID: 0,
			Room:   nil,
			Users:  nil,
		},
		options:   opts,
		writeChan: make(chan Message),
	}, nil
}

// Run begins the read/write loop of the puzzle session. This function will
// block until the passed context finishes.
func (s *Session) Run(ctx context.Context) {
	ctx2, cancel := context.WithCancel(ctx)
	s.ctx = ctx2
	s.done = cancel

	go s.readLoop()
	go s.writeLoop(s.ctx)

	<-s.ctx.Done()
}

func (s *Session) sendJoin() {
	var secret *string = nil
	if s.options.Secret != "" {
		secret = &s.options.Secret
	}

	s.writeMessage(&JoinMessage{
		UserName: s.options.UserName,
		Color:    s.options.UserColor,
		Room:     s.options.Room,
		Secret:   secret,
	})
}

// Exit closes the puzzle session.
func (s *Session) Exit() {
	_ = s.conn.Close()
	s.done()
}

func (s *Session) readLoop() {
	for {
		messageType, data, err := s.conn.ReadMessage()
		if err != nil {
			log.Printf("read: %v\n", err)
			return
		}

		// Pass message to correct handler
		if messageType == websocket.TextMessage {
			s.debug("recv %v\n", string(data))
			err = s.processJSON(data)
		} else if messageType == websocket.BinaryMessage {
			s.debug("recv %X\n", data)
			err = s.processBinary(data)
		}

		if err != nil {
			log.Printf("handle: %v\n", err)
			s.Exit()
			return
		}
	}
}

func (s *Session) writeLoop(ctx context.Context) {
	for {
		select {
		case m := <-s.writeChan:
			err := s.writeMessageDirect(m)
			if err != nil {
				log.Printf("write: %v\n", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *Session) processJSON(data []byte) error {
	var msg BaseJSONMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	switch msg.Type {
	case "version":
		if msg.Version == challengeVersion { // ignore in favor of challenge
			return nil
		}

		log.Printf("Joined Puzzle server version %s\n", msg.Version)
		if msg.Version != expectedVersion && !s.options.OverrideVersion {
			e := puzzleVersionError(msg.Version)
			return &e
		}
	case "me":
		s.state.UserID = msg.MeID
		log.Printf("Got my user ID of %d\n", msg.MeID)
	case "room":
		var room Room
		if err := json.Unmarshal(msg.Room, &room); err != nil {
			return err
		}
		s.state.Room = &room
	case "users":
		var users []User
		if err := json.Unmarshal(msg.Users, &users); err != nil {
			return err
		}
		s.state.Users = users
	}

	if !s.state.Joined {
		if s.state.UserID != 0 && s.state.Room != nil && s.state.Users != nil {
			s.state.Joined = true
			if s.OnJoined != nil {
				go s.OnJoined(&s.state)
			}
		}
	}

	return nil
}

func (s *Session) processBinary(data []byte) error {
	view := DataView(data)
	messageType := view.Uint8(0)

	var err error = nil
	switch messageType {
	case messagePing:
		s.writeMessage(pongMessage)
	case messageChallenge:
		err = s.handleChallenge(view)
	case messageUserID:
		s.handleUserID(view)
	}

	if err != nil {
		return fmt.Errorf("handle message type %d: %w", messageType, err)
	} else {
		return nil
	}
}

func (s *Session) handleChallenge(view DataView) error {
	pos := 3
	version, consumed := view.ReadString(pos)
	pos += consumed

	if version != expectedVersion && !s.options.OverrideVersion {
		e := puzzleVersionError(version)
		return &e
	}

	challenge, _ := view.ReadString(pos)
	matches := challengeRegex.FindStringSubmatch(challenge)
	if len(matches) != 4 {
		return fmt.Errorf("unexpected challenge form %q", challenge)
	}

	aCoeff, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return err
	}
	bCoeff, err := strconv.ParseUint(matches[2], 10, 32)
	if err != nil {
		return err
	}
	c, err := strconv.ParseUint(matches[3], 10, 32)
	if err != nil {
		return err
	}

	challengeResult := executeChallenge(uint32(aCoeff), uint32(bCoeff), uint32(c), fixedA, fixedB)
	challengeResponse := &ChallengeResponseMessage{Value: challengeResult}
	s.writeMessage(challengeResponse)
	s.sendJoin()

	return nil
}

func executeChallenge(aCoeff, bCoeff, c, a, b uint32) uint32 {
	return aCoeff*a + bCoeff*b - c - 1
}

func (s *Session) handleUserID(view DataView) {
	id := view.Uint16(1)
	s.state.UserID = id
	log.Printf("Got my user ID of %d\n", id)
}

func (s *Session) writeMessage(message Message) {
	s.writeChan <- message
}

func (s *Session) writeMessageDirect(message Message) error {
	s.debug("Sending %q\n", message.Name())
	kind, data, err := message.Encode(&s.state)
	if err != nil {
		return fmt.Errorf("encode %q: %w", message.Name(), err)
	}
	if kind == websocket.TextMessage {
		s.debug("send %v\n", string(data))
	} else {
		s.debug("send %X\n", data)
	}
	return s.conn.WriteMessage(kind, data)
}

func (s *Session) debug(format string, v ...any) {
	if s.options.Debug {
		log.Printf(format, v)
	}
}

type puzzleVersionError string

func (e *puzzleVersionError) Error() string {
	return fmt.Sprintf("Unexpected puzzle server version %s (wanted %s). Enable puzzle.Options.OverrideVersion to ignore.", *e, expectedVersion)
}
