package puzzle

func (s *Session) PickUpPiece(id uint16, x float32, y float32) {
	s.writeMessage(&PickUpPieceMessage{
		ID: id,
		X:  x,
		Y:  y,
	})
}

func (s *Session) MovePiece(id uint16, x float32, y float32) {
	s.writeMessage(&MovePieceMessage{
		ID: id,
		X:  x,
		Y:  y,
	})
}

func (s *Session) PutDownPiece(id uint16, x float32, y float32) {
	s.writeMessage(&PutDownPieceMessage{
		ID: id,
		X:  x,
		Y:  y,
	})
}

func (s *Session) CombinePieces(id1 uint16, id2 uint16, x float32, y float32) {
	s.writeMessage(&CombinePiecesMessage{
		FirstID:  id1,
		SecondID: id2,
		X:        x,
		Y:        y,
	})
}
