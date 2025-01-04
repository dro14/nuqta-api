package meili

func (m *Meili) Ping() error {
	_, err := m.client.Health()
	return err
}
