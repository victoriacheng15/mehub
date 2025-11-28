{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  packages = with pkgs; [
    nodejs_24
  ];

  # Optional: set environment variables
  TF_IN_AUTOMATION = "1";

  # Optional: show shell info
  shellHook = ''
    echo "Node.js:   $(node --version)"
    echo "npm:       $(npm --version)"
  '';
}