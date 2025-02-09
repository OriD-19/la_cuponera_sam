package database

// Add all the methods supported by each of the stores

import (
	"OriD19/webdev2/types"
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBStore struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBClient(ctx context.Context, tableName string) *DynamoDBStore {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &DynamoDBStore{
		client:    client,
		tableName: tableName,
	}
}

// ************************************************************
// COUPON METHODS
// ************************************************************

func (d *DynamoDBStore) GetAllCoupons(ctx context.Context, nextToken *string) (types.CouponRange, error) {

	couponRange := types.CouponRange{
		Coupons: []types.Coupon{},
	}

	input := &dynamodb.QueryInput{
		TableName:              &d.tableName,
		Limit:                  aws.Int32(10), // for pagination purposes
		KeyConditionExpression: aws.String("entityType = :entityType"),
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":entityType": &ddbtypes.AttributeValueMemberS{
				Value: "coupon",
			},
		},
	}

	if nextToken != nil {
		input.ExclusiveStartKey = map[string]ddbtypes.AttributeValue{
			"id": &ddbtypes.AttributeValueMemberS{
				Value: *nextToken,
			},
		}
	}

	result, err := d.client.Query(ctx, input)

	if err != nil {
		// return empty coupon range and error
		return couponRange, err
	}

	// take the result from dynamodb, and put it into the coupon range
	err = attributevalue.UnmarshalListOfMaps(result.Items, &couponRange.Coupons)

	if err != nil {
		return couponRange, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	if len(result.LastEvaluatedKey) > 0 {
		if key, ok := result.LastEvaluatedKey["id"]; ok {
			couponRange.Next = &key.(*ddbtypes.AttributeValueMemberS).Value
		}
	}

	return couponRange, nil
}

func (d *DynamoDBStore) GetAllCouponsFromCategory(ctx context.Context, category string) (types.CouponRange, error) {
	couponRange := types.CouponRange{
		Coupons: []types.Coupon{},
	}

	// DONE: Check this thing, because Copilot is terrible
	input := &dynamodb.QueryInput{
		TableName:              &d.tableName,
		KeyConditionExpression: aws.String("entityType = :entityType"),
		// More read operations, maybe an index would be useful...
		FilterExpression: aws.String("#category = :category"), // filter only by category
		ExpressionAttributeNames: map[string]string{
			"#category": "category",
		},
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":entityType": &ddbtypes.AttributeValueMemberS{
				Value: "coupon",
			},
			":category": &ddbtypes.AttributeValueMemberS{
				Value: category,
			},
		},
	}

	result, err := d.client.Query(ctx, input)

	if err != nil {
		// return empty coupon range and error
		return couponRange, err
	}

	// take the result from dynamodb, and put it into the coupon range
	err = attributevalue.UnmarshalListOfMaps(result.Items, &couponRange.Coupons)

	if err != nil {
		return couponRange, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	return couponRange, nil
}

func (d *DynamoDBStore) GetCoupon(c context.Context, id string) (types.Coupon, error) {
	// query a single coupon with the GetItem API. Better resource (RCU) efficiency
	input := &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"entityType": &ddbtypes.AttributeValueMemberS{
				Value: "coupon",
			},
			"id": &ddbtypes.AttributeValueMemberS{
				Value: id,
			},
		},
	}

	result, err := d.client.GetItem(c, input)

	if err != nil {
		return types.Coupon{}, err
	}

	var coupon types.Coupon
	err = attributevalue.UnmarshalMap(result.Item, &coupon)

	if err != nil {
		return types.Coupon{}, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	return coupon, nil
}

