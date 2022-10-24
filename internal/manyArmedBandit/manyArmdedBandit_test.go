package manyArmedBandit

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	var myArStruct []BannerStruct
	var myRandProc []int
	for i := 1; i <= 10; i++ {
		myArStruct = append(myArStruct, BannerStruct{i, 0, 0})
		myRandProc = append(myRandProc, rand.Intn(100))
	}
	fmt.Println(myArStruct)
	fmt.Println(myRandProc)
	myBandit := New(BanditConfig{100, 100, 10})
	fmt.Print(myBandit)
	require.Equal(t, 1, myBandit.GetBannerId(myArStruct))

}
