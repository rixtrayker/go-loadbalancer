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
            go_1_21
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

          vendorHash = "sha256-0000000000000000000000000000000000000000000000000000";

          meta = with pkgs.lib; {
            description = "A Go-based load balancer";
            homepage = "https://github.com/amr/go-loadbalancer";
            license = licenses.mit;
            maintainers = with maintainers; [ ];
          };
        };
      });
} 