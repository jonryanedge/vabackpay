# VA Disability Backpay Calculator

A web-based calculator for veterans to estimate their VA disability backpay based on effective date, disability rating, and dependents.

## Features

- Calculate VA disability backpay from any start date (2000-2025)
- Support for all disability ratings (10%-100%)
- Dependent calculations (for ratings 30% and above)
- Year-by-year breakdown with monthly rates
- Dark mode UI
- Email results to user
- Responsive design

## Tech Stack

- **Backend**: Go (Golang)
- **Frontend**: HTML, CSS, HTMX
- **Email**: gomail v2

## Prerequisites

- Go 1.21 or higher
- Git

## Setup

1. Clone the repository:
```bash
git clone https://github.com/jonryanedge/vabackpay.git
cd vabackpay
```

2. Copy the example environment file:
```bash
cp example.env .env
```

3. Edit `.env` with your configuration (see Configuration section below)

4. Install dependencies:
```bash
go mod download
```

## Running

### Development
```bash
go run main.go
```

The server will start at `http://localhost:8080`

### Production Build
```bash
go build -o vabackpay main.go
./vabackpay
```

## Configuration

Edit the `.env` file to configure email and application settings:

| Variable | Description | Default |
|----------|-------------|---------|
| `DEBUG` | Set to `true` to test without sending emails | `true` |
| `SMTP_HOST` | SMTP server hostname | - |
| `SMTP_PORT` | SMTP server port | `587` |
| `SMTP_USERNAME` | SMTP authentication username | - |
| `SMTP_PASSWORD` | SMTP authentication password | - |
| `FROM_EMAIL` | Sender email address | - |
| `FROM_NAME` | Sender display name | `VA Backpay Calculator` |
| `APP_URL` | Application URL for email links | `http://localhost:8080` |

### Email Testing

When `DEBUG=true`, emails are logged but not sent. To test email functionality:

1. Set `DEBUG=true` in `.env`
2. Run the application
3. Submit an email through the UI
4. Check server logs for email output

### Production Email

For production use:
1. Set `DEBUG=false`
2. Configure SMTP credentials (Gmail, SendGrid, Mailgun, etc.)
3. Ensure `APP_URL` points to your production URL

## Project Structure

```
vabackpay/
├── main.go              # Application entry point
├── .env                 # Local configuration (gitignored)
├── example.env          # Configuration template
├── go.mod              # Go module definition
├── static/
│   └── style.css       # Stylesheets
├── templates/
│   ├── index.html      # Main page template
│   └── results.html    # Results partial template
└── README.md
```

## VA Disability Rates

The calculator includes historical VA disability compensation rates from 2000-2026. Rates are based on official VA compensation tables and include annual COLA adjustments.

### Dependent Add-ons

The calculator uses simplified dependent calculations:
- **30%-60% ratings**: $109/month per dependent
- **70%-100% ratings**: $153/month per dependent

Note: This is a simplified calculation. Actual VA rates may vary based on dependent status (spouse, children, parents, Aid & Attendance).

## Deployment

### Cloud Run (Recommended)

See [DEPLOY_CLOUDRUN.md](./DEPLOY_CLOUDRUN.md) for detailed instructions.

Quick start:
```bash
gcloud run deploy va-backpay-calc \
  --source . \
  --region us-central1 \
  --allow-unauthenticated
```

### Docker

```bash
docker build -t vabackpay .
docker run -p 8080:8080 vabackpay
```

## License

MIT License

## Disclaimer

This calculator provides estimates based on publicly available VA compensation rates. Actual backpay amounts may vary. This tool is for informational purposes only and should not be considered financial or legal advice. Consult with a VA-accredited representative for accurate calculations.
