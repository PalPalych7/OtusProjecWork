package manyarmedbandit

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func isClick(myRandProc int) int {
	var res int
	curRand := rand.Intn(101) //nolint

	if myRandProc >= curRand {
		res = 1
	}
	return res
}

func TestMABandit(t *testing.T) {
	var myArStruct []BannerStruct
	var myRandProc []int // сгенерированная вероятность кликоа на баннер
	genCount := 50000    // количество запросов
	bannerCount := 50    // кол-во баннеров

	rand.Seed(time.Now().UTC().UnixNano())
	var minProc float32 = 100
	var maxProc float32

	myBandit := New(BanditConfig{250, 500, 10})

	for i := 1; i <= bannerCount; i++ { // генерим вероятность клика для каждого баннера
		myArStruct = append(myArStruct, BannerStruct{i, 0, 0})
		myRandProc = append(myRandProc, rand.Intn(101)) //nolint
	}
	for i := 1; i <= genCount; i++ { // вызов метода заданное кол-во раз
		curNum := myBandit.GetBannerNum(myArStruct)
		myArStruct[curNum].ShowCount++
		if isClick(myRandProc[curNum]) == 1 { // определения кликнули ли по баннеру на основании сгенерированной вероятности
			myArStruct[curNum].ClickCount++ // увеличение счётчмка кликов
		}
	}
	fmt.Println("myArStruct=", myArStruct)
	for i := 0; i < bannerCount; i++ { // определение процента показа для самого популярного и самого редкого баннера
		if float32(myArStruct[i].ShowCount)/float32(genCount)*100 > maxProc {
			maxProc = float32(myArStruct[i].ShowCount) / float32(genCount) * 100
		}
		if float32(myArStruct[i].ShowCount)/float32(genCount)*100 < minProc {
			minProc = float32(myArStruct[i].ShowCount) / float32(genCount) * 100
		}
	}
	fmt.Println("minProc=", minProc)
	fmt.Println("maxProc=", maxProc)
	// максимальный должен показываться минимум в 2 раза чаще среднестатистического (для 50 баннеров >4%)
	require.LessOrEqual(t, 100/float32(bannerCount)*2, maxProc)
	// минимальный должен показываться минимум в 2 раза реже среднестатистического (для 50 баннеров <1%)
	require.LessOrEqual(t, minProc, 100/float32(bannerCount)/2)
	// минимальный должен показываться чаще чем 1/20 от среднестатистического (для 50 баннеров >0,1%)
	require.LessOrEqual(t, 100/float32(bannerCount)/100, minProc)
}
