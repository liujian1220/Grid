package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"strconv"
)

type Grid_simple struct{

}
type Node struct {
	NodeName   string `json:"name"`
	NodeType  string `json:"type"`
	NodePowerNum string `json:"powernum"`
	NodeMoneyNum string `json:"moneynum"`
	NodeOwner  string `json:"owner"`
	NodeCredit string `json:"credit"`
}
func (g *Grid_simple)InitNode(stub shim.ChaincodeStubInterface){
	n1:=Node{"peer0.org1.com","power","100","1000","org1","1"}
    n2:=Node{"peer0.org2.com","load","40","500","org2","1"}
    n1bytes,_:=json.Marshal(n1)
	n2bytes,_:=json.Marshal(n2)
	stub.PutState(n1.NodeName,n1bytes)
	stub.PutState(n2.NodeName,n2bytes)
}
func (g *Grid_simple)Init(stub shim.ChaincodeStubInterface)peer.Response{
	_,args:=stub.GetFunctionAndParameters()
	if len(args)!=1{
		fmt.Println("the init args must be one")
	}
	fmt.Println("init successfil!!!")
	return shim.Success([]byte("success"))
}
func (g *Grid_simple)Invoke(stub shim.ChaincodeStubInterface)peer.Response{
	function,args:=stub.GetFunctionAndParameters()
	if function=="get"{
		if len(args)!=1{
			fmt.Println("the function GetElectricityAndMoney args must be one")
		}
		_,_=GetElectricityAndMoney(stub,args[0])
		return shim.Success(nil)
	}else if function=="init"{
		if  len(args)!=0{
			fmt.Println("the function BuyToken args must be zero")
		}
	}else if function=="trans"{
		_=Trans(stub,args[0],args[1],args[2])
	}
	return shim.Success([]byte(""))
}

func Trans(stub shim.ChaincodeStubInterface,a,b,c string)error{
	//a->b the money c
	a_electricnum:=GetNode(stub,a).NodeMoneyNum
	b_electricnum:=GetNode(stub,b).NodeMoneyNum
	if a_electricnum<c{
		fmt.Printf("the %s money is%s,less than %s",a,a_electricnum,c)
		return nil
	}
	a_int,_:=strconv.Atoi(a_electricnum)
	b_int,_:=strconv.Atoi(b_electricnum)
	c_int,_:=strconv.Atoi(c)
	num1:=strconv.Itoa(a_int-c_int)
	num2:=strconv.Itoa(b_int+c_int)
	SetMoneyNum(stub,num1,GetNode(stub,a).NodeName)
	SetMoneyNum(stub,num2,GetNode(stub,b).NodeName)
	fmt.Printf("%strans to %s the money as %s",a,b,c)
	return nil

}
func SetMoneyNum(stub shim.ChaincodeStubInterface,a,b string)error{
	_,err:=strconv.Atoi(a)
	if err!=nil{
		fmt.Println("the SetMoney must be number")
		return err
	}

	GetNode(stub,b).NodeMoneyNum=a
	NodeBytes,_:=json.Marshal(GetNode(stub,b))
	stub.PutState(b,NodeBytes)
	return nil
}


func main() {
	if err := shim.Start(new(Grid_simple)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}

}
func GetElectricityAndMoney(stub shim.ChaincodeStubInterface,s string)(a,b string){
	nodeBytes,err:=stub.GetState(s)
	if err!=nil{
		fmt.Println("Failed read the world state")
		return "",""
	}
	if nodeBytes==nil{
		fmt.Printf("%sis not exist\n",nodeBytes)
	}
	node:=new(Node)
	_=json.Unmarshal(nodeBytes,node)
	fmt.Printf("the electric of %s i s%s,the money of %s is %s\n",s,a,s,b)
	return node.NodePowerNum,node.NodeMoneyNum

}
func GetNode(stub shim.ChaincodeStubInterface,s string)(*Node){
	nodeBytes,err:=stub.GetState(s)
	if err!=nil{
		fmt.Println("Failed read the world state")
		return nil
	}
	if nodeBytes==nil{
		fmt.Printf("%sis not exist\n",nodeBytes)
	}
	node:=new(Node)
	_=json.Unmarshal(nodeBytes,node)
	return node
}
