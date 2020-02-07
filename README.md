<p align="center">
    <svg width="120" style="enable-background:new 0 0 512 512;" version="1.1" viewBox="0 0 512 512" xml:space="preserve" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
        <g>
            <path class="st0" fill="#6F707E" d="M477.6,196.4H82.2C82.9,101,160.7,23.5,256.3,23.5c61.2,0,116.8,31.2,148.7,83.4   c3.5,5.8,11.1,7.6,16.8,4c5.7-3.5,7.5-11,4-16.8C389.4,34.7,326-0.8,256.3-0.8c-109,0-197.8,88.4-198.5,197.2H35   c-6.7,0-12.2,5.5-12.2,12.2V499c0,6.7,5.5,12.2,12.2,12.2h442.6c6.7,0,12.2-5.5,12.2-12.2V208.6   C489.8,201.8,484.3,196.4,477.6,196.4z M465.4,486.8H47.2v-266h418.2V486.8z" />
            <path class="st0" fill="#6F707E" d="M244.1,368.1v56.3c0,6.7,5.5,12.2,12.2,12.2s12.2-5.5,12.2-12.2v-56.3   c21.4-5.5,37.3-24.7,37.3-47.8c0-27.3-22.2-49.5-49.5-49.5s-49.5,22.2-49.5,49.5C206.8,343.4,222.7,362.7,244.1,368.1z    M256.3,295.3c13.8,0,25.1,11.3,25.1,25.1c0,13.8-11.3,25.1-25.1,25.1c-13.8,0-25.1-11.2-25.1-25.1   C231.2,306.6,242.4,295.3,256.3,295.3z" />
        </g>
    </svg>
</p>

<p align="center">Utility for obtaining database credentials from <a href="https://github.com/hashicorp/consul">Consul</a> and <a href="https://github.com/hashicorp/vault">Vault</a>.</p>

## Getting started

### Installation

Download and install using go get:

```sh
go get -u github.com/romantomjak/knock-knock
```

or grab a binary from [releases](knock-knock/releases/latest) section!

### Configuration

Configuration file is a [Go template](https://golang.org/pkg/html/template/) with [TOML](https://en.wikipedia.org/wiki/TOML) syntax and by default is searched in `~/.knock-knock.toml`

```toml
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

```
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
