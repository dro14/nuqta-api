package memcached

func (m *Memcached) Ping() error {
	return m.client.Ping()
}
