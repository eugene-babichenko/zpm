workflow "Test" {
  on = "push"
  resolves = ["GoTests"]
}

workflow "Release" {
  on = "push"
  resolves = ["GoReleaser"]
}

action "IsTag" {
  uses = "actions/bin/filter@master"
  args = "tag"
}

action "GoTests" {
  uses = "docker://golang"
  env = {
    GO111MODULE = "on"
  }
  runs = "go test ./..."
}

action "GoReleaser" {
  uses = "docker://goreleaser/goreleaser"
  needs = ["GoTests", "IsTag"]
  args = "release"
  secrets = ["GITHUB_TOKEN"]
  env = {
    GO111MODULE = "on"
  }
}
