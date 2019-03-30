package main

import (
    "encoding/json"
    "net/http"
    "time"
    "html/template"
	"crypto/rand"
	"crypto/rsa"
	"os"
    "strconv"
    "strings"
)

/* index */
func (a *App) index(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
    var html HtmlStruct
    err := a.DB.One("Id", 1, &html.Company)
    if (err != nil) {
        a.Log.Printf("%s\t%s\t%s\t%s",r.Method,r.RequestURI,err.Error(),time.Since(start))
    }

    err = a.DB.All( &html.Quotes)
    if err != nil { 
        a.Log.Printf("%s\t%s\t%s\t%s",r.Method,r.RequestURI,err.Error(),time.Since(start))
    }
    err = a.DB.All( &html.Pricing)
    if err != nil { 
        a.Log.Printf("%s\t%s\t%s\t%s",r.Method,r.RequestURI,err.Error(),time.Since(start))
    }
    err = a.DB.All( &html.Portfolio)
    if err != nil { 
        a.Log.Printf("%s\t%s\t%s\t%s",r.Method,r.RequestURI,err.Error(),time.Since(start))
    }
    err = a.DB.All( &html.Services)
    if err != nil { 
        a.Log.Printf("%s\t%s\t%s\t%s",r.Method,r.RequestURI,err.Error(),time.Since(start))
    }

    t, errtmpl := template.ParseFiles(a.Conf.Htmlfile)
    if errtmpl != nil { 
        a.Log.Printf("%s\t%s\t%s\t%s",r.Method,r.RequestURI,errtmpl.Error(),time.Since(start))
    }
    errtmpl =   t.Execute(w, html) //step 2
    if errtmpl != nil { 
        a.Log.Printf("%s\t%s\t%s\t%s",r.Method,r.RequestURI,errtmpl.Error(),time.Since(start))
    }
    a.Log.Printf("%s\t%s\t%s",r.Method,r.RequestURI,time.Since(start))
}
/* Company */
func (a *App) updateCompany(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
    var c Company
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&c); err != nil {
	    a.Log.Printf("%s\t%s\t%s\t%s",r.Method,r.RequestURI,err.Error(),time.Since(start))
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()
    if err := c.updateCompany(a.DB); err != nil {
	    a.Log.Printf("%s\t%s\t%s\t%s",r.Method,r.RequestURI,err.Error(),time.Since(start))
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    a.Log.Printf("%s\t%s\t%s",r.Method,r.RequestURI,time.Since(start))
    respondWithJSON(w, http.StatusOK, c)
}

/* Places */
func (a *App) createPlace(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
    var p Place
    var erreur error
    decoder := json.NewDecoder(r.Body)
    if erreur = decoder.Decode(&p); erreur != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()
    p.Hash = getToken(16)
    p.Created_at = time.Now();
    p.Updated_at = time.Now();
    p.Picture,erreur = WriteBase64ToPNG(p.Encoded,a.Conf.Pngdir)
    if erreur != nil {
        a.Log.Printf("Product error %s",erreur.Error())
        respondWithError(w, http.StatusInternalServerError, erreur.Error())
        return
    }
    /* Created by */
    var u User
    erreur = a.DB.One("Hash", p.Referenced_by.Hash, &u)
    if (erreur != nil) {
        a.Log.Printf("%s\t%s\t%s\t%s",r.Method,r.RequestURI,erreur.Error(),time.Since(start))
    }
    /* Located at */
    p.Referenced_by = u
    erreur = p.createPlace(a.DB) 
    if erreur != nil {
        respondWithError(w, http.StatusInternalServerError, erreur.Error())
        return
    }
    /* end subset */
    respondWithJSON(w, http.StatusCreated, p)
    a.Log.Printf("%s\t%s\t%s",r.Method,r.RequestURI,time.Since(start))
}
/* Users */

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
    var u User
    var err error
    decoder := json.NewDecoder(r.Body)
    err = decoder.Decode(&u) 
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()
    
    u.Hash = getToken(16)
    u.Created_at = time.Now();
    u.Updated_at = time.Now();
	u.PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
        a.Log.Printf("%s\t%s\t%s\t%s",r.Method,r.RequestURI,err.Error(),time.Since(start))
		os.Exit(1)
	}

	u.PublicKey = &u.PrivateKey.PublicKey
    a.Log.Printf("%s\t%s\t Public Key (%s)\t%s",r.Method,r.RequestURI,u.PublicKey,time.Since(start))
    err = u.createUser(a.DB)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJSON(w, http.StatusCreated,u)
    a.Log.Printf("%s\t%s\t%s",r.Method,r.RequestURI,time.Since(start))
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
    users,err := getUsers(a.DB)
    if err != nil {
	    a.Log.Printf("%s\t%s\t%s\t%s",r.Method,r.RequestURI,err.Error(),time.Since(start))
        respondWithError(w, http.StatusInternalServerError, err.Error())
    }
    respondWithJSON(w, http.StatusOK, users)
    a.Log.Printf("%s\t%s\t%s",r.Method,r.RequestURI,time.Since(start))
}
func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
    var err error
    hash := r.URL.Query().Get("hash")
    user := User{Hash: hash}
    err = user.getUser(a.DB)
    if err != nil {
	    a.Log.Printf("%s\t%s\t%s\t%s",r.Method,r.RequestURI,err.Error(),time.Since(start))
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return 
    }
    respondWithJSON(w, http.StatusOK, user)
    a.Log.Printf("%s\t%s\t%s",r.Method,r.RequestURI,time.Since(start))
}

