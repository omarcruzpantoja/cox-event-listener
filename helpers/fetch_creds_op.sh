#!/bin/bash

# ---- TENANT ACCOUNT ----
TENANT_USERNAME=$(op item get $TENANT_ACCOUNT_ID --vault "$OP_VAULT" --fields TENANT_USERNAME --reveal)
TENANT_PRIVATE_SSH_KEY=$(op item get $TENANT_ACCOUNT_ID --vault "$OP_VAULT" --fields TENANT_PRIVATE_SSH_KEY | tr -d '"')
TENANT_ADDRESS=$(op item get $TENANT_ACCOUNT_ID --vault "$OP_VAULT" --fields TENANT_ADDRESS)

DISCORD_APPLICATION_ID=$(op item get $TENANT_ACCOUNT_ID --vault "$OP_VAULT" --fields DISCORD_APPLICATION_ID --reveal)
DISCORD_PUBLIC_KEY=$(op item get $TENANT_ACCOUNT_ID --vault "$OP_VAULT" --fields DISCORD_PUBLIC_KEY --reveal)
DISCORD_BOT_TOKEN=$(op item get $TENANT_ACCOUNT_ID --vault "$OP_VAULT" --fields DISCORD_BOT_TOKEN --reveal)

# ---- BACKEND ACCOUNT ----

# Define file paths
KEY_FILE="key.pem"
INVENTORY_FILE="ansible/inventory/hosts.ini"

# Create ansible inventory directory and temporary file
mkdir -p ansible/inventory
touch $INVENTORY_FILE
echo "[lxc-discord-bot-service]" >> "$INVENTORY_FILE"
echo $TENANT_ADDRESS >> "$INVENTORY_FILE"

# Create PRIVATE_SSH_KEY temporary file
touch $KEY_FILE
chmod 600 "$KEY_FILE"
echo -e "$TENANT_PRIVATE_SSH_KEY" >> "$KEY_FILE"


# Create tmporary environment variables file
touch .tmpenvs
echo "export TENANT_USERNAME=$TENANT_USERNAME" >> ".tmpenvs"

echo "export DISCORD_APPLICATION_ID=\"$DISCORD_APPLICATION_ID\"" >> ".tmpenvs"
echo "export DISCORD_PUBLIC_KEY=\"$DISCORD_PUBLIC_KEY\"" >> ".tmpenvs"
echo "export DISCORD_BOT_TOKEN=\"$DISCORD_BOT_TOKEN\"" >> ".tmpenvs"