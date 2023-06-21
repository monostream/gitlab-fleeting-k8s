package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.com/gitlab-org/fleeting/fleeting/plugin"

	"fleeting_plugin_k8s"
)

var (
	showVersion = flag.Bool("version", false, "Show version information and exit")
)

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Println(fleeting_plugin_k8s.Version.Full())
		os.Exit(0)
	}

	ig := fleeting_plugin_k8s.InstanceGroup{}

	// testing

	//ctx := context.Background()
	//
	//pi, err := ig.Init(ctx, hclog.L(), provider.Settings{})
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//_ = pi
	//
	//iid := "default-default-0"
	//
	//ig.Update(ctx, func(instance string, state provider.State) {
	//	if state == provider.StateRunning {
	//		iid = instance
	//	}
	//})
	//
	//if iid == "" {
	//	log.Println("instance not found")
	//	return
	//}

	//d, err := ig.Decrease(ctx, []string{iid})
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//_ = d

	//ci, err := ig.ConnectInfo(ctx, iid)
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//_ = ci

	//n, err := ig.Increase(ctx, 1)
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//_ = n

	// testing end

	//return

	plugin.Serve(&ig)
}
