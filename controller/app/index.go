package app

import (
	"coder/packs/beego"
	"coder/utils"
	"fmt"
	"os"
	"os/exec"
)

type IndexController struct {
	baseController
}

func GetGraphSetting(reqSet map[string]string) map[string]string {
	ret := map[string]string{
		"shape":	"record",
		"style":	"filled",
		"color":	"black",
		"fillcolor":"white",
		"bgcolor":	"white",
	}

	var val string
	var exist bool

	val, exist = reqSet["shape"]
	if exist {
		ret["shape"] = val
	}

	val, exist = reqSet["style"]
	if exist {
		ret["style"] = val
	}

	val, exist = reqSet["color"]
	if exist {
		ret["color"] = val
	}

	val, exist = reqSet["fillcolor"]
	if exist {
		ret["fillcolor"] = val
	}

	val, exist = reqSet["bgcolor"]
	if exist {
		ret["bgcolor"] = val
	}

	return ret
}

func (this *IndexController)Index() {
	this.Data["title"]	=	"Coder"
	this.Data["shapes"]	= 	utils.GetGraphShapeSetting()
	this.Data["colors"]	= 	utils.GetGraphColorSetting()
	this.display("app/index")
}

func (this *IndexController)Review() {
	this.displayErrorPage()
}

func (this *IndexController)Analyse() {
	ret := displayJson{"status":1, "msg":"success", "data":""}
	graphSetting := map[string]string{
		"bgcolor": this.GetString("bgcolor"),
		"color": this.GetString("color"),
		"style": this.GetString("style"),
		"fillcolor": this.GetString("fillcolor"),
		"shape": this.GetString("shape"),
	}

	gPath, err := this.makeGraph(GetGraphSetting(graphSetting))
	if nil != err {
		ret["status"] = 0
		ret["msg"] = fmt.Sprintf("%s", err)
	}else{
		ret["data"] = gPath
	}
	this.Data["json"] = ret

	this.ServeJSON()
}

func getGraphSavePath(dotName string) string {
	staticDir := beego.BConfig.WebConfig.StaticDir["/static"]

	return fmt.Sprintf("%s/tmp/%s.png", staticDir, utils.SetStringDatePrefix(dotName))
}

func getGvTmpPath() string {
	return beego.AppConfig.String("TMP_GVFILE_PATH")
}

func getSolTmpDir() string {
	return beego.AppConfig.String("TMP_SOLFILE_DIR")
}

func DrawGvGraph(fileName string, format string) (string, error) {
	dotFilePath   := getGvTmpPath()
	graphSavePath := getGraphSavePath(fileName)
	//var out bytes.Buffer
	//var stderr bytes.Buffer

	os.Remove(graphSavePath)
	f, _ := os.Create(graphSavePath)
	f.Close()
	cmd := exec.Command("dot", fmt.Sprintf("-T%s", format), "-o", graphSavePath, dotFilePath)
	//cmd.Stdout = &out
	//cmd.Stderr = &stderr

	runError := cmd.Run()

	//f nil != runError {
	//	fmt.Println(fmt.Sprint(runError) + ": " + stderr.String())
	//

	//fmt.Println("Result: " + out.String())

	return graphSavePath, runError
}


/**
 * @retur: return the graph's save path or error description.
 */
func (this *IndexController)makeGraph(graphSetting map[string]string) (string, error) {
	var ret string
	fileName, saveFilePath, uploadErr := utils.SaveUploadFileWithPath(this.Ctx.Request, "codefile", getSolTmpDir(), true)
	if nil != uploadErr {
		return ret, uploadErr
	}

	// defer func() {
	// 	os.Remove(saveFilePath)
	// }()

	funcs, modifysolErr := utils.ModifySol(saveFilePath)
	if nil != modifysolErr {
		return ret, modifysolErr
	}
	for _, fname := range funcs {
		utils.OwPush(fname)
		utils.OwPop(fname)
	}
	utils.EndDot(getGvTmpPath(), graphSetting)
	// defer os.Remove(getGvTmpPath())

	ret, drawErr := DrawGvGraph(fileName, "png")

	return ret, drawErr
}