#!/bin/bash

cd cli
go build .

echo "RUNNING PROGRAM WITH FIRST DATASET (AIRPORTS)..."
./cli -d ../data/dataset1.csv -f 1

echo "RUNNING PROGRAM WITH FIRST DATASET (CITIES)..."
./cli -d ../data/dataset2.csv -f 2

rm cli
cd ../