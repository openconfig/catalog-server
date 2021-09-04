# Copyright 2021 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

if [ -z "$CLOUDSDK_CORE_PROJECT" ]; then
  echo 'Required environment variable $CLOUDSDK_CORE_PROJECT not set.'
  exit 1
fi

# Note: this command is creating a service account, and should only be run once per project.
# It's possible that you might encounter errors when running this command.
# As long as you obain the key after running this script, you are all set.
gcloud iam service-accounts create sa-claims \
	--description="Service account for claims admin" \
	--display-name="Claims Service account"
gcloud projects add-iam-policy-binding $CLOUDSDK_CORE_PROJECT \
	--member serviceAccount:sa-claims@$CLOUDSDK_CORE_PROJECT.iam.gserviceaccount.com \
	--role roles/firebase.sdkAdminServiceAgent
gcloud iam service-accounts keys create sa-claims.key \
	--iam-account sa-claims@$CLOUDSDK_CORE_PROJECT.iam.gserviceaccount.com
