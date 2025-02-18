# EzMail - AWS SES Emailing made eazy

This project aims to simplify the integration with AWS Simple Email Service (SES) by providing a REST API for managing SES configurations, email templating, sending, and scheduling.

The project consists of two main components:

- **REST API**: Handles domain identity registration, email template creation and management, as well as immediate or scheduled email sending.
- **Event Processor**: A background service responsible for processing scheduled emails, handling retries, and managing domain validation.

This solution streamlines email operations with AWS SES, reducing complexity and improving efficiency.

