package profiler

type Profiler interface {
	AgifyPerson(name string) (int, error)
	GenderizePerson(name string) (string, error)
	NationalizePerson(name string) (string, error)
}
