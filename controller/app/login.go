package app

import (
	"coder/packs/beego"
	"coder/packs/beego/config"
	"time"
)

var UserConfig config.Configer

type LoginController struct {
	beego.Controller
}

func init() {
	UserConfig, _ = config.NewConfig("ini", "conf/user.conf")
}

func (this *LoginController)Login() {
	this.Data["flag"] = false
	if "" != this.GetString("flag") {
		this.Data["flag"] = true
	}
	this.TplName = "app/login" + beego.AppConfig.String("TplPostfix")
}

func (this *LoginController)DoLogin() {
	ret := displayJson{"status":1,"msg":"success","data":""}

	reqUname := this.GetString("uname")
	reqUpass := this.GetString("upass")
	act := this.GetString("act")

	if "" == reqUname || "" == reqUpass {
		ret["status"] = 0
		ret["msg"] = "账号或密码不能为空"
	}else{
		pass := UserConfig.String(reqUname)
		if "login" == act {
			if "" != pass {
				reqUpass := this.GetString("upass")
				if reqUpass == pass {
					this.SetSession("uname", reqUname)
					if "true" == this.GetString("remember") {
						this.Ctx.SetCookie("CODER", reqUname, time.Second*86400)
					}
				}else{
					ret["status"] = 0
					ret["msg"] = "密码错误"
				}
			}else{
				ret["status"] = 0
				ret["msg"] = "账号不存在"
			}
		}else if "reg" == act {
			if "" == pass {
				UserConfig.Set(reqUname, reqUpass)
				UserConfig.SaveConfigFile("conf/user.conf")
				this.SetSession("uname", reqUname)
				if "true" == this.GetString("remember") {
					this.Ctx.SetCookie("CODER", reqUname, time.Second*86400)
				}
			}else{
				ret["status"] = 0
				ret["msg"] = "账号已存在"
			}
		}else{
			ret["status"] = 0
			ret["msg"] = "非法操作"
		}
	}



	this.Data["json"] = ret
	this.ServeJSON()
}

func (this *LoginController)Register() {
	ret := displayJson{"status":1,"msg":"success","data":""}

	uname := this.GetString("uname")
	upass := this.GetString("upass")

	if "" == uname || "" == upass {
		ret["status"] = 0
		ret["msg"] = "账户和密码不能为空"
	}else{
		UserConfig.Set(uname, upass)
	}

	this.Data["json"] = ret
	this.ServeJSON()
}

func (this *LoginController)Logout() {
	this.DelSession("uname")
	this.Ctx.SetCookie("CODER", "", -1)
	this.Redirect("/login", 302)
}