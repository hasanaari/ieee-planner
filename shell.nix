{pkgs ? import <nixpkgs> {}}:
pkgs.mkShell rec {
  allowUnfree = true;

  buildInputs = with pkgs; [
    bashInteractive
    nodejs
    nodePackages.pnpm
  ];
}
