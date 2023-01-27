package app

import (
	"encoding/json"
	"net/http"

	"k8s.io/klog"
)

// jsonPrint prints output in json format
func jsonPrint(w http.ResponseWriter, code int, res any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		klog.Errorf("cannot encode response: %v", err)
	}
}

// jsonPrintError error log to server console and prints out error in json format
func jsonPrintError(w http.ResponseWriter, code int, errMsj, consoleMsj string) {
	klog.Errorf(consoleMsj+": %v", errMsj)
	jsonPrint(w, code, map[string]string{"error": errMsj})
}
