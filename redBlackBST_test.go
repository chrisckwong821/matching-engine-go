package main

import (
	"math/rand"
	"testing"
)

func TestRedBlackEmpty(t *testing.T) {
	rb := NewRedBlackBST()
	if rb.Size() != 0 || !rb.IsEmpty() {
		t.Errorf("Red Black BST should be empty")
	}
}

func TestRedBlackBasic(t *testing.T) {
	st := NewRedBlackBST()
	keys := make([]float32, 0)
	for i := 0; i < 10; i += 1 {
		k := rand.Float32()
		keys = append(keys, k)
		st.Put(k, nil)
	}

	if st.Size() != 10 {
		t.Errorf("size should equals 10, got %d", st.Size())
	}
	if st.IsEmpty() {
		t.Errorf("st should not be empty")
	}

	for _, k := range keys {
		if !st.Contains(k) {
			t.Errorf("st should contain the key %0.8f", k)
		}
	}
}

func TestRedBlackHeight(t *testing.T) {
	st := NewRedBlackBST()
	n := 100000
	for i := 0; i < n; i += 1 {

		st.Put(float32(i), nil)
	}

	if st.Size() != n {
		t.Errorf("size should equals %d, got %d", n, st.Size())
	}
	if st.IsEmpty() {
		t.Errorf("st should not be empty")
	}

	height := st.Height()
	if height < 17 || height > 34 {
		t.Errorf("red black bst height should be in range lgN <= height <= 2*lgN, in our case from 17 to 34, but we got %d", height)
	}
}

func TestRedBlackMinMax(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 10; i += 1 {
		st.Put(float32(10-i), nil)
	}

	min := float32(1)
	if st.Min() != min {
		t.Errorf("min %0.8f != %0.8f", st.Min(), min)
	}

	max := float32(10)
	if st.Max() != max {
		t.Errorf("min %0.8f != %0.8f", st.Max(), max)
	}
}

func TestRedBlackMinMaxCachedOnDelete(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 100; i += 1 {
		st.Put(float32(100-i), nil)
	}

	min := float32(1)
	if st.Min() != min {
		t.Errorf("min %0.8f != %0.8f", st.Min(), min)
	}

	max := float32(100)
	if st.Max() != max {
		t.Errorf("min %0.8f != %0.8f", st.Max(), max)
	}

	st.DeleteMin()
	st.DeleteMin()
	for i := 3; i < 20; i += 1 {
		st.Delete(float32(i))
	}
	st.DeleteMax()
	st.DeleteMax()
	for i := 98; i > 70; i -= 1 {
		st.Delete(float32(i))
	}

	min = float32(20)
	if st.Min() != min {
		t.Errorf("min %0.8f != %0.8f", st.Min(), min)
	}

	max = float32(70)
	if st.Max() != max {
		t.Errorf("min %0.8f != %0.8f", st.Max(), max)
	}
}

func TestRedBlackFloor(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 10; i += 1 {
		k := float32(20 - 2*i)
		st.Put(k, nil)
	}

	keymiss := float32(3)
	flmiss := float32(2)
	if st.Floor(keymiss) != flmiss {
		t.Errorf("floor != %0.8f", st.Floor(keymiss))
	}

	keyhit := float32(10)
	if st.Floor(keyhit) != keyhit {
		t.Errorf("floor != %0.8f", st.Floor(keyhit))
	}
}

func TestRedBlackCeiling(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 10; i += 1 {
		k := float32(20 - 2*i)
		st.Put(k, nil)
	}

	keymiss := float32(3)
	clmiss := float32(4)
	if st.Ceiling(keymiss) != clmiss {
		t.Errorf("ceiling != %0.8f", st.Ceiling(keymiss))
	}

	keyhit := float32(10)
	if st.Ceiling(keyhit) != keyhit {
		t.Errorf("ceiling != %0.8f", st.Ceiling(keyhit))
	}
}

func TestRedBlackSelect(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 10; i += 1 {
		k := float32(10 - i)
		st.Put(k, nil)
	}

	key := float32(3)
	if st.Select(2) != key {
		t.Errorf("element with rank=2 should be %0.8f", key)
	}

	key = float32(10)
	if st.Select(9) != key {
		t.Errorf("element with rank=9 should be %0.8f", key)
	}
}

