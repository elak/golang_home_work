package hw04_lru_cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {

	toSlice := func(l List) []int {
		elems := make([]int, 0, l.Len())

		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}

		return elems
	}

	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		require.Equal(t, []int{10, 20}, toSlice(l))

		l.PushBack(30) // [10, 20, 30]
		require.Equal(t, 3, l.Len())
		require.Equal(t, []int{10, 20, 30}, toSlice(l))

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())
		require.Equal(t, []int{10, 30}, toSlice(l))

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]
		require.Equal(t, []int{80, 60, 40, 10, 30, 50, 70}, toSlice(l))

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		require.Equal(t, []int{80, 60, 40, 10, 30, 50, 70}, toSlice(l))

		l.MoveToFront(l.Back()) // [70, 80, 60, 40, 10, 30, 50]

		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, toSlice(l))
	})
}
