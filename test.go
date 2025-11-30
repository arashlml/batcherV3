package batcherV3

import (
	"testing"
	"time"
)

// ------------------------------
// ✅ Test 1: Flush by Size
// ------------------------------

func TestBatcher_FlushBySize(t *testing.T) {
	b := NewBatcher[int](3, time.Second)

	b.Add(1)
	b.Add(2)
	b.Add(3)

	select {
	case batch := <-b.OutChan():
		if len(batch) != 3 {
			t.Fatalf("expected batch size 3, got %d", len(batch))
		}
		if batch[0] != 1 || batch[1] != 2 || batch[2] != 3 {
			t.Fatalf("unexpected batch data: %v", batch)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for batch flush by size")
	}

	b.Close()
}

// ------------------------------
// ✅ Test 2: Flush by Time
// ------------------------------

func TestBatcher_FlushByTime(t *testing.T) {
	b := NewBatcher[int](10, 200*time.Millisecond)

	b.Add(10)
	b.Add(20)

	select {
	case batch := <-b.OutChan():
		if len(batch) != 2 {
			t.Fatalf("expected batch size 2, got %d", len(batch))
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for batch flush by time")
	}

	b.Close()
}

// ------------------------------
// ✅ Test 3: Close triggers Flush
// ------------------------------

func TestBatcher_CloseFlushesRemaining(t *testing.T) {
	b := NewBatcher[int](10, time.Second)

	b.Add(100)
	b.Add(200)

	b.Close()

	select {
	case batch := <-b.OutChan():
		if len(batch) != 2 {
			t.Fatalf("expected batch size 2 after close, got %d", len(batch))
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for batch flush on close")
	}
}

// ------------------------------
// ✅ Test 4: Multiple Batches
// ------------------------------

func TestBatcher_MultipleFlushes(t *testing.T) {
	b := NewBatcher[int](2, time.Second)

	b.Add(1)
	b.Add(2)

	batch1 := <-b.OutChan()
	if len(batch1) != 2 {
		t.Fatalf("expected batch1 size 2, got %d", len(batch1))
	}

	b.Add(3)
	b.Add(4)

	batch2 := <-b.OutChan()
	if len(batch2) != 2 {
		t.Fatalf("expected batch2 size 2, got %d", len(batch2))
	}

	b.Close()
}
