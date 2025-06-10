{
  description = "gysmo flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "gysmo";
          version = "0.2.2";
          src = self;
          vendorHash = null;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = [ pkgs.go ];
        };

        apps.default = {
          type = "app";
          program = "${self.packages.default}/bin/gysmo";
        };
      }
    );
}
