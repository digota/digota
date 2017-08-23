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
	"github.com/digota/digota/product/productpb"
	"github.com/digota/digota/sdk"
	"github.com/icrowley/fake"
	"golang.org/x/net/context"
	"log"
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

	// Charge amount
	log.Println(productpb.NewProductServiceClient(c).New(context.Background(), &productpb.NewRequest{
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
	}))

}
