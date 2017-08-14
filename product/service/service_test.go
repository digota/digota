//   Copyright 2015 Digota Inc.
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
	_ "github.com/digota/digota/sku/service"
)

import (
	"github.com/digota/digota/config"
	"github.com/digota/digota/locker"
	"github.com/digota/digota/product/productpb"
	"github.com/digota/digota/storage"
	"github.com/icrowley/fake"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"os"
	"testing"
	"time"
)

var service = &productService{}
var db = "testing-product-" + uuid.NewV4().String()

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

func TestProducts_GetNamespace(t *testing.T) {
	p := Products{}
	if p.GetNamespace() != "product" {
		t.Fatal()
	}
}

func TestProduct_GetNamespace(t *testing.T) {
	p := Product{}
	if p.GetNamespace() != "product" {
		t.Fatal()
	}
}

func TestProduct_SetCreated(t *testing.T) {
	p := Product{}
	ti := time.Now().Unix()
	p.SetCreated(ti)
	if p.Created != ti {
		t.Fatal()
	}
}

func TestProduct_SetId(t *testing.T) {
	p := Product{}
	uid := uuid.NewV4().String()
	p.SetId(uid)
	if p.GetId() != uid {
		t.Fatal()
	}
}

func TestProduct_SetUpdated(t *testing.T) {
	p := Product{}
	ti := time.Now().Unix()
	p.SetUpdated(ti)
	if p.Updated != ti {
		t.Fatal()
	}
}

func TestProductService_New(t *testing.T) {

	// ok
	if _, err := service.New(context.Background(), &productpb.NewRequest{
		Name:        fake.Brand(),
		Active:      true,
		Attributes:  []string{"size"},
		Description: fake.Paragraph(),
		Images: []string{
			"http://digota.com/image1.jpg",
			"http://digota.com/image2.jpg",
		},
		Metadata: map[string]string{
			"myAppExtraData": "12345",
		},
		Shippable: true,
		Url:       "http://digota.com",
	}); err != nil {
		t.Fatal(err)
	}

	// validation fail
	if _, err := service.New(context.Background(), &productpb.NewRequest{
		Url: "http://digotacom",
	}); err == nil {
		t.Fatal()
	}

}

func TestProductService_Delete(t *testing.T) {

	// ok
	p, err := service.New(context.Background(), &productpb.NewRequest{
		Name:        fake.Brand(),
		Active:      true,
		Attributes:  []string{"size"},
		Description: fake.Paragraph(),
		Images: []string{
			"http://digota.com/image1.jpg",
			"http://digota.com/image2.jpg",
		},
		Metadata: map[string]string{
			"myAppExtraData": "12345",
		},
		Shippable: true,
		Url:       "http://digota.com",
	})

	if err != nil {
		t.Fatal(err)
	}

	// ok
	if _, err := service.Delete(context.Background(), &productpb.DeleteRequest{
		Id: p.GetId(),
	}); err != nil {
		t.Fatal()
	}

	// validation fail
	if _, err := service.Delete(context.Background(), &productpb.DeleteRequest{
		Id: uuid.NewV4().String(),
	}); err == nil {
		t.Fatal()
	}

}

func TestProductService_Get(t *testing.T) {

	// ok
	p, err := service.New(context.Background(), &productpb.NewRequest{
		Name:        fake.Brand(),
		Active:      true,
		Attributes:  []string{"size"},
		Description: fake.Paragraph(),
		Images: []string{
			"http://digota.com/image1.jpg",
			"http://digota.com/image2.jpg",
		},
		Metadata: map[string]string{
			"myAppExtraData": "12345",
		},
		Shippable: true,
		Url:       "http://digota.com",
	})

	if err != nil {
		t.Fatal(err)
	}

	// ok
	if _, err := service.Get(context.Background(), &productpb.GetRequest{
		Id: p.GetId(),
	}); err != nil {
		t.Fatal()
	}

	// validation fail
	if _, err := service.Get(context.Background(), &productpb.GetRequest{
		Id: "notvaliduuid",
	}); err == nil {
		t.Fatal()
	}

	// not found
	if _, err := service.Get(context.Background(), &productpb.GetRequest{
		Id: uuid.NewV4().String(),
	}); err == nil {
		t.Fatal()
	}

}

func TestProductService_List(t *testing.T) {

	storage.Handler().DropCollection(db, &Product{})

	// create few products
	for _, v := range []string{"one", "two", "three"} {
		_, err := service.New(context.Background(), &productpb.NewRequest{
			Name:        v,
			Active:      true,
			Attributes:  []string{"size"},
			Description: fake.Paragraph(),
			Images: []string{
				"http://digota.com/image1.jpg",
				"http://digota.com/image2.jpg",
			},
			Metadata: map[string]string{
				"myAppExtraData": "12345",
			},
			Shippable: true,
			Url:       "http://digota.com",
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// ok
	list, err := service.List(context.Background(), &productpb.ListRequest{
		Page:  0,
		Limit: 3,
	})

	if err != nil {
		t.Fatal(err)
	}

	if list.Total != 3 || list.Products[0].Name != "one" || list.Products[1].Name != "two" || list.Products[2].Name != "three" {
		t.Fatal()
	}

}
func TestProductService_Update(t *testing.T) {

	// ok
	p, err := service.New(context.Background(), &productpb.NewRequest{
		Name:        fake.Brand(),
		Active:      true,
		Attributes:  []string{"size"},
		Description: fake.Paragraph(),
		Images: []string{
			"http://digota.com/image1.jpg",
			"http://digota.com/image2.jpg",
		},
		Metadata: map[string]string{
			"myAppExtraData": "12345",
		},
		Shippable: true,
		Url:       "http://digota.com",
	})

	if err != nil {
		t.Fatal(err)
	}

	// ok
	_, err1 := service.Update(context.Background(), &productpb.UpdateRequest{
		Id:          p.GetId(),
		Name:        fake.Brand(),
		Active:      true,
		Attributes:  []string{"size"},
		Description: fake.Paragraph(),
		Images: []string{
			"http://digota.com/image1.jpg",
			"http://digota.com/image2.jpg",
		},
		Metadata: map[string]string{
			"myAppExtraData": "12345",
		},
		Shippable: false,
		Url:       "http://digota.com",
	})

	if err1 != nil {
		t.Fatal(err)
	}

	// validation fail
	_, err2 := service.Update(context.Background(), &productpb.UpdateRequest{
		Url: "http://digotacom",
	})

	if err2 == nil {
		t.Fatal(err)
	}

	// not found
	_, err3 := service.Update(context.Background(), &productpb.UpdateRequest{
		Id: uuid.NewV4().String(),
	})

	if err3 == nil {
		t.Fatal(err)
	}

}
