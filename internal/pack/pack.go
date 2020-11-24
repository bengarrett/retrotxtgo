package pack

//go:generate go run generator.go

type embed struct {
	storage map[string][]byte
}

func new() *embed {
	return &embed{storage: make(map[string][]byte)}
}

func (e *embed) Add(name string, content []byte) {
	e.storage[name] = content
}

func (e *embed) Get(name string) []byte {
	if n, ok := e.storage[name]; ok {
		return n
	}
	return nil
}

func (e *embed) Valid(name string) bool {
	if _, ok := e.storage[name]; ok {
		return true
	}
	return false
}

var pack = new()

// Add a named file to pack.
func Add(name string, content []byte) {
	pack.Add(name, content)
}

// Get a named file from pack.
func Get(name string) []byte {
	return pack.Get(name)
}

// Valid pack.
func Valid(name string) bool {
	return pack.Valid(name)
}
