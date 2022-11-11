package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/replit/database-go"
	"github.com/xojoc/useragent"
)

var tmpl = template.Must(template.ParseFiles("Templates/index.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func NewToken(w http.ResponseWriter, r *http.Request) {
	bytes := make([]byte, 20)
	rand.Read(bytes)
	database.Set(fmt.Sprintf("%X", bytes), "")
	w.Write([]byte(fmt.Sprintf(`{"token":"%X"}`, bytes)))
}

func GetInfo(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	info, err := database.Get(q["token"][0])
	if err != nil {
		w.Write([]byte(`{"Ok":false, "Error":"Token Does Not Exist"}`))
		return
	} else if info == "" {
		w.Write([]byte(`{"Ok": false, "Error": "The victim has not opened the email yet."}`))
		return
	}

	w.Write([]byte(info))
}

func Image(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header["User-Agent"][0], "GoogleImageProxy") {
		info := GenerateInfoStruct(r)
		j, _ := json.Marshal(info)
		q := r.URL.Query()
		if len(q["token"]) == 0 {
			return
		}else if q["token"][0] == "" {
			return
		}
		key:= q["token"][0]
		database.Set(key, string(j))
	}

	file, _ := os.ReadFile("image.png")
	w.Header().Set("Content-Type", "image/png")
	w.Write(file)
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}
	return ""
}

func GenerateInfoStruct(r *http.Request) (infoPack Info) {
	var ipAddr IP
	ua := r.Header.Get("User-Agent")
	parsedua := useragent.Parse(ua)
	ip := getIP(r)
	resp, _ := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", url.QueryEscape(ip)))
	b, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	json.Unmarshal(b, &ipAddr)

	geoLocation := fmt.Sprintf("\nCountry: %s\nRegion: %s\nCity: %s\nZip: %s\nISP: %s",
		ipAddr.Country,
		ipAddr.RegionName,
		ipAddr.City,
		ipAddr.Zip,
		ipAddr.ISP,
	)

	device := fmt.Sprintf("\nName: %s\nOS: %s\nMobile: %v\nTablet: %v",
		parsedua.Name,
		parsedua.OS,
		parsedua.Mobile,
		parsedua.Tablet,
	)
	infoPack = Info{IP: ipAddr.Query, UserAgent: ua, Device: device, GeoLocation: geoLocation}
	return infoPack
}
