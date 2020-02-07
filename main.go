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

	err = tmpl.Execute(consul, nil)
	if err != nil {
		panic(err)
	}

	println(tmpl.Contents())
}
