::     Digota <http://digota.com> - eCommerce microservice
::     Copyright (C) 2017  Yaron Sumel <yaron@digota.com>. All Rights Reserved.
::
::     This program is free software: you can redistribute it and/or modify
::     it under the terms of the GNU Affero General Public License as published
::     by the Free Software Foundation, either version 3 of the License, or
::     (at your option) any later version.
::
::     This program is distributed in the hope that it will be useful,
::     but WITHOUT ANY WARRANTY; without even the implied warranty of
::     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
::     GNU Affero General Public License for more details.
::
::     You should have received a copy of the GNU Affero General Public License
::     along with this program.  If not, see <http://www.gnu.org/licenses/>.

::     pb.bat is protoBuff generator helper for windows, will pause and
::     echo any error. Usage: "> pb.bat"

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