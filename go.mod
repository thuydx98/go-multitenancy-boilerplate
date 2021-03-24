module go-multitenancy-boilerplate

go 1.12

replace go-multitenancy-boilerplate => ../go-multitenancy-boilerplate

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/gorilla/sessions v1.2.1
	github.com/jinzhu/gorm v1.9.16
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/wader/gormstore v0.0.0-20210319162436-2b0cf73a0321
	golang.org/x/crypto v0.0.0-20191205180655-e7c4368fe9dd
)
