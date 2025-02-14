package linesegment

import (
	"embed"
	"encoding/json"
	"github.com/mikenye/geom2d/options"
	"io/fs"
	"testing"
)

// Embed testdata files
//
//go:embed testdata/testsegments_25x25_10.json
//go:embed testdata/testsegments_100x100_100.json
//go:embed testdata/testsegments_224x224_500.json
var testFiles embed.FS

func BenchmarkFindIntersectionsFast(b *testing.B) {
	err := fs.WalkDir(testFiles, "testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			b.Fatalf("WalkDir error: %v", err)
			return err
		}
		if !d.IsDir() {

			// Read and unmarshal once per file
			jsonBytes, err := testFiles.ReadFile(path)
			if err != nil {
				b.Fatalf("File read error (%s): %v", path, err)
			}

			var segments []LineSegment[int]
			if err := json.Unmarshal(jsonBytes, &segments); err != nil {
				b.Fatalf("Unmarshal error (%s): %v", path, err)
			}

			// Run the benchmark FindIntersectionsFast
			b.Run(path, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = FindIntersectionsFast(segments, options.WithEpsilon(1e-8))
				}
			})
		}
		return nil
	})
	if err != nil {
		b.Fatalf("WalkDir failed: %v", err)
	}
}

func BenchmarkFindIntersectionsSlow(b *testing.B) {
	err := fs.WalkDir(testFiles, "testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			b.Fatalf("WalkDir error: %v", err)
			return err
		}
		if !d.IsDir() {

			// Read and unmarshal once per file
			jsonBytes, err := testFiles.ReadFile(path)
			if err != nil {
				b.Fatalf("File read error (%s): %v", path, err)
			}

			var segments []LineSegment[int]
			if err := json.Unmarshal(jsonBytes, &segments); err != nil {
				b.Fatalf("Unmarshal error (%s): %v", path, err)
			}

			// Run the benchmark FindIntersectionsFast
			b.Run(path, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = FindIntersectionsSlow(segments, options.WithEpsilon(1e-8))
				}
			})
		}
		return nil
	})
	if err != nil {
		b.Fatalf("WalkDir failed: %v", err)
	}
}

//func TestFindIntersectionsFastFromFiles(t *testing.T) {
//	err := fs.WalkDir(testFiles, "testdata", func(path string, d fs.DirEntry, err error) error {
//		if err != nil {
//			return err
//		}
//		if !d.IsDir() {
//
//			// read and unmarshal
//			t.Log("read and unmarshal test data")
//			b, err := testFiles.ReadFile(path)
//			require.NoError(t, err, "file read error")
//			var segments []LineSegment[int]
//			err = json.Unmarshal(b, &segments)
//			require.NoError(t, err, "unmarshal error")
//
//			t.Run(path, func(t *testing.T) {
//
//				if testing.Short() {
//					t.Skip("Skipping long-running test in short mode.")
//				}
//				t.Parallel() // Runs in parallel with other tests
//
//				// run sweepline
//				t.Log("run FindIntersectionsFast")
//				resultsFast := FindIntersectionsFast(segments, options.WithEpsilon(1e-8))
//				t.Logf("FindIntersectionsFast computed %d intersections", len(resultsFast))
//
//				// run naive
//				t.Log("run FindIntersectionsSlow")
//				resultsSlow := FindIntersectionsSlow(segments, options.WithEpsilon(1e-8))
//				t.Logf("FindIntersectionsSlow computed %d intersections", len(resultsSlow))
//
//				// compare
//				t.Log("comparing results")
//				require.True(t, InterSectionResultsEq(resultsFast, resultsSlow, options.WithEpsilon(1e-8)))
//
//			})
//		}
//		return nil
//	})
//	require.NoError(t, err, "walkdir error")
//}
