package store

// TODO: implement the store interface

type Store struct {
	// memory store
	// key is the module and Bolt bucket name
	// value is the slice of Data
	data map[string][]Data
}

func NewMemoryStore() *Store {
	return &Store{
		data: make(map[string][]Data),
	}
}

// Implement the Storer interface

// Read reads records for the given module from the database
func (s *Store) Read(module string) (data []Data, err error) {
	if _, ok := s.data[module]; !ok {
		return nil, ErrRecordNotFound
	}
	return s.data[module], nil
}

// Write writes the data to the database
func (s *Store) Write(data Data) error {
	if _, ok := s.data[data.Module]; !ok {
		s.data[data.Module] = []Data{}
	}
	s.data[data.Module] = append(s.data[data.Module], data)
	return nil
}

// View returns a map of topics and their values for the given module
// The map is sorted by DateTime and structured as follows:
// map[Topic]map[DateTime]Value
func (s *Store) View(module string) (data map[string]map[string]string, err error) {

	data = make(map[string]map[string]string)

	// select distinct topics from module
	for _, d := range s.data[module] {
		data[d.Topic] = make(map[string]string)
	}

	// select all records from module and fill the map
	for _, d := range s.data[module] {
		data[d.Topic][d.DateTime] = d.Value
	}

	return
}
