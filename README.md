# btcli

CLI client for the Bigtable with auto-completion.

## Features

btcli is the cloud bigtable client tool.  Connect to your bigtable instances and any read items with auto-completion.

__write commands are not implemented yet__

#### Compare to cbt

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

#### Connect to bigtable

```
btcli -project <GCP_PROJECT_NAME> -instance <BIGTABLE_INSTANCE_ID> -creds <GCP_CREDENTIAL_FILE>
```

__-creds e.g. `~/.config/gcloud/application_default_credentials.json`__

#### Interactive shell

```
ls
count <TABLE_NAME>
lookup <TABLE_NAME> <KEY> [options...]
read <TABLE_NAME> [options...]
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
