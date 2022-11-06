//go:build integration
// +build integration

package integration_test

import (
	"fmt"
	"testing"
)

func TestService(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		fmt.Println("integration test statrt")
	})
}

/*
 //+build integration
*/
