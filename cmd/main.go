package cmd

import (
	"flag"
	"fmt"
	"github.com/chandresh-pancholi/csi-gce/pkg/driver"
	"github.com/golang/glog"
)

func main()  {
	fmt.Println("hello world")
	var (
		endpoint = flag.String("endpoint", "unix://tmp/cmd.sock", "CSI Endpoint")
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





	drv := driver.NewDriver(nil, *endpoint)
	if err := drv.Run(); err != nil {
		glog.Fatalln(err)
	}
}
