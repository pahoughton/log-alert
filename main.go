/* 2019-01-21 (cc) <paul4hough@gmail.com>
   log-alert application entry point
*/
package main

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/pahoughton/log-alert/config"
	"github.com/pahoughton/log-alert/watcher"

	"gopkg.in/alecthomas/kingpin.v2"

)

type CommandArgs struct {
	ConfigFn	*string
	Debug		*bool
}

func main() {

	app := kingpin.New(path.Base(os.Args[0]),
		"log scanner").
			Version("0.1.1")

	args := CommandArgs{
		ConfigFn: app.Flag("config-fn","config filename").
			Default("log-alert.yml").String(),
		Debug:		app.Flag("debug","debug output to stdout").
			Default("true").Bool(),
	}

	kingpin.MustParse(app.Parse(os.Args[1:]))

	fmt.Println(os.Args[0]," starting")
	fmt.Println("loading ",*args.ConfigFn)

	cfg, err := config.LoadFile(*args.ConfigFn)
	if err != nil {
		panic(err)
	}

	watcher := new.Watcher(cfg,*args.Debug)

	fmt.Println("INFO: ",os.Args[0]," listening on ",cfg.Global.ListenAddr)
	amgrcnt := 0
	for _, amgr := range cfg.Amgrs {
		for _, targ := range amgr.SConfigs.Targets {
			url := fmt.Sprintf("%s://%s",amgr.Scheme, targ)
			fmt.Println("INFO: sending alerts to ",url)
			amgrcnt += 1
		}
	}
	if amgrcnt < 1 {
		fmt.Println("FATAL: no alertmanagers configured")
		os.Exit(2)
	}

	http.Handle("/metrics",promh.Handler())

	fmt.Println("FATAL: ",
		http.ListenAndServe(cfg.Global.ListenAddr,nil).
			Error())
	os.Exit(1)
}
