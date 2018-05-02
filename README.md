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
