package main

import "os"

func main() {
	a := NewApp()
	//a.loadSavedImageFiles()
	a.addImageFiles(os.Args[1:])
	a.run()
}
