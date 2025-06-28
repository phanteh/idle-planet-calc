package main

func main() {
	data := loadData()
	app := NewApp(data)
	app.Run()
}
