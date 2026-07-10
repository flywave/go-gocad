package gocad

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestOpenCubeTS(t *testing.T) {
	f, err := os.Open("testdata/cube.ts")
	if err != nil {
		t.Skip("test data not found:", err)
	}
	defer f.Close()

	ts, err := Parse(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(ts.Vertices) != 8 {
		t.Errorf("expected 8 vertices, got %d", len(ts.Vertices))
	}
	if len(ts.Triangles) != 12 {
		t.Errorf("expected 12 triangles, got %d", len(ts.Triangles))
	}
	if ts.Metadata["name"] != "unit-cube" {
		t.Errorf("expected name unit-cube, got %s", ts.Metadata["name"])
	}
	if ts.Color[0] != 0.8 || ts.Color[1] != 0.8 || ts.Color[2] != 0.8 {
		t.Errorf("expected color 0.8 0.8 0.8, got %v", ts.Color)
	}
}

func TestOpenLinePL(t *testing.T) {
	f, err := os.Open("testdata/line.pl")
	if err != nil {
		t.Skip("test data not found:", err)
	}
	defer f.Close()

	pl, err := ParsePLine(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(pl.Vertices) != 4 {
		t.Errorf("expected 4 vertices, got %d", len(pl.Vertices))
	}
	if pl.Metadata["name"] != "fault-trace" {
		t.Errorf("expected name fault-trace, got %s", pl.Metadata["name"])
	}
}

func TestRoundTripFile(t *testing.T) {
	orig, err := os.ReadFile("testdata/cube.ts")
	if err != nil {
		t.Skip("test data not found:", err)
	}
	// parse, write, parse again
	ts1, err := Parse(bytes.NewReader(orig))
	if err != nil {
		t.Fatal(err)
	}
	var buf strings.Builder
	if err := ts1.Write(&buf); err != nil {
		t.Fatal(err)
	}
	ts2, err := Parse(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatal(err)
	}
	if len(ts1.Vertices) != len(ts2.Vertices) {
		t.Errorf("vertex count mismatch: %d vs %d", len(ts1.Vertices), len(ts2.Vertices))
	}
	if len(ts1.Triangles) != len(ts2.Triangles) {
		t.Errorf("triangle count mismatch: %d vs %d", len(ts1.Triangles), len(ts2.Triangles))
	}
}

func TestExportToOBJFromFile(t *testing.T) {
	f, err := os.Open("testdata/cube.ts")
	if err != nil {
		t.Skip("test data not found:", err)
	}
	defer f.Close()

	ts, err := Parse(f)
	if err != nil {
		t.Fatal(err)
	}
	var buf strings.Builder
	if err := ts.ToOBJ(&buf); err != nil {
		t.Fatal(err)
	}
	obj := buf.String()
	if len(obj) == 0 {
		t.Error("OBJ output should not be empty")
	}
}
