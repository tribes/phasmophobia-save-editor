package phasmophobia

// This may change in future updates
var salt = []byte("CHANGE ME TO YOUR OWN RANDOM STRING")

// StringSaveEntry is a configuration item that contains a value of type string
type StringSaveEntry struct {
	Key   string
	Value string
}

// IntSaveEntry is a configuration item that contains a value of type integer
type IntSaveEntry struct {
	Key   string
	Value int64
}

// FloatSaveEntry is a configuration item that contains a value of type float
type FloatSaveEntry struct {
	Key   string
	Value float64
}

// BoolSaveEntry is a configuration item that contains a value of type boolean
type BoolSaveEntry struct {
	Key   string
	Value bool
}

// Save is the struct that hold the whole phasmophobia
type Save struct {
	Path       string `json:"-"`
	StringData []*StringSaveEntry
	IntData    []*IntSaveEntry
	FloatData  []*FloatSaveEntry
	BoolData   []*BoolSaveEntry
}
