package main

import (
	"encoding/json"
	"fmt"
	"github.com/josemrobles/robification-go"
	"io/ioutil"
	"net/http"
)

func main() {

	http.HandleFunc("/", indexAction)
	err := http.ListenAndServe(":8008", nil)
	if err != nil {
		panic(err)
	}

}

func indexAction(res http.ResponseWriter, req *http.Request) {

	config := getConfig("config.json")

	decoder := json.NewDecoder(req.Body)
	var p gitPayload
	err := decoder.Decode(&p)

	if err != nil {
		println(err)
	} else {
		if p.Action == "opened" {

			message := fmt.Sprint(p.Pull_Request.Html_Url, " To: ", p.Pull_Request.Head.Repo.Name, " by: ", p.Pull_Request.User.Login)
			post := robification.NewFdChat(string(config.Fd_Token), string(message))
			err = robification.Send(post)
			if err != nil {
				println(err)
			}

			res.Header().Set("Content-Type", "application/json")
			b, _ := json.Marshal(p)
			fmt.Fprintf(res, string(b))
		}
	}
}

func getConfig(jsonFile string) (config *JSONConfigData) {
	config = &JSONConfigData{}
	J, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(J), &config)
	return config
}