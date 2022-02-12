#! /bin/bash
BIN="./bin/cr"

exec > >(trap "" INT TERM; sed 's/^/server: /') $BIN server --config ./configs/sample.toml &
exec > >(trap "" INT TERM; sed 's/^/stats: /') $BIN update_stats --config ./configs/sample.toml