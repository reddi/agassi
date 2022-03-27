package service

type DataStore struct {
}

func (s *DataStore) Hello() (helloStr string) {
	return "Hello!"
}

func (s *DataStore) Close() error {
	return nil
}
