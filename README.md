# btcli

CLI client for the Bigtable with auto-completion.

## Features

btcli is the cloud bigtable client tool.  Connect to your bigtable instances and any read items with auto-completion.

_write commands are not implemented yet_

### Compare to cbt

`cbt` is an official bigtable client tool

- btcli has auto-completion
- btcli can decode a big-endian values
- btcli has a filter for the version and family
- A print format that same as the cbt

## Installation

```
go get github.com/takashabe/btcli
```

## Usage

### Connect to bigtable

```
btcli -project <GCP_PROJECT_NAME> -instance <BIGTABLE_INSTANCE_ID> -creds <GCP_CREDENTIAL_FILE>
```

_-creds e.g. `~/.config/gcloud/application_default_credentials.json`_

### Interactive shell

- ls

List tables and column families

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
  family    Read only columns family with <columns_family>
  version   Read only latest <n> columns
```

- read

Read rows

```
read <table> [start=<row>] [end=<row>] [prefix=<prefix>] [family=<column_family>] [version=<n>]
  start     Start reading at this row
  end       Stop reading before this row
  prefix    Read rows with this prefix
  family    Read only columns family with <columns_family>
  version   Read only latest <n> columns
```

## Support commands

### Read commands

- [x] ls
- [x] count
- [x] lookup
    - [x] version
    - [x] family
- [x] read
    - [x] start
    - [x] end
    - [x] prefix
    - [x] version
    - [x] family

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

- [ ] help
