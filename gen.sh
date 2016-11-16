#!/bin/sh

python -m peachpy.x86_64 decode.py -S -o decode_amd64.s -mabi=goasm
# peachpy doesn't quite support what we need, so we modify the output here
sed -i 's/dst_len+8(FP)/·sseMasks+0(SB)/' decode_amd64.s
sed -i '
1i// +build !noasm\n
' decode_amd64.s

# generate the tables we need
go run gen.go -table=ssemasks |gofmt >ssemasks.go
go run gen.go -table=bytesused |sed 's/ $//' >bytesused.go
