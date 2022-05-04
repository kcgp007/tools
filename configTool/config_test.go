package configTool

import "testing"

func TestAdd(t *testing.T) {
	type dataConfig struct {
		S string
		I int
		F float64
	}

	var data = dataConfig{
		S: "aaa",
		I: 10,
		F: 10.5,
	}

	Add(&data)
	t.Log(data)

	AddWithKey("data", &data)
	t.Log(data)
}
