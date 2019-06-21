workflow "CI" {
  resolves = ["Go Tests"]
  on = "push"
}

action "Go Tests" {
  uses = "docker://golang"
  env = {
    GO111MODULE = "on"
  }
  runs = "go test ./..."
}
