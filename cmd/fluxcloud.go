package main

import (
	"log"
	"strings"

	"github.com/topfreegames/fluxcloud/pkg/apis"
	"github.com/topfreegames/fluxcloud/pkg/config"
	"github.com/topfreegames/fluxcloud/pkg/exporters"
	"github.com/topfreegames/fluxcloud/pkg/formatters"
)

func initExporter(config config.Config) (exporter []exporters.Exporter) {
	exporterType := config.Optional("exporter_type", "slack")

	exporterTypes := strings.Split(exporterType, ",")

	for _, v := range exporterTypes {
		if v == "webhook" {
			webhook, err := exporters.NewWebhook(config)
			if err != nil {
				log.Fatal(err)
			}
			exporter = append(exporter, webhook)
		}

		if v == "matrix" {
			matrix, err := exporters.NewMatrix(config)
			if err != nil {
				log.Fatal(err)
			}
			exporter = append(exporter, matrix)
		}

		if v == "msteams" {
			msteams, err := exporters.NewMSTeams(config)
			if err != nil {
				log.Fatal(err)
			}
			exporter = append(exporter, msteams)
		}

		if v == "slack" {
			slack, err := exporters.NewSlack(config)
			if err != nil {
				log.Fatal(err)
			}
			exporter = append(exporter, slack)
		}
		if v == "datadog" {
			dd, err := exporters.NewDatadog(config)
			if err != nil {
				log.Fatal(err)
			}
			exporter = append(exporter, dd)
		}
	}

	for _, e := range exporter {
		log.Printf("Using %s exporter", e.Name())
	}

	return exporter
}

func main() {
	log.SetFlags(0)

	cfg := &config.DefaultConfig{}

	fmtTemplates := formatters.ReadTemplates()
	fmtCfg := config.NewChain(config.MapConfig(fmtTemplates), cfg)
	formatter, err := formatters.NewDefaultFormatter(fmtCfg)
	if err != nil {
		log.Fatal(err)
	}

	apiConfig := apis.NewAPIConfig(formatter, initExporter(cfg), cfg)

	apis.HandleWebsocket(apiConfig)
	apis.HandleV6(apiConfig)
	log.Fatal(apiConfig.Listen(cfg.Optional("listen_address", ":3031")))
}
