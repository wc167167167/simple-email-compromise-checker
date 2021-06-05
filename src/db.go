package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var ses, _ = session.NewSession(&aws.Config{
	Endpoint:    aws.String("http://host.docker.internal:8990"),
	Region:      aws.String("ap-southeast-1"),
	Credentials: credentials.NewStaticCredentials("test", "test", "test"),
},
)

var db = dynamodb.New(ses)
var tableName = "CompromisedEmails"

func createTable() {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("prefix"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("email"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("prefix"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("email"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("CompromisedEmails"),
	}

	result, err := db.CreateTable(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceInUseException:
				fmt.Println(dynamodb.ErrCodeResourceInUseException, aerr.Error())
			case dynamodb.ErrCodeLimitExceededException:
				fmt.Println(dynamodb.ErrCodeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

type Item struct {
	Prefix string `dynamodbav:"prefix"`
	Email  string `dynamodbav:"email"`
}

func fillTestingData() {
	for i := 0; i < 100; i++ {
		item := Item{
			Prefix: fmt.Sprint("jx", i),
			Email:  fmt.Sprint("jx", i, "compromised@fake_email.com"),
		}

		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			fmt.Println(err.Error())
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}
		db.PutItem(input)
	}
}

func Init() {
	createTable()
	fillTestingData()
}

func IsExist(email string) bool {

	if len(email) < 3 {
		return false
	}

	prefix := email[0:3]

	primaryKey := map[string]string{
		"prefix": prefix,
		"email":  email,
	}

	key, _ := dynamodbattribute.MarshalMap(primaryKey)

	proj := expression.NamesList(expression.Name("prefix"), expression.Name("email"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		return false
	}

	input := &dynamodb.GetItemInput{
		TableName:                aws.String(tableName),
		Key:                      key,
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
	}

	result, _ := db.GetItem(input)

	return len(result.Item) != 0
}
