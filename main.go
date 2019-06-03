package main

import "fmt"

func init() {
	CreateEnv()
}
func main() {
	defer SaveStack()
	r := GrlsReader{Url: "https://grls.rosminzdrav.ru/pricelims.aspx", Added: 0}
	Logging("start")
	r.reader()
	Logging(fmt.Sprintf("Added %d elements", r.Added))
	Logging("end")
}
