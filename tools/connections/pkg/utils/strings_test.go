package utils

import "testing"

func TestRandomString(t *testing.T) {
	type testCase struct {
		length     int
		wantlength int
	}
	var tcs = []testCase{
		{-20, 0},
		{-1, 0},
		{0, 0},
		{1, 1},
		{20, 20},
	}
	for _, tc := range tcs {
		if len(RandomString(tc.length)) != tc.wantlength {
			t.Errorf("len(RandomString(%d)) != %d", tc.length, tc.wantlength)
			t.Fail()
		}
	}
}
