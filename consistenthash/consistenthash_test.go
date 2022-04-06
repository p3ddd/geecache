package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// 哈希值: 2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.Add("6", "4", "2")
	// fmt.Printf("%+v\n", hash)
	// &{hash:0x77f740 replicas:3
	//	 keys:[2 4 6 12 14 16 22 24 26]
	//   hashMap:map[2:2 4:4 6:6 12:2 14:4 16:6 22:2 24:4 26:6]}

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	hash.Add("8")
	// fmt.Printf("%+v\n", hash)
	// &{hash:0x5af740 replicas:3
	//	 keys:[2 4 6 8 12 14 16 18 22 24 26 28]
	//   hashMap:map[2:2 4:4 6:6 8:8 12:2 14:4 16:6 18:8 22:2 24:4 26:6 28:8]}

	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// hash.Add("9", "10", "11", "12", "233")
	// fmt.Printf("%+v\n", hash)
	// &{hash:0xf1f740 replicas:3
	// keys:[2 4 6 8 9 10 11 12 12 14 16 18 19 22 24 26 28 29 110 111 112 210 211 212]
	// hashMap:map[2:2 4:4 6:6 8:8 9:9 10:10 11:11 12:12 14:4 16:6 18:8 19:9 22:2 24:4 26:6 28:8 29:9 110:10 111:11 112:12 210:10 211:11 212:12]}
}
