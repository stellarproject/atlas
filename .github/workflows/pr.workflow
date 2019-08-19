workflow "PR" {
  on = "pull-request"
}

action "Set up Go" {
  id = "go"
  uses = "actions/setup-go@v1"
}

action "Check out code" {
  uses = "actions/checkout@v1"
}

action "Test" {
  runs = "go test -cover -v ./..."
}
