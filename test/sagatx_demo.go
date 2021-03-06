package main

import (
	"github.com/jeremyxu2010/matrix-saga-go"
	"fmt"
	"github.com/jeremyxu2010/matrix-saga-go/utils"
	"time"
	"os/signal"
	"syscall"
	"os"
	"errors"
)

var (
	BALANCES map[string]int
	TransferMoneySagaStartDecorated func() error
	TransferOutCompensableDecorated func(from string, amount int) error
	TransferInCompensableDecorated func(to string, amount int) error
)

func init()  {
	err := saga.DecorateSagaStartMethod(&TransferMoneySagaStartDecorated, TransferMoney, 20)
	if err != nil {
		panic(err)
	}
	err = saga.DecorateCompensableMethod(&TransferOutCompensableDecorated, TransferOut, CancelTransferOut, 5)
	if err != nil {
		panic(err)
	}
	err = saga.DecorateCompensableMethod(&TransferInCompensableDecorated, TransferIn, CancelTransferIn, 5)
	if err != nil {
		panic(err)
	}

	initDatas()
}

func initDatas(){
	BALANCES = make(map[string]int, 0)
	BALANCES["foo"] = 500
	BALANCES["bar"] = 500
}

func TransferMoney() error {
	err := TransferOutCompensableDecorated("foo", 100)
	if err != nil {
		return err
	}
	err = TransferInCompensableDecorated("bar", 100)
	if err != nil {
		return err
	}
	return nil
}

func TransferOut(from string, amount int) error {
	//ctx, _ := sagactx.GetSagaAgentContext()
	//fmt.Println(ctx.GlobalTxId, ctx.LocalTxId)

	oldAmount, _ := BALANCES[from]
	BALANCES[from] = oldAmount - amount
	return nil
}

func CancelTransferOut(from string, amount int) error {
	//ctx, _ := sagactx.GetSagaAgentContext()
	//fmt.Println(ctx.GlobalTxId, ctx.LocalTxId)

	oldAmount, _ := BALANCES[from]
	BALANCES[from] = oldAmount + amount
	return nil
}

func TransferIn(to string, amount int) error {
	//ctx, _ := sagactx.GetSagaAgentContext()
	//fmt.Println(ctx.GlobalTxId, ctx.LocalTxId)

	//oldAmount, _ := BALANCES[to]
	//BALANCES[to] = oldAmount + amount
	//return nil

	return errors.New("xx")
}

func CancelTransferIn(to string, amount int) error {
	oldAmount, _ := BALANCES[to]
	BALANCES[to] = oldAmount - amount
	return nil
}

func main() {
	utils.DisableHttpProxy()
	saga.InitSagaAgent("saga-go-demo", "10.12.142.216:30571", nil)
	TransferMoneySagaStartDecorated()
	stopped := false
	go func() {
		s := make(chan os.Signal)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
		<-s
		stopped = true
	}()
	for !stopped {
		fmt.Println(BALANCES["foo"], BALANCES["bar"])
		time.Sleep(time.Second * 3)
	}
}
