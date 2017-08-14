//    Copyright 2015 Digota Inc.
//
//    This program is free software: you can redistribute it and/or  modify
//    it under the terms of the GNU Affero General Public License, version 3,
//    as published by the Free Software Foundation.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU Affero General Public License for more details.
//
//    You should have received a copy of the GNU Affero General Public License
//    along with this program.  If not, see <https://www.gnu.org/licenses/agpl-3.0.en.html>.

package service

import (
	_ "github.com/digota/digota/product/service"
)

import (
	"github.com/digota/digota/config"
	"github.com/digota/digota/locker"
	"github.com/digota/digota/payment/paymentpb"
	"github.com/digota/digota/product"
	"github.com/digota/digota/product/productpb"
	"github.com/digota/digota/sku"
	"github.com/digota/digota/sku/skupb"
	"github.com/digota/digota/storage"
	"github.com/icrowley/fake"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"os"
	"reflect"
	"testing"
	"time"
)

var db = "testing-sku-" + uuid.NewV4().String()

func TestMain(m *testing.M) {
	// storage
	if err := storage.New(config.Storage{
		Address:  []string{"localhost"},
		Handler:  "mongodb",
		Database: db,
	}); err != nil {
		panic(err)
	}
	// locker
	if err := locker.New(config.Locker{
		Handler: "zookeeper",
		Address: []string{"localhost"},
	}); err != nil {
		panic(err)
	}
	retCode := m.Run()
	storage.Handler().DropDatabase(db)
	// teardown
	os.Exit(retCode)
}

func TestSkus_GetNamespace(t *testing.T) {
	p := Skus{}
	if p.GetNamespace() != "sku" {
		t.Fatal()
	}
}

func TestProduct_GetNamespace(t *testing.T) {
	p := Sku{}
	if p.GetNamespace() != "sku" {
		t.Fatal()
	}
}

//func TestSKUInventory_GetNamespace(t *testing.T) {
//	p := SKUInventory{}
//	if p.GetNamespace() != "skuinventory" {
//		t.Fatal()
//	}
//}

func TestProduct_SetCreated(t *testing.T) {
	p := Sku{}
	ti := time.Now().Unix()
	p.SetCreated(ti)
	if p.Created != ti {
		t.Fatal()
	}
}

func TestProduct_SetId(t *testing.T) {
	p := Sku{}
	uid := uuid.NewV4().String()
	p.SetId(uid)
	if p.GetId() != uid {
		t.Fatal()
	}
}

func TestProduct_SetUpdated(t *testing.T) {
	p := Sku{}
	ti := time.Now().Unix()
	p.SetUpdated(ti)
	if p.Updated != ti {
		t.Fatal()
	}
}

