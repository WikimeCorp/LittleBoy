package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	conflib "github.com/JeremyLoy/config"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type configStruct struct {
	VkClientID     string `config:"REACT_APP_CLIENT_ID"`
	VkAuthURI      string `config:"REACT_APP_VKAUTH_URI"`
	VkClientSecret string `config:"REACT_APP_CLIENT_SECRET"`
	RedirectURI    string `config:"REACT_APP_REDIRECT_URI"`
	Port           string `config:"REACT_APP_LITTLE_BOY_PORT"`
}

var config = configStruct{}

func getToken(code string) string {
	url := fmt.Sprintf("%s/access_token?grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s",
		config.VkAuthURI, code, config.RedirectURI, config.VkClientID, config.VkClientSecret)
	req, _ := http.NewRequest("POST", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("ERR: ", err)
		return ""
	}
	defer resp.Body.Close()
	token := struct {
		AccessToken string `json:"access_token"`
	}{}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ERR BODY READING ", err)
	}
	log.Println("vk answer ", string(bytes))
	err = json.Unmarshal(bytes, &token)
	if err != nil {
		log.Println("ERR UNMARSHAL ", err)
	}
	return token.AccessToken
}

func main() {
	configPath := flag.String("configPath", "empty.env", "Path to config file.")

	err := conflib.From(*configPath).FromEnv().To(&config)
	if err != nil {
		log.Fatal("Config read error:", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/access_token", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/access_token enter")
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
		fmt.Println("/access_token exit")
	})

	// router.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
	// 	url := fmt.Sprintf("https://oauth.vk.com/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s", config.VkClientID, config.RedirectURI, "12345")
	// 	http.Redirect(w, r, url, http.StatusSeeOther)
	// })

	// router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte(r.URL.Query().Get("code")))
	// })

	fmt.Println(config)
	fmt.Println("Server start on " + "0.0.0.0:" + config.Port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+config.Port, middleware(handlers.CORS()(router))))
}

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("middleware enter")
		h.ServeHTTP(w, r)
		fmt.Println("middleware end")
	})
}
