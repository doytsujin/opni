package main

import (
	"fmt"
	"os"

	"github.com/rancher/opni/pkg/test"
	"github.com/spf13/pflag"
	"github.com/ttacon/chalk"
)

func main() {
	var enableGateway, enableEtcd, enableCortex bool
	var remoteGatewayAddress string

	pflag.BoolVar(&enableGateway, "enable-gateway", true, "enable gateway")
	pflag.BoolVar(&enableEtcd, "enable-etcd", true, "enable etcd")
	pflag.BoolVar(&enableCortex, "enable-cortex", true, "enable cortex")
	pflag.StringVar(&remoteGatewayAddress, "remote-gateway-address", "", "remote gateway address")
	pflag.Parse()

	defaultAgentOpts := []test.StartAgentOption{}
	if remoteGatewayAddress != "" {
		fmt.Fprintln(os.Stdout, chalk.Blue.Color("disabling gateway and cortex since remote gateway address is set"))
		enableGateway = false
		enableCortex = false
		defaultAgentOpts = append(defaultAgentOpts, test.WithRemoteGatewayAddress(remoteGatewayAddress))
	}

	test.StartStandaloneTestEnvironment(
		test.WithEnableGateway(enableGateway),
		test.WithEnableEtcd(enableEtcd),
		test.WithEnableCortex(enableCortex),
		test.WithDefaultAgentOpts(defaultAgentOpts...),
	)
}
