package config

type BucketConfig struct {
	Capacity  float64 `yaml:"capacity"`
	RefillRPS float64 `yaml:"refill_rps"`
}

type RateLimiterConfig struct {
	DefaultBucket BucketConfig `yaml:"default_bucket"`
}
