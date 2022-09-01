package main

import (
	"flag"
	"fmt"
	"github.com/7yyo/sunflower/prompt"
	"github.com/pingcap/log"
	"go.uber.org/zap/zapcore"
	"keep/components"
	pd "keep/components/pd"
	"keep/components/tidb"
	"keep/util/color"
	"keep/util/printer"
	"os"
)

var mainMenu = []interface{}{
	"cdc",
	"tidb",
	"pd",
}

func main() {
	defer prompt.Close()
	printer.Welcome()
	endpoints := args()
	etcd := components.NewEtcd(endpoints)

	s := components.Server{
		Etcd:            etcd,
		PlacementDriver: pd.NewPlacementDriver(endpoints[0]),
		TiDBCluster:     tidb.NewTiDBCluster(etcd),
	}
	log.SetLevel(zapcore.PanicLevel)

	se := prompt.Select{
		Option: mainMenu,
	}
	_, c, _ := se.Run()

	if err := s.Run(c.(string)); err != nil && err.Error() != "" {
		if prompt.IsBackSpace(err) {
			os.Exit(0)
		}
		fmt.Println(color.Red(fmt.Sprintf("%s", err.Error())))
	}
}

func args() []string {
	var ep string
	flag.StringVar(&ep, "pd", "", "")
	flag.Parse()
	if ep == "" {
		ep = "172.16.5.149:2379"
	}
	return []string{ep}
}
