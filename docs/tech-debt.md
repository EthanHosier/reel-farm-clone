## Infra

Set up vpc enpdoints for accessing ecr and s3 instead of going through gateway

## Server

For the github actions deploy, currently using my laptops administrator tokens -> should create seperate iam role
Define actual openapi error types for the different cases
For the github actions, the health check is checking the alb url, not the cloudfront url

## Database

Enable RLS
Add separate DEV database
