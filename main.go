package main

//go build  -tags bdebug -ldflags "-H windowsgui" -o SynyiQcLocalproxy.exe  打包成后台应用，带devtools
//双击exe启动
//http://localhost:3000 查看是否运行
//http://localhost:3000/inpat?inpatId=ZY230417002&emplCode=99999&hospCode=0101 打开/刷新/切换患者
//http://localhost:3000/close 关闭质控医生端

import (
	"fmt"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/del-xiong/miniblink"
	"github.com/go-chi/chi"
	"github.com/go-vgo/robotgo"
	"github.com/joho/godotenv"
)

var (
	view *miniblink.WebView

	isDebug          bool
	qcWinTitle       string
	qcAddress        string
	qcLocalproxyPort string
	qcWidth          int
	qcHeight         int
	qcTop            int
	qcLeft           int

	emplCode string
	hospCode string
)

func main() {
	//读取配置
	readConfig()

	//初始化 blink
	initBlink()

	//启动server
	startServe()
}

func readConfig() {
	// 获取屏幕大小
	width, _ := robotgo.GetScreenSize()

	if e := godotenv.Load(); e != nil {
		log.Fatal(e)
	}

	isDebug, _ = strconv.ParseBool(os.Getenv("is_debug"))
	qcWinTitle = os.Getenv("qc_win_title")

	qcAddress = os.Getenv("qc_address")
	qcLocalproxyPort = os.Getenv("qc_localproxy_port")

	qcWidth, _ = strconv.Atoi(os.Getenv("qc_width"))
	qcHeight, _ = strconv.Atoi(os.Getenv("qc_height"))
	qcTop, _ = strconv.Atoi(os.Getenv("qc_top"))
	var right, _ = strconv.Atoi(os.Getenv("qc_right"))
	qcLeft = width - right - qcWidth
}

func initBlink() {
	//设置调试模式
	miniblink.SetDebugMode(isDebug)
	//初始化miniblink模块
	err := miniblink.InitBlink()
	if err != nil {
		log.Fatal(err)
	}
}

func startServe() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this is local proxy serve for synyi quality control system"))
	})
	r.Get("/inpat", inpat)
	r.Get("/close", closeView)
	if isDebug {
		r.Use(middleware.Logger)
	}
	log.Fatal(http.ListenAndServe(":"+qcLocalproxyPort, r))
}

func inpat(rw http.ResponseWriter, r *http.Request) {
	inpatId := r.URL.Query().Get("inpatId")
	urlEmplCode := r.URL.Query().Get("emplCode")
	urlHospCode := r.URL.Query().Get("hospCode")
	isSetUser := (len(emplCode) > 0 && len(hospCode) > 0) && (urlEmplCode != emplCode || urlHospCode != hospCode)
	emplCode = urlEmplCode
	hospCode = urlHospCode

	if view != nil {
		if view.IsDestroy {
			go newView(inpatId)
		} else {
			go refreshView(isSetUser, inpatId)
		}
	} else {
		go newView(inpatId)
	}

	rw.Write([]byte("success"))
}

func closeView(rw http.ResponseWriter, _ *http.Request) {
	if view != nil {
		view.DestroyWindow()
	}

	rw.Write([]byte("success"))
}

func newView(inpatId string) {
	// 启动浏览器
	view = miniblink.NewWebView(false, qcWidth, qcHeight, qcLeft, qcTop)
	// 启动浏览器(只有web界面会显示)
	//view := miniblink.NewWebView(false, qc_width, qc_height, qc_left, qc_top)
	view.LoadURL(qcAddress)

	setUser()

	view.LoadURL(qcAddress + "?inpatId=" + inpatId)

	// 显示窗口
	view.ShowWindow()

	//debug
	view.HideDockIcon()
	setWinTitle()
	view.MostTop(true)
}

func refreshView(isSetUser bool, inpatId string) {
	if isSetUser {
		setUser()
	}

	view.LoadURL(qcAddress + "?inpatId=" + inpatId)

	setWinTitle()

	view.RestoreWindow()
}

func setUser() {
	var user string = `{"orgCode":"` + hospCode + `","emplCode":"` + emplCode + `"}`
	var setUserJS string = `localStorage.setItem('user','` + user + `')`
	_, err := view.Eval(setUserJS)
	if err != nil {
		fmt.Println("setUser err:", err)
	}
}

func setWinTitle() {
	// 设置窗体标题(会被web页面标题覆盖)
	view.SetWindowTitle(qcWinTitle)
	view.DisableAutoTitle()
}
