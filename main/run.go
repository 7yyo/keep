package main

import (
	"flag"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/pingcap/log"
	"go.uber.org/zap/zapcore"
	"keep/components"
	pd "keep/components/pd"
	"keep/components/tidb"
	"keep/util/color"
)

func main() {

	fmt.Println("\n" +
		" __   ___  _______   _______    _______   \n" +
		"|/\"| /  \")/\"     \"| /\"     \"|  |   __ \"\\  \n" +
		"(: |/   /(: ______)(: ______)  (. |__) :) \n" +
		"|    __/  \\/    |   \\/    |    |:  ____/  \n" +
		"(// _  \\  // ___)_  // ___)_   (|  /      \n" +
		"|: | \\  \\(:      \"|(:      \"| /|__/ \\     \n" +
		"(__|  \\__)\\_______) \\_______)(_______)    \n" +
		"                                          ")

	defer fmt.Println("\ngoodbye!")
	endpoints := args()

	etcd := components.NewEtcd(endpoints)
	s := components.Server{
		Etcd:            etcd,
		PlacementDriver: pd.NewPlacementDriver(endpoints[0]),
		TiDBCluster:     tidb.NewTiDBCluster(etcd),
	}
	log.SetLevel(zapcore.PanicLevel)

	prompt := promptui.Select{
		Label: "what can i do for you",
		Items: []string{
			"cdc",
			"tidb",
			"pd",
		},
	}
	_, c, _ := prompt.Run()

	if err := s.Run(c); err != nil && err.Error() != "" {
		fmt.Println(color.Red(fmt.Sprintf("%s", err.Error())))
	}
}

func args() []string {
	var endpoint string
	flag.StringVar(&endpoint, "pd", "", "")
	flag.Parse()
	if endpoint == "" {
		endpoint = "172.16.5.133:2379"
	}
	endpoints := make([]string, 0)
	endpoints = append(endpoints, endpoint)
	return endpoints
}
