/*@Author Manasvini Banavara Suryanarayana
*SJSU ID : 010102040
*CMPE 273 Assignment #1
*/
package main

import (
    "fmt"
    "log"
    "net"
 	"net/rpc/jsonrpc"
    "os"
    "strconv"
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


func main() {
   
    conn, err := net.Dial("tcp", "localhost:8061")
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    c := jsonrpc.NewClient(conn)
    //fmt.Println(len(os.Args))
    commandinput := os.Args
    if (len(commandinput) == 3){
        stockstr := commandinput[1]
        budget, err6 := strconv.ParseFloat(commandinput[2], 64)
        if err6 !=nil {
            fmt.Println("error:",err6)
        }
        fmt.Println("Request : {stockSymbolAndPercentage : "+stockstr+", budget : "+commandinput[2]+" }")
        args := Args1{stockstr, budget}
        var reply Buyresponse
        err = c.Call("Arith.Buy", args, &reply)
        if err != nil {
            log.Fatal("arith error:", err)
        }
        halt := fmt.Sprintf("%.2f",reply.Rem)
        //fmt.Printf("%d  %s  %s\n", reply.Tid, reply.Stock, halt)
        fmt.Println("Response : {tradeid : "+strconv.Itoa(reply.Tid)+"  stocks : "+reply.Stock+", unvestedamount : "+halt+" }")
    
    } else if (len(commandinput) == 2){
        tid, err7 := strconv.Atoi(commandinput[1])
        if err7 !=nil {
            fmt.Println("error:",err7)
        }
        fmt.Println("Request : {tradeid : "+commandinput[1]+" }")
        args2 := Args2{tid}
        var reply2 Portfolioresponse
        err = c.Call("Arith.Portfolio", args2, &reply2)
        if err != nil {
            log.Fatal("arith error:", err)
        }
        halt2 := fmt.Sprintf("%.2f",reply2.Rem)
        halt3 := fmt.Sprintf("%.2f",reply2.Currentstockval)
     
        fmt.Println("Response : {stocks : "+reply2.Stock+"  currentstockvalue : "+halt3+"  unvestedamount : "+halt2+" }")

    } else {
        println(" Please use this command:\n go run client.go <stockstring> <budget>\n OR \n go run client.go <tradeid>")
    }
   
}