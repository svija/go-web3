package falcon

import (
	"encoding/base64"
	"encoding/hex"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"strconv"
	"io/ioutil"
	"fmt"
	"net/http"
	"io"
	"strings"
	"github.com/garyburd/redigo/redis"
	"github.com/svija/redisCliPool"
    "unsafe"
    "time"
    "os"
	"os/exec"
	"bytes"
	ini "gopkg.in/ini.v1"
)


func F_bool_to_str(in bool) string{
    if in {
        return "true"
    }
    return "false"
}

func F_int(str string)(rst int){
    fl,_:=strconv.Atoi(str)
    return fl
}

func F_int_to_str(in int)(rst string){
    return strconv.Itoa(in)
}

func F_float(str string)(rst float64){
    fl,_:=strconv.ParseFloat(str, 64)
    return fl
}

func F_float_to_str(str float64)(rst string){
    return strconv.FormatFloat(str, 'f', 12, 64)
}


func F_float_to_str8(str float64)(rst string){
    return strconv.FormatFloat(str, 'f', 8, 64)
}

func F_uint64_to_str(in uint64)(rst string){
    rst=strconv.FormatUint(in, 10)
    return rst
}

func F_int64_to_str(in int64)(rst string){
    rst=strconv.FormatInt(in, 10)
    return rst
}



const (
	base64Table        = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-"
	hashFunctionHeader = ""
	hashFunctionFooter = ""
)


func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s)) 
	return hex.EncodeToString(h.Sum(nil))
}


func GetSHA1String(s string) string {
	t := sha1.New()
	t.Write([]byte(s))
	return hex.EncodeToString(t.Sum(nil))

}


func GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

var coder = base64.NewEncoding(base64Table)


func Base64Encode(str string) string {
	var src []byte = []byte(hashFunctionHeader + str + hashFunctionFooter)
	return string([]byte(coder.EncodeToString(src)))
}

func Base64Decode(str string) (string, error) {
	var src []byte = []byte(str)
	by, err := coder.DecodeString(string(src))
	return strings.Replace(strings.Replace(string(by), hashFunctionHeader, "", -1), hashFunctionFooter, "", -1), err
}


func tkmon(reqstr string)(token string){
	TimeMon:=time.Now().Format("2006-01") 
	return GetSHA1String(reqstr+TimeMon+"C4b1a7f7ddf688a1d9a546e1ff63cdbeae6ef4309")
}

func tkday(reqstr string)(token string){
	TimeDay:=time.Now().Format("2006-01-02")
	return GetSHA1String(reqstr+TimeDay+"A830cea5b478ca075a321f0b4339cf4d6208dd68e")
}

func tkhour(reqstr string)(token string){
	TimeHour:=time.Now().Format("2006-01-02 15")
	return GetSHA1String(reqstr+TimeHour+"B2659f2d739db05ff12a0003cf4d361ec9ce398ba")
}

func tkua(reqstr string)(token string){
	TimeHour:=time.Now().Format("2006-01-02 15")
	return GetSHA1String(reqstr+TimeHour+"32323a")
}

func CmdBash(commandName string) *exec.Cmd {
	cmd := exec.Command("/bin/bash", "-c", commandName)
	go func() {
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}()
	return cmd
}


func findstr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n+len(start):])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}


func repstr(str,start,end,new string)string {
	if len(str)==0{
		return ""
	}else{
		spls:=strings.ReplaceAll(str,findstr(str,start,end),new)
		splx:=strings.Replace(spls,start,"",-1)
		return strings.Replace(splx,end,"",-1)
	}

}


func httpurl(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	return client.Do(request)
}

func httpgeturl(url string)([]byte,error)  {
	resp, err := httpurl(url)
	if err != nil{
		return nil,err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d of %s", resp.StatusCode, url)
	}
	return ioutil.ReadAll(resp.Body)
}



