//go:build integration || ignore || (тест && ignore) || только || при || поднятой || БД
// +build integration ignore тест,ignore только при поднятой БД

package sqlstorage

import (
	"context"
	"fmt"
	"testing"

	manyarmedbandit "github.com/PalPalych7/OtusProjectWork/internal/manyArmedBandit"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		myBandit := manyarmedbandit.New(manyarmedbandit.BanditConfig{500, 500, 10})
		storage := New(context.Background(), "otusfinalproj", "testuser", "123456", myBandit)

		// DB connect
		err := storage.Connect()
		require.NoError(t, err)

		bammerId, err := storage.GetBannerForSlot(1, 1)
		require.NoError(t, err)
		require.Greater(t, bammerId, 0)
		fmt.Println(bammerId)

		err2 := storage.AddBannerSlot(2, 2)
		require.NoError(t, err2)

		err3 := storage.DelBannerSlot(2, 2)
		require.NoError(t, err3)

		err4 := storage.BannerClick(1, 1, 1)
		require.NoError(t, err4)

		// close DB
		err = storage.Close()
		require.NoError(t, err)
	})
}
