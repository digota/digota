// Digota <http://digota.com> - eCommerce microservice
// Copyright (c) 2018 Yaron Sumel <yaron@digota.com>
//
// MIT License
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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

	// in-memory locker
	locker.New(config.Locker{})

	retCode := m.Run()
	storage.Handler().DropDatabase(db)
	// teardown
	os.Exit(retCode)
}

func TestProducts_GetNamespace(t *testing.T) {
	p := products{}
	if p.GetNamespace() != "product" {
		t.Fatal()
	}
}

func TestProduct_GetNamespace(t *testing.T) {
	p := product{}
	if p.GetNamespace() != "product" {
		t.Fatal()
	}
}

func TestProduct_SetCreated(t *testing.T) {
	p := product{}
	ti := time.Now().Unix()
	p.SetCreated(ti)
	if p.Created != ti {
		t.Fatal()
	}
}

func TestProduct_SetId(t *testing.T) {
	p := product{}
	uid := uuid.NewV4().String()
	p.SetId(uid)
	if p.GetId() != uid {
		t.Fatal()
	}
}

func TestProduct_SetUpdated(t *testing.T) {
	p := product{}
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

	storage.Handler().DropCollection(db, &product{})

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
