package main

func main() {
	tmpl, err := NewTemplate("/Users/romantomjak/.knock-knock.toml")
	if err != nil {
		panic(err)
	}

	consul, err := NewConsulClient()
	if err != nil {
		panic(err)
	}

	vault, err := NewVaultClient()
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(consul, vault)
	if err != nil {
		panic(err)
	}

	println(tmpl.Contents())
}
