{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    gnumake
  ];

  # Use a writable GOPATH location in the home directory or project directory
  # You can specify a specific folder in your project as needed
  shellHook = ''
    export GOPATH="$HOME/go"  # Or "./go" if you prefer it in the project directory
    export GO111MODULE="on"

    # Copy .env if it doesn't exist
    if [ ! -f ".env" ]; then
      echo "Copying .env..."
      cp ../shared-resources/.env .
    fi

    # Load variables from .env file into the shell environment
    set -o allexport
    source .env || true
    set +o allexport

    echo "Go development environment ready"
    
    # Start zsh shell
    exec zsh
  '';
}

