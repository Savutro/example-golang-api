package controllers

import (
	"encoding/json"
	"net/http"
)

type AdminInfo struct {
	Info string `json:"info"`
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	adminInfo := AdminInfo{
		Info: "Hello Admin",
	}
	res, _ := json.Marshal(adminInfo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
