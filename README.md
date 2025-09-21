## Cox event listener

A discord bot to listen ConquerX events and send messages with roles.

### Setup

1. Install go
2. Install dependencies 

   ```
   go mod download
   ```

3. Set up ENV variables
   ```
   export DISCORD_APPLICATION_ID=
   export DISCORD_PUBLIC_KEY=
   export DISCORD_BOT_TOKEN=
   export DISCORD_ACCOUNT_TOKEN=
   ```

4. Run the bot
   ```
   make run
   ```

### What to look into?

`parsers.go` contains most of the logic related to parsing the messages. This is the meat of the code. For the most part I don't expect much code changes here, for the most part you'd want to look over it to understand what's happening

`constants.go` This one is VERY important, in here you'll add the IDs for the roles, messages, channels...etc. You must configure this according to your server objects. If you don't configure this, the bot won't do anything.

`handlers` folder will contain logic for adding handlers to events.