package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dnsge/puzzle-bot"
	"log"
	"time"
)

const (
	defaultDelay = time.Microsecond * 500
)

var (
	userNameFlag  = flag.String("name", "Puzzle Bot", "User display name")
	userColorFlag = flag.String("color", "#00ff00", "User display color")
	roomFlag      = flag.String("room", "", "puzzle.aggie.io room code")
	secretFlag    = flag.String("secret", "", "puzzle.aggie.io room secret")

	edgesFlag    = flag.Bool("edges", false, "Solve edges")
	completeFlag = flag.Bool("complete", false, "Solve puzzle completely")
	regionFlag   = flag.String("region", "", "Solve region from (row,col):(row2,col2)")

	delayFlag           = flag.Duration("delay", defaultDelay, "Delay between actions")
	debugFlag           = flag.Bool("debug", false, "Show debug information")
	overrideVersionFlag = flag.Bool("override-version", false, "Override the server version check")
	forceFlag           = flag.Bool("force", false, "Force operation even if root is not found or locked")
)

func init() {
	flag.Parse()
}

func main() {
	if *roomFlag == "" {
		log.Fatalln("You must specify --room")
	}

	if !*edgesFlag && !*completeFlag && *regionFlag == "" {
		log.Fatalln("You must specify one of --edges, --complete, or --region (row,col):(row:col)")
	}

	var reg region
	if *regionFlag != "" {
		var err error
		reg, err = parseRegion(*regionFlag)
		if err != nil {
			log.Fatalf("Failed to parse region specifier %q: must be in format \"(row,col):(row,col)\"", *regionFlag)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session, err := puzzle.NewSession(ctx, puzzle.Options{
		Room:            *roomFlag,
		Secret:          *secretFlag,
		UserName:        *userNameFlag,
		UserColor:       *userColorFlag,
		Debug:           *debugFlag,
		OverrideVersion: *overrideVersionFlag,
	})

	if err != nil {
		log.Fatalln(err)
	}

	session.OnJoined = func(state *puzzle.SessionState) {
		defer func() {
			// some time to make sure any pending writes probably go through
			time.Sleep(time.Millisecond * 250)
			session.Exit()
			cancel()
		}()

		if *completeFlag {
			solveComplete(session, state)
		} else if *edgesFlag {
			solveEdges(session, state)
		} else if *regionFlag != "" {
			combineRegionCenter(session, state, reg, *forceFlag)
		}
	}

	session.Run(ctx)
}

func parseRegion(regionStr string) (region, error) {
	var reg region
	_, err := fmt.Sscanf(regionStr, "(%d,%d):(%d,%d)", &reg.startRow, &reg.startCol, &reg.endRow, &reg.endCol)
	if err != nil {
		return region{}, err
	}
	return reg, nil
}
