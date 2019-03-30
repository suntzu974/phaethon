package main
import (
    "os"
    "log"
    "net/http"
	"github.com/asdine/storm"
    "github.com/gorilla/mux"
    "time"
)

func (a *App) Initialize(conf *Configuration) {
    start := time.Now()
    a.Log.Printf("%s\t%s\t%s\t%s","Init","Init","Initialize",time.Since(start))

    db, err := storm.Open(conf.Dbname) 
    a.DB = db
    a.Conf = conf
    if err != nil {
        a.Log.Fatal(err)
    }

    var c Company
    c.Id=1
    c.Header ="Goyav Company"
    c.Logo = "GOYAV Ltd" 
    c.Name =  "Goyav Corp" 
    c.Introduction =  "Cette structure a pour vocation de vous proposer les meilleurs solutions "
    c.Description = " Cette société est une société de service très répandue dans le monde"
    c.Content= "Concept de la frugalité , une dimension écologique"
    c.Address = "6 Place Hermann Melville, 97420 LE PORT"
    c.Phone = "069209997"
    c.Email = "jeannick.grodnin@goayv.com"
    c.Latitude = -20.947254  
    c.Longitude = 55.298622 
    c.BgColor = "#f4511e"

    c.About.Header = "A propos "
    c.Service.Header = "Services"
    c.Portfolio.Header = "Portfolio"
    c.Portfolio.Description = "Ce que nous avons crées"
    c.Pricing.Header = "Tarif"
    c.Contact.Header = "Contact "
    c.Contact.Description = "Contactez nous . Délai de réponse 24 Heures."

    c.Mission.Key = "Mission"
    c.Mission.Header = "Notre Mission"
    c.Mission.Content = "Notre mission est primordiale"

    c.Vision.Key = "Value"
    c.Vision.Header = "Nos valeurs"
    c.Vision.Content = "Nos valeurs sont fondamentales"
    
    errdb := db.One("Id", 1, &c)
    if (errdb != nil) {
        errdb = db.Save(&c)
        if (errdb != nil) {
            a.Log.Printf("%s\t%s\t%s\t%s","Save DB","Save DB First","First time",time.Since(start))
        }
    }
    a.Log.Printf("%s\t%s\t%s\t%s","Save DB","Save DB First","First time",time.Since(start))
    a.Router = mux.NewRouter()
    a.initializeRoutes()
}

func (a *App) Run() {
    a.Log.Fatal(http.ListenAndServeTLS(a.Conf.Addr, a.Conf.Crt, a.Conf.Key , a.Router))
}

func (a *App) initializeRoutes() {
    a.Router.HandleFunc("/", a.index)

    a.Router.HandleFunc("/company", a.updateCompany).Methods("PUT")
    a.Router.HandleFunc("/user", a.createUser).Methods("POST")
    a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
    a.Router.HandleFunc("/user", a.getUser).Methods("GET")
    a.Router.HandleFunc("/places", a.getPlaces).Methods("GET")
    a.Router.HandleFunc("/place", a.getPlace).Methods("GET")
    a.Router.HandleFunc("/place", a.createPlace).Methods("POST")
//    a.Router.HandleFunc("/place/{key}", a.getPlace).Methods("GET")
    
}
func (a *App) NewLog(logpath string) {
    file, err := os.Create(logpath)	
    //file, err := os.OpenFile(logpath,os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	a.Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
}
