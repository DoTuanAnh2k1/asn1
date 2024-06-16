package main

import (
	"encoding/asn1"
)

type ExampleData struct {
	booleanValue     bool           `asn1:""`
	integerValue     int            `asn1:""`
	bitstringValue   asn1.BitString `asn1:""`
	octetstringValue []byte         `asn1:""`
	nullValue        asn1.RawValue  `asn1:""`
}
type ExampleSequence struct {
	sequenceField1 int    `asn1:""`
	sequenceField2 string `asn1:""`
}
type ExampleSet struct {
	setField1 int    `asn1:""`
	setField2 string `asn1:""`
}
type ExampleChoice struct {
	choiceField1 int    `asn1:""`
	choiceField2 string `asn1:""`
}
