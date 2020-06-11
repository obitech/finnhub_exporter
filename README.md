# [finnhub_exporter](https://github.com/obitech/finnhub_exporter)

The exporter allows for exposing financial data from 
[Finnhub.io](https://finnhub.io) as Prometheus metrics. 

## Installation

### From source

Clone the repository:

```
git clone git@github.com:obitech/finnhub_exporter.git
```

Build it:

```
make all
```

By default the binary will be installed into `./bin`:

```
± ./bin/finnhub_exporter -h
Export financial data from finnhub.io

Usage:
  run [flags]

Flags:
  -h, --help                        help for run
  -l, --log.level string            log level (debug, info, warn, error). Empty or invalid values will fallback to info (default "info")
  -a, --web.listen-address string   The address to listen on for HTTP requests. (default ":9780")
```

## Running it

The exporter follows the idea of the 
[blackbox_exporter](https://github.com/prometheus/blackbox_exporter) or 
[snmp_exporter](https://github.com/prometheus/snmp_exporter) where modules
can be applied to targets. In this case queries, which correspond to Finnhub.io
endpoints, can be applied to stock symbols.

Start it from the command line, while passing your token as an environment 
variable:

```
$ export FINNHUB_API_KEY=abcd
$ ./finnhub_exporter
```

With a request to the `/query` path an with passing additional paramaters, a
query can be made.

Parameter|Description
---|---
`endpoint`|The Finnhub endpoint to query. See below for supported endpoints.
`symbol`|The stock symbol such as `AAPL`, `MSFT`, etc.

### Supported Endpoints

Endpoint|Finnhub.io Endpoint|Description
---|---|---
`companyprofile2`|[/stock/profile2](https://finnhub.io/docs/api#company-profile2)|Get general information of a company.
`quote`|[/quote](https://finnhub.io/docs/api#quote)|Get real-time quote data for US stocks.

## Example

Getting the latest quote from `AAPL`:

```
$ curl "localhost:9780/query?endpoint=quote&symbol=AAPL"
# HELP finnhub_query_duration Returns how long a query to the Finnhub API took to complete in seconds
# TYPE finnhub_query_duration gauge
finnhub_query_duration 0.417014498
# HELP finnhub_query_success Displays whether a query to the Finnhub API was successful
# TYPE finnhub_query_success gauge
finnhub_query_success 1
# HELP finnhub_quote_current
# TYPE finnhub_quote_current gauge
finnhub_quote_current{symbol="AAPL"} 353.5299987792969
# HELP finnhub_quote_high
# TYPE finnhub_quote_high gauge
finnhub_quote_high{symbol="AAPL"} 353.8399963378906
# HELP finnhub_quote_low
# TYPE finnhub_quote_low gauge
finnhub_quote_low{symbol="AAPL"} 346.0899963378906
# HELP finnhub_quote_open
# TYPE finnhub_quote_open gauge
finnhub_quote_open{symbol="AAPL"} 347.8999938964844
# HELP finnhub_quote_prev_close
# TYPE finnhub_quote_prev_close gauge
finnhub_quote_prev_close{symbol="AAPL"} 343.989990234375
```

## License

[MIT](https://choosealicense.com/licenses/mit/)