package main

func main() {
	tmpl, err := NewTemplate("/Users/romantomjak/.knock-knock.toml")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute()
	if err != nil {
		panic(err)
	}

	println(tmpl.Contents())
}
