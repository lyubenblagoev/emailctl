# emailctl

emailctl is a command line interface (CLI) for the [Postfix Rest Server](https://github.com/lyubenblagoev/postfix-rest-server) V1 API. 

emailctl is a work in progress and currently supports only the domain and account API.

```
emailctl is a command line interface (CLI) to the Postfix Rest Server

Usage:
  emailctl [command]

Available Commands:
  account     Account commands
  domain      Domain commands
  help        Help about any command
  version     Prints the version number of emailctl

Flags:
      --config string   config file (default is $HOME/.emailctl.yaml)
  -h, --help            help for emailctl

Use "emailctl [command] --help" for more information about a command.
```

## Installation

### Install from source

If you have a Go environment configured, you can install the development version of `emailctl` from the command line like so: 

```
go get github.com/lyubenblagoev/emailctl
```

## Configuration

By default `emailctl` will load a configuration file from `$HOME/.emailctl.yaml`. 

### Configuration options

* `https` - Boolean setting for whether to use HTTPS for connection to the [Postfix Rest Server](https://github.com/lyubenblagoev/postfix-rest-server).
* `host` - The hostname or IP address of the machine on which the [Postfix Rest Server](https://github.com/lyubenblagoev/postfix-rest-server) is running.
* `port` - The port on which the [Postfix Rest Server](https://github.com/lyubenblagoev/postfix-rest-server) is listening.

Example: 

```yaml
https: false
host: localhost
port: 8080
```

The above values are the defaults. You can omit options that don't change the default values.

## Examples

Below are a few usage examples:

### Domain API

* List all domains on your server:

```
emailctl domain list
```

* Show information for specific domain:

```
emailctl domain show example.com
```

* Add a new domain: 

```
emailctl domain add example.com
```

* Delete a domain: 

```
emailctl domain delete example.com
```

* Disable a domain on your server:

```
emailctl domain disable example.com
```

* Rename a domain

```
emailctl domain rename example.com example.net
```

### Account API

* List accounts for domain:

```
emailctl account list example.com
```

* Show information for account in specific domain:

```
emailctl account show example.com user1
```

* Delete account:

```
emailctl account delete example.com user1
```

* Disable account: 

```
emailctl account disable example.com user1
```

* Change password for account:

```
emailctl account password example.com user1
Password: 
Confirm password: 
```

## More information

To learn more about the features and commands available run

```
emailctl help [command [sub-command]...]
```

or 

```
emailctl [command [sub-command]...] --help
```
