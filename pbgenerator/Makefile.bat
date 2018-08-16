:: Digota <http://digota.com> - eCommerce microservice
:: Copyright (c) 2018 Yaron Sumel <yaron@digota.com>

:: MIT License
:: Permission is hereby granted, free of charge, to any person obtaining a copy
:: of this software and associated documentation files (the "Software"), to deal
:: in the Software without restriction, including without limitation the rights
:: to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
:: copies of the Software, and to permit persons to whom the Software is
:: furnished to do so, subject to the following conditions:

:: The above copyright notice and this permission notice shall be included in all
:: copies or substantial portions of the Software.

:: THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
:: IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
:: FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
:: AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
:: LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
:: OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
:: SOFTWARE.

:: @echo off
cd /d "%GOPATH%\src\github.com\digota\digota"

:: sku
DEL "sku\skupb\sku.pb.go" || pause
DEL "sku\skupb\sku.pb.gw.go" || pause

(protoc ^
-I=. ^
-I=../../../ ^
-I=../../gogo/protobuf/protobuf ^
--gofast_out=plugins=grpc:../../../ ^
sku/skupb/sku.proto || pause)

:: order
DEL "order\orderpb\order.pb.go" || pause
DEL "order\orderpb\order.pb.gw.go" || pause

(protoc ^
-I=. ^
-I=../../../ ^
-I=../../gogo/protobuf/protobuf ^
--gofast_out=plugins=grpc:../../../ ^
order/orderpb/order.proto || pause)

:: payment
DEL "payment\paymentpb\payment.pb.go" || pause
DEL "payment\paymentpb\payment.pb.gw.go" || pause

(protoc ^
-I=. ^
-I=../../../ ^
-I=../../gogo/protobuf/protobuf ^
--gofast_out=plugins=grpc:../../../ ^
payment/paymentpb/payment.proto || pause)

:: product
DEL "product\productpb\product.pb.go" || pause
DEL "product\productpb\product.pb.gw.go" || pause

(protoc ^
-I=. ^
-I=../../../ ^
-I=../../gogo/protobuf/protobuf ^
--gofast_out=plugins=grpc:../../../ ^
product\productpb\product.proto || pause)

:: pause
exit