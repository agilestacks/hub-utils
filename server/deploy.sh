#!/bin/bash

NUMERIC_ID=$(gcloud projects list --filter="$PROJECT" --format="value(PROJECT_NUMBER)")
SECRET_FQN="projects/$NUMERIC_ID/secrets/$CLOUD_SQL_CREDENTIALS_SECRET/versions/latest"

gcloud functions deploy "$NAME" \
    --region "$REGION" \
    --runtime nodejs16 \
    --trigger-http \
    --allow-unauthenticated \
    --entry-point "$ENTRY_POINT" \
    --set-env-vars INSTANCE_CONNECTION_NAME="$INSTANCE_CONNECTION_NAME",DB_USER="$DB_USER",DB_NAME="$DB_NAME",CLOUD_SQL_CREDENTIALS_SECRET="$SECRET_FQN"
