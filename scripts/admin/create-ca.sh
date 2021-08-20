# Set this variable to the name of your project.
PROJECT_ID="disco-idea-817"

# Note: this command is creating an account 
gcloud iam service-accounts create sa-claims \
	--description="Service account for claims admin" \
	--display-name="Claims Service account"
gcloud projects add-iam-policy-binding $PROJECT_ID \
	--member serviceAccount:sa-claims@$PROJECT_ID.iam.gserviceaccount.com \
	--role roles/firebase.sdkAdminServiceAgent
gcloud iam service-accounts keys create sa-claims.key \
	--iam-account sa-claims@$PROJECT_ID.iam.gserviceaccount.com