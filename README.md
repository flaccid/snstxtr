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

```
POST /
{"phone": "+61406650430", "msg": "hello, there."}
```

If you use `--allow-get` then you can lazily send SMS:

```
GET /?phone=%2B61406650430&msg=Hello%2C%20there.
```

### Pingdom Webhooks

```
POST /pingdom-webhook/?recipients=%2B61406650430
{
  "check_id": 12345,
  "check_name": "Name of HTTP check",
  "check_type": "HTTP",
  "check_params": {
    "basic_auth": false,
    "encryption": true,
    "full_url": "https://www.example.com/path",
    "header": "User-Agent:Pingdom.com_bot",
    "hostname": "www.example.com",
    "ipv6": false,
    "port": 443,
    "url": "/path"
  },
  "tags": [
    "example_tag"
  ],
  "previous_state": "UP",
  "current_state": "DOWN",
  "importance_level": "HIGH",
  "state_changed_timestamp": 1451610061,
  "state_changed_utc_time": "2016-01-01T01:01:01",
  "long_description": "Long error message",
  "description": "Short error message",
  "first_probe": {
    "ip": "123.4.5.6",
    "ipv6": "2001:4800:1020:209::5",
    "location": "Stockholm, Sweden"
  },
  "second_probe": {
    "ip": "123.4.5.6",
    "ipv6": "2001:4800:1020:209::5",
    "location": "Austin, US",
    "version": 1
  }
}
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
