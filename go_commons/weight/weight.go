// Package weight provides utilities for handling weights and unit conversions.
package weight

import (
	error2 "github.com/omniful/go_commons/error"
	"golang.org/x/exp/maps"
)

type (
	UOM string
)

const (
	Kg   UOM = "kg"
	Lbs  UOM = "lbs"
	G    UOM = "g"
	L    UOM = "l"
	Ml   UOM = "ml"
	EA   UOM = "ea"
	Pack UOM = "pack"
	Oz   UOM = "oz"
)

type MinimumSellableWeight struct {
	UOM   UOM
	Value uint64
}

type MinimumSellableWeightInfo struct {
	UOM        UOM
	multiplier float64
	isUnitVal  bool
}

type ConversionInfo struct {
	multiplier float64
}

type Weight struct {
	UOM   UOM     `json:"uom"`
	Value float64 `json:"value"`
}

type WeightInfo struct {
	minimumSellableWeight MinimumSellableWeightInfo
	conversionInfo        map[UOM]ConversionInfo
	isWeighed             bool
}

var weightInfoMap = map[UOM]WeightInfo{
	Kg: {
		minimumSellableWeight: MinimumSellableWeightInfo{UOM: G, multiplier: 1000},
		conversionInfo: map[UOM]ConversionInfo{
			G:   {multiplier: 1000},
			Lbs: {multiplier: 2.20462},
			Oz:  {multiplier: 35.27396},
		},
		isWeighed: true,
	},
	Lbs: {
		minimumSellableWeight: MinimumSellableWeightInfo{UOM: G, multiplier: 453.59237},
		conversionInfo: map[UOM]ConversionInfo{
			G:  {multiplier: 453.59237},
			Kg: {multiplier: 0.453592},
			Oz: {multiplier: 16},
		},
		isWeighed: true,
	},
	G: {
		minimumSellableWeight: MinimumSellableWeightInfo{UOM: G, multiplier: 1},
		conversionInfo: map[UOM]ConversionInfo{
			Kg:  {multiplier: 0.001},
			Lbs: {multiplier: 0.00220462},
			Oz:  {multiplier: 0.03527396},
		},
		isWeighed: true,
	},
	L: {
		minimumSellableWeight: MinimumSellableWeightInfo{UOM: Ml, multiplier: 1000},
		conversionInfo: map[UOM]ConversionInfo{
			Ml: {multiplier: 1000},
		},
		isWeighed: true,
	},
	Ml: {
		minimumSellableWeight: MinimumSellableWeightInfo{UOM: Ml, multiplier: 1},
		conversionInfo: map[UOM]ConversionInfo{
			L: {multiplier: 0.001},
		},
		isWeighed: true,
	},
	EA: {
		minimumSellableWeight: MinimumSellableWeightInfo{UOM: EA, isUnitVal: true},
		conversionInfo: map[UOM]ConversionInfo{
			Pack: {multiplier: 1},
		},
	},
	Pack: {
		minimumSellableWeight: MinimumSellableWeightInfo{UOM: Pack, isUnitVal: true},
		conversionInfo: map[UOM]ConversionInfo{
			EA: {multiplier: 1},
		},
	},
	Oz: {
		minimumSellableWeight: MinimumSellableWeightInfo{UOM: G, multiplier: 28.3495},
		conversionInfo: map[UOM]ConversionInfo{
			G:   {multiplier: 28.3495},
			Kg:  {multiplier: 0.0283495},
			Lbs: {multiplier: 0.0625},
		},
		isWeighed: true,
	},
}

func ConvertToMinimumSellableWeight(uom UOM, weight float64, quantity uint64) (w MinimumSellableWeight, err error2.CustomError) {
	data, ok := weightInfoMap[uom]

	if !ok {
		err = error2.NewCustomError(error2.BadRequestError, "INVALID UOM")
		return
	}

	wInfo := data.minimumSellableWeight

	w.UOM = wInfo.UOM

	if wInfo.isUnitVal {
		w.Value = quantity
		return
	}

	w.Value = uint64(weight * float64(quantity) * wInfo.multiplier)

	return
}

func ConvertUOM(weight Weight, conversionUOM UOM) (w Weight, err error2.CustomError) {
	if weight.UOM == conversionUOM {
		w = weight
		return
	}

	data, ok := weightInfoMap[weight.UOM]

	if !ok {
		err = error2.NewCustomError(error2.BadRequestError, "INVALID UOM")
		return
	}

	if !conversionUOM.In(maps.Keys(data.conversionInfo)) {
		err = error2.NewCustomError(error2.BadRequestError, "INVALID UOM")
		return
	}

	w = weight

	conversion := data.conversionInfo[conversionUOM]

	w.Value = weight.Value * conversion.multiplier
	w.UOM = conversionUOM
	return
}

func IsUOMWeighted(input UOM) bool {
	data, ok := weightInfoMap[input]
	if !ok {
		return false
	}
	return data.isWeighed
}

func IsValidUOM(inputUOM UOM, skuUOM UOM) bool {
	data, ok := weightInfoMap[inputUOM]
	if !ok {
		return false
	}
	if inputUOM != skuUOM && !skuUOM.In(maps.Keys(data.conversionInfo)) {
		return false
	}
	return true
}

func GetMinimumSellableUOM(input UOM) UOM {
	data, ok := weightInfoMap[input]
	if !ok {
		return EA
	}
	return data.minimumSellableWeight.UOM
}

func (uom UOM) In(uoms []UOM) bool {
	for _, u := range uoms {
		if u == uom {
			return true
		}
	}
	return false
}
