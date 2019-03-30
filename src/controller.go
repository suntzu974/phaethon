package main

import (
    "github.com/asdine/storm"
    "reflect"
    "strings"
	"github.com/astaxie/beego/validation"
	"errors"
	"fmt"
)

func isEmpty(object interface{}) bool {
    //First check normal definitions of empty
    if object == nil {
        return true
    } else if object == "" {
        return true
    } else if object == false {
        return true
    }

    //Then see if it's a struct
    if reflect.ValueOf(object).Kind() == reflect.Struct {
        // and create an empty copy of the struct object to compare against
        empty := reflect.New(reflect.TypeOf(object)).Elem().Interface()
        if reflect.DeepEqual(object, empty) {
            return true
        }
    }
    return false
}
func (c *Company) updateCompany(db *storm.DB)(error) {
    // Update multiple fields
    err := db.Update(c)
    if (err != nil) {
        return nil
    }
    return nil
}
/* Users */
func (u *User) createUser(db *storm.DB) error {
	valid := validation.Validation{}
	if v := valid.Required(u.Username, "username"); !v.Ok {
		return errors.New(v.Error.Key+":"+v.Error.Message)
	}
	if v := valid.Required(u.Phone, "phone"); !v.Ok {
		return errors.New(v.Error.Key+":"+v.Error.Message)
	}
	if v := valid.Required(u.Email, "email"); !v.Ok {
		return errors.New(v.Error.Key+":"+v.Error.Message)
	}

    err := db.Save(u)
    if (err != nil) {
        return err
    }
    return nil
}
func (u *User) getUser(db *storm.DB) (error) {
	err := db.One("Hash", u.Hash, u)
	if err != nil {
		return err
	}
	return nil
}

func getUsers(db *storm.DB) ([]User,error) {
    var users []User
    err := db.All(&users)
    if err != nil {
        return users,err
    }
    return users,nil
}
/* Place */
func (p *Place) createPlace(db *storm.DB) error {
    err := db.Save(p)
    if (err != nil) {
        return nil
    }
    return nil
}
// -20.961963, 55.323329 Saint Paul
//https://www.goyav.com/api/v2/search.json?latitude=-20.947422&longitude=55.298410&category=1,2,3,4&start_date=1972-02-02&distance=0.100
func (p *Place) getPlaceByLocation(db *storm.DB) ([]Place,error) {
    var places []Place
    var places_founded []Place
    var distance float64
    err := db.All(&places)
    if (err != nil) {
        return places,err
    }
    for _,item := range places {
        distance = Distance(item.Located_at.Latitude,item.Located_at.Longitude,p.Query.Latitude,p.Query.Longitude)
        for _,k := range p.Query.Category {
            if item.Category == k {
                if distance < p.Query.Distance {
                    item.Located_at.Distance = distance
                    fmt.Println("Distance : %s from item %d",distance, item.Located_at.Distance)
                    places_founded = append(places_founded,item)
                }
            }
            
        }
    }
    return places_founded,nil
}
func (p *Place) getPlaceByName(db *storm.DB) ([]Place,error) {
    var places []Place
    var places_founded []Place
    err := db.All(&places)
    if (err != nil) {
        return places,err
    }
    for _,item := range places {
        if strings.Contains(strings.ToLower(item.Name),strings.ToLower(p.Query.Name)) {
            places_founded = append(places_founded,item)
        } 
    }
    return places_founded,nil
}

func (p *Place) getPlaceByAddress(db *storm.DB) ([]Place,error) {
    var places []Place
    var places_founded []Place
    err := db.All(&places)
    if (err != nil) {
        return places,err
    }
    for _,item := range places {
        if strings.Contains(strings.ToLower(item.Located_at.Address),strings.ToLower(p.Query.Address)) {
            places_founded = append(places_founded,item)
        } 
    }
    return places_founded,nil
}
func (p *Place) getPlaceByCountry(db *storm.DB) ([]Place,error) {
    var places []Place
    var places_founded []Place
    err := db.All(&places)
    if (err != nil) {
        return places,err
    }
    for _,item := range places {
        if strings.Contains(strings.ToLower(item.Located_at.Country),strings.ToLower(p.Query.Country)) {
            places_founded = append(places_founded,item)
        } 
    }
    return places_founded,nil
}
func (p *Place) getPlaceByTown(db *storm.DB) ([]Place,error) {
    var places []Place
    var places_founded []Place
    err := db.All(&places)
    if (err != nil) {
        return places,err
    }
    for _,item := range places {
        if strings.Contains(strings.ToLower(item.Located_at.Town),strings.ToLower(p.Query.Town)) {
            places_founded = append(places_founded,item)
        } 
    }
    return places_founded,nil
}
func (p *Place) getPlaceByZipcode(db *storm.DB) ([]Place,error) {
    var places []Place
    var places_founded []Place
    err := db.All(&places)
    if (err != nil) {
        return places,err
    }
    for _,item := range places {
        if strings.Contains(strings.ToLower(item.Located_at.Zipcode),strings.ToLower(p.Query.Zipcode)) {
            places_founded = append(places_founded,item)
        } 
    }
    return places_founded,nil
}

func getPlaces(db *storm.DB) ([]Place,error) {
    var places []Place
    err := db.All(&places)
    if err != nil {
        return places,err
    }
    return places,nil
}
