# validator-monitor

The service is used to calculate the effectiveness of operators.

## Usage

### Step 1. Installation

```console
git clone https://github.com/stakewise/validator-monitor
cd ./validator-monitor
```

### Step 2. Configuration

Specify the address of the Graph (`GRAPH_NODE_URL`), ETH2 nodes (`BEACON_NODE_URL`) and comma separated list of operator wallets `OPERATOR_WALLETS` as environment variables


### Step 3. Start the service

```console
go run .
```

or

```
go build .
./validator-monitor
```

### Operator Effectiveness Environment Settings

| Variable                       | Description                                                                      | Required | Default                                                                 |
|--------------------------------|----------------------------------------------------------------------------------|----------|-------------------------------------------------------------------------|
| BIND_ADDRESS             | Which port/interface to listen on                        | No       | 0.0.0.0:0000                                      |
| GRAPH_NODE_URL           | Graph Node Endpoint                   | Yes       | https://api.thegraph.com/subgraphs/name/stakewise/stakewise-mainnet |
| BEACON_NODE_URL            | ETH2 Node Endpoint                 | Yes       | -                                                                       |
| OPERATOR_WALLETS         | List of operators wallets              | Yes       | 0x6a5b7d2,0x6a5b7d3,0x6a5b7d4 ...                                                                       |                                                              |
