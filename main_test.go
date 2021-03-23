package main

import (
	"fmt"
	"testing"
	"time"
)

type IBreadInterface interface {
	BreadName() string
	Price() int
}

type ButteredBread struct{}

func (ButteredBread) BreadName() string {
	return "奶油面包"
}

func (ButteredBread) Price() int {
	return 10
}

/*
func (ButteredBread) Color() string {
	return "白色"
}

func (ButteredBread) ProdDate() string {
	return "2019-06-19"
}*/

func BenchmarkIfaceToType(b *testing.B) {
	b.Run("iface-to-type", func(b *testing.B) {
		var iface IBreadInterface = ButteredBread{}
		for i := 0; i < b.N; i++ {

			iface.BreadName()
			iface.Price()
		}
	})
	b.Run("direct-call", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bread := ButteredBread{}
			bread.BreadName()
			bread.Price()
		}
	})
}

func TestChan(t *testing.T) {
	done := make(chan struct{})
	go func() {
		fmt.Println("start", time.Now())
		time.Sleep(3 * time.Second)
		close(done)
	}()
	<-done
	fmt.Println("end", time.Now())
}
