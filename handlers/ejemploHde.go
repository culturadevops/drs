package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
func (t *Ejemplo) List() fiber.Handler {
	return func(c *fiber.Ctx) error {
		AppConfig.Configure("./configs", "app")
		//tableName := "mitabla"
		//Region := "us-east-1"
		//buckets3 := "lima-aws-cicd-pipeline"
		Region := AppConfig.Region
		buckets3 := AppConfig.Bucket
		tableName := c.Params("tabla")
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		svc := dynamodb.New(sess)
		params := &dynamodb.ScanInput{
			TableName: aws.String(tableName),
		}
		result, err := svc.Scan(params)
		if err != nil {
			log.Fatalf("Query API call failed: %s", err)
			return c.Status(http.StatusInternalServerError).JSON(err)
		}
		var SUPERMAP []map[string]string
		for _, i := range result.Items {
			var item map[string]string
			err = dynamodbattribute.UnmarshalMap(i, &item)
			if err != nil {
				log.Fatalf("Got error unmarshalling: %s", err)
				return c.Status(http.StatusInternalServerError).JSON("Got error unmarshalling")
			}
			SUPERMAP = append(SUPERMAP, item)
		}
		jsonfilename := tableName + ".json"
		writeFile(jsonfilename, SUPERMAP)
		if UploadMultipartformToS3(jsonfilename, buckets3, "", Region) != nil {
			log.Fatalf("s3 upload error: %s", err)
			return c.Status(http.StatusInternalServerError).JSON("s3 upload error")
		}
		s3 := new(s3.S3Client)
		s3.NewSession(Region)
		fmt.Println("------------------------")
		fmt.Println(s3.GenerateUrlForDownload(buckets3, jsonfilename))
		return c.Status(http.StatusOK).JSON("ok")
	}
}

/*
func (t *Ejemplo) Add() fiber.Handler {
	return func(c *fiber.Ctx) error {
		u := new(ejemplorqt)
		err := c.BodyParser(u)
		data := ""
		if err != nil {
			// falta el agregar model
			return c.Status(http.StatusInternalServerError).JSON(err)
		}
		return c.Status(http.StatusOK).JSON(data)
	}
}

func (t *Ejemplo) List() fiber.Handler {
	return func(c *fiber.Ctx) error {
		//data := models.Ejemplo.List()
		return c.Status(http.StatusOK).JSON(data)
	}
}
func (t *Ejemplo) Show() fiber.Handler {
	return func(c *fiber.Ctx) error {
		//id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
		//data, _ := models.Ejemplo.Info(uint(id))
		return c.Status(http.StatusOK).JSON(data)
	}
}
func (t *Ejemplo) Del() fiber.Handler {
	return func(c *fiber.Ctx) error {
		//id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
		//data, _ := models.Ejemplo.Del(uint(id))
		return c.Status(http.StatusOK).JSON(data)
	}
}
func (t *Ejemplo) Update() fiber.Handler {
	return func(c *fiber.Ctx) error {
		//id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
		u := new(ejemplorqt)
		if err := c.BodyParser(r); err != nil {
			return err
		}
		//models.Ejemplo.Update(uint(id),u)
		return c.Status(http.StatusOK).JSON("Registro actualizado")

	}
}*/
