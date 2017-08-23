//     Digota <http://digota.com> - eCommerce microservice
//     Copyright (C) 2017  Yaron Sumel <yaron@digota.com>. All Rights Reserved.
//
//     This program is free software: you can redistribute it and/or modify
//     it under the terms of the GNU Affero General Public License as published
//     by the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU Affero General Public License for more details.
//
//     You should have received a copy of the GNU Affero General Public License
//     along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"github.com/digota/digota/sdk"
	"github.com/digota/digota/sku/skupb"
	"golang.org/x/net/context"
	"log"
	"math/rand"
	"time"
)

func main() {

	c, err := sdk.NewClient("localhost:3051", &sdk.ClientOpt{
		InsecureSkipVerify: false,
		ServerName:         "server.com",
		CaCrt:              "out/ca.crt",
		Crt:                "out/client.com.crt",
		Key:                "out/client.com.key",
	})

	if err != nil {
		panic(err)
	}

	defer c.Close()

	rand.Seed(time.Now().Unix())

	// Charge amount
	log.Println(skupb.NewSkuServiceClient(c).Update(context.Background(), &skupb.UpdateRequest{
		Id: "af350ecc-56c8-485f-8858-74d4faffa9cb",
		//Name:      fake.Brand(),
		Active: true,
		Price:  uint64(rand.Int31n(10001)),
		//Currency:  paymentpb.Currency_USD,
		//Parent: "cb379ae1-8729-4b32-ba7a-3119dc2bd212",
		//Metadata: map[string]string{
		//	"key": "val",
		//},
		//Image: "http://sadf.124.com",
		//PackageDimensions: &skupb.PackageDimensions{
		//	Weight: 1.1,
		//	Length: 1.2,
		//	Height: 1.3,
		//	Width:  1.4,
		//},
		//Inventory: &skupb.Inventory{
		//	Quantity: 1111,
		//	Type:     skupb.Inventory_Finite,
		//},
		//Attributes: map[string]string{
		//	"color": "red",
		//},
	}))

}
