package app

import (
	"coder/packs/beego"
)

type baseController struct {
	beego.Controller
	moduleName			string
	controllerName		string
	actionName			string
	tplPostfix			string
	errorPageName		string
}

type displayJson map[string]interface{}

func (this *baseController)GetLoginName() string {
	loginName := ""

	loginName, _ = this.GetSession("uname").(string)

	if "" == loginName {
		loginName = this.Ctx.GetCookie("CODER")
		this.SetSession("uname", loginName)
	}

	return loginName
}

func (this *baseController)Prepare() {
	uname := this.GetLoginName()
	if "" == uname {
		this.Redirect("/login", 302)
	}else{
		this.Data["uname"] = uname
	}
	this.Data["uname"] = uname
	controllerName, actionName := this.GetControllerAndAction()
	this.moduleName = "app"
	this.actionName = actionName
	this.tplPostfix = beego.AppConfig.String("TplPostfix")
	this.controllerName = controllerName
	this.errorPageName  = "error"
}

func (this *baseController)display(template string) {
	this.TplName = template + this.tplPostfix
}

func (this *baseController)displayErrorJson() {

}

/**
 * @description: you can use this.Data["errorpage_status"] = int to define the error page's status showing
 * @description: you can use this.Data["errorpage_message"] = string to define the error page's message showing
 */
func (this *baseController)displayErrorPage() {
	this.TplName = this.errorPageName + this.tplPostfix
}