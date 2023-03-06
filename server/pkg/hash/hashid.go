package hash

import (
	"fmt"
	"github.com/speps/go-hashids/v2"
)

type ID uint

const salt = "aptksnhjsdf"
const minLength = 20

func GetHashIds() (*hashids.HashID, error) {
	var hd = hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func IdToHash(id uint) (string, error) {
	h, err := GetHashIds()
	if err != nil {
		return "", err
	}
	hash, err := h.Encode([]int{int(id)})
	if err != nil {
		return "", err
	}
	return hash, nil
}

func HashToId(hash string) (uint, error) {
	h, err := GetHashIds()
	if err != nil {
		return 0, err
	}
	ids, err := h.DecodeWithError(hash)
	if err != nil {
		return 0, err
	}
	if len(ids) == 0 {
		return 0, fmt.Errorf("no id found")
	}
	return uint(ids[0]), err
}

func (I *ID) UnmarshalJSON(bytes []byte) error {
	str := string(bytes)
	hash := str[1 : len(str)-1]
	id, err := HashToId(hash)
	if err != nil {
		return err
	}
	*I = ID(id)
	return nil
}

func (I ID) MarshalJSON() ([]byte, error) {
	hash, err := IdToHash(uint(I))
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf(`"%s"`, hash)), nil
}

func (I ID) Uint() uint {
	return uint(I)
}
