package main

func main() {
	if err := StartServer("127.0.0.1:29998", "/Users/danslimmon/i"); err != nil {
		panic(err.Error())
	}
	if err := StopServer(); err != nil {
		panic(err.Error())
	}
}
