SHELL := /bin/bash
COMPOSE := docker compose -f docker-compose.yaml

.PHONY: help network-up network-down channel-create chaincode-deploy seed-data benchmark stats clean

help:
	@echo "Targets: network-up network-down channel-create chaincode-deploy seed-data benchmark stats clean"

network-up:
	./network/generate.sh
	$(COMPOSE) up -d
	./network/up.sh

network-down:
	./network/down.sh
	$(COMPOSE) down -v

channel-create:
	./network/create-channel.sh

chaincode-deploy:
	./network/deploy-chaincode.sh

seed-data:
	cd data/synthetic && go run generator.go --records 2000 --out ../../benchmark/results/records.json

benchmark:
	./benchmark/run_all.sh

stats:
	python3 benchmark/stats.py --input benchmark/results/raw.csv --output benchmark/results/table6.csv

clean:
	rm -rf crypto-config channel-artifacts organizations benchmark/results/*.csv