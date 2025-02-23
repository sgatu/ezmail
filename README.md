# EzMail - AWS SES Emailing made eazy

This project aims to simplify the integration with AWS Simple Email Service (SES) by providing a REST API for managing SES configurations, email templating, sending, and scheduling.

The project consists of two main components:

- **REST API**: Handles domain identity registration, email template creation and management, as well as immediate or scheduled email sending.
- **Event Processor**: A background service responsible for processing scheduled emails, handling retries, and managing domain validation.

This solution streamlines email operations with AWS SES, reducing complexity and improving efficiency.



## API Documentation

### Base URL

```
http://localhost:3000
```

### Authentication

All requests require an `Authorization` header with a Bearer token:

```
Authorization: Bearer YOUR_ACCESS_TOKEN
```

### Endpoints

#### Domains

| Method | Endpoint                        | Description              |
| ------ | ------------------------------- | ------------------------ |
| GET    | `/domain`                       | Retrieve all domains     |
| GET    | `/domain/{domain_id}`           | Retrieve a single domain |
| POST   | `/domain`                       | Create a new domain      |
| POST   | `/domain/{domain_id}/refresh`   | Refresh a domain         |
| DELETE | `/domain/{domain_id}?full=true` | Delete a domain          |

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
    "to": "user@example.com",
    "template_id": "TEMPLATE_ID",
    "variables": {
        "FIRST_NAME": "John",
        "COMPANY_NAME": "Acme Inc."
    }
}
```