func TestSKUService_New(t *testing.T) {

	skuService := skuService{}

	// validation fail
	if _, err := skuService.New(context.Background(), &skupb.NewRequest{}); err == nil {
		t.Fatal(err)
	}

	newReq := &skupb.NewRequest{
		Name:     "sku name 2123",
		Active:   true,
		Price:    10001,
		Currency: paymentpb.Currency_EUR,
		Parent:   "cb379ae1-8729-4b32-ba7a-3119dc2bd211",
		Metadata: map[string]string{
			"key": "val",
		},
		Image: "http://sadf.124.com",
		PackageDimensions: &skupb.PackageDimensions{
			Weight: 1.1,
			Length: 1.2,
			Height: 1.3,
			Width:  1.4,
		},
		Inventory: &skupb.Inventory{
			Quantity: 1111,
			Type:     skupb.Inventory_Finite,
		},
		Attributes: map[string]string{
			"color": "red", "size": "M",
		},
	}

	// validation pass
	// product id fail
	if _, err := skuService.New(context.Background(), newReq); err == nil {
		t.Fatal(err)
	}
	// create demo product
	p, err := product.Service().New(context.Background(), &productpb.NewRequest{
		Active:      true,
		Name:        fake.Sentences(),
		Description: fake.Sentences(),
		Shippable:   false,
		Url:         "http://" + fake.Characters() + ".com",
		Images:      []string{"http://" + fake.Characters() + ".com", "http://" + fake.Characters() + ".com"},
		Metadata: map[string]string{
			"key": "val",
		},
		Attributes: []string{"color"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// set new product id to pass product check
	newReq.Parent = p.GetId()

	// should pass
	if _, err := skuService.New(context.Background(), newReq); err != nil {
		t.Fatal(err)
	}
}

func TestSKUService_Get(t *testing.T) {

	s := skuService{}

	// create demo product
	p, err := product.Service().New(context.Background(), &productpb.NewRequest{
		Active:      true,
		Name:        fake.Sentences(),
		Description: fake.Sentences(),
		Shippable:   false,
		Url:         "http://" + fake.Characters() + ".com",
		Images:      []string{"http://" + fake.Characters() + ".com", "http://" + fake.Characters() + ".com"},
		Metadata: map[string]string{
			"key": "val",
		},
		Attributes: []string{"color"},
	})

	if err != nil {
		t.Fatal(err)
	}

	// create sku
	sku, err := s.New(context.Background(), &skupb.NewRequest{
		Name:     "sku name 2123",
		Active:   true,
		Price:    10001,
		Currency: paymentpb.Currency_EUR,
		Parent:   p.GetId(),
		Metadata: map[string]string{
			"key": "val",
		},
		Image: "http://sadf.124.com",
		PackageDimensions: &skupb.PackageDimensions{
			Weight: 1.1,
			Length: 1.2,
			Height: 1.3,
			Width:  1.4,
		},
		Inventory: &skupb.Inventory{
			Quantity: 1111,
			Type:     skupb.Inventory_Finite,
		},
		Attributes: map[string]string{
			"color": "red", "size": "M",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	// pass
	if _, err := s.Get(context.Background(), &skupb.GetRequest{Id: sku.GetId()}); err != nil {
		t.Fatal(err)
	}

	// validation fail
	if _, err := s.Get(context.Background(), &skupb.GetRequest{Id: "notuuid"}); err == nil {
		t.Fatal(err)
	}

	// not found
	if _, err := s.Get(context.Background(), &skupb.GetRequest{Id: uuid.NewV4().String()}); err == nil {
		t.Fatal(err)
	}

	// lock fail
	unlock, err := locker.Handler().Lock(&Sku{Sku: skupb.Sku{Id: sku.GetId()}})
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()
	if _, err := s.Get(context.Background(), &skupb.GetRequest{Id: sku.GetId()}); err == nil {
		t.Fatal(err)
	}

}

func TestSKUService_Update(t *testing.T) {

	s := skuService{}

	// create demo product
	p, err := product.Service().New(context.Background(), &productpb.NewRequest{
		Active:      true,
		Name:        fake.Sentences(),
		Description: fake.Sentences(),
		Shippable:   false,
		Url:         "http://" + fake.Characters() + ".com",
		Images:      []string{"http://" + fake.Characters() + ".com", "http://" + fake.Characters() + ".com"},
		Metadata: map[string]string{
			"key": "val",
		},
		Attributes: []string{"color"},
	})

	if err != nil {
		t.Fatal(err)
	}

	// create sku
	skuItem, err := s.New(context.Background(), &skupb.NewRequest{
		Name:     "sku name 2123",
		Active:   true,
		Price:    10001,
		Currency: paymentpb.Currency_EUR,
		Parent:   p.GetId(),
		Metadata: map[string]string{
			"key": "val",
		},
		Image: "http://sadf.124.com",
		PackageDimensions: &skupb.PackageDimensions{
			Weight: 1.1,
			Length: 1.2,
			Height: 1.3,
			Width:  1.4,
		},
		Inventory: &skupb.Inventory{
			Quantity: 1111,
			Type:     skupb.Inventory_Finite,
		},
		Attributes: map[string]string{
			"color": "red", "size": "M",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	// validation fail
	if _, err := s.Update(context.Background(), &skupb.UpdateRequest{Id: skuItem.Id + "notvalid"}); err == nil {
		t.Fatal(err)
	}

	// get fail
	if _, err := s.Update(context.Background(), &skupb.UpdateRequest{Id: uuid.NewV4().String()}); err == nil {
		t.Fatal(err)
	}

	func() {
		// lock fail
		unlock, err := locker.Handler().Lock(&Sku{Sku: skupb.Sku{Id: skuItem.GetId()}})
		if err != nil {
			t.Fatal(err)
		}
		defer unlock()

		// try lock fail
		if _, err := s.Update(context.Background(), &skupb.UpdateRequest{Id: skuItem.GetId()}); err == nil {
			t.Fatal(err)
		}
	}()

	// fake product id
	if _, err := s.Update(context.Background(), &skupb.UpdateRequest{
		Id:     skuItem.GetId(),
		Parent: uuid.NewV4().String(),
	}); err == nil {
		t.Fatal(err)
	}

	// valid product id
	// attrs from req
	if sku1, err := s.Update(context.Background(), &skupb.UpdateRequest{
		Id:     skuItem.GetId(),
		Parent: p.GetId(),
		Attributes: map[string]string{
			"color": "red",
			"size":  "m",
		},
	}); err != nil {
		t.Fatal(err)
	} else if sku1.Attributes["color"] != "red" {
		t.Fatal()
	}

	// valid product id
	// attrs from sku
	if _, err := s.Update(context.Background(), &skupb.UpdateRequest{
		Id:     skuItem.GetId(),
		Parent: p.GetId(),
	}); err != nil {
		t.Fatal(err)
	}

	// update price
	sku2, err := s.Update(context.Background(), &skupb.UpdateRequest{
		Id:       skuItem.GetId(),
		Parent:   p.GetId(),
		Price:    1000,
		Currency: paymentpb.Currency_USD,
		Metadata: map[string]string{
			"key": "val",
		},
		Name: "123",
		PackageDimensions: &skupb.PackageDimensions{
			Width:  1,
			Height: 1,
			Length: 1,
			Weight: 1,
		},
		Inventory: &skupb.Inventory{
			Type:     skupb.Inventory_Finite,
			Quantity: 20,
		},
		Image: "http://image.com/image.png",
	})
	if err != nil {
		t.Fatal(err)
	}

	if sku2.Parent != p.GetId() {
		t.Fatal()
	}

	if sku2.Price != 1000 && sku2.Currency != paymentpb.Currency_USD {
		t.Fatal()
	}

	if len(sku2.Metadata) != 1 && sku2.Metadata["key"] != "val" {
		t.Fatal()
	}

	if sku2.Name != "123" {
		t.Fatal()
	}

	if !reflect.DeepEqual(sku2.PackageDimensions, &skupb.PackageDimensions{
		Width:  1,
		Height: 1,
		Length: 1,
		Weight: 1,
	}) {
		t.Fatal()
	}

	if !reflect.DeepEqual(sku2.Inventory, &skupb.Inventory{
		Type:     skupb.Inventory_Finite,
		Quantity: 20,
	}) {
		t.Fatal()
	}

	if sku2.Image != "http://image.com/image.png" {
		t.Fatal()
	}

}

func TestSKUService_Delete(t *testing.T) {

	s := skuService{}

	// create demo product
	p, err := product.Service().New(context.Background(), &productpb.NewRequest{
		Active:      true,
		Name:        fake.Sentences(),
		Description: fake.Sentences(),
		Shippable:   false,
		Url:         "http://" + fake.Characters() + ".com",
		Images:      []string{"http://" + fake.Characters() + ".com", "http://" + fake.Characters() + ".com"},
		Metadata: map[string]string{
			"key": "val",
		},
		Attributes: []string{"color"},
	})

	if err != nil {
		t.Fatal(err)
	}

	// create sku
	sku, err := s.New(context.Background(), &skupb.NewRequest{
		Name:     "sku name 2123",
		Active:   true,
		Price:    10001,
		Currency: paymentpb.Currency_EUR,
		Parent:   p.GetId(),
		Metadata: map[string]string{
			"key": "val",
		},
		Image: "http://sadf.124.com",
		PackageDimensions: &skupb.PackageDimensions{
			Weight: 1.1,
			Length: 1.2,
			Height: 1.3,
			Width:  1.4,
		},
		Inventory: &skupb.Inventory{
			Quantity: 1111,
			Type:     skupb.Inventory_Finite,
		},
		Attributes: map[string]string{
			"color": "red", "size": "M",
		},
	})

	if _, err := s.Delete(context.Background(), &skupb.DeleteRequest{Id: sku.GetId() + "notvalid"}); err == nil {
		t.Fatal(err)
	}

	if _, err := s.Delete(context.Background(), &skupb.DeleteRequest{Id: sku.GetId()}); err != nil {
		t.Fatal(err)
	}

	if _, err := s.Delete(context.Background(), &skupb.DeleteRequest{Id: uuid.NewV4().String()}); err == nil {
		t.Fatal(err)
	}

	func() {
		// lock fail
		unlock, err := locker.Handler().Lock(&Sku{skupb.Sku{Id: sku.GetId()}})
		if err != nil {
			t.Fatal(err)
		}
		defer unlock()
		// try lock fail
		if _, err := s.Delete(context.Background(), &skupb.DeleteRequest{Id: sku.GetId()}); err == nil {
			t.Fatal(err)
		}
	}()

}

func TestSKUService_GetWithInventoryLock(t *testing.T) {

	s := skuService{}

	// create demo product
	p, err := product.Service().New(context.Background(), &productpb.NewRequest{
		Active:      true,
		Name:        fake.Sentences(),
		Description: fake.Sentences(),
		Shippable:   false,
		Url:         "http://" + fake.Characters() + ".com",
		Images:      []string{"http://" + fake.Characters() + ".com", "http://" + fake.Characters() + ".com"},
		Metadata: map[string]string{
			"key": "val",
		},
		Attributes: []string{"color"},
	})

	if err != nil {
		t.Fatal(err)
	}

	skuItem, err := s.New(context.Background(), &skupb.NewRequest{
		Name:     "sku name 2123",
		Active:   true,
		Price:    10001,
		Currency: paymentpb.Currency_EUR,
		Parent:   p.GetId(),
		Metadata: map[string]string{
			"key": "val",
		},
		Image: "http://sadf.124.com",
		PackageDimensions: &skupb.PackageDimensions{
			Weight: 1.1,
			Length: 1.2,
			Height: 1.3,
			Width:  1.4,
		},
		Inventory: &skupb.Inventory{
			Quantity: 1111,
			Type:     skupb.Inventory_Finite,
		},
		Attributes: map[string]string{
			"color": "red", "size": "M",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	sku2, unlock, update, err := s.GetWithInventoryLock(context.Background(), &sku.GetWithInventoryLockRequest{
		Id:       skuItem.GetId(),
		Duration: time.Second,
	})

	if err != nil {
		t.Fatal(err)
	}

	sku2.Image = "updateimage"
	if err := update(); err != nil {
		t.Fatal(err)
	}

	unlock()

	func() {
		// lock fail
		unlock, err := locker.Handler().Lock(&Sku{skupb.Sku{Id: skuItem.GetId()}})
		if err != nil {
			t.Fatal(err)
		}
		defer unlock()
		if _, _, _, err := s.GetWithInventoryLock(context.Background(), &sku.GetWithInventoryLockRequest{
			Id:       skuItem.GetId(),
			Duration: time.Second,
		}); err == nil {
			t.Fatal(err)
		}
	}()

}

func TestSKUService_List(t *testing.T) {

	// clean collection
	storage.Handler().DropCollection(db, &Sku{})

	s := skuService{}

	// create demo product
	p, err := product.Service().New(context.Background(), &productpb.NewRequest{
		Active:      true,
		Name:        fake.Sentences(),
		Description: fake.Sentences(),
		Shippable:   false,
		Url:         "http://" + fake.Characters() + ".com",
		Images:      []string{"http://" + fake.Characters() + ".com", "http://" + fake.Characters() + ".com"},
		Metadata: map[string]string{
			"key": "val",
		},
		Attributes: []string{"color"},
	})

	if err != nil {
		t.Fatal(err)
	}

	// create few products
	for _, v := range []string{"one", "two", "three"} {
		// create sku
		_, err := s.New(context.Background(), &skupb.NewRequest{
			Name:     v,
			Active:   true,
			Price:    10001,
			Currency: paymentpb.Currency_EUR,
			Parent:   p.GetId(),
			Metadata: map[string]string{
				"key": "val",
			},
			Image: "http://sadf.124.com",
			Inventory: &skupb.Inventory{
				Quantity: 1111,
				Type:     skupb.Inventory_Finite,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	l, err := s.List(context.Background(), &skupb.ListRequest{
		Page:  0,
		Limit: 10,
	})

	if err != nil {
		t.Fatal()
	}

	if l.Total != 3 {
		t.Fatal()
	}

}
