# EzMail - AWS SES Emailing made eazy

This project aims to simplify the integration with AWS Simple Email Service (SES) by providing a REST API for managing SES configurations, email templating, sending, and scheduling.

The project consists of two main components:

- **REST API**: Handles domain identity registration, email template creation and management, as well as immediate or scheduled email sending.
- **Event Processor**: A background service responsible for processing scheduled emails, handling retries, and managing domain validation.

This solution streamlines email operations with AWS SES, reducing complexity and improving efficiency.

## How to build

### To build the api binary execute

```bash
make build-api
```

This will create a binary in folder dist/api called **ezmail**


### To build the event processor execute

```bash
make build-executor
```

This will create a binary in folder dist/executor called **executor**


## Environment Variables

| Variable                           | Description                                                      | Required | Default                               |
|------------------------------------|------------------------------------------------------------------|----------|---------------------------------------|
| `MYSQL_DSN`                        | MySQL connection string                                          | Yes      | -                                     |
| `REDIS`                            | Redis server address                                             | Yes      | `localhost:6379`                      |
| `PORT`                             | API server port                                                  | No       | `3000`                                |
| `NODE_ID`                          | Unique node identifier, used for snowflake ID generation         | No       | `Random generated between 0 and 1023` |
| `AUTH_TOKEN`                       | API authentication token                                         | No       | -                                     |
| `REDIS_EVENTS_MAX_LEN`             | Max length for Redis events queue                                | No       | `2500`                                |
| `EVENTS_TOPIC`                     | Topic for email events                                           | No       | `topic:email.events`                  |
| `SCHEDULING_KEY`                   | Key for scheduled tasks queue, if not set scheduling is disabled | No       | -                                     |
| `RESCHEDULE_RETRIES`               | Number of retry attempts for failed emails                       | No       | -                                     |
| `RESCHEDULE_TIME_MS`               | Time in milliseconds between retries                             | No       | -                                     |
| `LOG_LEVEL`                        | Logging level (e.g., INFO)                                       | No       | `INFO`                                |
| `REFRESH_DOMAIN_RETRIES`           | Number of retries for domain refresh                             | No       | `12`                                  |
| `REFRESH_DOMAIN_RETRY_SEC_BETWEEN` | Time in seconds between domain refresh attempts                  | No       | `1800`                                |

Info: Both RESCHEDULE_RETRIES and RESCHEDULE_TIME_MS must be set to enable email retrying 

## API Documentation


### Authentication

If AUTH_TOKEN env variable is defined then all requests require an `Authorization` header with a Bearer token:

```
Authorization: Bearer YOUR_ACCESS_TOKEN
```

### Endpoints

#### Domains

| Method | Endpoint                        | Description                        |
| ------ | ------------------------------ | ---------------------------------- |
| GET    | `/domain`                       | Retrieve all domains               |
| GET    | `/domain/{domain_id}`           | Retrieve a single domain           |
| POST   | `/domain`                       | Create a new domain                |
| POST   | `/domain/{domain_id}/refresh`   | Refresh a domain status            |
| DELETE | `/domain/{domain_id}?full=true` | Delete a domain                    |

#### Create Domain

**Request Body:**

```json
{
    "name": "domain.tld",
    "region": "eu-west-1"
}
```

---

#### Templates

| Method | Endpoint    | Description       |
| ------ | ----------- | ----------------- |
| POST   | `/template` | Create a template |

#### Create Template

**Request Body:**

```json
{
    "html": "<div>Hello [[FIRST_NAME]]</div>",
    "text": "Hello [[FIRST_NAME]]",
    "subject": "Salutations from [[COMPANY_NAME]]"
}
```

---

#### Emails

| Method | Endpoint            | Description       |
| ------ | ------------------- | ----------------- |
| POST   | `/email`            | Send an email     |
| GET    | `/email/{email_id}` | Get email details |

#### Send Email

**Request Body:**

```json
{
    "from": "Name <source@domain.tld>"
    "to": [ "user@example.com" ],
    "template_id": "TEMPLATE_ID",
    "context": {
        "FIRST_NAME": "John",
        "COMPANY_NAME": "Acme Inc."
    },
    "reply_to": "email_to_reply@domain.com",
    "bcc": [ "bcc_email@somedomain.tld" ],
    "when": "2024/10/12 20:37:00"
}
```

Email fields follow the RFC 5322 format and can include a display-name. Example:

- name@domain.com : email without a display-name
- Name Surname &lt;name@domain.com&gt; : email with a display-name

**bcc**, **reply_to** and **when** fields are optional

Using the when field you can schedule an email for later, dates are UTC
