package utils

import (
	"bufio"
	"coder/packs/collections/stack"
	"coder/packs/gographviz"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var lock sync.Mutex
var graph *gographviz.Graph

type data struct {
	funName string
	start int64
	end int64
	father string
	startPc int64
	endPc int64
	gasStart uint64
	gasEnd uint64
}

var s = stack.New()
var dot  = stack.New()
var pc int64
var sum uint64

func NewDot(bgcolor string) {
	lock.Lock()
	digraphStr := fmt.Sprintf("digraph G{bgcolor=%s}", bgcolor)
	graphAst, _ := gographviz.Parse([]byte(digraphStr))
	graph = gographviz.NewGraph()
	gographviz.Analyse(graphAst, graph)

	lock.Unlock()
}

func OwPush(funName string) {
	pc  = 0
	sum = 0
	s.Push(data{
		funName: funName,
		start: time.Now().UnixNano(),
		startPc: pc,
		gasStart: sum,
	})
}

func OwPop(funName string)  {
	d := *new(data)
	if s.Len() != 0{
		d =s.Peek().(data)
		if d.funName == funName{
			s.Pop()
			d.end = time.Now().UnixNano()
			d.endPc = pc
			d.gasEnd = sum
			if s.Len() != 0 {
				d.father=s.Peek().(data).funName
			} else {
				d.father="nil"
			}
			dot.Push(d)
			//fmt.Println(d.father)
			//fmt.Println(d.funName)
			//graph.AddNode("G",d.funName,nil)
			//if d.father!="nil"{
			//	graph.AddEdge(d.father, d.funName, true, nil)
			//}
		}
	}
}

func ReadLine(fileName string) ([]string){
	f, err := os.Open(fileName)
	if err != nil {
		return nil
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result []string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {   //读取结束，会报EOF
				return result
			}
			return nil
		}
		result = append(result,line)
	}
	return result
}

func stackNil() {
	l := s.Len()
	for i := 0; i < l; i++ {
		s.Pop()
	}
}

func writeDot(setting map[string]string)  {
	NewDot(setting["bgcolor"])
	node :=make(map[string]string)
	node["shape"]="record,height=.1"
	for dot.Len()!=0 {
		dotdata :=dot.Pop().(data)
		//时间
		//fmt.Println(dotdata.end-dotdata.start)
		//字节码个数
		//fmt.Println("数值")
		//fmt.Println(dotdata.endPc-dotdata.startPc)
		c := `"`+dotdata.funName+"|time:"+strconv.FormatInt(dotdata.end-dotdata.start,10)+"ms,BytecodeNumber:"+strconv.FormatInt(dotdata.endPc-dotdata.startPc,10)+",gasCost:"+strconv.FormatUint(dotdata.gasEnd-dotdata.gasStart,10)+`"`
		node["label"]=c
		node["shape"]=setting["shape"]
		node["color"]=setting["color"]
		node["style"]=setting["style"]
		node["fillcolor"]=setting["fillcolor"]
		graph.AddNode("G",dotdata.funName,node)
		if dotdata.father != "nil"{
			graph.AddEdge(dotdata.father, dotdata.funName, true, nil)
		}
	}
}

func AddPc(cost uint64){
	pc++
	sum += cost
	fmt.Println("gasSum:",sum)
}

func EndDot(saveFilePath string, setting map[string]string) {
	stackNil()
	writeDot(setting)
	//fmt.Println(graph.String())
	SaveFile(saveFilePath, graph.String())
}

func ModifySol(fileName string) ([]string, error) {
	var file string
	var funName string
	var funNames []string
	fi, err := os.Open(fileName)
	if err != nil {
		return funNames, err
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		r := regexp.MustCompile("(?U)function(.*)\\(")
		z:=r.FindString(string(a))
		if z != ""{
			funName =z[9:len(z)-1]
			funNames = append(funNames, funName + "\n")
			res, _ := regexp.Compile("(?U)function(.*){")
			file=file+res.ReplaceAllString(string(a),string(a)+"\n"+"bytes32 funName="+ "'"+z[9:len(z)-1]+"';")+"\n"
		}else {
			file=file+string(a)+"\n"
		}
	}

	return funNames, nil
}

func SaveFile(fileNameNew string,wireteString string)  {
	os.Remove(fileNameNew)
	f, err := os.Create(fileNameNew)
	defer f.Close()
	CheckErr(err)
	_, err = io.WriteString(f, wireteString)
	//fmt.Printf("写入 %d 个字节\n", n)
}