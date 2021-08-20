# Set this variable to the name of your project.
PROJECT_ID="disco-idea-817"

# Note: this command is creating a service account, and should only be run once per project.
# It's possible that you might encounter errors when running this command.
# As long as you obain the key after running this script, you are all set.
gcloud iam service-accounts create sa-claims \
	--description="Service account for claims admin" \
	--display-name="Claims Service account"
gcloud projects add-iam-policy-binding $PROJECT_ID \
	--member serviceAccount:sa-claims@$PROJECT_ID.iam.gserviceaccount.com \
	--role roles/firebase.sdkAdminServiceAgent
gcloud iam service-accounts keys create sa-claims.key \
	--iam-account sa-claims@$PROJECT_ID.iam.gserviceaccount.com