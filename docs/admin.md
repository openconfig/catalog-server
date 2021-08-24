### How to grant accounts with access to different organizations

+ First, please ensure that your catalog server is running on GCP correctly.
+ Enable identity platform on GCP following this [instruction](https://cloud.google.com/identity-platform/docs/quickstart-email-password). It also gives instruction on how to sign in users with email and password, but it's also doable to use other signing in methods supported by identity platform, including signing in using Google account, LinkedIn account.
+ Replace the boilerplate codes in [update.html](../frontend/static/update.html) from `line 19` to `line 28` with that from your identity platform. See this [instruction](https://cloud.google.com/identity-platform/docs/quickstart-email-password) for more details.
+ See [scripts/admin](../scripts/admin) directory for details on how to grant access to different accounts.