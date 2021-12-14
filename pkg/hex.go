package cheek

import (
	"math/rand"
)

type sGIOrigin2800 []string

var hexComp sGIOrigin2800 = []string{
	"++?????++ Out of Cheese Error. Redo From Start.",
	"Mr. Jelly! Mr. Jelly! Error at Address Number 6, Treacle Mine Road.",
	"Melon melon melon",
	"+++Wahhhhhhh! Mine!+++",
	"+++ Divide By Cucumber Error. Please Reinstall Universe And Reboot +++",
	"+++Whoops! Here comes the cheese! +++",
	"+++ I Regret This +++",
}

func (hq sGIOrigin2800) Poke() string {
	return hq[rand.Intn(len(hq))]
}
