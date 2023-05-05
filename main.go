package main

import (
	"fmt"
	"log"

	"github.com/del-xiong/miniblink"
)

func main() {
	//设置调试模式
	miniblink.SetDebugMode(true)
	//初始化miniblink模块
	err := miniblink.InitBlink()
	if err != nil {
		log.Fatal(err)
	}
	// 启动1366x920普通浏览器
	view := miniblink.NewWebView(false, 350, 850, 1500, 10)
	// 启动1366x920透明浏览器(只有web界面会显示)
	//view := miniblink.NewWebView(true, 1366, 920)
	view.LoadURL("http://10.1.72.55:31090/third/control")
	orgCode := "0101"
	emplCode := "99999"
	var user string = `{"orgCode":"` + orgCode + `","emplCode":"` + emplCode + `"}`
	var setUserJS string = `localStorage.setItem('user','` + user + `')`
	fmt.Println(user)
	view.Eval(setUserJS)
	view.LoadURL("http://10.1.72.55:31090/third/control?inpatId=ZY230427002")
	// 设置窗体标题(会被web页面标题覆盖)
	view.SetWindowTitle("森亿电子病历内涵质控系统")
	view.DisableAutoTitle()

	// 显示窗口
	view.ShowWindow()

	<-make(chan bool)
}
