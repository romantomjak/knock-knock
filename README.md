<p align="center">
    <img src="logo.svg" alt="knock-knock" title="knock-knock" />
</p>

<p align="center">Utility for obtaining database credentials from <a href="https://github.com/hashicorp/consul">Consul</a> and <a href="https://github.com/hashicorp/vault">Vault</a>.</p>

## Getting started

### Installation

Download and install using go get:

```sh
go get -u github.com/romantomjak/knock-knock
```

or grab a binary from [releases](https://github.com/romantomjak/knock-knock/releases/latest) section!

### Configuration

Configuration file is a [Go template](https://golang.org/pkg/html/template/) with [TOML](https://en.wikipedia.org/wiki/TOML) syntax and by default is searched in `~/.knock-knock.toml`

```hcl
[myservice]
host = {{ key "services/myservice/db/host" }}
port = 5432
username = {{ with secret "secret/services/myservice/db" }}{{ .Data.username }}{{ end }}
password = {{ with secret "secret/services/myservice/db" }}{{ .Data.password }}{{ end }}
dbname = {{ key "services/myservice/db/database" }}
```

TOML sections are your service names. `key` retrieves values from Consul and
likewise `secret` is for retrieving secrets from Vault.

#### Vault K/V version 2 backend

To access a versioned secret value:

```hcl
password = {{ with secret "secret/services/myservice/db" }}{{ .Data.data.password }}{{ end }}
```

Note the nested `.Data.data` syntax when referencing the secret value. For more
information about using the K/V v2 backend, see the [Vault Documentation](https://www.vaultproject.io/docs/secrets/kv/kv-v2/).

### Usage

Running the application requires to specify a service name from the TOML file:

```sh
$ knock-knock myservice
host = myexampledb.a1b2c3d4wxyz.us-west-2.rds.amazonaws.com
port = 5432
username = awsuser
password = awssecretpassword
dbname = awsdatabase
```

Magic! :sparkles:

## Contributing

You can contribute in many ways and not just by changing the code! If you have
any ideas, just open an issue and tell me what you think.

Contributing code-wise - please fork the repository and submit a pull request.

## Credits

Logo made by Ely Wahib from [http://wahib.me](http://wahib.me)

## License

MIT
