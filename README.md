# La Cuponera - AWS SAM Architecture

Serverless Architecture Modeol for La Cuponera, a fictional enterprise for managing coupons and discounts.
This whole project is meant to be deployed and integrated into the AWS ApiGateway service, alongside with 
AWS Lambda for functionality computing and DynamoDB for entity storage.

## CDK Usage

For deploying the proyect, simply compile the programs inside `lambda/functions/*` and zip them into a 
.zip file. This generates the asset for the Lambda function.
The whole infrastructure is defined using the AWS CDK for Go, just for convenience in the deployment.

## Endpoints

For a full list of endpoints, refer to the AWS ApiGateway documentation. The hierarchy looks something like 
the following:

![Resource Hierarchy displayed in the AWS ApiGateway panel](./resource-hierarchy.PNG)

### Requisites

For building this project, the following programs and their specific versions were used:
- Go (1.23.6)
- AWS CDK (2.178.1)
- Operating System: Windows (Linux for Lambda binaries)