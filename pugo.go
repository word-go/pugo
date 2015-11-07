package main

import (
	"github.com/codegangsta/cli"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/inconshreveable/log15.v2/ext"
	"pugo/builder"
	"pugo/server"
)

const (
	VERSION  = "1.0"
	VER_DATE = "2015-11-05"

	SRC_DIR = "source"   // source contents dir
	TPL_DIR = "template" // template dir
	DST_DIR = "dest"     // destination dir
)

var (
	app = cli.NewApp()
)

func init() {
	app.Name = "pugo"
	app.Usage = "a beautiful site generator"
	app.Author = "https://github.com/fuxiaohei"
	app.Email = "fuxiaohei@vip.qq.com"
	app.Version = VERSION + "(" + VER_DATE + ")"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Value: "0.0.0.0:9899",
			Usage: "pugo's http server address",
		},
		cli.StringFlag{
			Name:  "theme",
			Value: "default",
			Usage: "pugo's theme to display",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "debug mode flag",
		},
	}
	app.Action = action
}

func action(ctx *cli.Context) {
	lv := log15.LvlInfo
	if ctx.Bool("debug") {
		lv = log15.LvlDebug
	}
	log15.Root().SetHandler(log15.LvlFilterHandler(lv, ext.FatalHandler(log15.StderrHandler)))

	log15.Debug("Dir.Source./" + SRC_DIR)
	log15.Debug("Dir.Template./" + TPL_DIR)
	log15.Debug("Dir.Destination./" + DST_DIR)

	// builder
	b := builder.New(SRC_DIR, TPL_DIR, ctx.String("theme"), ctx.Bool("debug"))
	if b.Error != nil {
		panic(b.Error)
	}

	b.Build(DST_DIR)

	// server
	staticDir := b.Renders().Current().StaticDir()
	static := server.NewStatic()
	static.RootPath = staticDir
	s := server.NewServer(ctx.String("addr"))
	s.Static = static
	s.Helper = server.NewHelper(b, DST_DIR)
	s.Run()
}

func main() {
	app.RunAndExitOnError()
}
