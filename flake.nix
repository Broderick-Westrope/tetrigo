{
  description = "Play Tetris in your terminal.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    utils = {
      url = "github:NewDawn0/nixUtils";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, utils, ... }: {
    overlays.default = final: prev: {
      tetrigo = self.packages.${prev.system}.default;
    };

    packages = utils.lib.eachSystem { } (pkgs:
      let
        # Build Go binary
        builder = pkgs.buildGoModule {
          name = "tetrigo";
          src = ./.;
          vendorHash = "sha256-wXafPSj64mUa/Onim4zunUO5mrguUh9BylTuOG0MfTk=";
        };
        # Config TOML helpers
        mkCfg = {
          fromToml =
            builtins.fromTOML (builtins.readFile ./example.config.toml);
          toToml = config:
            (pkgs.formats.toml { }).generate "config.toml"
            (mkCfg.fromToml // config);
        };
        # Static helper package derivation
        tetrigoHelper = { config, dbPath }:
          pkgs.stdenvNoCC.mkDerivation {
            name = "tetrigo-helper";
            src = null;
            dontUnpack = true;
            runtimeInputs = [ pkgs.sqlite ];
            buildPhase = if builtins.isPath dbPath then
              ''cp "${dbPath}" tetrigo.db''
            else
              ''${pkgs.sqlite}/bin/sqlite3 tetrigo.db ""'';
            installPhase = ''
              install -D     ${builder}/bin/tetrigo $out/bin/tetrigo
              install -Dm644 tetrigo.db             $out/share/tetrigo.db
              install -Dm644 ${mkCfg.toToml config} $out/share/config.toml
            '';
          };
        # Final runtime wrapper
        tetrigoWrapped = { config ? mkCfg.fromToml, dbPath ? null }:
          let t = tetrigoHelper { inherit config dbPath; };
          in pkgs.writeShellApplication {
            name = "tetrigo";
            runtimeInputs = [ pkgs.coreutils pkgs.sqlite ];
            text = ''
              TETRIGO_HOME="$HOME/.config/tetrigo"
              mkdir -p "$TETRIGO_HOME"
              [ -f "$TETRIGO_HOME/tetrigo.db" ] || install -Dm644 "${t}/share/tetrigo.db" "$TETRIGO_HOME/tetrigo.db"
              ${t}/bin/tetrigo --config="${t}/share/config.toml" --db="$TETRIGO_HOME/tetrigo.db"
            '';
            meta = {
              description = "Play Tetris in your terminal.";
              longDescription = ''
                A Golang implementation of Tetris, following the official 2009 Tetris Design Guideline.
              '';
              homepage = "https://github.com/Broderick-Westrope/tetrigo";
              license = pkgs.lib.licenses.gpl3;
              maintainers = [ "Broderick-Westrope" ];
            };
          };
      in {
        default = pkgs.lib.makeOverridable tetrigoWrapped {
          config = mkCfg.fromToml;
          dbPath = null;
        };
      });
  };
}
