package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/fglo/particles-rules-of-attraction/pkg/particlelifesim/game"
)

func run() {
	g := game.New()

	if err := ebiten.RunGame(g); err != nil {
		if err == game.Terminated {
			return
		}
		log.Fatal(err)
	}
}

func main() {
	run()
}
