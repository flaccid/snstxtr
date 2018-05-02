# snstxtr

Send SMS using SNS by CLI or a simple web service.

## Usage

    $ snstxtr --help

### Environment Variables

These can be used as an alternative or in conjunction with the applicable CLI options.

- `$PHONE` - phone number to send the SMS message to
- `$AWS_REGION` - AWS region for the SNS service
- `$AWS_ACCESS_KEY_ID` - AWS access key id
- `$AWS_SECRET_ACCESS_KEY` - AWS secret access key

## Daemon Mode

Example:

```
GET /?phone=%2B61406650430&msg=Hello%2C%20there.
```

## IAM

Example policy:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Deny",
            "Action": [
                "sns:Publish"
            ],
            "Resource": "arn:aws:sns:*:*:*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "sns:Publish"
            ],
            "Resource": "*"
        }
    ]
}
```

### Terraform

Example that also stores the creds in vault:

```
resource "aws_iam_user" "sns-sms" {
  name = "${var.environment}-sns-sms"
}

resource "aws_iam_access_key" "sns-sms" {
  user = "${aws_iam_user.sns-sms.name}"
}

resource "vault_generic_secret" "sns-sms-key-id" {
  path = "secret/aws/${var.environment}/sns-sms/AWS_ACCESS_KEY_ID"

  data_json = <<EOT
{ "value": "${aws_iam_access_key.sns-sms.id}" }
EOT
}

resource "vault_generic_secret" "sns-sms-secret" {
  path = "secret/aws/${var.environment}/sns-sms/AWS_SECRET_ACCESS_KEY"

  data_json = <<EOT
{ "value": "${aws_iam_access_key.sns-sms.secret}" }
EOT
}

resource "aws_iam_user_policy" "sns-sms" {
  name = "${var.environment}-sns-sms"
  user = "${aws_iam_user.sns-sms.name}"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Deny",
            "Action": [
                "sns:Publish"
            ],
            "Resource": "arn:aws:sns:*:*:*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "sns:Publish"
            ],
            "Resource": "*"
        }
    ]
}
EOF
}
```
