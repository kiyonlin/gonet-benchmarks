package main

var delimiter = []byte("\r\n\r\n")

func main() {
	go runNet()

	go runGnet()

	go runGev()

	select {}
}