func (a *App) getPlaces(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
    places,err := getPlaces(a.DB)
    if err != nil {
	    a.Log.Printf("%s\t%s\t%s\t%s",r.Method,r.RequestURI,err.Error(),time.Since(start))
        respondWithError(w, http.StatusInternalServerError, err.Error())
    }
    respondWithJSON(w, http.StatusOK, places)
    a.Log.Printf("%s\t%s\t%s",r.Method,r.RequestURI,time.Since(start))
}
func (a *App) getPlace(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var t2 = []int{}
    var query QueryPlace
    var erreur error
    query.Address = r.URL.Query().Get("address")
    query.Name = r.URL.Query().Get("name")
    query.Zipcode = r.URL.Query().Get("zipcode")
    query.Town = r.URL.Query().Get("town")
    query.Country = r.URL.Query().Get("country")
    query.Latitude,erreur    = strconv.ParseFloat(r.URL.Query().Get("latitude"),64)
    query.Longitude,erreur   = strconv.ParseFloat(r.URL.Query().Get("longitude"),64)
    categories := strings.Split(r.URL.Query().Get("category"), ",")
    for _, i := range  categories {
        j, erreur := strconv.Atoi(i)
        if erreur != nil {
            panic(erreur)
        }
        t2 = append(t2, j)
    }
    query.Category   = t2
    query.Distance,erreur    = strconv.ParseFloat(r.URL.Query().Get("distance"),64)
    if erreur != nil { 
        //
    }

    a.Log.Printf("Lat:%f\t lng:%f \t category:%s\t distance:%f",query.Latitude,query.Longitude,query.Category,query.Distance)

    var place Place
    var places []Place
    // Use method by parameters
    place.Query = query
    if len(strings.TrimSpace(query.Town)) != 0 {
        places,erreur = place.getPlaceByTown(a.DB)
    }
    if len(strings.TrimSpace(query.Country)) != 0 {
        places,erreur = place.getPlaceByCountry(a.DB)
    }
    if len(strings.TrimSpace(query.Zipcode)) != 0 {
        places,erreur = place.getPlaceByZipcode(a.DB)
    }
    if len(strings.TrimSpace(query.Name)) != 0 {
        places,erreur = place.getPlaceByName(a.DB)
    }
    if query.Latitude != 0 {
        places,erreur = place.getPlaceByLocation(a.DB)
    }

    if erreur != nil {
	    a.Log.Printf("%s\t%s\t%s\t%s",r.Method,r.RequestURI,erreur.Error(),time.Since(start))
        respondWithError(w, http.StatusInternalServerError, erreur.Error())
        return 
    }
    respondWithJSON(w, http.StatusOK, places)
    a.Log.Printf("%s\t%s\t%s",r.Method,r.RequestURI,time.Since(start))
}
