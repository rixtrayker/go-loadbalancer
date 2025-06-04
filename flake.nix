{
  description = "Go Load Balancer Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_24
            gopls
            go-outline
            gocode
            gopkgs
            godef
            golint
            delve
          ];

          shellHook = ''
            export GO111MODULE=on
            export GOPATH=$HOME/go
            export PATH=$GOPATH/bin:$PATH
          '';
        };

        packages.default = pkgs.buildGoModule {
          pname = "go-loadbalancer";
          version = "0.1.0";
          src = ./.;
          vendorHash = self.inputs.flake-utils.lib.flakeVendorHash ./.;

          meta = with pkgs.lib; {
            description = "A Go-based load balancer";
            homepage = "https://github.com/amr/go-loadbalancer";
            license = licenses.mit;
            maintainers = with maintainers; [ ];
          };
        };
      });

  nixConfig = {
    extra-substituters = [
      "https://cache.nixos.org"
      "https://nix-community.cachix.org"
    ];
    extra-trusted-public-keys = [
      "cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY="
      "nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs="
    ];
  };
} 