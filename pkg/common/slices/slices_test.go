package slices

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapErr(t *testing.T) {
	var errCounter int
	countErr := errors.New("count error")

	errGenerator := func(number int) (res int32, err error) {
		res = int32(number)
		if errCounter == 3 {
			err = countErr
		}
		errCounter++
		return res, err
	}

	fourElemSlice := []int{1, 2, 3, 4}
	resultSlice, err := MapErr(fourElemSlice, errGenerator)

	assert.Error(t, err)
	assert.Equal(t, countErr, err)
	assert.Len(t, resultSlice, 0)

	twoElemSlice := []int{1, 2}
	resultSlice, err = MapErr(twoElemSlice, errGenerator)

	assert.NoError(t, err)
	assert.Len(t, resultSlice, 2)
}

func TestFilter_EmptySlice(t *testing.T) {
	var emptySlice []int

	resultSlice := Filter(emptySlice, func(_ int) bool {
		return true
	})

	assert.Len(t, resultSlice, 0)
}

func TestFilterErr_FilterElement(t *testing.T) {
	fourElemSlice := []int{1, 2, 3, 4}

	fourNumberExcluder := func(number int) (bool, error) { return number != 4, nil }

	resultSlice, err := FilterErr(fourElemSlice, fourNumberExcluder)

	assert.NoError(t, err)
	assert.Len(t, resultSlice, len(fourElemSlice)-1)
}

func TestFilterErr_ReturnErr(t *testing.T) {
	filterErr := errors.New("filter err")

	fourElemSlice := []int{1, 2, 3, 4}

	errGenerator := func(number int) (bool, error) {
		if number == 3 {
			return true, filterErr
		}
		return true, nil
	}

	resultSlice, err := FilterErr(fourElemSlice, errGenerator)

	assert.Error(t, err)
	assert.Equal(t, filterErr, err)
	assert.Len(t, resultSlice, 0)

	twoElemSlice := []int{1, 2}
	resultSlice, err = FilterErr(twoElemSlice, errGenerator)

	assert.NoError(t, err)
	assert.Len(t, resultSlice, 2)
}
