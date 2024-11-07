.PHONY: help
help: # Lists commands
	@awk 'BEGIN { print "\033[32mAvaiable Commands:\033[0m"; } \
		 /^##/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 3); next; } \
		 /^[a-zA-Z0-9_-]+:/ { \
			 split($$0, parts, /:.*#/); \
			 cmd = parts[1]; \
			 sub(/^[ \t]+/, "", cmd); \
			 desc = substr($$0, index($$0, "#") + 1); \
			 if (desc != "") \
				 printf "\033[36m%-30s\033[0m %s\n", cmd, desc; \
		 }' $(MAKEFILE_LIST)



.PHONY: t
t: # Runs the golang testing suite 
	go test ./... -cover

.PHONY: b
b: # Builds the golang app
	go build .  

.PHONY: r
r: # Runs the golang app 
	go run .  
