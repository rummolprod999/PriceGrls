package main

func init() {
	CreateEnv()
}
func main() {
	defer SaveStack()
	r := GrlsReader{Url: "https://grls.rosminzdrav.ru/pricelims.aspx"}
	r.reader()
}