func (d *DynamoDBStore) PutCoupon(c context.Context, coupon types.Coupon) error {
	coupon.EntityType = "coupon"
	av, err := attributevalue.MarshalMap(coupon)

	if err != nil {
		return fmt.Errorf("failed to marshal coupon, %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      av,
	}

	_, err = d.client.PutItem(c, input)

	if err != nil {
		return fmt.Errorf("failed to put coupon, %v", err)
	}

	return nil
}

func (d *DynamoDBStore) PutGeneratedOffer(c context.Context, offer types.GeneratedOffer) error {
	offer.EntityType = "generatedOffer"
	av, err := attributevalue.MarshalMap(offer)

	if err != nil {
		return fmt.Errorf("failed to marshal generated offer, %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      av,
	}

	_, err = d.client.PutItem(c, input)

	if err != nil {
		return fmt.Errorf("failed to put generated offer, %v", err)
	}

	// we need to decrease the count of available coupons

	// get the coupon associated with the generated offer
	coupon, err := d.GetCoupon(c, offer.CouponId)

	if err != nil {
		return fmt.Errorf("failed to get coupon, %v", err)
	}

	// decrease the count of available coupons
	coupon.AvailableCoupons--

	// update the coupon in the database
	err = d.PutCoupon(c, coupon)

	if err != nil {
		return fmt.Errorf("failed to update coupon availability, %v", err)
	}

	return nil
}

func (d *DynamoDBStore) GenerateId(c context.Context, enterpriseId string) (string, error) {
	// generate a random ID for the generated offer

	// 7-digit random number for the code
	randInt, err := rand.Int(rand.Reader, big.NewInt(9999999))

	if err != nil {
		return "", fmt.Errorf("failed to generate random number, %v", err)
	}

	// get the enterprise associated with the coupon
	enterprise, err := d.GetEnterprise(c, enterpriseId)

	if err != nil {
		return "", fmt.Errorf("failed to get enterprise, %v", err)
	}

	return fmt.Sprintf("%s%d", enterprise.EnterpriseCode, randInt), nil

}

// only clients can buy a coupon
func (d *DynamoDBStore) BuyCoupon(c context.Context, couponId string, userId string) (types.GeneratedOffer, error) {

	coupon, err := d.GetCoupon(c, couponId)

	if err != nil {
		return types.GeneratedOffer{}, fmt.Errorf("failed to get coupon, %v", err)
	}

	// check if the coupon is still available
	if coupon.AvailableCoupons <= 0 {
		return types.GeneratedOffer{}, fmt.Errorf("coupon is not available")
	}

	// get the user associated with the coupon
	user, err := d.GetClient(c, userId)

	if err != nil {
		return types.GeneratedOffer{}, fmt.Errorf("failed to get user, %v", err)
	}

	var newGenOffer types.GeneratedOffer

	// generate a new ID for the generated offer
	generatedId, err := d.GenerateId(c, coupon.EnterpriseId)

	if err != nil {
		return types.GeneratedOffer{}, fmt.Errorf("failed to generate ID, %v", err)
	}

	newGenOffer.Id = generatedId
	newGenOffer.UserId = user.Email
	newGenOffer.CouponId = coupon.Id
	newGenOffer.GeneratedAt = time.Now()
	newGenOffer.ExpirationDate = coupon.ValidUntil
	newGenOffer.Redeemed = false

	// save the generated offer in the database
	err = d.PutGeneratedOffer(c, newGenOffer)

	if err != nil {
		return types.GeneratedOffer{}, fmt.Errorf("failed to put generated offer, %v", err)
	}

	return newGenOffer, nil
}

// get the user ID from a route parameter
func (d *DynamoDBStore) GetUserOffers(c context.Context, userId string) (types.OfferRange, error) {
	// query all the generated offers for a given user
	offers := types.OfferRange{
		Offers: []types.GeneratedOffer{},
	}
	input := &dynamodb.QueryInput{
		TableName:              &d.tableName,
		KeyConditionExpression: aws.String("entityType = :entityType"),
		FilterExpression:       aws.String("#userId = :userId"),
		ExpressionAttributeNames: map[string]string{
			"#userId": "userId",
		},
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":entityType": &ddbtypes.AttributeValueMemberS{
				Value: "generatedOffer",
			},
			":userId": &ddbtypes.AttributeValueMemberS{
				Value: userId,
			},
		},
	}

	result, err := d.client.Query(c, input)

	if err != nil {
		return offers, err
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, &offers)

	if err != nil {
		return offers, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	return offers, nil
}

func (d *DynamoDBStore) GetGeneratedOffer(c context.Context, id string) (types.GeneratedOffer, error) {
	// query a single generated offer with the GetItem API. Better resource (RCU) efficiency
	input := &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"entityType": &ddbtypes.AttributeValueMemberS{
				Value: "generatedOffer",
			},
			"id": &ddbtypes.AttributeValueMemberS{
				Value: id,
			},
		},
	}

	result, err := d.client.GetItem(c, input)

	if err != nil {
		return types.GeneratedOffer{}, err
	}

	var offer types.GeneratedOffer
	err = attributevalue.UnmarshalMap(result.Item, &offer)

	if err != nil {
		return types.GeneratedOffer{}, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	return offer, nil
}

func (d *DynamoDBStore) RedeemCoupon(c context.Context, id string) error {
	// get the generated offer
	offer, err := d.GetGeneratedOffer(c, id)

	if err != nil {
		return fmt.Errorf("failed to get generated offer, %v", err)
	}

	// check if the offer is still valid
	if offer.ExpirationDate.Before(time.Now()) {
		return fmt.Errorf("offer is expired")
	}

	// check if the offer is already redeemed
	if offer.Redeemed {
		return fmt.Errorf("offer is already redeemed")
	}

	// mark the offer as redeemed
	offer.Redeemed = true

	// update the offer in the database
	err = d.PutGeneratedOffer(c, offer)

	if err != nil {
		return fmt.Errorf("failed to update generated offer, %v", err)
	}

	return nil
}

// ************************************************************
// USER METHODS
// ************************************************************

// !IMPORTANT: REGISTER METHODS ALSO UPDATES IF THE VALUE ALREADY EXISTS
func (d *DynamoDBStore) RegisterClient(c context.Context, client types.Client) error {
	client.EntityType = "client"
	av, err := attributevalue.MarshalMap(client)

	if err != nil {
		return fmt.Errorf("failed to marshal client, %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      av,
	}

	_, err = d.client.PutItem(c, input)

	if err != nil {
		return fmt.Errorf("failed to put client, %v", err)
	}

	return nil
}

func (d *DynamoDBStore) RegisterEnterprise(c context.Context, enterprise types.Enterprise) error {
	enterprise.EntityType = "enterprise"
	av, err := attributevalue.MarshalMap(enterprise)

	if err != nil {
		return fmt.Errorf("failed to marshal enterprise, %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      av,
	}

	_, err = d.client.PutItem(c, input)

	if err != nil {
		return fmt.Errorf("failed to put enterprise, %v", err)
	}

	return nil
}

func (d *DynamoDBStore) RegisterAdministrator(c context.Context, administrator types.Administrator) error {
	administrator.EntityType = "administrator"
	av, err := attributevalue.MarshalMap(administrator)

	if err != nil {
		return fmt.Errorf("failed to marshal administrator, %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      av,
	}

	_, err = d.client.PutItem(c, input)

	if err != nil {
		return fmt.Errorf("failed to put administrator, %v", err)
	}

	return nil
}

func (d *DynamoDBStore) RegisterEmployee(c context.Context, employee types.Employee) error {
	employee.EntityType = "employee"
	av, err := attributevalue.MarshalMap(employee)

	if err != nil {
		return fmt.Errorf("failed to marshal employee, %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      av,
	}

	_, err = d.client.PutItem(c, input)

	if err != nil {
		return fmt.Errorf("failed to put employee, %v", err)
	}

	return nil
}

// we're using an username as the ID field inside of the database to make login easier
func (d *DynamoDBStore) GetClient(c context.Context, username string) (types.Client, error) {
	// query a single client with the GetItem API. Better resource (RCU) efficiency
	input := &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"entityType": &ddbtypes.AttributeValueMemberS{
				Value: "client",
			},
			"id": &ddbtypes.AttributeValueMemberS{
				Value: username,
			},
		},
	}

	result, err := d.client.GetItem(c, input)

	if err != nil {
		return types.Client{}, err
	}

	if len(result.Item) == 0 {
		return types.Client{}, fmt.Errorf("client not found")
	}

	var client types.Client
	err = attributevalue.UnmarshalMap(result.Item, &client)

	if err != nil {
		return types.Client{}, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	return client, nil
}

func (d *DynamoDBStore) GetEnterprise(c context.Context, id string) (types.Enterprise, error) {
	// query a single enterprise with the GetItem API. Better resource (RCU) efficiency
	input := &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"entityType": &ddbtypes.AttributeValueMemberS{
				Value: "enterprise",
			},
			"id": &ddbtypes.AttributeValueMemberS{
				Value: id,
			},
		},
	}

	result, err := d.client.GetItem(c, input)

	if err != nil {
		return types.Enterprise{}, err
	}

	var enterprise types.Enterprise
	err = attributevalue.UnmarshalMap(result.Item, &enterprise)

	if err != nil {
		return types.Enterprise{}, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	return enterprise, nil
}

func (d *DynamoDBStore) GetAdministrator(c context.Context, id string) (types.Administrator, error) {
	// query a single administrator with the GetItem API. Better resource (RCU) efficiency
	input := &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"entityType": &ddbtypes.AttributeValueMemberS{
				Value: "administrator",
			},
			"id": &ddbtypes.AttributeValueMemberS{
				Value: id,
			},
		},
	}

	result, err := d.client.GetItem(c, input)

	if err != nil {
		return types.Administrator{}, err
	}

	var administrator types.Administrator
	err = attributevalue.UnmarshalMap(result.Item, &administrator)

	if err != nil {
		return types.Administrator{}, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	return administrator, nil
}

func (d *DynamoDBStore) GetEmployee(c context.Context, id string) (types.Employee, error) {
	// query a single employee with the GetItem API. Better resource (RCU) efficiency
	input := &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"entityType": &ddbtypes.AttributeValueMemberS{
				Value: "employee",
			},
			"id": &ddbtypes.AttributeValueMemberS{
				Value: id,
			},
		},
	}

	result, err := d.client.GetItem(c, input)

	if err != nil {
		return types.Employee{}, err
	}

	var employee types.Employee
	err = attributevalue.UnmarshalMap(result.Item, &employee)

	if err != nil {
		return types.Employee{}, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	return employee, nil
}
