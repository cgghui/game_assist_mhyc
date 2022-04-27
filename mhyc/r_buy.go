package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// Buy 购买道具
func Buy(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	defer t1.Stop()
	f1 := func() time.Duration {
		go func() {
			_ = CLI.ShopBuy(&C2SShopBuy{GoodsId: 616, Num: 5})
		}()
		_ = Receive.Wait(&S2CShopBuy{}, s3)
		go func() {
			_ = CLI.ShopBuy(&C2SShopBuy{GoodsId: 617, Num: 1})
		}()
		_ = Receive.Wait(&S2CShopBuy{}, s3)
		return TomorrowDuration(RandMillisecond(30000, 30600))
	}
	for {
		select {
		case <-t1.C:
			t1.Reset(f1())
		case <-ctx.Done():
			return
		}
	}
}

////////////////////////////////////////////////////////////

// ShopBuy 商城购物
func (c *Connect) ShopBuy(goods *C2SShopBuy) error {
	body, err := proto.Marshal(goods)
	if err != nil {
		return err
	}
	return c.send(432, body)
}

func (x *S2CShopBuy) ID() uint16 {
	return 433
}

func (x *S2CShopBuy) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ShopBuy] tag=%d %v", x.Tag, x)
}
