# tiramolla

Utility for file transferring from/to a remote server, possibly at multiple hops distance

## Feature

tiramolla is a cli that allows easy transfer of files from/to remote servers, even when the target server is at one or multiple hops distance.
User is assumed to have access to Gateway and Target server, possibly using different username on each.

```sh
+------+        +---------+       +--------+
| Host | <--->  | Gateway | <---> | Target |
+------+        +---------+       +--------+
```

## Commands

* show - print configured server names or details for a specific server
* copy - download or upload a file

## Usage

### Configuration

The servers that tiramolla can reach are configured in a file named .tiramolla (or .tiramolla.yaml) in the home directory.
An example of .tiramolla.yaml is included in this repository (see [tiramolla.example.yaml](tiramolla.example.yaml)).

### tiramolla command usage

#### general usage
```sh
$ tiramolla --help
Utility for file transferring from/to a remote server, possibly at multiple hops distance

Usage:
  tiramolla [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  copy        download or upload a file
  help        Help about any command
  show        print configured server names or details for a specific server

Flags:
  -h, --help   help for tiramolla

Use "tiramolla [command] --help" for more information about a command.
$
```

#### show
```sh
$ tiramolla show --help
Shows configuration details of a specific server, if serverName is provided.
Shows list of servers if 'servers' argument is passed, list can be matched against a regular expression if the flag is provided.

Usage:
  tiramolla show {servers|serverName} [flags]

Flags:
  -h, --help           help for show
      --match string   expression to match server names
$
```

#### copy
```sh
$ tiramolla copy --help
Downloads or uploads a file to or from a remote server.

Use absolute paths to avoid unexpected behaviour.
Downloading and uploading is set by the mode flag.

Usage:
  tiramolla copy /path/to/file /path/to/dest [flags]

Flags:
  -h, --help            help for copy
      --mode string     down or up
      --server string   target server
$
```

## Install

You have [Go installed](https://go.dev/doc/install).
```sh
git clone https://github.com/kantonop/tiramolla.git ~/tiramolla
cd ~/tiramolla
make
```

tiramolla executable will be installed in `$GOPATH`, `$GOBIN` or `~/go/bin`, make sure to include the appropriate folder to your `$PATH`

## Uninstall

```sh
# clone again the repository if necessary and navigate to the folder
make uninstall
```

## License

Apache License 2.0, see [LICENSE](LICENSE).
