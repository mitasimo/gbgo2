package mapstruct

import (
	"testing"
)

func TestMapStruct(t *testing.T) {

	type Person struct {
		Name string
		Age  int
	}

	const (
		personName = "John Lennon"
		personAge  = 81
	)

	// шаг 1, все ключи мапы соответствуют полям структуры

	// создать структуру
	s1 := new(Person)

	// создать мапу
	cfg1 := make(map[string]interface{})
	cfg1["Name"] = personName
	cfg1["Age"] = 81

	// заполнить структуру из мапы
	// возврат должен быть nil
	err := MapStruct(s1, cfg1)
	if err != nil {
		t.Errorf("%v", err)
	}

	// сравнить значения структуры с референстными
	if s1.Name != personName || s1.Age != personAge {
		t.Errorf("Шаг 1: s1.Name == %s, must be a %s; s.Age == %d, must be a %d", s1.Name, personName, s1.Age, personAge)
	}

	// Шаг 2: только один ключ в мапе соответствует полю в структуре

	// создать структуру
	s2 := new(Person)

	// создать мапу
	cfg2 := make(map[string]interface{})
	cfg2["Name"] = personName
	cfg2["Band"] = "Beatls"

	// заполнить структуру из мапы
	// возврат должен быть nil
	err = MapStruct(s2, cfg2)
	if err != nil {
		t.Errorf("%v", err)
	}

	// сравнить значения структуры с референстными
	if s2.Name != personName || s2.Age != 0 {
		t.Errorf("Шаг 2: s2.Name == %s, must be a %s; s.Age == %d, must be a %d", s2.Name, personName, s2.Age, 0)
	}

	// Шаг 3: только один ключ в мапе соответствует полю в структуре

	// создать структуру
	s3 := new(Person)

	// создать мапу
	cfg3 := make(map[string]interface{})
	cfg3["Name"] = true

	// заполнить структуру из мапы
	MapStruct(s3, cfg3)
	if s3.Name != "" || s3.Age != 0 {
		t.Errorf("Шаг 3: s3.Name == %s, must be a %s; s.Age == %d, must be a %d", s3.Name, "", s3.Age, 0)
	}

}
