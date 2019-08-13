package customerimporter

import "testing"

func TestReadFile(t *testing.T) {
	_, err := ReadFile("../testdata/notfound.csv")
	if err == nil {
		t.Error()
	}

	mapOfDomains, err := ReadFile("../testdata/customers.csv")
	if err != nil {
		t.Error(err)
	}

	if len(mapOfDomains) != 500 {
		t.Error()
	}

	if mapOfDomains["github.io"] != 8 || mapOfDomains["google.com.br"] != 5 {
		t.Error()
	}

	if mapOfDomains["test.com"] != 0 {
		t.Error()
	}
}

func BenchmarkReadFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReadFile("../testdata/customers.csv")
	}
}

func BenchmarkReadFileLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReadFile("../testdata/customersLarge.csv")
	}
}

func BenchmarkReadFileSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReadFile("../testdata/customersSmall.csv")
	}
}
