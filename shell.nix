{pkgs ? import <nixpkgs> {}}:
pkgs.mkShell rec {
  allowUnfree = true;

  buildInputs = with pkgs; [
    zsh
    bashInteractive
    nodejs
    nodePackages.pnpm
    docker
  ];
}
