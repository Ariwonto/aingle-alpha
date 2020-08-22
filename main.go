package main

import (
	"github.com/iotaledger/hive.go/node"

	"github.com/Ariwonto/aingle-alpha/pkg/config"
	"github.com/Ariwonto/aingle-alpha/pkg/toolset"
	"github.com/Ariwonto/aingle-alpha/plugins/autopeering"
	"github.com/Ariwonto/aingle-alpha/plugins/cli"
	"github.com/Ariwonto/aingle-alpha/plugins/coordinator"
	"github.com/Ariwonto/aingle-alpha/plugins/dashboard"
	"github.com/Ariwonto/aingle-alpha/plugins/database"
	"github.com/Ariwonto/aingle-alpha/plugins/gossip"
	"github.com/Ariwonto/aingle-alpha/plugins/gracefulshutdown"
	"github.com/Ariwonto/aingle-alpha/plugins/metrics"
	"github.com/Ariwonto/aingle-alpha/plugins/mqtt"
	"github.com/Ariwonto/aingle-alpha/plugins/peering"
	"github.com/Ariwonto/aingle-alpha/plugins/pow"
	"github.com/Ariwonto/aingle-alpha/plugins/profiling"
	"github.com/Ariwonto/aingle-alpha/plugins/prometheus"
	"github.com/Ariwonto/aingle-alpha/plugins/snapshot"
	"github.com/Ariwonto/aingle-alpha/plugins/spammer"
	"github.com/Ariwonto/aingle-alpha/plugins/tangle"
	"github.com/Ariwonto/aingle-alpha/plugins/urts"
	"github.com/Ariwonto/aingle-alpha/plugins/warpsync"
	"github.com/Ariwonto/aingle-alpha/plugins/webapi"
	"github.com/Ariwonto/aingle-alpha/plugins/zmq"
)

func main() {
	cli.HideConfigFlags()
	cli.PrintVersion()
	cli.ParseConfig()
	toolset.HandleTools()
	cli.PrintConfig()

	plugins := []*node.Plugin{
		cli.PLUGIN,
		gracefulshutdown.PLUGIN,
		profiling.PLUGIN,
		database.PLUGIN,
		autopeering.PLUGIN,
		webapi.PLUGIN,
	}

	if !config.NodeConfig.GetBool(config.CfgNetAutopeeringRunAsEntryNode) {
		plugins = append(plugins, []*node.Plugin{
			pow.PLUGIN,
			gossip.PLUGIN,
			tangle.PLUGIN,
			peering.PLUGIN,
			warpsync.PLUGIN,
			urts.PLUGIN,
			metrics.PLUGIN,
			snapshot.PLUGIN,
			dashboard.PLUGIN,
			zmq.PLUGIN,
			mqtt.PLUGIN,
			spammer.PLUGIN,
			coordinator.PLUGIN,
			prometheus.PLUGIN,
		}...)
	}

	node.Run(node.Plugins(plugins...))
}
