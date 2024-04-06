package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

const (
	PORT = 4000
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type ShortenedLink struct {
	ShortParam string
	ActualLink string
}

func getRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func main(){
	shortenedLinks := []ShortenedLink{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "src/index.html")
	})

	http.HandleFunc("/add-link", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		link := r.Form.Get("link")
		parsedUrl, err := url.Parse(link)
		if err != nil {
			w.Write([]byte("Failed! WRONG URL"))
		}else{
			redirectUrl := parsedUrl.String()
			shortenedUrl := strings.ToLower(getRandomString(5))
			fmt.Printf("%s -> %s\n", shortenedUrl, redirectUrl)
			sl := ShortenedLink{
				ShortParam: shortenedUrl,
				ActualLink: redirectUrl,
			}
			shortenedLinks = append(shortenedLinks, sl)
			w.Write([]byte(fmt.Sprintf("http://%s.localhost:4000/go", shortenedUrl)))
		}
	})

	http.HandleFunc("/go", func(w http.ResponseWriter, r *http.Request) {
		shortLinkString := strings.Split(r.Host, ".")[0]
		redirectUrl := ""
		for _, sl := range shortenedLinks{
			if sl.ShortParam == shortLinkString{
				redirectUrl = sl.ActualLink
			}
		}
		if redirectUrl == ""{
			w.Write([]byte("Wrong short link"))
		}else{
			w.Header().Set("Location", redirectUrl)
			w.WriteHeader(302)
		}
	})

	fmt.Printf("Started listening on server :%d\n", PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil); err!= nil{
		fmt.Printf("Error while listening: %s\n", err)
	}
}