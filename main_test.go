package main

import (
	"fmt"
	"testing"
)

type IBreadInterface interface {
    BreadName() string
    Price() int
}

type ButteredBread struct {
}

func (ButteredBread) BreadName() string {
    return "奶油面包"
}

func (ButteredBread) Price() int {
    return 10
}

 func (ButteredBread) Color() string {
     return "白色"
 }

func (ButteredBread) ProdDate() string {
    return "2019-06-19"
}

func BenchmarkIfaceToType(b *testing.B) {
    b.Run("iface-to-type", func(b *testing.B) {
        
        for i := 0; i < b.N; i ++ {
			var iface IBreadInterface = ButteredBread{}
            iface.BreadName()
            iface.Price()
        }
    })
    b.Run("direct-call", func(b *testing.B) {
        
        for i := 0; i < b.N; i ++ {
			var bread = ButteredBread{}
            bread.BreadName()
            bread.Price()
        }
    })
fmt.Println("mk")
}