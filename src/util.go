package main

import (
    "encoding/json"
    "crypto/rand"
    "encoding/base64"
    "encoding/base32"
    "fmt"
    "os"
    "bytes"
    "image/png"
    "bufio"
    "time"
   "gopkg.in/yaml.v2"
   "io/ioutil"
   "math"
       	"net/http"
)

func loadConfig(filename string) (Configuration, error) {  
    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return Configuration{}, err
    }

    var c Configuration
    err = yaml.Unmarshal(bytes, &c)
    if err != nil {
        return Configuration{}, err
    }

    return c, nil
}

func getTS() (ts string) {
    t:= time.Now()
    ts = t.Format( time.RFC3339)
    return
}
func getToken(length int) string {
    randomBytes := make([]byte, 64)
    _, err := rand.Read(randomBytes)
    if err != nil {
        panic(err)
    }
    return base64.StdEncoding.EncodeToString(randomBytes)[:length]
}
func getFile(length int) string {
    randomBytes := make([]byte, 32)
    _, err := rand.Read(randomBytes)
    if err != nil {
        panic(err)
    }
    return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}

func WriteBase64ToPNG(b64 string, pathb64 string) (string,error) {

	token := getFile(16)
	var path = pathb64+token+".png"
	fo, err := os.Create(path)
	if err != nil {
		return path ,err
	}
	defer fo.Close()
	unbased, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
    	panic("Cannot decode b64")
	}

	r := bytes.NewReader(unbased)
	im, err := png.Decode(r)
	if err != nil {
    	panic("Bad png")
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
    	panic("Cannot open file")
	}
	//var opt jpeg.Options
    //opt.Quality = 80
	png.Encode(f, im)
	fmt.Println("==> done creating file", path)
	return path ,nil
}
func PNG2B64(source string) (string) {
	imgFile, err := os.Open(source) // a QR code image
    if err != nil {
    	imgFile, err = os.Open("/home/suntzu974/public/default.png")
    }

    defer imgFile.Close()

    // create a new buffer base on file size
    fInfo, _ := imgFile.Stat()
    var size int64 = fInfo.Size()
    buf := make([]byte, size)

    // read file content into buffer
    fReader := bufio.NewReader(imgFile)
    fReader.Read(buf)

    // convert the buffer bytes to base64 string - use buf.Bytes() for new image
    return base64.StdEncoding.EncodeToString(buf)

}

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
  // must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)
	return 2 * r * math.Asin(math.Sqrt(h))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}
