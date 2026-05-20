# SecretServer E2E Setup

The SecretServer e2e suite creates disposable secrets in a real Secret Server
instance, syncs them through ESO, and deletes them after each test. The target
Secret Server instance and its IDs are environment-specific. Do not assume IDs
from a developer sandbox also work in QA or upstream CI.

## Required Access

Use a dedicated test account with the minimum permissions needed to:

- authenticate to the Secret Server API
- read the selected secret template
- create secrets in the selected folder
- read created secrets
- delete created secrets

The selected folder must be safe for disposable test data. Test secret names use
the `e2e-tests-eso-secretserver-` prefix.

## Environment Variables

Set these variables before running the suite:

```sh
export SECRETSERVER_URL="https://..."
export SECRETSERVER_USERNAME="..."
export SECRETSERVER_PASSWORD="..."

export SECRETSERVER_SITE_ID="..."
export SECRETSERVER_FOLDER_ID="..."
export SECRETSERVER_TEMPLATE_ID="..."
export SECRETSERVER_DATA_FIELD_ID="..."
export SECRETSERVER_REQUIRED_PASSWORD_FIELD_ID="..."
export SECRETSERVER_REQUIRED_PASSWORD_VALUE="external-secrets-e2e"
```

`SECRETSERVER_REQUIRED_PASSWORD_FIELD_ID` may be omitted or set to `0` only if
the selected template has no required password field. If the template has any
required fields, the e2e helper fills them with
`SECRETSERVER_REQUIRED_PASSWORD_VALUE`.

## Selecting IDs

Choose IDs for the specific Secret Server instance under test:

- `SECRETSERVER_SITE_ID`: site where test secrets are created
- `SECRETSERVER_FOLDER_ID`: writable folder for disposable test secrets
- `SECRETSERVER_TEMPLATE_ID`: readable template used to create test secrets
- `SECRETSERVER_DATA_FIELD_ID`: non-file template field where the suite stores
  the JSON test payload
- `SECRETSERVER_REQUIRED_PASSWORD_FIELD_ID`: required password field, if the
  template has one

A simple working template is a password-style template with one general text
field and one required password field. For example, in one developer sandbox the
working values were:

```sh
SECRETSERVER_SITE_ID=1
SECRETSERVER_FOLDER_ID=14
SECRETSERVER_TEMPLATE_ID=2
SECRETSERVER_DATA_FIELD_ID=60
SECRETSERVER_REQUIRED_PASSWORD_FIELD_ID=7
```

These values are examples only. Re-discover or provision equivalent values for
QA and CI.

## Discovering IDs

With `curl` and `jq`, a tester can inspect accessible records without exposing
credentials in logs:

```sh
token=$(
  curl -sk -X POST "$SECRETSERVER_URL/oauth2/token" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    --data-urlencode "username=$SECRETSERVER_USERNAME" \
    --data-urlencode "password=$SECRETSERVER_PASSWORD" \
    --data-urlencode "grant_type=password" \
  | jq -r .access_token
)

curl -sk -H "Authorization: Bearer $token" \
  "$SECRETSERVER_URL/api/v1/secrets?paging.take=30&paging.skip=0" \
| jq -r '.records[] | [.id, .name, .secretTemplateId, .folderId, .siteId] | @tsv'
```

Use a visible secret in the intended folder to identify candidate
`secretTemplateId`, `folderId`, and `siteId`. Then inspect the candidate
template:

```sh
curl -sk -H "Authorization: Bearer $token" \
  "$SECRETSERVER_URL/api/v1/secret-templates/$SECRETSERVER_TEMPLATE_ID" \
| jq -r '
  "template=\(.name) id=\(.id)",
  (.fields[] | [
    .secretTemplateFieldId,
    .fieldSlugName,
    .displayName,
    .isRequired,
    .isPassword,
    .isFile
  ] | @tsv)
'
```

Pick a non-file field for `SECRETSERVER_DATA_FIELD_ID`. If the template has a
required password field, set `SECRETSERVER_REQUIRED_PASSWORD_FIELD_ID` to that
field ID.

## Running Locally

Run only the SecretServer provider specs:

```sh
IMAGE_NAME=ghcr.io/external-secrets/external-secrets \
VERSION=manual-secretserver \
TEST_SUITES=provider \
GINKGO_LABELS=secretserver \
make -C e2e test
```

Expected selection:

```text
Will run 12 of 334 specs
```

A successful run ends with:

```text
SUCCESS! -- 12 Passed | 0 Failed | 0 Pending | 322 Skipped
```

## QA And CI Notes

For QA, create or identify a QA-owned Secret Server instance and repeat the ID
selection process above. Record the final environment variables in the QA test
run notes or the secure CI secret store, not in source control.

For upstream ESO CI, the same variables must exist as GitHub Actions secrets in
the upstream repository before maintainers run `/ok-to-test`. The CI Secret
Server instance may use different IDs than local or QA. The PR or work item
should document which environment was used and confirm that all required
variables were populated.

If every selected spec fails with `API_AccessDenied` at
`provider.go:70`, the account cannot read the configured template. Re-check
`SECRETSERVER_TEMPLATE_ID` and template permissions. If creation fails later,
re-check folder write/delete permissions and the field IDs.

