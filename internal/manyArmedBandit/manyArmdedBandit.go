package manyArmedBandit

import (
	"math"
	"math/rand"
	"time"
)

type BanditConfig struct {
	FullLearnigCount     int // количество запросов в режиме "полного обучения"
	PartialLearningCount int // количество запросов в режиме "чаcтичного обучения"
	FinalRandomPecent    int // вероятность случайного выбора после обучения (в процентах)
}

type BannerStruct struct {
	BannerId   int
	ShowCount  int
	ClickCount int
}

type MyBandit interface {
	GetBannerNum(arrStruct []BannerStruct) int
}

type banditStruct struct {
	BanditConfig BanditConfig
	Сend         float32 // цена деления (на какой процент уменьшаем "случайную величину" за 1 показ)
}

func kvadrProc(arrStruct []BannerStruct) int {
	var arrSumKvProc []float64
	var curKvProc float64
	var sumKvProc float64 = 0
	for _, v := range arrStruct {
		curKvProc = math.Pow(float64(v.ClickCount)/float64(v.ShowCount)*100, 2)
		sumKvProc = sumKvProc + curKvProc
		arrSumKvProc = append(arrSumKvProc, sumKvProc)
	}
	sumKvProcInt := int(sumKvProc)
	randVal := rand.Intn(sumKvProcInt + 1)
	res := 0
	for i, v := range arrSumKvProc {
		if float64(randVal) <= v {
			res = i
			return res
		}
	}
	return res
}

func (b banditStruct) GetBannerNum(arrStruct []BannerStruct) int {
	rand.Seed(time.Now().UTC().UnixNano())
	showSum := 0
	var res int
	for _, v := range arrStruct {
		showSum = showSum + v.ShowCount
	}
	if showSum <= b.BanditConfig.FullLearnigCount { // режим обучения
		res = rand.Intn(len(arrStruct))
	} else {
		rand100 := rand.Intn(101) // случайная величина для определения алгоритма
		var randomPecent float32
		if showSum <= b.BanditConfig.FullLearnigCount+b.BanditConfig.PartialLearningCount {
			// вычисляем "вероятностный процент" линейно от 100 до минимального
			randomPecent = float32(100) - float32(showSum-b.BanditConfig.FullLearnigCount)*b.Сend
		} else {
			// "вероятностный процент" - минимальный из конфига
			randomPecent = float32(b.BanditConfig.FinalRandomPecent)
		}
		if float32(rand100) < randomPecent {
			res = rand.Intn(len(arrStruct))
		} else {
			// подбор согласно квадратичным вероятностям просмотров
			res = kvadrProc(arrStruct)
		}
	}
	return res
}

func New(bc BanditConfig) MyBandit {
	cend := float32(100-bc.FinalRandomPecent) / float32(bc.PartialLearningCount)
	return banditStruct{bc, cend}
}
