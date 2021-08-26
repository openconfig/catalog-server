### How to deploy catalog server on GCP

+ Follow [instructions](https://cloud.google.com/sql/docs/postgres/quickstart) to set up your postgres database on GCP.
+ Create database tables by running commands in [schema](schema) directory in your postgres database. See this detailed [instruction](https://cloud.google.com/sql/docs/postgres/connect-overview) on how to connect to your database. Our recommended way is to use a Cloud SQL Auth proxy, please see this [instruction](https://cloud.google.com/sql/docs/postgres/connect-admin-proxy).
+ Deploy catalog server on GCP following this [instruction](https://cloud.google.com/run/docs/quickstarts/build-and-deploy/go). This run would fail due to that you haven't set up related environment variables and connection to postgres database.
+ Set up connection from your launched cloud run instance to your postgres database following this [instruction](https://cloud.google.com/sql/docs/postgres/connect-run).
+ Set up environment variables that are required in `pkg/db` in the cloud run instance you have just launched following this [instruction](https://cloud.google.com/run/docs/configuring/environment-variables). That includes `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PWD`, `DB_NAME`. See [pkg/db/db.go](../pkg/db/db.go)'s comments for more details about these variables.
+ Change `CLOUD_RUN_URL` in both [query.html](../frontend/static/query.html), [update.html](../frontend/static/update.html) to the URL of your launched cloud run instance.
+ The catalog server should be running after all stpes above.