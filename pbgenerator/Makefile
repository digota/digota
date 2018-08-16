# Digota <http://digota.com> - eCommerce microservice
# Copyright (c) 2018 Yaron Sumel <yaron@digota.com>
#
# MIT License
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

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
