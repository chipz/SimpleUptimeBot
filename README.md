# SimpleUptimeBot
Simple Telegram Bot to monitor uptime of your service made with Go

## Installation:
1. Clone this project
2. Create a [telegram bot](https://core.telegram.org/bots#creating-a-new-bot)
3. Get your bot token, it will be in format something like this
```
379261111:11111xwyq1CO1JwZNeeOqSq1AygnhVL11111
```
4. start it as:
SimpleUptimeBot "token:forBot"
or
BOT_TOKEN="token:forBot" SimpleUptimeBot

## Motivation:
I have services that needs to the uptime should be monitored, and give me a report
once the service is **down**. I want the report will be in Telegram. Using
the existing services will be overkill for something as simple as this.

##TODO:
- Use a persistent database to manage the data
- Add mode to monitor certain port
- Add checking for valid url
- Add command to `list` and `remove` the added urls
