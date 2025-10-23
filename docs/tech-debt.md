## Infra

Set up vpc enpdoints for accessing ecr and s3 instead of going through gateway

## Server

For the github actions deploy, currently using my laptops administrator tokens -> should create seperate iam role
Define actual openapi error types for the different cases
For the github actions, the health check is checking the alb url, not the cloudfront url
Am I using the transactions correctly for the use repo sqlc?
Add idempotency and race condition protection for spending credits
Use signed urls for the ai avatar thumbnails and vids?
Refac the ai-generate video service to move the signing of urls into the service layer and the db types shouldnt leak into the handler?

## Stripe

Race conditions during payment account creation

## Database

Enable RLS
Add separate DEV database
