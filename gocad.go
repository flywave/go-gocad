package gocad

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Vertex [3]float64
type Triangle [3]int

type TriSurf struct {
	Name       string
	Properties map[string][]float64
	Vertices   []Vertex
	Triangles  []Triangle
	Color      [3]float64
	Metadata   map[string]string
}

type PLine struct {
	Name     string
	Vertices []Vertex
	Closed   bool
	Metadata map[string]string
}

func Parse(r io.Reader) (*TriSurf, error) {
	scanner := bufio.NewScanner(r)
	ts := &TriSurf{
		Properties: make(map[string][]float64),
		Metadata:   make(map[string]string),
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "GOCAD":
			if len(fields) > 2 {
				ts.Metadata["format"] = fields[2]
			}

		case "HEADER":
			for scanner.Scan() {
				hdr := strings.TrimSpace(scanner.Text())
				if hdr == "}" {
					break
				}
				if colonIdx := strings.Index(hdr, ":"); colonIdx > 0 {
					key := strings.TrimSpace(hdr[:colonIdx])
					val := strings.TrimSpace(hdr[colonIdx+1:])
					ts.Metadata[key] = val
				}
			}

		case "PVRTX":
			if len(fields) >= 5 {
				x, _ := strconv.ParseFloat(fields[2], 64)
				y, _ := strconv.ParseFloat(fields[3], 64)
				z, _ := strconv.ParseFloat(fields[4], 64)
				ts.Vertices = append(ts.Vertices, Vertex{x, y, z})
			} else if len(fields) >= 4 {
				x, _ := strconv.ParseFloat(fields[1], 64)
				y, _ := strconv.ParseFloat(fields[2], 64)
				z, _ := strconv.ParseFloat(fields[3], 64)
				ts.Vertices = append(ts.Vertices, Vertex{x, y, z})
			}

		case "TRGL":
			if len(fields) >= 4 {
				i1, _ := strconv.Atoi(fields[1])
				i2, _ := strconv.Atoi(fields[2])
				i3, _ := strconv.Atoi(fields[3])
				ts.Triangles = append(ts.Triangles, Triangle{i1 - 1, i2 - 1, i3 - 1})
			}

		case "END":
			// end of object

		case "*solid*color:":
			if len(fields) >= 4 {
				r, _ := strconv.ParseFloat(fields[1], 64)
				g, _ := strconv.ParseFloat(fields[2], 64)
				b, _ := strconv.ParseFloat(fields[3], 64)
				ts.Color = [3]float64{r, g, b}
			}
		}
	}
	return ts, scanner.Err()
}

func ParsePLine(r io.Reader) (*PLine, error) {
	scanner := bufio.NewScanner(r)
	pl := &PLine{
		Metadata: make(map[string]string),
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "GOCAD":
			if len(fields) > 2 {
				pl.Metadata["format"] = fields[2]
			}

		case "HEADER":
			for scanner.Scan() {
				hdr := strings.TrimSpace(scanner.Text())
				if hdr == "}" {
					break
				}
				if colonIdx := strings.Index(hdr, ":"); colonIdx > 0 {
					key := strings.TrimSpace(hdr[:colonIdx])
					val := strings.TrimSpace(hdr[colonIdx+1:])
					pl.Metadata[key] = val
				}
			}

		case "PVRTX":
			if len(fields) >= 5 {
				x, _ := strconv.ParseFloat(fields[2], 64)
				y, _ := strconv.ParseFloat(fields[3], 64)
				z, _ := strconv.ParseFloat(fields[4], 64)
				pl.Vertices = append(pl.Vertices, Vertex{x, y, z})
			}

		case "SEG":
			pl.Closed = false

		case "END":
			break
		}
	}
	return pl, scanner.Err()
}

func (ts *TriSurf) Write(w io.Writer) error {
	fmt.Fprintf(w, "GOCAD TriSurf 1.0\n")
	fmt.Fprintf(w, "HEADER {\n")
	fmt.Fprintf(w, "name: %s\n", ts.Name)
	for k, v := range ts.Metadata {
		if k != "format" {
			fmt.Fprintf(w, "%s: %s\n", k, v)
		}
	}
	fmt.Fprintf(w, "}\n")

	if ts.Color != [3]float64{0, 0, 0} {
		fmt.Fprintf(w, "*solid*color: %.6f %.6f %.6f\n", ts.Color[0], ts.Color[1], ts.Color[2])
	}

	for i, v := range ts.Vertices {
		fmt.Fprintf(w, "PVRTX %d %.6f %.6f %.6f\n", i+1, v[0], v[1], v[2])
	}

	for _, t := range ts.Triangles {
		fmt.Fprintf(w, "TRGL %d %d %d\n", t[0]+1, t[1]+1, t[2]+1)
	}

	fmt.Fprintf(w, "END\n")
	return nil
}

func (ts *TriSurf) ToOBJ(w io.Writer) error {
	for _, v := range ts.Vertices {
		fmt.Fprintf(w, "v %.6f %.6f %.6f\n", v[0], v[1], v[2])
	}
	for _, t := range ts.Triangles {
		fmt.Fprintf(w, "f %d %d %d\n", t[0]+1, t[1]+1, t[2]+1)
	}
	return nil
}
