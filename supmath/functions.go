package supmath

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

const (
	standardDicreteName = "standard"
	elitistDicreteName  = "elitist"
)

type Binarizer struct {
	transfer func([]float64)
	discrete func([]float64, interface{})
}

func NewBinarizer(transfer func([]float64), discrete func([]float64, interface{})) *Binarizer {
	return &Binarizer{transfer: transfer, discrete: discrete}
}

func (binarizer *Binarizer) Binarization(xVec []float64, items interface{}) {
	binarizer.transfer(xVec)
	binarizer.discrete(xVec, items)
}

func (binarizer *Binarizer) GetDiscrete() func([]float64, interface{}) {
	return binarizer.discrete
}

func GetTransferByStr(str string) (transfer func([]float64)) {
	str = strings.ToLower(str)
	if str == "s1" {
		transfer = TransferS1
	} else if str == "s2" {
		transfer = TransferS2
	} else if str == "s3" {
		transfer = TransferS3
	} else if str == "s12" {
		transfer = TransferS12
	} else if str == "s13" {
		transfer = TransferS13
	} else if str == "v1" {
		transfer = TransferV1
	} else if str == "v2" {
		transfer = TransferV2
	} else if str == "v3" {
		transfer = TransferV3
	} else if str == "v4" {
		transfer = TransferV4
	}
	return
}

func GetDiscreteByStr(name string) func([]float64, interface{}) {
	if strings.ToLower(name) == standardDicreteName {
		return StandardDiscrete
	} else if strings.ToLower(name) == elitistDicreteName {
		return ElitistDiscrete
	} else {
		panic("Error param for GetDiscreteByStr")
	}
}

func TransferV1(xVec []float64) {
	sq2 := 1.41421356237
	for i := 0; i < len(xVec); i++ {
		xVec[i] = math.Abs(math.Erf(sq2 * xVec[i] / math.Pi))
	}
}

func TransferV2(xVec []float64) {
	for i := 0; i < len(xVec); i++ {
		xVec[i] = math.Abs(math.Tanh(xVec[i]))
	}
}

func TransferV3(xVec []float64) {
	for i := 0; i < len(xVec); i++ {
		xVec[i] = math.Abs(xVec[i] / (math.Sqrt(1 + xVec[i]*xVec[i])))
	}
}

func TransferV4(xVec []float64) {
	for i := 0; i < len(xVec); i++ {
		xVec[i] = math.Abs((2.0 / math.Pi) * (math.Atan(xVec[i] * math.Pi / 2.0)))
	}
}

func TransferS1(xVec []float64) {
	for i := 0; i < len(xVec); i++ {
		xVec[i] = 1.0 / (1.0 + math.Exp(-1.0*xVec[i]))
	}
}

func TransferS2(xVec []float64) {
	for i := 0; i < len(xVec); i++ {
		xVec[i] = 1.0 / (1.0 + math.Exp(-2.0*xVec[i]))
	}
}

func TransferS3(xVec []float64) {
	for i := 0; i < len(xVec); i++ {
		xVec[i] = 1.0 / (1.0 + math.Exp(-3.0*xVec[i]))
	}
}

func TransferS12(xVec []float64) {
	for i := 0; i < len(xVec); i++ {
		xVec[i] = 1.0 / (1.0 + math.Exp(-xVec[i]/2.0))
	}
}

func TransferS13(xVec []float64) {
	for i := 0; i < len(xVec); i++ {
		xVec[i] = 1.0 / (1.0 + math.Exp(-xVec[i]/3.0))
	}
}

func CalcEuclidDist(xVec []float64, yVec []float64) (result float64, ok bool) {
	if len(xVec) != len(yVec) {
		ok = false
		return
	}
	ok = true
	for i := 0; i < len(xVec); i++ {
		result += math.Pow(xVec[i]-yVec[i], 2.0)
	}
	result = math.Sqrt(result)
	return
}

func CalcManhattanDist(xVec []float64, yVec []float64) (result float64, ok bool) {
	if len(xVec) != len(yVec) {
		ok = false
		return
	}
	ok = true

	for i := 0; i < len(xVec); i++ {
		result += math.Abs(xVec[i] - yVec[i])
	}
	return
}

func CalcMinkDist(xVec []float64, yVec []float64, p float64) (result float64, ok bool) {
	if len(xVec) != len(yVec) {
		ok = false
		return
	}
	ok = true
	for i := 0; i < len(xVec); i++ {
		result += math.Pow(xVec[i]-yVec[i], p)
	}
	result = math.Pow(result, 1.0/p)
	return
}

