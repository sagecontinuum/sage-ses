package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"log"
	"fmt"
	"strings"
	"time"
	"encoding/base64"
	"io/ioutil"

)

type goal struct {
	ErrorStruct `json:",inline"`
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name,omitempty"`
	Owner       string            `json:"owner,omitempty"`
	status		string			  `json:"status,omitempty"`
	TimeCreated *time.Time        `json:"time_created,omitempty"`
	TimeUpdated *time.Time        `json:"time_last_updated,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type ErrorStruct struct {
	Error string `json:"error,omitempty"`
}



func getSESstats(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, "SES stats")
	return
}

func listGoals(w http.ResponseWriter, r *http.Request) {
	
	respondJSON(w, http.StatusOK, "Goals!")
	return
}

func createGoal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	goalName, err := getQueryField(r, "name")
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, err.Error(), goalName)
		return
	}
	goalID := "0000000001"
	// newGoal, err := createGoalRequest(username, goalName, status)
	newGoal := goal{ID: goalID, Name: goalName, Owner: username, status:"submitted"}
	respondJSON(w, http.StatusOK, newGoal)
}

func getGoalStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	goalID := vars["goal"]
	log.Println("goalID: ", goalID)
	status := "temporary status"
	respondJSON(w, http.StatusOK, status)
}

func addGoalStatus(w http.ResponseWriter, r *http.Request) {
	
	status, err := getQueryField(r, "status")
	if err != nil {
		respondJSONError(w, http.StatusInternalServerError, err.Error(), status)
		return
	}

	respondJSON(w, http.StatusOK, status)
}

func getGoalmetrics(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	username := vars["username"]

	goalMetrics := username + "'s goals"

	respondJSON(w, http.StatusOK, goalMetrics)
}

func getQueryField(r *http.Request, fieldName string) (value string, err error) {
	query := r.URL.Query()
	dataTypeArray, ok := query[fieldName]

	if !ok {
		err = fmt.Errorf("Please specify data type via query field \"type\"")
		return
	}

	if len(dataTypeArray) == 0 {
		err = fmt.Errorf("Please specify data type via query field \"type\"")
		return
	}

	value = dataTypeArray[0]
	if value == "" {
		err = fmt.Errorf("Please specify data type via query field \"type\"")
		return
	}

	return
}

func authMW(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	vars := mux.Vars(r)
	vars["username"] = ""

	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		next(w, r)
		//respondJSONError(w, http.StatusInternalServerError, "Authorization header is missing")
		return
	}
	log.Printf("authorization: %s", authorization)
	authorizationArray := strings.Split(authorization, " ")
	if len(authorizationArray) != 2 {
		respondJSONError(w, http.StatusInternalServerError, "Authorization field must be of form \"sage <token>\"")
		return
	}

	if strings.ToLower(authorizationArray[0]) != "sage" {
		respondJSONError(w, http.StatusInternalServerError, "Only bearer \"sage\" supported")
		return
	}

	//tokenStr := r.FormValue("token")
	tokenStr := authorizationArray[1]
	log.Printf("tokenStr: %s", tokenStr)

	if disableAuth {
		if strings.HasPrefix(tokenStr, "user:") {
			username := strings.TrimPrefix(tokenStr, "user:")
			vars["username"] = username
		} else {
			vars["username"] = "user-auth-disabled"
		}

		next(w, r)
		return
	}

	url := tokenInfoEndpoint

	log.Printf("url: %s", url)

	payload := strings.NewReader("token=" + tokenStr)
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Print("NewRequest returned: " + err.Error())
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		respondJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	auth := tokenInfoUser + ":" + tokenInfoPassword
	//fmt.Printf("auth: %s", auth)
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+authEncoded)

	req.Header.Add("Accept", "application/json; indent=4")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		respondJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		respondJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if res.StatusCode != 200 {
		fmt.Printf("%s", body)
		//http.Error(w, fmt.Sprintf("token introspection failed (%d) (%s)", res.StatusCode, body), http.StatusInternalServerError)
		respondJSONError(w, http.StatusUnauthorized, fmt.Sprintf("token introspection failed (%d) (%s)", res.StatusCode, body))
		return
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		//fmt.Println(err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		respondJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	val, ok := dat["error"]
	if ok && val != nil {
		fmt.Fprintf(w, val.(string)+"\n")

		//http.Error(w, val.(string), http.StatusInternalServerError)
		respondJSONError(w, http.StatusInternalServerError, val.(string))
		return

	}

	isActiveIf, ok := dat["active"]
	if !ok {
		//http.Error(w, "field active was misssing", http.StatusInternalServerError)
		respondJSONError(w, http.StatusInternalServerError, "field active missing")
		return
	}
	isActive, ok := isActiveIf.(bool)
	if !ok {
		//http.Error(w, "field active is noty a boolean", http.StatusInternalServerError)
		respondJSONError(w, http.StatusInternalServerError, "field active is not a boolean")
		return
	}

	if !isActive {
		//http.Error(w, "token not active", http.StatusInternalServerError)
		respondJSONError(w, http.StatusUnauthorized, "token not active")
		return
	}

	usernameIf, ok := dat["username"]
	if !ok {
		//respondJSONError(w, http.StatusInternalServerError, "username is missing")
		respondJSONError(w, http.StatusInternalServerError, "username is missing")
		return
	}

	username, ok := usernameIf.(string)
	if !ok {
		respondJSONError(w, http.StatusInternalServerError, "username is not string")
		return
	}

	//vars := mux.Vars(r)

	vars["username"] = username

	next(w, r)

}

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	respondJSONError(w, http.StatusInternalServerError, "resource unknown")
	return
}

func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// json.NewEncoder(w).Encode(data)
	s, err := json.MarshalIndent(data, "", "  ")
	if err == nil {
		w.Write(s)
	}
}

func respondJSONError(w http.ResponseWriter, statusCode int, msg string, args ...interface{}) {
	errorStr := fmt.Sprintf(msg, args...)
	log.Printf("Reply to client: %s", errorStr)
	respondJSON(w, statusCode, ErrorStruct{Error: errorStr})
}
