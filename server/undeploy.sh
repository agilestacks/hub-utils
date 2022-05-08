#!/bin/bash

gcloud functions delete "$NAME" --region "$REGION" --quiet
