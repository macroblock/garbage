package sdf

import (
	"github.com/macroblock/garbage/sdf/application"
	"github.com/macroblock/garbage/sdf/rake"
)

var app rake.IApplication

// Application -
func Application() rake.IApplication {
	if app != nil {
		return app
	}

	app = application.New()
	return app
}
