package aws

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mugayoshi/currency_rate_tracker/helpers"
)

type CurrencyRateData struct {
	Date string  `dynamodbav:"date"`
	Rate float32 `dynamodbav:"rate"`
	Base string  `dynamodbav:"base"`
}
type CurrencyData struct {
	LastData CurrencyRateData `dynamodbav:"last-data"`
	Minimum  CurrencyRateData `dynamodbav:"minimum"`
}

type ResponseJson struct {
	JsonData CurrencyData `dynamodbav:"json-data"`
}

const TABLE_NAME = "currency-rate"

func GetDynamoDbClient() *dynamodb.Client {
	accessKey := helpers.GetEnvVariable("DYNAMO_API")
	secret := helpers.GetEnvVariable("DYNAMO_SECRET")
	region := helpers.GetEnvVariable("DYNAMO_REGION")
	credentials := credentials.NewStaticCredentialsProvider(accessKey, secret, "")
	client := dynamodb.NewFromConfig(aws.Config{
		Region:      region,
		Credentials: credentials,
	})
	return client
}

func GetCurrencyItem(client *dynamodb.Client, currency string) (CurrencyData, error) {
	item, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]types.AttributeValue{
			"name": &types.AttributeValueMemberS{Value: currency},
		},
	})
	if err != nil {
		var d CurrencyData
		return d, errors.New("api error")
	}
	if item.Item == nil {
		var d CurrencyData
		return d, errors.New("not found")
	}
	// fmt.Printf("%+v\n", item.Item)
	actual := ResponseJson{}
	e := attributevalue.UnmarshalMap(item.Item, &actual)
	if e != nil {
		var d CurrencyData
		return d, errors.New("unmarshal error")
	}

	return actual.JsonData, nil
}

func updateItemBase(client *dynamodb.Client, input *dynamodb.UpdateItemInput) (bool, error) {
	log.Println("update dynamo db")
	_, err := client.UpdateItem(context.TODO(), input)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return false, errors.New("update api error")
	}
	log.Println("update dynamo db is done")
	return true, nil
}

func UpdateLastData(client *dynamodb.Client, currency string, rate float32, date string) (bool, error) {
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]types.AttributeValue{
			"name": &types.AttributeValueMemberS{Value: currency},
		},
		UpdateExpression: aws.String("SET #json.#last.#rate = :rate, #json.#last.#date = :date"),
		ExpressionAttributeNames: map[string]string{
			"#json": "json-data",
			"#last": "last-data",
			"#rate": "rate",
			"#date": "date",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":rate": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", rate)},
			":date": &types.AttributeValueMemberS{Value: date},
		},
	}
	return updateItemBase(client, updateInput)
}

func UpdateMinimumRate(client *dynamodb.Client, currency string, rate float32, date string) (bool, error) {
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]types.AttributeValue{
			"name": &types.AttributeValueMemberS{Value: currency},
		},
		UpdateExpression: aws.String("SET #json.#min.#rate = :rate, #json.#min.#date = :date"),
		ExpressionAttributeNames: map[string]string{
			"#json": "json-data",
			"#min":  "minimum",
			"#rate": "rate",
			"#date": "date",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":rate": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", rate)},
			":date": &types.AttributeValueMemberS{Value: date},
		},
	}
	return updateItemBase(client, updateInput)
}
