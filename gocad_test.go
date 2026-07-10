package gocad

import (
	"strings"
	"testing"
)

func TestParseTriSurf(t *testing.T) {
	data := `GOCAD TriSurf 1.0
HEADER {
name: test-surface
}
PVRTX 1 0.0 0.0 0.0
PVRTX 2 1.0 0.0 0.0
PVRTX 3 1.0 1.0 0.0
PVRTX 4 0.0 1.0 0.0
TRGL 1 2 3
TRGL 1 3 4
END`

	ts, err := Parse(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if len(ts.Vertices) != 4 {
		t.Errorf("expected 4 vertices, got %d", len(ts.Vertices))
	}
	if len(ts.Triangles) != 2 {
		t.Errorf("expected 2 triangles, got %d", len(ts.Triangles))
	}
	if ts.Metadata["name"] != "test-surface" {
		t.Errorf("expected name test-surface, got %s", ts.Metadata["name"])
	}
}

func TestParseTriSurfWithColor(t *testing.T) {
	data := `GOCAD TriSurf 1.0
HEADER {
name: colored-surface
}
*solid*color: 0.5 0.6 0.7
PVRTX 1 0 0 0
PVRTX 2 1 0 0
PVRTX 3 0 1 0
TRGL 1 2 3
END`

	ts, err := Parse(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if ts.Color[0] != 0.5 || ts.Color[1] != 0.6 || ts.Color[2] != 0.7 {
		t.Errorf("unexpected color: %v", ts.Color)
	}
}

func TestWriteReadRoundTrip(t *testing.T) {
	orig := &TriSurf{
		Name: "roundtrip",
		Vertices: []Vertex{
			{0, 0, 0}, {1, 0, 0}, {1, 1, 0}, {0, 1, 0},
		},
		Triangles: []Triangle{
			{0, 1, 2}, {0, 2, 3},
		},
	}

	var buf strings.Builder
	if err := orig.Write(&buf); err != nil {
		t.Fatal(err)
	}

	parsed, err := Parse(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.Vertices) != 4 {
		t.Errorf("expected 4 vertices, got %d", len(parsed.Vertices))
	}
	if len(parsed.Triangles) != 2 {
		t.Errorf("expected 2 triangles, got %d", len(parsed.Triangles))
	}
}

func TestToOBJ(t *testing.T) {
	ts := &TriSurf{
		Vertices: []Vertex{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}},
		Triangles: []Triangle{{0, 1, 2}},
	}
	var buf strings.Builder
	if err := ts.ToOBJ(&buf); err != nil {
		t.Fatal(err)
	}
	obj := buf.String()
	if !strings.Contains(obj, "v 0.000000 0.000000 0.000000") {
		t.Errorf("OBJ missing vertex")
	}
	if !strings.Contains(obj, "f 1 2 3") {
		t.Errorf("OBJ missing face")
	}
}

func TestParsePLine(t *testing.T) {
	data := `GOCAD PLine 1.0
HEADER {
name: test-line
}
PVRTX 1 0 0 0
PVRTX 2 1 0 0
PVRTX 3 2 0 0
END`

	pl, err := ParsePLine(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if len(pl.Vertices) != 3 {
		t.Errorf("expected 3 vertices, got %d", len(pl.Vertices))
	}
	if pl.Metadata["name"] != "test-line" {
		t.Errorf("expected name test-line, got %s", pl.Metadata["name"])
	}
}
