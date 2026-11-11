package fuse

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestWritevEmpty(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "writev-empty")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer f.Close()

	n, err := writev(int(f.Fd()), nil)
	if err != nil || n != 0 {
		t.Fatalf("writev(nil) = (%d, %v), want (0, nil)", n, err)
	}

	n, err = writev(int(f.Fd()), [][]byte{{}, nil, {}})
	if err != nil || n != 0 {
		t.Fatalf("writev(empty slices) = (%d, %v), want (0, nil)", n, err)
	}
}

func TestWritevExceedsMaxIovecs(t *testing.T) {
	for _, numSlices := range []int{1, 64, 128, 129, 1024, 1025, 1500} {
		t.Run(fmt.Sprintf("slices_%d", numSlices), func(t *testing.T) {
			f, err := os.CreateTemp(t.TempDir(), "writev")
			if err != nil {
				t.Fatalf("CreateTemp: %v", err)
			}
			defer f.Close()

			var packet [][]byte
			var want bytes.Buffer
			for i := 0; i < numSlices; i++ {
				if i%7 == 0 {
					packet = append(packet, nil)
				}
				chunk := []byte(fmt.Sprintf("chunk-%04d;", i))
				packet = append(packet, chunk)
				want.Write(chunk)
			}

			n, err := writev(int(f.Fd()), packet)
			if err != nil {
				t.Fatalf("writev(%d non-empty slices) failed: %v", numSlices, err)
			}
			if n != want.Len() {
				t.Fatalf("writev returned n=%d, want %d", n, want.Len())
			}

			got, err := os.ReadFile(f.Name())
			if err != nil {
				t.Fatalf("ReadFile: %v", err)
			}
			if !bytes.Equal(got, want.Bytes()) {
				t.Fatalf("file content mismatch for %d slices (got %d bytes, want %d bytes)", numSlices, len(got), want.Len())
			}
		})
	}
}

func TestPutOutMessageResetsSglist(t *testing.T) {
	var c Connection
	om := c.getOutMessage()
	payload := []byte("payload-data")
	om.Append(payload)
	om.Grow(64)

	if len(om.Sglist) == 0 {
		t.Fatal("expected non-empty Sglist before putOutMessage")
	}
	backing := om.Sglist[:cap(om.Sglist)]

	c.putOutMessage(om)

	if len(om.Sglist) != 0 {
		t.Fatalf("len(om.Sglist) = %d after putOutMessage, want 0", len(om.Sglist))
	}
	for i, s := range backing {
		if s != nil {
			t.Fatalf("backing Sglist[%d] still holds slice reference after putOutMessage", i)
		}
	}
}
