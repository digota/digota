#     Digota <http://digota.com> - eCommerce microservice
#     Copyright (C) 2017  Yaron Sumel <yaron@digota.com>. All Rights Reserved.
#
#     This program is free software: you can redistribute it and/or modify
#     it under the terms of the GNU Affero General Public License as published
#     by the Free Software Foundation, either version 3 of the License, or
#     (at your option) any later version.
#
#     This program is distributed in the hope that it will be useful,
#     but WITHOUT ANY WARRANTY; without even the implied warranty of
#     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#     GNU Affero General Public License for more details.
#
#     You should have received a copy of the GNU Affero General Public License
#     along with this program.  If not, see <http://www.gnu.org/licenses/>.

go:

# clean old *.pb.go
	(rm -f payment/paymentpb/payment.pb.go \
	rm -f order/orderpb/order.pb.go \
	rm -f sku/skupb/sku.pb.go \
	rm -f product/productpb/product.pb.go )

# generate payment pb
	(protoc \
	-I=. \
	-I=../../../ \
	-I=../../gogo/protobuf/protobuf \
	 --gogofast_out=plugins=grpc:../../../ \
	payment/paymentpb/payment.proto)

# generate order pb
	(protoc \
	-I=. \
	-I=../../../ \
	-I=../../gogo/protobuf/protobuf \
	 --gogofast_out=plugins=grpc:../../../ \
	order/orderpb/order.proto)

# generate sku pb
	(protoc \
	-I=. \
	-I=../../../ \
	-I=../../gogo/protobuf/protobuf \
	 --gogofast_out=plugins=grpc:../../../ \
	sku/skupb/sku.proto)

# generate product pb
	(protoc \
	-I=. \
	-I=../../../ \
	-I=../../gogo/protobuf/protobuf \
	 --gogofast_out=plugins=grpc:../../../ \
	product/productpb/product.proto)

php:

# create _php folder
	(mkdir -p _php && rm -rf _php/* )

# generate payment pb
	(protoc \
	-I=. \
	-I=../../../ \
	-I=../../gogo/protobuf/protobuf \
	--php_out=_php \
	--plugin=protoc-gen-grpc=bins/opt/grpc_php_plugin \
	payment/paymentpb/payment.proto)

# generate order pb
	(protoc \
	-I=. \
	-I=../../../ \
	-I=../../gogo/protobuf/protobuf \
	--php_out=_php \
	--plugin=protoc-gen-grpc=bins/opt/grpc_php_plugin \
	order/orderpb/order.proto)

# generate sku pb
	(protoc \
	-I=. \
	-I=../../../ \
	-I=../../gogo/protobuf/protobuf \
	--php_out=_php \
	--plugin=protoc-gen-grpc=bins/opt/grpc_php_plugin \
	sku/skupb/sku.proto)

# generate product pb
	(protoc \
	-I=. \
	-I=../../../ \
	-I=../../gogo/protobuf/protobuf \
	--php_out=_php \
	--plugin=protoc-gen-grpc=bins/opt/grpc_php_plugin \
	product/productpb/product.proto)
