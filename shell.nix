{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go
  ];

  # Set GOPATH to a writable location (like a 'go' directory in the project)
  GOPATH = "${toString ./go}";

  GO111MODULE = "on";

  shellHook = ''
    if [ ! -f ".env" ]; then
      echo "Copying .env..."
      cp  ../shared-resources/.env .
    fi
    export $(grep -v '^#' .env | xargs)
    echo "Go development environment ready"

    exec zsh
  '';
}

