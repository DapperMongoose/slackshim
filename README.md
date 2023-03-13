# slackshim
An app that takes an event as a POST request and turns it in to a Slack message if it's of the type "SpamNotification".

The input parameters look like this:
```
{
  "RecordType": "string",
  "Type": "string",
  "TypeCode": int,
  "Name": "string",
  "Tag": "string",
  "MessageStream": "string",
  "Description": "string",
  "Email": "string",
  "From": "string",
  "BouncedAt": "string",
}
```


# Setup in Slack
You'll need to configure an app with an incoming webhook in Slack.  The documentation for this process can be found here: https://api.slack.com/messaging/webhooks


# To Deploy: 
Simply run go build and then execute the resulting "shim" binary with the Slack Webhook URL provided as the URL parameter.

IE: ```./shim -url https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX```
