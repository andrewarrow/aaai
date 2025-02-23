package diff

// Change represents a modification to a file
type Change struct {
	File   string
	Ranges []Range
}

// Range represents a specific change within a file
type Range struct {
	Start  int      // Start line number (0-based)
	End    int      // End line number (0-based)
	Before []string // Lines before the change
	After  []string // Lines after the change
}
