package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	BUCKET       = os.Getenv("S3_BUCKET_NAME")
	awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

	storageAccessControl = "private"
	storageClass         = "GLACIER"
)

func main() {
	http.HandleFunc("/", handleRequestAndRedirect)
	log.Println("server start at :9080")
	if err := http.ListenAndServe(":9080", nil); err != nil {
		panic(err)
	}
}

func generateSignature(secretToken, payloadBody string) string {
	h := hmac.New(sha1.New, []byte(secretToken))
	h.Write([]byte(payloadBody))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func handleRequestAndRedirect(response http.ResponseWriter, request *http.Request) {

	body, err := ioutil.ReadAll(request.Body)

	t := time.Now()
	today := t.Format("2006-01-02")
	cacheKey := fmt.Sprintf("%x", md5.Sum(body))

	requestUri := fmt.Sprintf("toonapp/%s/%s.req", today, cacheKey)
	awsBackend := fmt.Sprintf("http://%s.s3.amazonaws.com/%s", BUCKET, requestUri)

	client := http.Client{}
	req, err := http.NewRequest(http.MethodPut, awsBackend, bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
	}

	d := time.Now().UTC().Format("20060102T150405Z")
	amzHeaders := "x-amz-acl:" + storageAccessControl + "\nx-amz-date:" + d + "\nx-amz-storage-class:" + storageClass
	resource := "/" + BUCKET + "/" + requestUri

	httpContentMd5 := request.Header.Get("Content-MD5")
	httpContentType := request.Header.Get("Content-Type")

	stringToSign := "PUT\n" + httpContentMd5 + "\n" + httpContentType + "\n\n" + amzHeaders + "\n" + resource

	awsSignature := generateSignature(awsSecretKey, stringToSign)
	auth := "AWS " + awsAccessKey + ":" + awsSignature

	req.Header = request.Header.Clone()
	req.Header.Set("Authorization", auth)
	req.Header.Set("x-amz-storage-class", storageClass)
	req.Header.Set("x-amz-acl", storageAccessControl)
	req.Header.Set("x-amz-date", d)

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	bdy, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(bdy))
}
