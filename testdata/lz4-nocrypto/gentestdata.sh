"Mdp��  �(#!/bin/bash

ROOT_DIR=$(cd $(dirname $0); cd ..; pwd)

 @ROOT. �
pushd $ �/cmd/brbundle
go build
popd

KEY=12345678 H

./D A/brbM � --help

rm -rf testdata/br* Mlz4* N*aes nchacha �	.zip

# content-folder
#� Tz raw/ M     Atest� �raw-nocrypto� Draw- G
./cJB/brb@     pcontent I     Atesth Dbr-ng B tes Draw- O
./ch -z lz4 � Dlz4-R O tesh Eraw h �-c AES -k ${KEY}� Atest� praw-aes G  te @nocrNO
./c� @     @cont8O-c Ah Dbr-ag F tes~ Onocrh 
H-z l8O-c Ah Elz4-i Etest� Onocrh Hraw 8gchacha;EtestR chacha 80Ochach Bbr-c H   t~ Onocr� Olz4 � Elz4-i Etest� Anocrh �
# embeddede-z raw( M     �-p rawnoenc  -�A/embk O.go � F./cmI/brb�@embe� M     S-p br� E -o b Cbr-n� A/embB F.go " Draw-# O
./c� o-z lz4Nlz4nDlz4-b E/emb� Ntest�O./cmH-z r�L-c A8c-p raw1H-o tb Haes/ J    " AnocrO./cm� f      O-c A�  Bbrae�D-o t� Obr-a� O tesDlz4 � O -c � Blz4a� D-o tb Olz4-� Otest�*Ichac�d-p raw Go teb kchacha�Otest� K    �Lchac� Abrch E -o b Abr-c E/emb�O   t� K-z l�Lchac� Clz4c Co te� Blz4- G/emb� Ntest_
#zip�	�raw zip-8E     MtestZ F.zip� Oraw-Ozip-n Cbr-n[O.zipf$Ozip-n Ilz4-o Mtest� O
./c�J-z rJL-c A�EtestT Aaes.� I     OnocrRCzip-� O-c An Hbr-am F tes� Onocrn 
h-z lz4&O-c An Ilz4-o Etest� Onocrn Jraw � Ichac EtestT cchachaMEtest Onocrn 
o      n Hbr-cm O tesn M-z lJOchac� Blz4- K.zip&�-nocrypto