// Rand.Seed(time.Now().UnixNano())
func RandFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}

func RandBinFloat(n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = float64(rand.Intn(2))
	}
	return res
}

func StandardDiscrete(xVec []float64, items interface{}) {
	randVec := RandFloats(0.0, 1.0, len(xVec))
	for i := range randVec {
		if xVec[i] >= randVec[i] {
			xVec[i] = 1.0
		} else {
			xVec[i] = 0.0
		}
	}
}

func ElitistDiscrete(xVec []float64, items interface{}) {
	randVec := RandFloats(0.0, 1.0, len(xVec))
	switch value := items.(type) {
	case []float64:
		for i := range randVec {
			if xVec[i] >= randVec[i] {
				xVec[i] = value[i]
			} else {
				xVec[i] = 0.0
			}
		}
	default:
		fmt.Println("Problem, bad parameter [items]", fmt.Sprintf("%T", value))
	}
}

/* Binarization , parameters by yourself and copy
For complement:
	arr := []int{1, 2, 3}
	tmp := make([]int, len(arr))
	copy(tmp, arr)
	tmp - old version
For elistic:
	best solution
*/
func Binarization(xVec []float64, transfer func([]float64),
	discrete func([]float64, ...interface{}), items ...interface{}) {
	transfer(xVec)
	discrete(xVec, items...)
}

func MinFloatSlice(xVec []float64) (ind int, res float64) {
	res = xVec[0]
	for i, e := range xVec {
		if e < res {
			res = e
			ind = i
		}
	}
	return
}

func MaxFloatSlice(xVec []float64) (ind int, res float64) {
	res = xVec[0]
	for i, e := range xVec {
		if e > res {
			res = e
			ind = i
		}
	}
	return
}

func SumFloatSlice(xVec []float64) (res float64) {
	for i := range xVec {
		res += xVec[i]
	}
	return
}

func MeanFloat64(xVec []float64) float64 {
	res := 0.0

	for _, e := range xVec {
		res += e
	}
	return res / float64(len(xVec))
}

func MeanInt64(xVec []int64) int64 {
	var res int64 = 0

	for _, e := range xVec {
		res += e
	}
	return res / int64(len(xVec))
}

func MaxFloat64(xVec []float64) float64 {
	res := xVec[0]
	for _, e := range xVec {
		if e > res {
			res = e
		}
	}
	return res
}

func MinFloat64(xVec []float64) float64 {
	res := xVec[0]
	for _, e := range xVec {
		if e < res {
			res = e
		}
	}
	return res
}

//MSDFloat64 calc standart deviation
func MSDFloat64(xVec []float64, mean float64) (sigma float64) {
	for i := range xVec {
		sigma += (xVec[i] - mean) * (xVec[i] - mean)
	}
	sigma /= float64(len(xVec))
	sigma = math.Sqrt(sigma)
	return
}

func IntRange(min, max int) int {
	if min == max {
		return min
	}
	return rand.Intn(max-min) + min
}

func SimpleMutation(xVec []float64, count int) {
	for i := 0; i < count; i++ {
		ind := IntRange(0, len(xVec))
		if xVec[ind] == 0.0 {
			xVec[ind] = 1.0
		} else {
			xVec[ind] = 0.0
		}
	}
}

func ShuffleMutation(xVec []float64, _ int) {
	rand.Shuffle(len(xVec), func(first, second int) {
		xVec[first], xVec[second] = xVec[second], xVec[first]
	})
}

func LimitMutation(xVec []float64, count int, limit int) {
	for i := 0; i < count; i++ {
		ind := IntRange(0, limit)
		if xVec[ind] == 0.0 {
			xVec[ind] = 1.0
		} else {
			xVec[ind] = 0.0
		}
	}
}

func RandBinLimit(size int, count int) []float64 {
	count = IntRange(count-int(float64(count)*0.1), count+int(float64(count)*0.1))
	res := make([]float64, size, size)
	for i := 0; i < count; i++ {
		ind := IntRange(0, len(res))
		res[ind] = 1.0
	}
	return res
}

func NormilizeFloatSliceMean(vec []float64) {
	xMean := MeanFloat64(vec)
	for i := range vec {
		vec[i] /= xMean
	}
}

func NormilizeFloatSliceMax(vec []float64) {
	xMax := MaxFloat64(vec)
	for i := range vec {
		vec[i] /= xMax
	}
}
