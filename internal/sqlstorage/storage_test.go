package sqlstorage

import (
	"context"
	"fmt"
	"testing"

	"github.com/PalPalych7/OtusProjectWork/internal/manyArmedBandit"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		myBandit := manyArmedBandit.New(manyArmedBandit.BanditConfig{500, 500, 10})
		storage := New(context.Background(), "otusfinalproj", "testuser", "123456", myBandit)

		// DB connect
		err := storage.Connect()
		require.NoError(t, err)

		err2, bammerId := storage.GetBannerForSlot(1, 1)
		require.NoError(t, err2)
		fmt.Println(bammerId)
		/*
			err2 := storage.DelBannerSlot(2, 2)
			require.NoError(t, err2)

			err3 := storage.BannerClick(1, 1, 1)
			require.NoError(t, err3)
		*/
		// close DB
		err = storage.Close()
		require.NoError(t, err)

		/*
			// Get Event By Date
			//		myEventList, err2 := storage.GetEventByDate("2022-05-11")
			myEventList, err2 := storage.GetEventByDate("11.05.2022")
			require.NoError(t, err2)
			len1 := len(myEventList)

			// new event
			err = storage.CreateEvent("t1", "11.05.2022", "something", 1)
			require.NoError(t, err)

			// Get Event By Date2
			myEventList, err2 = storage.GetEventByDate("11.05.2022")
			require.NoError(t, err2)
			len2 := len(myEventList)
			fmt.Println(len2)
			require.Equal(t, len1+1, len2)

		*/
	})
}
