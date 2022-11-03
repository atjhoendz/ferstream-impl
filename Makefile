publisher:
ifndef type
	@go run main.go publisher
endif
ifdef type
	@go run main.go publisher --type $(type)
endif


subscriber:
	@go run main.go subscriber

.PHONY: publisher subscriber