func TestRedBlackRank(t *testing.T) {
	st := NewRedBlackBST()
	keys := make([]float32, 0)
	for i := 0; i < 10; i += 1 {
		k := float32(10 - i)
		keys = append(keys, k)
		st.Put(k, nil)
	}

	for i := range keys {
		k := st.Select(i)
		if st.Rank(k) != i {
			t.Errorf("rank of %0.8f != %d", k, i)
		}
	}

	if st.Rank(11) != len(keys) {
		t.Errorf("rank of new maximum should equal to the number of nodes in the tree")
	}

	if st.Rank(11) != st.Rank(12) {
		t.Errorf("rank of new maximum should not depend on the new maximum concrete value")
	}
}

func TestRedBlackKeys(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 10; i += 1 {
		k := float32(10 - i)
		st.Put(k, nil)
	}

	lo := float32(3)
	hi := float32(6)
	keys := st.Keys(lo, hi)
	if len(keys) != 4 {
		t.Errorf("keys len should equals 4, %+v", keys)
	}

	if keys[0] != lo {
		t.Errorf("first key should be %0.8f", lo)
	}

	if keys[len(keys)-1] != hi {
		t.Errorf("last key should be %0.8f", hi)
	}

	for i := 1; i < len(keys); i += 1 {
		if keys[i] < keys[i-1] {
			t.Errorf("non-decreasing keys order validation failed")
		}
	}
}

func TestRedBlackDeleteMin(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 10; i += 1 {
		k := float32(10 - i)
		st.Put(k, nil)
	}

	st.DeleteMin()
	if st.Size() != 9 {
		t.Errorf("tree size should shrink")
	}

	if st.Contains(1) {
		t.Errorf("minimum element should be removed from the tree")
	}

	if !st.IsRedBlack() {
		t.Errorf("certification failed")
	}
}

func TestRedBlackDeleteMax(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 10; i += 1 {
		k := float32(i)
		st.Put(k, nil)
	}

	st.DeleteMax()
	if st.Size() != 9 {
		t.Errorf("tree size should shrink")
	}

	if st.Contains(9.0) {
		t.Errorf("minimum element should be removed from the tree")
	}

	if !st.IsRedBlack() {
		t.Errorf("certification failed")
	}
}

func TestRedBlackDelete(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 10; i += 1 {
		k := float32(i)
		st.Put(k, nil)
	}

	key := float32(5)
	st.Delete(key)
	if st.Size() != 9 {
		t.Errorf("tree size should shrink")
	}

	if st.Contains(key) {
		t.Errorf("minimum element should be removed from the tree")
	}

	if !st.IsRedBlack() {
		t.Errorf("certification failed")
	}
}

func TestRedBlackPutLinkedListOrder(t *testing.T) {
	st := NewRedBlackBST()
	for i := 0; i < 100; i += 1 {
		k := rand.Float32()
		st.Put(k, nil)
	}

	min := st.MinPointer()
	for p := min; p != nil && p.Next != nil; p = p.Next {
		if p.Next.Key < p.Key {
			t.Errorf("incorrect keys order")
			break
		}
	}
}

func TestRedBlackPutDeleteLinkedListOrder(t *testing.T) {
	st := NewRedBlackBST()
	n := 1000
	for i := 0; i < n; i += 1 {
		k := rand.Float32()
		st.Put(k, nil)
	}

	// deleting from both ends and in the middle 90% of the nodes
	k := int(float32(n) * 0.3)
	for i := 0; i < k; i += 1 {
		st.DeleteMin()
		k := st.Select(rand.Intn(st.Size()))
		st.Delete(k)
		st.DeleteMax()
	}

	if st.Size() != n-3*k {
		t.Errorf("incorrect tree size %d", st.Size())
	}

	min := st.MinPointer()
	for p := min; p != nil && p.Next != nil; p = p.Next {
		if p.Next.Key < p.Key {
			t.Errorf("incorrect keys order")
			break
		}
	}
}
