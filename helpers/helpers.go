package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

func GetRandomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

func FetchUrl(url string) []byte {
	const EncodeKeySize = 16
	keyString := GetRandomString(EncodeKeySize)
	//res, err := http.Get(url)
	res, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(keyString))

	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(res.Body)

	_ = res.Body.Close()
	decryptedData := AesDecryptECB(data, []byte(keyString))

	//return data
	return decryptedData
}
