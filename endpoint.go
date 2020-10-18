package main

func main() {
	logSetup()
	if err := getServer().ListenAndServeTLS("", ""); err != nil {
		panic(err)
	}
}
