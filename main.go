package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	conflib "github.com/JeremyLoy/config"
)

type configStruct struct {
	VkClientID     string `config:"REACT_APP_CLIENT_ID"`
	VkAuthURI      string `config:"REACT_APP_VKAUTH_URI"`
	VkClientSecret string `config:"REACT_APP_CLIENT_SECRET"`
	RedirectURI    string `config:"REACT_APP_REDIRECT_URI"`
	Port           string `config:"LITTLE_BOY_PORT"`
}

var config = configStruct{}

func getToken(code string) string {
	url := fmt.Sprintf("%s/access_token?grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s",
		config.VkAuthURI, code, config.RedirectURI, config.VkClientID, config.VkClientSecret)
	req, _ := http.NewRequest("POST", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	token := struct {
		AccessToken string `json:"access_token"`
	}{}
	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &token)
	return token.AccessToken
}

func main() {
	configPath := flag.String("configPath", ".env", "Path to config file.")

	err := conflib.From(*configPath).FromEnv().To(&config)
	if err != nil {
		log.Fatal("Config read error:", err)
	}

	http.HandleFunc("/access_token", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "code is empty", http.StatusBadRequest)
			return
		}

		token := getToken(code)
		if token == "" {
			http.Error(w, "invalide code", http.StatusBadRequest)
			return
		}
		ans := struct {
			Token string `json:"token"`
		}{token}
		w.Header().Set("Content-Type", "application/json")
		ansBytes, _ := json.Marshal(ans)
		w.Write(ansBytes)
	})

	fmt.Println(config)
	fmt.Println("Server start")
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
