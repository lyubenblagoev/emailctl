# emailctl

emailctl is a command line interface (CLI) for the [Postfix Rest Server][1] V1 API. 

emailctl currently implements all [Postfix Rest Server][1] APIs (domains, accounts, aliases, automatic sender and recipient BCC).

```
emailctl is a command line interface (CLI) to the Postfix Rest Server

Usage:
  emailctl [command]

Available Commands:
  account       Account commands
  alias         Alias commands
  auth          Authentication commands
  domain        Domain commands
  help          Help about any command
  recipient-bcc recipient-bcc commands
  sender-bcc    sender-bcc commands
  version       Prints the version number of emailctl

Flags:
      --config string   config file (default is $HOME/.emailctl.yaml)
  -h, --help            help for emailctl

Use "emailctl [command] --help" for more information about a command.

```

## Installation

### Install from source

If you have a Go environment configured, you can install the latest development version of `emailctl` directly from my GitHub account from the command line like so: 

```
go get github.com/lyubenblagoev/emailctl/...
```

which will install the executable in your GOPATH, or you can build it with the following command and move the executable somewhere under PATH:

```
go build -o some-path/emailctl ./cmd/emailctl
```

## Configuration

By default `emailctl` will load a configuration file from `$HOME/.emailctl.yaml`. Make sure to set proper permissions to this file, so no one else can read it, as `emailctl` uses it to store the access and refresh tokens.

### Configuration options

* `https` - Boolean setting for whether to use HTTPS for connection to the [Postfix Rest Server][1].
* `host` - The hostname or IP address of the machine on which the [Postfix Rest Server][1] is running.
* `port` - The port on which the [Postfix Rest Server][1] is listening.

Example: 

```yaml
https: false
host: localhost
port: 8080
```

The above values are the defaults. You can omit options that don't change the default values.

## Examples

Below are a few usage examples:

### Authentication

* Log in

```
emailctl auth login admin@example.com
```

* Log out

```
emailctl auth logout
```

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

[1]: https://github.com/lyubenblagoev/postfix-rest-server "Postfix Rest Server"
