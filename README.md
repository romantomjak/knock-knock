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

Configuration by default is read from `~/.knock-knock.conf`. It is based on the [INI](https://en.wikipedia.org/wiki/INI_file) file format which is rendered by Go [template](https://golang.org/pkg/html/template/).

```ini
[myservice]
host = {{ key "services/myservice/db/host" }}
port = 5432
username = {{ with secret "secret/services/myservice/db" }}{{ .Data.username }}{{ end }}
password = {{ with secret "secret/services/myservice/db" }}{{ .Data.password }}{{ end }}
dbname = {{ key "services/myservice/db/database" }}
```

Sections are your service names. They appear on a line by itself, in square
brackets ([ and ]). `key` retrieves values from Consul and likewise `secret`
is for retrieving secrets from Vault.

#### Vault K/V version 2 backend

Here's how to access a versioned secret value:

```hcl
password = {{ with secret "secret/services/myservice/db" }}{{ .Data.data.password }}{{ end }}
```

Note the nested `.Data.data` syntax when referencing the secret value. For more
information about using the K/V v2 backend, see the [Vault Documentation](https://www.vaultproject.io/docs/secrets/kv/kv-v2/).

### Usage

Running the application requires you to specify a service name from the
configuration file:

```sh
export VAULT_AUTH_GITHUB_TOKEN=<mygithubtoken>
export VAULT_ADDR=http://127.0.0.1:8200
export CONSUL_HTTP_ADDR=127.0.0.1:8500
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
