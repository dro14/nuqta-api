package memcached

func (m *Memcached) Ping() (string, error) {
	err := m.client.Ping()
	if err != nil {
		return "", err
	}
	return "Pong", nil
}
