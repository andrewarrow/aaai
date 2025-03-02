package diff

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestApplyPatch(t *testing.T) {
	tests := []struct {
		name     string
		original string
		diff     string
		expected string
	}{
		{
			name: "Basic modification",
			original: `line 1
line 2
line 3
line 4
line 5
`,
			diff: `--- original.txt
+++ modified.txt
@@ -2,3 +2,4 @@
 line 2
-line 3
+modified line 3
+new line
 line 4
`,
			expected: `line 1
line 2
modified line 3
new line
line 4
line 5
`,
		},
		{
			name: "Multiple hunks",
			original: `line 1
line 2
line 3
line 4
line 5
line 6
line 7
line 8
line 9
line 10
`,
			diff: `--- original.txt
+++ modified.txt
@@ -2,3 +2,4 @@
 line 2
-line 3
+modified line 3
+inserted line
 line 4
@@ -7,3 +8,4 @@
 line 7
-line 8
+modified line 8
+another new line
 line 9
`,
			expected: `line 1
line 2
modified line 3
inserted line
line 4
line 5
line 6
line 7
modified line 8
another new line
line 9
line 10
`,
		},
		{
			name: "Add-only hunk",
			original: `line 1
line 2
line 3
line 4
line 5
`,
			diff: `--- original.txt
+++ modified.txt
@@ -3,0 +4,2 @@
+new line A
+new line B
`,
			expected: `line 1
line 2
line 3
new line A
new line B
line 4
line 5
`,
		},
		{
			name: "Delete-only hunk",
			original: `line 1
line 2
line 3
line 4
line 5
`,
			diff: `--- original.txt
+++ modified.txt
@@ -3,1 +3,0 @@
-line 3
`,
			expected: `line 1
line 2
line 4
line 5
`,
		},
		{
			name:     "Create new file",
			original: ``,
			diff: `--- /dev/null
+++ new_file.txt
@@ -0,0 +1,3 @@
+line 1
+line 2
+line 3
`,
			expected: `line 1
line 2
line 3
`,
		},
		{
			name: "Empty file after deletion",
			original: `line 1
line 2
line 3
`,
			diff: `--- original.txt
+++ modified.txt
@@ -1,3 +0,0 @@
-line 1
-line 2
-line 3
`,
			expected: ``,
		},
		{
			name: "Fuzzy position matching",
			original: `extra line 0
line 1
line 2
line 3
line 4
line 5
`,
			diff: `--- original.txt
+++ modified.txt
@@ -2,3 +2,4 @@
 line 2
-line 3
+modified line 3
+new line
 line 4
`,
			expected: `extra line 0
line 1
line 2
modified line 3
new line
line 4
line 5
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tmpDir, err := ioutil.TempDir("", "diff_test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			// Create original file
			origPath := filepath.Join(tmpDir, "original.txt")
			if err := ioutil.WriteFile(origPath, []byte(tt.original), 0644); err != nil {
				t.Fatal(err)
			}

			// Create diff file
			diffPath := filepath.Join(tmpDir, "diff.patch")
			if err := ioutil.WriteFile(diffPath, []byte(tt.diff), 0644); err != nil {
				t.Fatal(err)
			}

			// Apply patch
			ApplyPatch(origPath, tt.diff)

			// Read result
			result, err := ioutil.ReadFile(origPath)
			if err != nil {
				t.Fatal(err)
			}

			// Compare with expected
			if string(result) != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, string(result))
			}
		})
	}
}

func TestParseHunks(t *testing.T) {
	diff := `--- original.txt
+++ modified.txt
@@ -2,3 +2,4 @@
 line 2
-line 3
+modified line 3
+new line
 line 4
@@ -7,3 +8,4 @@
 line 7
-line 8
+modified line 8
+another new line
 line 9
`

	lines, _ := readLinesFromString(diff)
	hunks := parseHunks(lines)

	if len(hunks) != 2 {
		t.Errorf("Expected 2 hunks, got %d", len(hunks))
	}

	expected := []Hunk{
		{
			StartLine: 1,
			Length:    3,
			NewStart:  1,
			NewLength: 4,
			Lines: []string{
				" line 2",
				"-line 3",
				"+modified line 3",
				"+new line",
				" line 4",
			},
		},
		{
			StartLine: 6,
			Length:    3,
			NewStart:  7,
			NewLength: 4,
			Lines: []string{
				" line 7",
				"-line 8",
				"+modified line 8",
				"+another new line",
				" line 9",
				"",
			},
		},
	}

	for i, hunk := range hunks {
		if hunk.StartLine != expected[i].StartLine {
			t.Errorf("Hunk %d: expected StartLine %d, got %d", i, expected[i].StartLine, hunk.StartLine)
		}
		if hunk.Length != expected[i].Length {
			t.Errorf("Hunk %d: expected Length %d, got %d", i, expected[i].Length, hunk.Length)
		}
		if hunk.NewStart != expected[i].NewStart {
			t.Errorf("Hunk %d: expected NewStart %d, got %d", i, expected[i].NewStart, hunk.NewStart)
		}
		if hunk.NewLength != expected[i].NewLength {
			t.Errorf("Hunk %d: expected NewLength %d, got %d", i, expected[i].NewLength, hunk.NewLength)
		}

		// Only check line count to avoid newline issues in testing
		if len(hunk.Lines) != len(expected[i].Lines) {
			t.Errorf("Hunk %d: expected %d lines, got %d", i, len(expected[i].Lines), len(hunk.Lines))
			continue
		}

		// Check each line content separately
		for j, line := range hunk.Lines {
			expectedLine := expected[i].Lines[j]
			if strings.TrimSpace(line) != strings.TrimSpace(expectedLine) {
				t.Errorf("Hunk %d, Line %d: expected '%s', got '%s'", i, j, expectedLine, line)
			}
		}
	}
}

func TestFindHunkPosition(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		hunk     Hunk
		expected int
	}{
		{
			name: "Exact match",
			lines: []string{
				"line 1\n",
				"line 2\n",
				"line 3\n",
				"line 4\n",
			},
			hunk: Hunk{
				StartLine: 1,
				Length:    2,
				Lines: []string{
					" line 2",
					"-line 3",
				},
			},
			expected: 1,
		},
		{
			name: "Fuzzy match (shifted)",
			lines: []string{
				"extra line\n",
				"line 1\n",
				"line 2\n",
				"line 3\n",
				"line 4\n",
			},
			hunk: Hunk{
				StartLine: 1,
				Length:    2,
				Lines: []string{
					" line 2",
					"-line 3",
				},
			},
			expected: 2,
		},
		{
			name: "Add-only hunk",
			lines: []string{
				"line 1\n",
				"line 2\n",
				"line 3\n",
			},
			hunk: Hunk{
				StartLine: 1,
				Length:    0,
				Lines: []string{
					"+new line",
				},
			},
			expected: 1,
		},
		{
			name: "Delete-only hunk",
			lines: []string{
				"line 1\n",
				"line 2\n",
				"line 3\n",
			},
			hunk: Hunk{
				StartLine: 1,
				Length:    1,
				Lines: []string{
					"-line 2",
				},
			},
			expected: 1,
		},
		{
			name: "No match",
			lines: []string{
				"line 1\n",
				"line 2\n",
				"line 3\n",
			},
			hunk: Hunk{
				StartLine: 1,
				Length:    2,
				Lines: []string{
					" line 2",
					"-completely different line",
				},
			},
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := findHunkPosition(tt.lines, tt.hunk)
			if pos != tt.expected {
				t.Errorf("Expected position %d, got %d", tt.expected, pos)
			}
		})
	}
}

func TestApplyHunks(t *testing.T) {
	original := []string{
		"line 1\n",
		"line 2\n",
		"line 3\n",
		"line 4\n",
		"line 5\n",
	}

	hunks := []Hunk{
		{
			StartLine: 1,
			Length:    2,
			Lines: []string{
				" line 2",
				"-line 3",
				"+modified line 3",
			},
		},
	}

	expected := []string{
		"line 1\n",
		"line 2\n",
		"modified line 3\n",
		"line 4\n",
		"line 5\n",
	}

	result := applyHunks(original, hunks)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, result)
	}
}

func TestReadLinesFromString(t *testing.T) {
	content := "line 1\nline 2\nline 3"
	expected := []string{"line 1\n", "line 2\n", "line 3\n"}

	lines, err := readLinesFromString(content)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(lines, expected) {
		t.Errorf("Expected:\n%v\nGot:\n%v", expected, lines)
	}
}

func TestReadWriteLines(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "diff_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "test.txt")
	content := []string{"line 1\n", "line 2\n", "line 3\n"}

	err = writeLines(filePath, content)
	if err != nil {
		t.Fatal(err)
	}

	lines, err := readLines(filePath)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(lines, content) {
		t.Errorf("Expected:\n%v\nGot:\n%v", content, lines)
	}
}
