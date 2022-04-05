package routes

import (
	"chatroom/global"
	"chatroom/logic"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
)

func homeHandleFunc(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(global.RootDir + "/template/home.html")
	if err != nil {
		fmt.Fprintf(w, "模版解析出错！")
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Fprintf(w, "模板执行出错！")
		return
	}
}

func userListHandleFunc(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	userList := logic.Broadcasters.GetUserList()
	b, err := json.Marshal(userList)

	if err != nil {
		fmt.Fprint(w, `[]`)
	} else {
		fmt.Fprint(w, string(b))
	}
}
