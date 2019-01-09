package main

import (
	"flag"
	"fmt"
	"github.com/chandresh-pancholi/csi-gce/pkg/driver"
	"log"
)

func main()  {
	fmt.Println("hello world")
	var (
		endpoint = flag.String("endpoint", "unix://tmp/cmd.sock", "CSI Endpoint")
		nodeID          = flag.String("nodeid", "", "node id")
		//version  = flag.Bool("version", false, "Print the version and exit.")
	)
	flag.Parse()

	//if *version {
	//	info, err := driver.GetVersionJSON()
	//	if err != nil {
	//		glog.Fatalln(err)
	//	}
	//	fmt.Println(info)
	//	os.Exit(0)
	//}




	log.Printf("Endpoint: %s ", *endpoint)
	log.Printf("NodeId: %s ", *nodeID)

	drv := driver.NewDriver(nil, *endpoint, *nodeID)
	if err := drv.Run(); err != nil {
		log.Fatal(err)
	}

	//driver, err := driver.NewGCE(nil, *nodeID, *endpoint)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//driver.Run()
	//os.Exit(0)

}
