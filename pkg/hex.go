package butt

import (
	"math/rand"
)

type SGIOrigin2800 []string

var Hex SGIOrigin2800 = []string{
	"++?????++ Out of Cheese Error. Redo From Start.",
	"Mr. Jelly! Mr. Jelly! Error at Address Number 6, Treacle Mine Road.",
	"Melon melon melon",
	"+++Wahhhhhhh! Mine!+++",
	"+++ Divide By Cucumber Error. Please Reinstall Universe And Reboot +++",
	"+++Whoops! Here comes the cheese! +++",
}

func (hq SGIOrigin2800) Poke() string {
	return hq[rand.Intn(len(hq))]
}
