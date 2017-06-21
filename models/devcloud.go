package models

import (
	//"net/url"
	"os/exec"
	"regexp"
	"strings"
	//"time"

	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/logs"
	//"github.com/OwnLocal/goes"
)

type Result struct {
	ID      string
	_index  string
	_type   string
	result  string
	current string
	//next        string
	PACKAGE_URL string
	SONAR_URL   string
}

var resp map[string]string

/*
func getClient() (conn *goes.Client) {
	conn = goes.NewClient(beego.AppConfig.String("ES_HOST"), beego.AppConfig.String("ES_PORT"))
	return
}

func writeEs(opt, id string, ret interface{}) (string, string, string, error) {
	conn := getClient()
	var esid interface{}
	if id != "" {
		esid = id
	}
	d := goes.Document{
		Index: beego.AppConfig.String("Index"),
		Type:  opt,
		ID:    esid,
		Fields: map[string]interface{}{
			"msg":    ret,
			"intime": time.Now().Format("2006-01-02 15:04:05"),
		},
	}
	extraArgs := make(url.Values, 0)
	response, err := conn.Index(d, extraArgs)
	if err != nil {
		return "", "", "", err
	}
	return response.Index, response.Type, response.ID, nil
}
*/
func GetCreateResult(project_name, svn_url string) (resp map[string]string) {
	resp = make(map[string]string)
	var ret string
	out, err := exec.Command("/bin/bash", beego.AppConfig.String("create"), project_name, svn_url).Output()
	if err != nil {
		ret = string(out) + "|" + err.Error()
		resp["current"] = "N"
		//resp["next"] = "N"
	} else {
		ret = project_name + " Create sucess"
		resp["current"] = "Y"
		//resp["next"] = "Y"
	}
	//write es
	Index, Type, Id, err := writeEs("create", "", ret)
	if err != nil {
		resp["result"] = ret
		resp["es_result"] = err.Error()
		return
	}
	resp["PACKAGE_URL"] = ""
	resp["SONAR_URL"] = ""
	resp["_index"] = Index
	resp["_type"] = Type
	resp["ID"] = Id
	resp["result"] = ret

	return
}

func GetCheckOutResult(project_name string) (resp map[string]string) {
	resp = make(map[string]string)
	var ret string
	out, err := exec.Command("/bin/bash", beego.AppConfig.String("checkout"), project_name).Output()
	if err != nil {
		ret = string(out) + "|" + err.Error()
		resp["current"] = "N"
		//resp["next"] = "N"
	} else {
		ret = string(out)
		resp["current"] = "Y"
		//resp["next"] = "Y"
	}

	//write es
	Index, Type, Id, err := writeEs("checkout", "", ret)
	if err != nil {
		resp["result"] = ret
		resp["es_result"] = err.Error()
		return
	}
	resp["PACKAGE_URL"] = ""
	resp["SONAR_URL"] = ""
	resp["_index"] = Index
	resp["_type"] = Type
	resp["ID"] = Id
	resp["result"] = ret

	return
}

func GetCodeCheckResult(project_name string) (resp map[string]string) {
	resp = make(map[string]string)
	var ret string
	out, err := exec.Command("/bin/bash", beego.AppConfig.String("codecheck"), project_name).Output()
	if err != nil {
		ret = string(out) + "|" + err.Error()
		resp["current"] = "N"
		resp["SONAR_URL"] = ""
		//resp["next"] = "N"
	} else {
		reg := regexp.MustCompile("SONAR_URL:.*[a-z]")
		SONAR_URL := strings.Replace(reg.FindAllString(string(out), -1)[0], "SONAR_URL:", "", -1)
		resp["SONAR_URL"] = SONAR_URL
		ret = string(out)
		resp["current"] = "Y"
		//resp["next"] = "Y"
	}

	//write es
	Index, Type, Id, err := writeEs("codecheck", "", ret)
	if err != nil {
		resp["result"] = ret
		resp["es_result"] = err.Error()
		return
	}
	resp["PACKAGE_URL"] = ""
	resp["_index"] = Index
	resp["_type"] = Type
	resp["ID"] = Id
	resp["result"] = ret
	return
}

func GetCompileResult(project_name, jdk_version string) (resp map[string]string) {
	resp = make(map[string]string)
	var ret string
	out, err := exec.Command("/bin/bash", beego.AppConfig.String("compile"), project_name, jdk_version).Output()
	if err != nil {
		ret = string(out) + "|" + err.Error()
		resp["current"] = "N"
		//resp["next"] = "N"
	} else {
		//ret = "Compile sucess "+project_name
		ret = string(out)
		resp["current"] = "Y"
		//resp["next"] = "Y"
	}
	//write es
	Index, Type, Id, err := writeEs("compile", "", ret)
	if err != nil {
		resp["result"] = ret
		resp["es_result"] = err.Error()
		return
	}
	resp["PACKAGE_URL"] = ""
	resp["SONAR_URL"] = ""
	resp["_index"] = Index
	resp["_type"] = Type
	resp["ID"] = Id
	resp["result"] = ret
	return
}

func GetPackResult(project_name, version, isE string) (resp map[string]string) {
	resp = make(map[string]string)
	var ret string
	out, err := exec.Command("/bin/bash", beego.AppConfig.String("pack"), project_name, version, isE).Output()
	if err != nil {
		ret = string(out) + "|" + err.Error()
		resp["current"] = "N"
		resp["PACKAGE_URL"] = ""
		//resp["next"] = "N"
	} else {
		isok := strings.Contains(string(out), "ERROR")
		if !isok {
			ret = string(out)
			reg := regexp.MustCompile("PACKAGE_URL:.*[0-9]")
			PACKAGE_URL := strings.Replace(reg.FindAllString(string(out), -1)[0], "PACKAGE_URL:", "", -1)
			resp["PACKAGE_URL"] = PACKAGE_URL
			resp["current"] = "Y"
			//resp["next"] = "Y"
		} else {
			ret = string(out)
			resp["current"] = "N"
			resp["PACKAGE_URL"] = ""
			//resp["next"] = "N"
		}
	}

	//write es
	Index, Type, Id, err := writeEs("pack", "", ret)
	if err != nil {
		resp["result"] = ret
		resp["es_result"] = err.Error()
		return
	}
	resp["SONAR_URL"] = ""
	resp["_index"] = Index
	resp["_type"] = Type
	resp["ID"] = Id
	resp["result"] = ret
	return
}
