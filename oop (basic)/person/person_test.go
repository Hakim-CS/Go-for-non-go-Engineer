package person

import "testing"

func TestNewAddress(t *testing.T) {
	ta := NewAddress("turkey", "denizli", 222)
	if ta == nil {
		t.Errorf("Newaddress() returned nil")
	}
}

func TestNewPerson(t *testing.T) {
	addr := Address{
		country: "turkey",
		city:    "denizli",
		zip:     222,
	}

	_, err := NewPerson("hakim", 21, "hakim@gmail.com", addr)
	if err != nil {
		t.Errorf("NewPerson() returned error: %v", err)
	}

}
