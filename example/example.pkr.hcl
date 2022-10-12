packer {
  required_plugins {
    playground = {
      version = ">= 0.0.1"
      source  = "github.com/takumin/playground"
    }
  }
}

source "playground" "example" {
  rootfs_urls     = ["https://cdimage.ubuntu.com/ubuntu-base/jammy/daily/current/jammy-base-amd64.tar.gz"]
  rootfs_checksum = "file:https://cdimage.ubuntu.com/ubuntu-base/jammy/daily/current/SHA256SUMS"
}

build {
  sources = ["source.playground.example"]
}
