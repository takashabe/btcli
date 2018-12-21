[![GoDoc](https://godoc.org/github.com/takashabe/btcli?status.svg)](https://godoc.org/github.com/takashabe/btcli)
[![CircleCI](https://circleci.com/gh/takashabe/btcli.svg?style=shield)](https://circleci.com/gh/takashabe/btcli)
[![Go Report Card](https://goreportcard.com/badge/github.com/takashabe/btcli)](https://goreportcard.com/report/github.com/takashabe/btcli)

# btcli

CLI client for the Bigtable with auto-completion.

![tty](https://user-images.githubusercontent.com/681508/44247754-9ad3ea80-a221-11e8-9172-2cb709e1420a.gif)

## Features

btcli is the cloud bigtable client tool.  Connect to your bigtable instances and any read items with auto-completion.

_write commands are not implemented yet_

### Compare to cbt

`cbt` is an official bigtable client tool

- btcli has auto-completion
- btcli can decode a big-endian values
- btcli has a filter for `value`, `version` and `family`
- A print format that same as the cbt

## Installation

```sh
go get -u github.com/takashabe/btcli/cmd/btcli
```

_Require go1.11 or later and GO111MODULE environemnts_

## Usage

### Connect to bigtable

```sh
btcli -project <GCP_PROJECT_NAME> -instance <BIGTABLE_INSTANCE_ID> -creds <GCP_CREDENTIAL_FILE>
```

_-creds e.g. `~/.config/gcloud/application_default_credentials.json`_

### Subcommand and options

- ls

List tables

```
ls          List tables
```

- count

Count rows in a table

```
count <table>
```

- lookup

Read from a single row

```
lookup <table> <row> [family=<column_family>] [version=<n>]
        version        Read only latest <n> columns
        decode         Decode big-endian value
        decode-columns Decode big-endian value with columns. <column_name:<string|int|float>[,<column_name:...>]
```

- read

Read rows

```
read <table> [start=<row>] [end=<row>] [prefix=<prefix>] [family=<column_family>] [version=<n>]
        start          Start reading at this row
        end            Stop reading before this row
        prefix         Read rows with this prefix
        value          Read rows with has value
        family         Read only columns family with <columns_family>
        version        Read only latest <n> columns
        from           Read cells whose version is newer than or equal to this unixtime
        to             Read cells whose version is older than this unixtime
        decode         Decode big-endian value
        decode-columns Decode big-endian value with columns. <column_name:<string|int|float>[,<column_name:...>]
```

### Environments

| Env | Detail |
| --- | --- |
| BTCLI_DECODE_TYPE | set the default decoding type.<br>values: `string, int, float` |

## Support commands

### Read commands

- [x] ls
- [x] count
- [x] lookup
    - [x] version
    - [x] decode
    - [x] decode-columns
- [x] read
    - [x] start
    - [x] end
    - [x] prefix
    - [x] value
    - [x] family
    - [x] version
    - [x] from
    - [x] to
    - [x] decode
    - [x] decode-columns

### Write commands

- [ ] createfamily
- [ ] createtable
- [ ] deletecolumn
- [ ] deletefamily
- [ ] deleterow
- [ ] deletetable
- [ ] set
- [ ] setgcpolicy

### Others

- [x] help
