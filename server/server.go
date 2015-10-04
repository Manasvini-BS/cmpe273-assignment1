/*@Author Manasvini Banavara Suryanarayana
*SJSU ID : 010102040
*CMPE 273 Assignment #1
*/
package main

import (
    "fmt"
    "log"
    "net"
    "net/rpc"
    "net/rpc/jsonrpc"
    "net/http"
    "io/ioutil"
    "net/url"
    "encoding/json"
     "strconv"
    "strings"
    "math"
    "math/rand"
    
)


type Args1 struct {
    Symb string
	Bud float64
}
type Args2 struct {
    Tradeid int
    
}

type Buyresponse struct {
    Stock string
    Tid int
    Rem float64
}

type Portfolioresponse struct {
    Stock string
    Currentstockval float64
    Rem float64
}


type Arith int


var tidmap = map[int]Buyresponse{}


func (t *Arith) Buy(args *Args1, reply *Buyresponse) error {
    
    fmt.Println("Buy Received Request : ", args.Symb, " ", args.Bud)
    
    // logic of the first request served
    //var in string ="GOOG:70%,YHOO:10%,AAPL:20%"
    //var budget float64 = 1000.00
    var in string = args.Symb
    var budget float64 = args.Bud
    inputstring := strings.Split(in, ",")
    symbolpercentmap := make(map[string]string)

    for index:=0; index < len(inputstring) ; index++ {
            keyvaluepair := strings.Split(inputstring[index],":")
            keys := keyvaluepair[0]
            values := keyvaluepair[1]
            symbolpercentmap[keys]=values

     }
     var req string 
     for key := range symbolpercentmap {
        if len(req)==0{
            req = "\""+key+"\""
        } else {
            req = req + ",\""+key+"\""
        }
    }
    var x string = "select Symbol,Ask from yahoo.finance .quotes where symbol in ("+req+")"
    fmt.Println("Query sent to Yahoo API : " + x)
    var input string = url.QueryEscape(x)
    resp, err := http.Get("http://query.yahooapis.com/v1/public/yql?q="+input+"&format=json&env=http://datatables.org/alltables.env")

    if err != nil {
        fmt.Println("error occured")
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    var yahfin interface{}
    err1 := json.Unmarshal(body, &yahfin)

    m := yahfin.(map[string]interface{})
    i := m["query"]
    m1 := i.(map[string]interface{})
    j := m1["results"]
    m2 := j.(map[string]interface{})
    j1 := m2["quote"]
    m3 := j1.([]interface{})

    stockmap := make(map[string]string)
 
    for index:=0; index < len(m3) ; index++ {
        x1 := m3[index]
        x2 := x1.(map[string]interface{})
        y1 := x2["Ask"]
        var ex string = y1.(string)
    
    
        y2 := x2["Symbol"]
        var ex2 string = y2.(string)
    
         stockmap[ex2] = ex
    }
    if err1 !=nil {
        fmt.Println("error:",err1)
    }
    // calculating budget value 
    symbolstockcountmap := make(map[string]int)
    var unvestedamount float64 = 0
    for key, value := range symbolpercentmap {
        val := strings.Trim(value,"%")
        per, err5 := strconv.ParseFloat(val, 64)
        if err5 !=nil {
            fmt.Println("error:",err5)
        }
        var amount float64 = (per/100)*budget 
        unitprice := stockmap[key]
        price, err6 := strconv.ParseFloat(unitprice, 64)
        if err6 !=nil {
            fmt.Println("error:",err6)
        }
        var numberofstocks float64 = amount/price
        reminder := math.Mod(amount,price)
        symbolstockcountmap[key] = int(numberofstocks)
        unvestedamount = unvestedamount + reminder
        
    }
    //creating response stock string 
    var stockstr string
    for key, value := range stockmap {
        val1 := strconv.Itoa(symbolstockcountmap[key])

        if len(stockstr)==0{
            stockstr = "\""+key+":"+val1+":$"+value+"\""

        }else {
            stockstr = stockstr + ",\""+key+":"+val1+":$"+value+"\""
        }

    }
    fmt.Println("Response Sent back to client: ",stockstr, "  ",unvestedamount)
    fmt.Println("------------------------------------------------------")
    trackid := rand.Intn(1000)
    reply.Stock = stockstr
    reply.Tid = trackid
    reply.Rem = unvestedamount
    tidmap[trackid] = *reply
	return nil
}

//function definition of portfolio response

func (t *Arith) Portfolio(args *Args2, reply *Portfolioresponse) error {
    // mock input for portfolio
    var test Buyresponse
    test.Stock="\"GOOG:100:$500.25\",\"YHOO:200:$31.40\""
    test.Tid = 1
    test.Rem = 600
    tidmap[1] = test

    fmt.Println("Portfolio Request Received : ", args.Tradeid)
    tradeid := args.Tradeid
    tidrecord := tidmap[tradeid]
    stocks := tidrecord.Stock
    unvestedamount := tidrecord.Rem
    //split stock symbol to pass it to yql

    inputstring := strings.Split(stocks, ",")
    symbolstockcountmap := make(map[string]string)
    symbolpurchasepricemap := make(map[string]string)

    for index:=0; index < len(inputstring) ; index++ {
            keyvaluepair := strings.Split(inputstring[index],":")
          
            keys := strings.Trim(keyvaluepair[0],"\"")
            values := keyvaluepair[1]
            values2 := strings.Trim(keyvaluepair[2],"\"")
            symbolstockcountmap[keys]=values
            symbolpurchasepricemap[keys]=values2

     }
     var req string 
     for key := range symbolstockcountmap {
        if len(req)==0{
            req = "\""+key+"\""
        } else {
            req = req + ",\""+key+"\""
        }
    }
    var x string = "select Symbol,Ask from yahoo.finance .quotes where symbol in ("+req+")"
    fmt.Println("Query sent to Yahoo API : " + x)
    var input string = url.QueryEscape(x)
    resp, err := http.Get("http://query.yahooapis.com/v1/public/yql?q="+input+"&format=json&env=http://datatables.org/alltables.env")

    if err != nil {
        fmt.Println("error occured")
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    var yahfin interface{}
    err1 := json.Unmarshal(body, &yahfin)

    m := yahfin.(map[string]interface{})
    i := m["query"]
    m1 := i.(map[string]interface{})
    j := m1["results"]
    m2 := j.(map[string]interface{})
    j1 := m2["quote"]
    m3 := j1.([]interface{})

    stockmap := make(map[string]string)
 
    for index:=0; index < len(m3) ; index++ {
        x1 := m3[index]
        x2 := x1.(map[string]interface{})
        y1 := x2["Ask"]
        var ex string = y1.(string)
    
    
        y2 := x2["Symbol"]
        var ex2 string = y2.(string)
    
         stockmap[ex2] = ex
    }
    if err1 !=nil {
        fmt.Println("error:",err1)
    }
    var stockstr string
    var Currentstockval float64
    for key, value := range stockmap {
        val1 := symbolstockcountmap[key]
        val3 := symbolpurchasepricemap[key]
         val4 := strings.Trim(val3,"$")
         oldprice, err7 := strconv.ParseFloat(val4, 64)
        if err7 !=nil {
            fmt.Println("error:",err7)
        }
         newprice, err8 := strconv.ParseFloat(value, 64)
        if err8 !=nil {
            fmt.Println("error:",err8)
        }
        var profitsymbol string
        if oldprice > newprice {
            profitsymbol = "-"
        }else {
            profitsymbol = "+"
        }


        if len(stockstr)==0{
            stockstr = "\""+key+":"+val1+":"+profitsymbol+"$"+value+"\""

        }else {
            stockstr = stockstr + ",\""+key+":"+val1+":"+profitsymbol+"$"+value+"\""
        }
        count, err9 := strconv.ParseFloat(val1, 64)
        if err9 !=nil {
            fmt.Println("error:",err9)
        }
        Currentstockval = Currentstockval + (count * newprice)
    }
    fmt.Println("Response Sent back to client: ",stockstr, "  ", Currentstockval,"  ",unvestedamount)
    fmt.Println("------------------------------------------------------")
    reply.Stock = stockstr
    reply.Currentstockval = Currentstockval
    reply.Rem = unvestedamount
    return nil
}

func startServer() {
	
    arith := new(Arith)

    server := rpc.NewServer()
    server.Register(arith)

    server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
    fmt.Println("server started and listening at port 8061")
    l, e := net.Listen("tcp", ":8061")
    if e != nil {
        log.Fatal("listen error:", e)
    }

    for {
        conn, err := l.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go server.ServeCodec(jsonrpc.NewServerCodec(conn))
    }
    fmt.Println("Connection closed")
}

func main() {
    startServer()
    fmt.Println("Server ended")
    }


