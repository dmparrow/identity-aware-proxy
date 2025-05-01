package auth

import "os"

func IsDev() bool {
    return os.Getenv("IAP_ENV") == "dev"
}

func IsProd() bool {
    return !IsDev()
}