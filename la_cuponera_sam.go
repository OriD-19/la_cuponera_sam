package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type LaCuponeraSamStackProps struct {
	awscdk.StackProps
}

func NewLaCuponeraSamStack(scope constructs.Construct, id string, props *LaCuponeraSamStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	table := awsdynamodb.NewTable(stack, jsii.String("LaCuponeraSamTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("entityType"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName:     jsii.String("LaCuponeraTable"),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	// generate three lamdbas, one for each type of functionality in the API:
	// - Managing coupons and offers
	// - Managing users
	// - Managing the login system

	// Coupons and offers
	couponsLambda := awslambda.NewFunction(stack, jsii.String("LaCuponeraCouponsLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("lambda/functions/couponFunction/couponsFunction.zip"), nil),
		Environment: &map[string]*string{
			"TABLE_NAME": table.TableName(),
		},
	})

	// Users
	usersLambda := awslambda.NewFunction(stack, jsii.String("LaCuponeraUsersLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("lambda/functions/userFunction/usersFunction.zip"), nil),
		Environment: &map[string]*string{
			"TABLE_NAME": table.TableName(),
		},
	})

	// Login
	loginLambda := awslambda.NewFunction(stack, jsii.String("LaCuponeraLoginLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromAsset(jsii.String("lambda/functions/loginFunction/loginFunction.zip"), nil),
		Environment: &map[string]*string{
			"TABLE_NAME": table.TableName(),
		},
	})

	table.GrantReadWriteData(couponsLambda)
	table.GrantReadWriteData(usersLambda)
	table.GrantReadWriteData(loginLambda)

	// Finally, create the integration with the API Gateway

	api := awsapigateway.NewRestApi(stack, jsii.String("LaCuponeraApi"), &awsapigateway.RestApiProps{
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowOrigins: jsii.Strings("*"),
			AllowMethods: jsii.Strings("GET", "POST", "PUT", "DELETE"),
			// remember modifying the headers if necessary
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
		},
		DeployOptions: &awsapigateway.StageOptions{
			LoggingLevel: awsapigateway.MethodLoggingLevel_INFO,
			StageName:    jsii.String("v1"),
		},
		RestApiName: jsii.String("LaCuponeraApi"),
	})

	// Create integrations with the lambda functions
	couponsIntegration := awsapigateway.NewLambdaIntegration(couponsLambda, nil)
	usersIntegration := awsapigateway.NewLambdaIntegration(usersLambda, nil)
	loginIntegration := awsapigateway.NewLambdaIntegration(loginLambda, nil)

	// **********************
	// CREATE THE RESOURCES
	// **********************

	// Create the resources
	// GET /coupons
	couponsResource := api.Root().AddResource(jsii.String("coupons"), nil)
	couponsResource.AddMethod(jsii.String("GET"), couponsIntegration, nil)

	// TODO POST /coupons
	//couponsResource.AddMethod(jsii.String("POST"), couponsIntegration, nil)

	// GET /coupons/{id}
	couponsResource.AddResource(jsii.String("{id}"), nil).
		AddMethod(jsii.String("GET"), couponsIntegration, nil)

	// PUT /coupons/{id}
	couponsResource.GetResource(jsii.String("{id}")).
		AddMethod(jsii.String("PUT"), couponsIntegration, nil)

	// get coupons by category
	// GET /coupons/category/{category}
	couponsResource.AddResource(jsii.String("category"), nil).
		AddResource(jsii.String("{category}"), nil).
		AddMethod(jsii.String("GET"), couponsIntegration, nil)

	// buy a coupon
	// POST /coupons/{id}/buy
	couponsResource.GetResource(jsii.String("{id}")).
		AddResource(jsii.String("buy"), nil).
		AddMethod(jsii.String("POST"), couponsIntegration, nil)

	// Offers resources

	// since offers only work for a given user id, we can query them directly as a parameter path
	// GET /offers/{userId}
	offersResource := api.Root().AddResource(jsii.String("offers"), nil)
	offersResource.AddResource(jsii.String("{userId}"), nil).
		AddMethod(jsii.String("GET"), couponsIntegration, nil)

	// get offer details
	// GET /offers/{offerId}
	offersResource.AddResource(jsii.String("{offerId}"), nil).
		AddResource(jsii.String("{userId}"), nil).
		AddMethod(jsii.String("GET"), couponsIntegration, nil)

	// redeem a coupon
	// POST /offers/{offerId}/redeem
	offersResource.AddResource(jsii.String("{offerId}"), nil).
		AddResource(jsii.String("redeem"), nil).
		AddMethod(jsii.String("POST"), couponsIntegration, nil)

	// Users resources

	// GET /users
	usersResource := api.Root().AddResource(jsii.String("users"), nil)
	usersResource.AddMethod(jsii.String("GET"), usersIntegration, nil)

	// register a new user of type client
	// POST /users/client
	usersResource.AddResource(jsii.String("client"), nil).
		AddMethod(jsii.String("POST"), usersIntegration, nil)

	// GET /users/{id}
	usersResource.AddResource(jsii.String("{id}"), nil).
		AddMethod(jsii.String("GET"), usersIntegration, nil)

	// view profile for client
	// GET /users/{id}/profile
	usersResource.
		AddResource(jsii.String("{id}"), nil).
		AddResource(jsii.String("profile"), nil).
		AddMethod(jsii.String("GET"), usersIntegration, nil)

	// update profile for client
	// PUT /users/{id}/profile
	usersResource.
		AddResource(jsii.String("{id}"), nil).
		AddResource(jsii.String("profile"), nil).
		AddMethod(jsii.String("PUT"), usersIntegration, nil)

	// login resources
	// POST /login/client
	loginResource := api.Root().AddResource(jsii.String("login"), nil)
	loginResource.
		AddResource(jsii.String("client"), nil).
		AddMethod(jsii.String("POST"), loginIntegration, nil)

	// POST login/employee
	loginResource.
		AddResource(jsii.String("employee"), nil).
		AddMethod(jsii.String("POST"), loginIntegration, nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewLaCuponeraSamStack(app, "LaCuponeraSamStack", &LaCuponeraSamStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
