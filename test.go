package main

import (
	"bufio"
	"fmt"
	"os"
	"storyteller/engine"
	"storyteller/game"

	"github.com/logrusorgru/aurora"
)

func mainConsole() {
	context := game.Sample()
	fmt.Println(aurora.Faint(context.Intro()))

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "restart" {
			context = game.Sample()
			fmt.Println(aurora.Faint(context.Intro()))
			continue
		}
		msg := engine.Process(context, scanner.Text())
		fmt.Println(aurora.Faint(msg + "\n"))
	}

}