func httpGet(urllink string )(rsp string) {
	timeout := time.Duration(1888 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	{

		res, err := client.Get(urllink)
		if err != nil {
			rsp:=fmt.Sprintf("Pasa Status: %s", err)
			return rsp
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			//log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
			rsp:=fmt.Sprintf("Pasa Status Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
			return rsp
		}
		if err != nil {
			//log.Fatal(err)
			rsp:=fmt.Sprintf("Pasa Status: %s", err)
			return rsp
		}
		rsp:=fmt.Sprintf("%s", body)
		return rsp
	}
}


 
func rdslist_get(list ,index_start,index_end string )(data string){

	c:= redisCliPool.Clipool.Get()
	values,_:= redis.Values(c.Do("lrange",list,index_start,index_end))
	defer c.Close()
	result:=""
	i:=0
	for _,v := range values {
		i+=1
		result+="\""+strconv.Itoa(i)+"\":"+"\""+string(v.([]byte))+"\","
	}
	defer c.Close()
	return  strings.TrimRight(result, ",")
}

func rdslist_get_v2(list ,index_start,index_end string )(data string){

	c:= redisCliPool.Clipool.Get()
	values,_:= redis.Values(c.Do("lrange",list,index_start,index_end))
	defer c.Close()
	result:=""
	i:=0
	for _,v := range values {
		i+=1
		result+="\""+strconv.Itoa(i)+"\":"+string(v.([]byte))+","
	}
	defer c.Close()
	return  strings.TrimRight(result, ",")
}



func rdslist_get_len(listin string )(len int) {
	c := redisCliPool.Clipool.Get()
	rep, err := redis.Int(c.Do("LLEN",listin))
	defer c.Close()
	if err == nil {
		return rep
	} else {
		return 0
	}
}

func rdslist_get_rpop(list string )(result string){
	c:= redisCliPool.Clipool.Get()
	redis.Values(c.Do("rpop",list))
	defer c.Close()
	return  "ok"
}


func rdslist_delete(list ,removecoin string )(data string){

	c:= redisCliPool.Clipool.Get()
	redis.Values(c.Do("lrem",list,"0",removecoin))
	defer c.Close()
	return  "removed"
}


func rdslist_contain(list,index_start,index_end,contain string )(data string){
	allstr:=rdslist_get(list,index_start,index_end)
	if  strings.Contains(allstr, contain) ==true  {
		return "true"
	}else{
		return "false"
	}
}



func rdslist_set(list,list_value string )(rsp string){
	c:= redisCliPool.Clipool.Get()
	data,err := redis.String(c.Do("lpush",list,list_value))
	for err != nil {
		return("Set Fail")
		defer c.Close()
		break
	}
	defer c.Close()
	return data
}

func BytesToString(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}


func StringToBytes(data string) []byte {
	return *(*[]byte)(unsafe.Pointer(&data))
}


func keygetgroup(node string  )(rsp string){
	cfg, err := ini.Load("pasacmd.ini")
	sr:=fmt.Sprint(cfg.Section(node))
	if err != nil {
		//fmt.Printf("Fail to read file: %v", err)
		rsp:=fmt.Sprintf("Fail to read file: %v", err)
		return rsp
		os.Exit(1)
	}
	if ( sr==""){
		return "None Section Value Found!"
	}else{
		return findstr(sr,"] map[","}")
	}

}




func keyget(node string , key string )(rsp string){
	cfg, err := ini.Load("filestarlink.ini")
	if err != nil {
		//fmt.Printf("Fail to read file: %v", err)
		rsp:=fmt.Sprintf("Fail to read file: %v", err)
		return rsp
		os.Exit(1)
	}
	if (cfg.Section(node).Key(key).String()==""){
		return "None Value Find!"
	}else{
		return  cfg.Section(node).Key(key).String()
	}

}

func keygetgroupsq(nodebt string,nodexc string ,btbase string,xcbase string )(rsp string){

	btstr := []string{}
	xcstr := []string{}
	for  i:=1;i<300;i++{
		if keygetraw(nodebt,btbase+ strconv.Itoa(i))!=""{
			btstr = append(btstr, "{"+btbase+strconv.Itoa(i)+"}="+keygetraw(nodebt,btbase+ strconv.Itoa(i)))
			xcstr=append(xcstr,"{"+xcbase+strconv.Itoa(i)+"}="+keygetraw(nodexc,xcbase+ strconv.Itoa(i)))
			//log.Print(btstr,i)
		}
	}
	return "<$<"+fmt.Sprint(btstr)+">$>;(*("+fmt.Sprint(xcstr)+")*)"
}


func keyset(node string , key string ,key_value string )(rsp string){
	cfg, err := ini.Load("pasa.ini")
	if err != nil {
		//fmt.Printf("Fail to read file: %v", err)
		rsp:=fmt.Sprintf("Fail to read file: %v", err)
		return rsp
		os.Exit(1)
	}
	cfg.Section(node).Key(key).SetValue(key_value)
	cfg.SaveTo("pasa.ini")
	return  "[Key]:"+key_value+"  write to cfg done!"
}


func miset(node string , key string ,key_value string )(rsp string){
	cfg, err := ini.Load("milontrol.ini")
	if err != nil {
		//fmt.Printf("Fail to read file: %v", err)
		rsp:=fmt.Sprintf("Fail to read file: %v", err)
		return rsp
		os.Exit(1)
	}
	cfg.Section(node).Key(key).SetValue(key_value)
	cfg.SaveTo("milontrol.ini")
	return  "[Key]:"+key_value+"  write to cfg done!"
}



func miget(node string , key string )(rsp string){
	cfg, err := ini.Load("milontrol.ini")
	if err != nil {
		//fmt.Printf("Fail to read file: %v", err)
		rsp:=fmt.Sprintf("Fail to read file: %v", err)
		return rsp
		os.Exit(1)
	}
	if (cfg.Section(node).Key(key).String()==""){
		return "None Value Find in Milontrol!"
	}else{
		return  cfg.Section(node).Key(key).String()
	}
}


func tarGet(node string , key string,filename string )(rsp string){
	cfg, err := ini.Load(filename)
	if err != nil {
		//fmt.Printf("Fail to read file: %v", err)
		rsp:=fmt.Sprintf("Fail to read file: %v", err)
		return rsp
		os.Exit(1)
	}
	if (cfg.Section(node).Key(key).String()==""){
		return "None Value Find in Milontrol!"
	}else{
		return  cfg.Section(node).Key(key).String()
	}
}





func keygetraw(node string , key string )(rsp string){
	cfg, err := ini.Load("pasacmd.ini")
	if err != nil {
		//fmt.Printf("Fail to read file: %v", err)
		rsp:=fmt.Sprintf("Fail to read file: %v", err)
		return rsp
		os.Exit(1)
	}
	return  cfg.Section(node).Key(key).String()
}

func keygetbyfile(node string ,key string,file string)(rsp string){
	cfg, err := ini.Load(file+".ini")
	if err != nil {
		//fmt.Printf("Fail to read file: %v", err)
		//rsp:=fmt.Sprintf("Fail to read file: %v", err)
		//return rsp
		//os.Exit(1)
		return "ErrorCode:8268;Status:This Service or Solution Developing;"
	}
	return  cfg.Section(node).Key(key).String()
}


func tarset(node string , key string ,key_value string,filename string )(rsp string){
	cfg, err := ini.Load(filename)
	if err != nil {
		//fmt.Printf("Fail to read file: %v", err)
		rsp:=fmt.Sprintf("Fail to read file: %v", err)
		return rsp
		os.Exit(1)
	}
	cfg.Section(node).Key(key).SetValue(key_value)
	cfg.SaveTo(filename)
	return  "[Key]:"+key_value+"  write to cfg done!"
}


func rdsget(key string )(value string){

	c:= redisCliPool.Clipool.Get()
	dataget,err := redis.String(c.Do("get",key))
	for err != nil {
		return("No Value")
		defer c.Close()
		break
	}
	defer c.Close()
	return dataget
}


func rdsdel(key string )(value string){

	c:= redisCliPool.Clipool.Get()
	dataget,err := redis.String(c.Do("del",key))
	for err != nil {
		return("No Value")
		defer c.Close()
		break
	}
	defer c.Close()
	return dataget
}

func rdsset(key,keyvalue string )(rsp string){

	c:= redisCliPool.Clipool.Get()
	data,err := redis.String(c.Do("set",key,keyvalue))
	for err != nil {
		return("Set Fail")
		defer c.Close()
		break
	}
	defer c.Close()
	return data
}




func CmdAndChangeDir(dir string, commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	fmt.Println("CmdAndChangeDir", dir, cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}


func Cmd(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	fmt.Println("Cmd", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}

func shell(s string) (string, error){
	cmd := exec.Command("/bin/bash", "-c", s)


	var out bytes.Buffer
	cmd.Stdout = &out

	
	err := cmd.Run()
	return out.String(), err
}

func wshell(s string) (string, error){
	cmd := exec.Command("cmd", "/C", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}


func SubStr(str string, start int, length int) (result string) {
	s := []rune(str)
	total := len(s)
	if total == 0 {
		return
	}

	if start < 0 {
		start = total + start
		if start < 0 {
			return
		}
	}
	if start > total {
		return
	}

	if length < 0 {
		length = total
	}

	end := start + length
	if end > total {
		result = string(s[start:])
	} else {
		result = string(s[start:end])
	}

	return
}

/**##########################################################
* 模块：常用包
* 说明：
* 备注
 *##########################################################*/
	//"bytes"
	//"encoding/base64"
	//"encoding/hex"
	//"encoding/json"
	//"github.com/tidwall/gjson"
	//"strconv"
	//"github.com/garyburd/redigo/redis"
	//"os"
	//"os/exec"
	//"github.com/ethereum/go-ethereum/common/hexutil"
    //"github.com/ethereum/go-ethereum/crypto"
