package customerimporter

import "testing"

func TestReadFileWithGoroutines(t *testing.T) {
	_, err := ReadFileWithGoroutines("../testdata/notfound.csv")
	if err == nil {
		t.Error()
	}

	mapOfDomains, err := ReadFileWithGoroutines("../testdata/customers.csv")
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

func BenchmarkReadFileWithGoroutines(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReadFileWithGoroutines("../testdata/customers.csv")
	}
}

func BenchmarkReadFileWithGoroutinesLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReadFileWithGoroutines("../testdata/customersLarge.csv")
	}
}

func BenchmarkReadFileWithGoroutinesSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReadFileWithGoroutines("../testdata/customersSmall.csv")
	}
}
