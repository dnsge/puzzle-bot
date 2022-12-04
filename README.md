# puzzle-bot

A library for creating bots for the interactive jigsaw puzzle https://puzzle.aggie.io.

## Usage

```bash
go run ./cmd/bot \
    --name "Puzzle Guy" \
    --color "#00ff00" \
    --room "ABC123" \
    [action flag]
```

Where `[action flag]` is one of the following:
```
--edges:                  Complete the edges of the puzzle
--complete:               Complete the entire puzzle
--region (r1,c1):(r2,c2)  Complete the region of the puzzle defined 
                          by the top left and bottom right corners
```