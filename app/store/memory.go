package store

// TODO: implement the store interface

type Store struct {
	messages map[int]*Message
}

func NewMemoryStore() *Store {
	return &Store{
		messages: make(map[int]*Message),
	}
}

func (m *Store) Read(id int) (*Message, error) {
	elem, ok := m.messages[id]
	if !ok {
		return nil, ErrRecordNotFound
	}
	return elem, nil
}

func (m *Store) Write(msg Message) (int, error) {
	msg.ID = len(m.messages) + 1
	m.messages[msg.ID] = &msg
	return msg.ID, nil
}

func (m *Store) CandelizePreviousMinute(sensor string) error {
	return ErrPretendToCandelize
}
