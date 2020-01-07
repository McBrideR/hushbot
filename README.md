# hushbot
List of errors that occured while trying to push this bot to cloud foundary.

  Error: Bot failed to push to cloud foundary because of dependency issues.
  Solution: install go's dependency management tool (go dep) and get the depedencies.
  1. 'go get -u github.com/golang/dep/cmd/dep'
  2. 'godep save'

  Error: Bot fails to start with 'exit description: Instance never healthy after 1m0s: Failed to make TCP connection to port 8080: connection refused'
  Solution: add --no-route to your cf push command ie: cf push hushbot --no-route or put no-route:true in the manifest.yml file.



