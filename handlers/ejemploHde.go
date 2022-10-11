package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/culturadevops/drs/extra/appConfigs"
	"github.com/culturadevops/drs/s3"

	"github.com/gofiber/fiber/v2"
)

var AppConfig appConfigs.Web

type Ejemplo struct {
}

type ejemplorqt struct {
	//	Account  string `json:"account" form:"account" query:"account"`
}

func writeFile(filename string, jsoncontent interface{}) {
	file, _ := json.MarshalIndent(jsoncontent, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)
}

func UploadMultipartformToS3(filename string, bucket string, RutaFinalEnS3 string, Region string) error {
	s3 := new(s3.S3Client)
	s3.NewSession(Region)
	RutaFinalEnS3 = RutaFinalEnS3 + filename
	return s3.Upload(filename, bucket, RutaFinalEnS3, "text/plain")
}
func Post(url string, s3ruta string, tabla string) (string, string, string, string) {
	client := &http.Client{}

	//fmt.Println(url)
	jsonBody := []byte(`{ "@type": "MessageCard","@context": "http://schema.org/extensions",
	"themeColor": "0076D7","summary": "Archivo de respuesta",
	"sections": [{"activityTitle": "dynamo scan", "activitySubtitle": "dev lima",
	"activityImage": "https://teamsnodesample.azurewebsites.net/static/img/image5.png",
	"facts": [{ "name": "tabla", "value": "` + tabla + `"},
	{"name": "archivo","value": "` + s3ruta + `"}],
	"markdown": true
	}],"potentialAction": [
	{ "@type": "OpenUri","name": "Learn More","targets": [{"os": "default","uri": "` + s3ruta + `"  }]}]}`)

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Errored when sending request to the server")

	} else {
		defer resp.Body.Close()
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		type Build struct {
			Name    string
			Version string
		}
		type Objt struct {
			Build Build
		}
		var jsonvars Objt

		json.Unmarshal([]byte(string(responseBody)), &jsonvars)
		if resp.Status == "200 OK" {
			if responseBody != nil {
				return resp.Status, "microname", jsonvars.Build.Name, jsonvars.Build.Version
			} else {
				return resp.Status, "", "", ""
			}
		}
	}
	return resp.Status, "microname", "", ""
}
func Getenv(key string) string {
	value, defined := os.LookupEnv(key)
	if !defined {
		fmt.Println("falta variable de entorno" + key)
		os.Exit(1)
	}
	return value
}
func (t *Ejemplo) List() fiber.Handler {
	return func(c *fiber.Ctx) error {

		//tableName := "mitabla"
		//Region := "us-east-1"
		//buckets3 := "lima-aws-cicd-pipeline"
		//Region := AppConfig.Region
		//buckets3 := AppConfig.Bucket
		Region := Getenv("REGION")
		buckets3 := Getenv("BUCKET")
		tableName := c.Params("tabla")
		Webhook := Getenv("WEBHOOK")
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		svc := dynamodb.New(sess)
		params := &dynamodb.ScanInput{
			TableName: aws.String(tableName),
		}
		result, err := svc.Scan(params)
		if err != nil {
			return TraditionalResponse(c, err, nil, "error scan")
		}
		var SUPERMAP []map[string]interface{}
		for _, i := range result.Items {
			var item map[string]interface{}
			err = dynamodbattribute.UnmarshalMap(i, &item)
			if err != nil {
				log.Fatalf("Got error unmarshalling: %s", err)
				return TraditionalResponse(c, err, nil, "Got error unmarshalling")

			}
			SUPERMAP = append(SUPERMAP, item)
		}
		jsonfilename := tableName + ".json"
		writeFile(jsonfilename, SUPERMAP)
		if UploadMultipartformToS3(jsonfilename, buckets3, "", Region) != nil {
			return TraditionalResponse(c, err, nil, "s3 upload error")
		}
		s3 := new(s3.S3Client)
		s3.NewSession(Region)
		fmt.Println("------------------------")
		s3rutalfinal := s3.GenerateUrlForDownload(buckets3, jsonfilename)

		Post(Webhook, s3rutalfinal, tableName)
		return c.Status(http.StatusOK).JSON("ok")
	}
}
