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
--edges                   Complete the edges of the puzzle
--complete                Complete the entire puzzle
--region (r1,c1):(r2,c2)  Complete the region of the puzzle defined 
                          by the top left and bottom right corners
```

## All Options
```
Join Options
    --name string           User display name (default "Puzzle Bot")
    --color <hex string>    User display color (default "#00ff00")
    --room string           puzzle.aggie.io room code
    --secret string         puzzle.aggie.io room secret
    
Actions
    --edges                 Solve edges
    --complete              Solve puzzle completely
    --region string         Solve region from (row,col):(row2,col2)

Misc
    --delay <duration>      Delay between actions (default 500µs)
    --force                 Force operation even if root is not found or locked
    --override-version      Override the server version check
    --debug                 Show debug information
```