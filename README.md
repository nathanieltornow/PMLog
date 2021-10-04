# PMLog

**Requirements**
- [Go](https://golang.org/) (version >= 1.15). If you're using Linux, you can use `./scripts/install_go.sh && source /etc/profile`.
- [Make](https://www.gnu.org/software/make/)


### Start a benchmark (single shard) on Cloudlab

1. Start the sequencer on one node
```shell
go run cmd/start_sequencer/start_sequencer.go -root
```

2. Startup the shard
```shell
# for ever replica
go run cmd/server/shard/main.go -order <sequencer-IP> -peers <peer-IPs (comma seperated)>
```

3. Start a client for a replica
   1. Modify benchmark/benchmark.config.yaml, that `enpoint: "<replica-IP>:4000"`
   2. Start the benchmark, which will start at the next minute
   ```shell
   go run benchmark/shared_log/shared_log.go -config benchmark/benchmark.config.yaml 
   ```

