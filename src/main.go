// option de compilation pour réduire la taille
// go build -ldflags="-s" -o programm_name ./src/*.go
package main


func main() {
  
  c,err := loadConfig("/usr/bin/configuration.yaml") 
//  c,err := loadConfig("configuration.yaml") 
  if err != nil {
    return
  }
  a := App{}
	a.NewLog(c.Logfile)
  a.Initialize(&c)
  a.Run()
}
