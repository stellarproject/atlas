workflow "Build and Release" {
  on = "push"
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

action "Build and Package" {
  needs = ["Build"]
  secrets = ["AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"]
  env = {
    BUILD = ""
  }
  runs = "./script/release.sh"
}
