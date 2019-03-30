package main

import (
    "log"
    "time"
    "github.com/gorilla/mux"
    "github.com/asdine/storm"
    "crypto/rsa"
)

type App struct {
    Router  *mux.Router
    DB  	*storm.DB
    Conf	*Configuration
	Log     *log.Logger
}
type PhaethonError struct {
	Key string
	Message string
}

type Apikey struct {
    Id int `storm:"id,increment"`
    Token string
}
type QueryPlace struct {
    Latitude float64
    Longitude float64
    Category []int
    Distance float64
    Address string
    Zipcode string
    Town string
    Country string
    Name string
}
type HtmlStruct struct {
    Company Company
    Pricing []Price
    Quotes []Quote
    Portfolio []Portfolio
    Services []Service

}
type Company struct {
    Id int `storm:"id"`
    Header string 
    Logo string
    Name string
    Introduction string
    Description string
    Content string
    Address string
    Phone string
    Email string
    Latitude float64
    Longitude float64
    BgColor string
    Mission Menu
    Vision Menu
    About Menu
    Service Menu
    Portfolio Menu
    Pricing Menu
    Contact Menu
}
type Quote struct {
    Id int `storm:"id,increment"`
    Hash string    
    Sentence string
    Author string
    Created_at time.Time
}
type Portfolio struct {
    Id int `storm:"id,increment"`
    Hash string
    Key string
    Description string
    Content string
    Picture string
    Created_at time.Time
}
type Service struct {
    Id int `storm:"id,increment"`
    Hash string
    Key string
    Description string
    Content string
    Picture string
    Created_at time.Time
}

type Price struct {
    Id int `storm:"id,increment"`
    Hash string
    Header string
    Price float64
    Unit string
    Periodic string
    Items []Item
    Created_at time.Time
}
type Menu struct {
    Key string
    Header string
    Description string
    Content string
}
type Item struct {
    Header string
    Key string
    Value float64
    Unit string
}
type Configuration struct {
    Htmldir string `yaml:"htmldir"`
    Htmlfile string `yaml:"htmlfile"`
    AppDir string `yaml:"appdir"`
    Dbname string `yaml:"dbname"`
    Logfile string `yaml:"logfile"`
    Pngdir string `yaml:"pngdir"`
    Addr string `yaml:"addr"`
    Key string `yaml:"key"`
    Crt string `yaml:"crt"`
    Pem string `yaml:"pem"`
    Host string `yaml:"host"`
    User string `yaml:"user"`
    Password string `yaml:"password"`

    
}

type Location struct {
    Address string
    Zipcode string
    Town string
    Country string `storm:"index"`
    Latitude float64
    Longitude float64
    Distance float64
}

type User struct {
    Id int `storm:"id,increment"`
    Username string `storm:"index"`
    Phone string `storm:"unique"`
    Email string `storm:"unique"`
    Hash string `storm:"index"`
    Model string `storm:"index"`
    Version string `storm:"index"`
    Application string `storm:"index"`
    PublicKey *rsa.PublicKey
    PrivateKey *rsa.PrivateKey
    Created_at time.Time 
    Updated_at time.Time
    Located_at Location
    Isonline bool
}

type Place struct {
    Id  int `storm:"id,increment"`
    Hash  string  `storm:"index"`
    Name  string  `storm:"index"`
    Referenced_by  User  `storm:"index"`
    Located_at Location  `storm:"index"`
    Category int  `storm:"index"`
    Picture string 
    Encoded string
    Evaluations []Evaluation
    Comments []Comment
    Created_at time.Time 
    Updated_at time.Time 
    Query QueryPlace
}
type Evaluation struct {
    Id  int `storm:"id,increment"`
    Hash  string  `storm:"index"`
    Place Place `storm:"index"`
    Evaluated_by User `storm:"index"`
    Price int
    Service int
    Quality int
}
type Comment struct {
    Id  int `storm:"id,increment"`
    Hash  string  `storm:"index"`
    Place Place `storm:"index"`
    Evaluated_by User `storm:"index"`
    Content string
    Picture string
}
