package constants

import (
	"regexp"
	"time"
)

const LocalSource = "local"
const LocalFreeFormPath = "./configs/config.yaml"
const RemoteFreeformProfile = "config"
const DeployedEnv = "env"
const CloudwatchNamespace = "AppConfig"
const CloudwatchErrorMetric = "ConfigurationError"
const CloudwatchErrorDimension = "Application"
const CloudwatchPutMetricInterval = 1 * time.Minute
const HeaderXOmnifulRequestID = "X-Omniful-Request-Id"
const HeaderXOmnifulCorrelationID = "X-Omniful-Correlation-ID"
const HeaderXClientService = "X-Client-Service"
const HeaderUserAgent = "User-Agent"
const HeaderAuthorization = "Authorization"
const Env = "env"
const Consistency = "consistency"
const DBPreference = "db_preference"
const SlaveDB = "slave_db"
const EventualConsistency = "eventual"
const StrongConsistency = "strong"
const Config = "config"
const DefaultTimeout = 1 * time.Second
const AuthorizationTokenPrefix = "Token"
const HealthCheck = "Health Check"
const Healthy = "healthy"
const Unhealthy = "unhealthy"
const Newrelic = "newrelic"
const JWTHeader = "x-jwt-header"
const JWTUserDetails = "jwt_user_details"
const PrivateUserDetails = "private_user_details"
const PublicUserDetails = "public_user_details"
const PublicCustomerDetails = "public_customer_details"
const PublicUserPlan = "public_user_plan"
const PrivateUserPlan = "private_user_plan"
const Limit = "limit"
const Page = "page"
const AllMessageAttributes = "All"
const KafkaConsumerTransactionName = "kafka_consumer_name"
const ProfilingEnabledKey = "profiling.enabled"
const Omniful = "Omniful"
const Compression = "x-compression"
const AttributeTenantID = "attr-omniful-tenant-id"
const AttributeTenantName = "attr-omniful-tenant-name"
const AttributeTenantEnvironment = "attr-omniful-tenant-env"
const AttributeUserID = "attr-omniful-user-id"
const ConfigPath = "CONFIG_PATH"

// Regex Patterns
var (
	// Pattern to match HTML script tags and their content
	scriptTagPattern = `<script[\s\S]*?>[\s\S]*?</script>`
	// Pattern to extract language code from localization files
	localizationPattern = `messages\.(?P<lang_code>.*)\.json`
)

// Regex Compilations
var (
	ScriptTagRegex    = regexp.MustCompile(scriptTagPattern)
	LocalizationRegex = regexp.MustCompile(localizationPattern)
)
