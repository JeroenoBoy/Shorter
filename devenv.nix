{ pkgs, lib, config, inputs, ... }:

let
  postgres_address = "127.0.0.1";
  postgres_port = 5432;
  postgres_database = "shorter";
  postgres_ssl = false;
in
{
  cachix.enable = false;

  # https://devenv.sh/basics/
  env.GREET = "Shorter";
  env.POSTGRES_ADDRESS = postgres_address;
  env.POSTGRES_PORT = postgres_port;
  env.POSTGRES_DATABASE = postgres_database;
  env.POSTGRES_USER = "root";
  env.POSTGRES_USESSL = postgres_ssl;

  # https://devenv.sh/packages/
  packages = with pkgs; [ templ ];

  # https://devenv.sh/scripts/
  scripts.hello.exec = "echo hello from $GREET";

  enterShell = ''
    hello
    git --version
  '';

  # https://devenv.sh/tests/
  enterTest = ''
    echo "Running tests"
    git --version | grep "2.42.0"
  '';

  # https://devenv.sh/services/
  # services.postgres.enable = true;
  services = {
    postgres = {
      enable = true;
      initialDatabases = [{ name = postgres_database; }];
    };
  };

  # https://devenv.sh/languages/
  # languages.nix.enable = true;
  languages.go.enable = true;

  # https://devenv.sh/pre-commit-hooks/
  pre-commit.hooks.shellcheck.enable = true;

  # https://devenv.sh/processes/
  # processes.ping.exec = "ping example.com";

  # See full reference at https://devenv.sh/reference/options/
}
