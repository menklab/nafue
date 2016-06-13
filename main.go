package nafue

import (
	"github.com/menkveldj/nafue/config"
)

var C config.Config

func Init(c config.Config) {
	C = c
}